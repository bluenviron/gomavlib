//autogenerated:yes
//nolint:revive,misspell,govet,lll
package avssuas

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Data for filling the OpenDroneID Authentication message. The Authentication Message defines a field that can provide a means of authenticity for the identity of the UAS (Unmanned Aircraft System). The Authentication message can have two different formats. For data page 0, the fields PageCount, Length and TimeStamp are present and AuthData is only 17 bytes. For data page 1 through 15, PageCount, Length and TimeStamp are not present and the size of AuthData is 23 bytes.
type MessageOpenDroneIdAuthentication = common.MessageOpenDroneIdAuthentication
