//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package development

import (
	"fmt"
	"strconv"
)

// Possible parameter transaction actions.
type PARAM_TRANSACTION_ACTION uint32

const (
	// Commit the current parameter transaction.
	PARAM_TRANSACTION_ACTION_START PARAM_TRANSACTION_ACTION = 0
	// Commit the current parameter transaction.
	PARAM_TRANSACTION_ACTION_COMMIT PARAM_TRANSACTION_ACTION = 1
	// Cancel the current parameter transaction.
	PARAM_TRANSACTION_ACTION_CANCEL PARAM_TRANSACTION_ACTION = 2
)

var labels_PARAM_TRANSACTION_ACTION = map[PARAM_TRANSACTION_ACTION]string{
	PARAM_TRANSACTION_ACTION_START:  "PARAM_TRANSACTION_ACTION_START",
	PARAM_TRANSACTION_ACTION_COMMIT: "PARAM_TRANSACTION_ACTION_COMMIT",
	PARAM_TRANSACTION_ACTION_CANCEL: "PARAM_TRANSACTION_ACTION_CANCEL",
}

var values_PARAM_TRANSACTION_ACTION = map[string]PARAM_TRANSACTION_ACTION{
	"PARAM_TRANSACTION_ACTION_START":  PARAM_TRANSACTION_ACTION_START,
	"PARAM_TRANSACTION_ACTION_COMMIT": PARAM_TRANSACTION_ACTION_COMMIT,
	"PARAM_TRANSACTION_ACTION_CANCEL": PARAM_TRANSACTION_ACTION_CANCEL,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e PARAM_TRANSACTION_ACTION) MarshalText() ([]byte, error) {
	if name, ok := labels_PARAM_TRANSACTION_ACTION[e]; ok {
		return []byte(name), nil
	}
	return []byte(strconv.Itoa(int(e))), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *PARAM_TRANSACTION_ACTION) UnmarshalText(text []byte) error {
	if value, ok := values_PARAM_TRANSACTION_ACTION[string(text)]; ok {
		*e = value
	} else if value, err := strconv.Atoi(string(text)); err == nil {
		*e = PARAM_TRANSACTION_ACTION(value)
	} else {
		return fmt.Errorf("invalid label '%s'", text)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e PARAM_TRANSACTION_ACTION) String() string {
	if name, ok := labels_PARAM_TRANSACTION_ACTION[e]; ok {
		return name
	}
	return strconv.Itoa(int(e))
}
