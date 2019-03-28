package gomavlib

import (
	"time"
)

type nodeChecker struct {
	n         *Node
	terminate chan struct{}
	done      chan struct{}
}

func newNodeChecker(n *Node) *nodeChecker {
	h := &nodeChecker{
		n:         n,
		terminate: make(chan struct{}, 1),
		done:      make(chan struct{}),
	}
	go h.do()
	return h
}

func (h *nodeChecker) close() {
	h.terminate <- struct{}{}
	<-h.done
}

func (h *nodeChecker) do() {
	defer func() { h.done <- struct{}{} }()

	ticker := time.NewTicker(nodeCheckPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			func() {
				h.n.nodeMutex.Lock()
				defer h.n.nodeMutex.Unlock()
				for i, t := range h.n.nodes {
					// delete nodes after a period of inactivity
					if time.Since(t) >= nodeDisappearAfter {
						delete(h.n.nodes, i)
						h.n.eventChan <- &EventNodeDisappear{i}
					}
				}
			}()

		case <-h.terminate:
			return
		}
	}
}
