package gomavlib

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

const (
	commandLongID  = 76
	commandIntID   = 75
	commandAckID   = 77
	commandLongCRC = 152
	commandIntCRC  = 158
	commandAckCRC  = 143

	defaultCommandTimeout = 5 * time.Second
)

var (
	//nolint:revive
	ErrCommandTimeout = errors.New("command timeout")
	//nolint:revive
	ErrCommandRejected = errors.New("command rejected")
	//nolint:revive
	ErrCommandFailed = errors.New("command failed")
	//nolint:revive
	ErrCommandUnsupported = errors.New("command unsupported")
	//nolint:revive
	ErrCommandDenied = errors.New("command denied")
	//nolint:revive
	ErrCommandCancelled = errors.New("command cancelled")
	//nolint:revive
	ErrNodeTerminated = errors.New("node terminated")
)

// CommandOptions configures command sending behavior.
type CommandOptions struct {
	// Channel to send the command on.
	Channel *Channel

	// Timeout for waiting for COMMAND_ACK response.
	// Defaults to 5 seconds if not specified or zero.
	Timeout time.Duration

	// OnProgress is called when MAV_RESULT_IN_PROGRESS ACKs are received.
	// The progress value ranges from 0-100, or 255 if unknown.
	OnProgress func(progress uint8)
}

// CommandResponse contains the response from a command.
type CommandResponse struct {
	// Result is the MAV_RESULT enum value
	Result uint64

	// Progress percentage (0-100, or 255 if unknown)
	Progress uint8

	// ResultParam2 contains additional error information (if any)
	ResultParam2 int32

	// ResponseTime is how long it took to get the response
	ResponseTime time.Duration

	// InProgress indicates if this was an intermediate progress update
	InProgress bool
}

// commandKey uniquely identifies a pending command
type commandKey struct {
	channel      *Channel
	targetSystem uint8
	targetComp   uint8
	commandID    uint32
}

// pendingCommand tracks a command awaiting response
type pendingCommand struct {
	key        commandKey
	ctx        context.Context
	cancel     context.CancelFunc
	responseCh chan *CommandResponse
	progressCh chan uint8
	sentAt     time.Time
	timeout    time.Duration
}

// sendCommandReq is an internal request to send a command
type sendCommandReq struct {
	channel      *Channel
	msg          message.Message
	targetSystem uint8
	targetComp   uint8
	commandID    uint32
	timeout      time.Duration
	progressCh   chan uint8
	responseCh   chan *CommandResponse
	errorCh      chan error
}

// nodeCommand manages the command protocol
type nodeCommand struct {
	node *Node

	// Message templates from dialect (validated during init)
	msgCommandLong message.Message
	msgCommandInt  message.Message
	msgCommandAck  message.Message

	// State management
	pendingMutex sync.RWMutex
	pending      map[commandKey]*pendingCommand

	// Control channels
	chRequest chan *sendCommandReq
	terminate chan struct{}
	done      chan struct{}
}

// initialize validates dialect and sets up the command manager
func (nc *nodeCommand) initialize() error {
	// Dialect must be present
	if nc.node.Dialect == nil {
		return errSkip
	}

	// Find and validate COMMAND_LONG message
	nc.msgCommandLong = func() message.Message {
		for _, m := range nc.node.Dialect.Messages {
			if m.GetID() == commandLongID {
				return m
			}
		}
		return nil
	}()
	if nc.msgCommandLong == nil {
		return errSkip
	}

	// Validate COMMAND_LONG CRC
	rwLong := &message.ReadWriter{Message: nc.msgCommandLong}
	err := rwLong.Initialize()
	if err != nil || rwLong.CRCExtra() != commandLongCRC {
		return errSkip
	}

	// Find and validate COMMAND_INT message
	nc.msgCommandInt = func() message.Message {
		for _, m := range nc.node.Dialect.Messages {
			if m.GetID() == commandIntID {
				return m
			}
		}
		return nil
	}()
	if nc.msgCommandInt == nil {
		return errSkip
	}

	// Validate COMMAND_INT CRC
	rwInt := &message.ReadWriter{Message: nc.msgCommandInt}
	err = rwInt.Initialize()
	if err != nil || rwInt.CRCExtra() != commandIntCRC {
		return errSkip
	}

	// Find and validate COMMAND_ACK message
	nc.msgCommandAck = func() message.Message {
		for _, m := range nc.node.Dialect.Messages {
			if m.GetID() == commandAckID {
				return m
			}
		}
		return nil
	}()
	if nc.msgCommandAck == nil {
		return errSkip
	}

	// Validate COMMAND_ACK CRC
	rwAck := &message.ReadWriter{Message: nc.msgCommandAck}
	err = rwAck.Initialize()
	if err != nil || rwAck.CRCExtra() != commandAckCRC {
		return errSkip
	}

	// Initialize state
	nc.pending = make(map[commandKey]*pendingCommand)
	nc.chRequest = make(chan *sendCommandReq, 16)
	nc.terminate = make(chan struct{})
	nc.done = make(chan struct{})

	return nil
}

