//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
)

// Tune formats (used for vehicle buzzer/tone generation).
type TUNE_FORMAT uint32

const (
	// Format is QBasic 1.1 Play: https://www.qbasic.net/en/reference/qb11/Statement/PLAY-006.htm.
	TUNE_FORMAT_QBASIC1_1 TUNE_FORMAT = 1
	// Format is Modern Music Markup Language (MML): https://en.wikipedia.org/wiki/Music_Macro_Language#Modern_MML.
	TUNE_FORMAT_MML_MODERN TUNE_FORMAT = 2
)

var labels_TUNE_FORMAT = map[TUNE_FORMAT]string{
	TUNE_FORMAT_QBASIC1_1:  "TUNE_FORMAT_QBASIC1_1",
	TUNE_FORMAT_MML_MODERN: "TUNE_FORMAT_MML_MODERN",
}

var values_TUNE_FORMAT = map[string]TUNE_FORMAT{
	"TUNE_FORMAT_QBASIC1_1":  TUNE_FORMAT_QBASIC1_1,
	"TUNE_FORMAT_MML_MODERN": TUNE_FORMAT_MML_MODERN,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e TUNE_FORMAT) MarshalText() ([]byte, error) {
	if name, ok := labels_TUNE_FORMAT[e]; ok {
		return []byte(name), nil
	}
	return []byte(strconv.Itoa(int(e))), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *TUNE_FORMAT) UnmarshalText(text []byte) error {
	if value, ok := values_TUNE_FORMAT[string(text)]; ok {
		*e = value
	} else if value, err := strconv.Atoi(string(text)); err == nil {
		*e = TUNE_FORMAT(value)
	} else {
		return fmt.Errorf("invalid label '%s'", text)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e TUNE_FORMAT) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
