//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package marsh

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// MAV FTP error codes (https://mavlink.io/en/services/ftp.html)
type MAV_FTP_ERR = common.MAV_FTP_ERR

const (
	// None: No error
	MAV_FTP_ERR_NONE MAV_FTP_ERR = common.MAV_FTP_ERR_NONE
	// Fail: Unknown failure
	MAV_FTP_ERR_FAIL MAV_FTP_ERR = common.MAV_FTP_ERR_FAIL
	// FailErrno: Command failed, Err number sent back in PayloadHeader.data[1].
	// This is a file-system error number understood by the server operating system.
	MAV_FTP_ERR_FAILERRNO MAV_FTP_ERR = common.MAV_FTP_ERR_FAILERRNO
	// InvalidDataSize: Payload size is invalid
	MAV_FTP_ERR_INVALIDDATASIZE MAV_FTP_ERR = common.MAV_FTP_ERR_INVALIDDATASIZE
	// InvalidSession: Session is not currently open
	MAV_FTP_ERR_INVALIDSESSION MAV_FTP_ERR = common.MAV_FTP_ERR_INVALIDSESSION
	// NoSessionsAvailable: All available sessions are already in use
	MAV_FTP_ERR_NOSESSIONSAVAILABLE MAV_FTP_ERR = common.MAV_FTP_ERR_NOSESSIONSAVAILABLE
	// EOF: Offset past end of file for ListDirectory and ReadFile commands
	MAV_FTP_ERR_EOF MAV_FTP_ERR = common.MAV_FTP_ERR_EOF
	// UnknownCommand: Unknown command / opcode
	MAV_FTP_ERR_UNKNOWNCOMMAND MAV_FTP_ERR = common.MAV_FTP_ERR_UNKNOWNCOMMAND
	// FileExists: File/directory already exists
	MAV_FTP_ERR_FILEEXISTS MAV_FTP_ERR = common.MAV_FTP_ERR_FILEEXISTS
	// FileProtected: File/directory is write protected
	MAV_FTP_ERR_FILEPROTECTED MAV_FTP_ERR = common.MAV_FTP_ERR_FILEPROTECTED
	// FileNotFound: File/directory not found
	MAV_FTP_ERR_FILENOTFOUND MAV_FTP_ERR = common.MAV_FTP_ERR_FILENOTFOUND
)
