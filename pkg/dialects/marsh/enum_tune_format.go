//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package marsh

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Tune formats (used for vehicle buzzer/tone generation).
type TUNE_FORMAT = common.TUNE_FORMAT

const (
	// Format is QBasic 1.1 Play: https://www.qbasic.net/en/reference/qb11/Statement/PLAY-006.htm.
	TUNE_FORMAT_QBASIC1_1 TUNE_FORMAT = common.TUNE_FORMAT_QBASIC1_1
	// Format is Modern Music Markup Language (MML): https://en.wikipedia.org/wiki/Music_Macro_Language#Modern_MML.
	TUNE_FORMAT_MML_MODERN TUNE_FORMAT = common.TUNE_FORMAT_MML_MODERN
)