func (nc *nodeCommand) close() {
	close(nc.terminate)
	<-nc.done
}

func (nc *nodeCommand) run() {
	defer close(nc.done)
	defer nc.node.wg.Done()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case req := <-nc.chRequest:
			nc.handleSendCommand(req)

		case <-ticker.C:
			nc.cleanupTimeouts()

		case <-nc.terminate:
			nc.cancelAllPending()
			return
		}
	}
}

func (nc *nodeCommand) handleSendCommand(req *sendCommandReq) {
	key := commandKey{
		channel:      req.channel,
		targetSystem: req.targetSystem,
		targetComp:   req.targetComp,
		commandID:    req.commandID,
	}

	timeout := req.timeout
	if timeout == 0 {
		timeout = defaultCommandTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	pending := &pendingCommand{
		key:        key,
		ctx:        ctx,
		cancel:     cancel,
		responseCh: req.responseCh,
		progressCh: req.progressCh,
		sentAt:     time.Now(),
		timeout:    timeout,
	}

	nc.pendingMutex.Lock()
	nc.pending[key] = pending
	nc.pendingMutex.Unlock()

	// Send the command message
	err := nc.node.WriteMessageTo(req.channel, req.msg)
	if err != nil {
		nc.removePending(key)
		req.errorCh <- err
		cancel()
		return
	}

	req.errorCh <- nil

	// Wait for timeout in background
	go nc.waitForResponse(pending)
}

func (nc *nodeCommand) waitForResponse(pending *pendingCommand) {
	<-pending.ctx.Done()

	// Only send timeout response if it was actually a timeout, not a cancellation
	if err := pending.ctx.Err(); err == context.DeadlineExceeded {
		// Check if still pending (might have been removed by successful response)
		nc.pendingMutex.Lock()
		_, stillPending := nc.pending[pending.key]
		if stillPending {
			delete(nc.pending, pending.key)
		}
		nc.pendingMutex.Unlock()

		// Only send timeout response if it was still pending
		if stillPending {
			select {
			case pending.responseCh <- &CommandResponse{
				Result:       0, // Could use a specific timeout result
				ResponseTime: time.Since(pending.sentAt),
			}:
			default:
			}
		}
	}
}

func (nc *nodeCommand) cleanupTimeouts() {
	now := time.Now()

	nc.pendingMutex.Lock()
	defer nc.pendingMutex.Unlock()

	for key, pending := range nc.pending {
		if now.Sub(pending.sentAt) > pending.timeout {
			delete(nc.pending, key)
			pending.cancel()
		}
	}
}

func (nc *nodeCommand) removePending(key commandKey) {
	nc.pendingMutex.Lock()
	defer nc.pendingMutex.Unlock()

	if pending, ok := nc.pending[key]; ok {
		pending.cancel()
		delete(nc.pending, key)
	}
}

func (nc *nodeCommand) cancelAllPending() {
	nc.pendingMutex.Lock()
	defer nc.pendingMutex.Unlock()

	for key, pending := range nc.pending {
		pending.cancel()
		delete(nc.pending, key)
	}
}

// onEventFrame is called when a COMMAND_ACK is received
func (nc *nodeCommand) onEventFrame(evt *EventFrame) {
	if evt.Message().GetID() != commandAckID {
		return
	}

	// Extract fields using reflection (dialect-agnostic)
	ackMsg := reflect.ValueOf(evt.Message()).Elem()

	commandID := uint32(ackMsg.FieldByName("Command").Uint())
	result := ackMsg.FieldByName("Result").Uint()

	// Extension fields (may not exist in V1)
	progress := uint8(255)
	resultParam2 := int32(0)

	if field := ackMsg.FieldByName("Progress"); field.IsValid() {
		progress = uint8(field.Uint())
	}
	if field := ackMsg.FieldByName("ResultParam2"); field.IsValid() {
		resultParam2 = int32(field.Int())
	}

	// Match based on who SENT the ACK (should be who we sent the command TO)
	key := commandKey{
		channel:      evt.Channel,
		targetSystem: evt.SystemID(),    // Who sent the ACK
		targetComp:   evt.ComponentID(), // Who sent the ACK
		commandID:    commandID,
	}

	nc.pendingMutex.RLock()
	pending, exists := nc.pending[key]
	nc.pendingMutex.RUnlock()

	if !exists {
		return
	}

	response := &CommandResponse{
		Result:       result,
		Progress:     progress,
		ResultParam2: resultParam2,
		ResponseTime: time.Since(pending.sentAt),
		InProgress:   result == 5, // MAV_RESULT_IN_PROGRESS
	}

	if response.InProgress {
		// Send progress update
		if pending.progressCh != nil {
			select {
			case pending.progressCh <- progress:
			default:
			}
		}
		return
	}

	// Final result
	nc.removePending(key)

	select {
	case pending.responseCh <- response:
	default:
	}
}
