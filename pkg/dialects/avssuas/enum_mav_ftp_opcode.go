//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package avssuas

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// MAV FTP opcodes: https://mavlink.io/en/services/ftp.html
type MAV_FTP_OPCODE = common.MAV_FTP_OPCODE

const (
	// None. Ignored, always ACKed
	MAV_FTP_OPCODE_NONE MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_NONE
	// TerminateSession: Terminates open Read session
	MAV_FTP_OPCODE_TERMINATESESSION MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_TERMINATESESSION
	// ResetSessions: Terminates all open read sessions
	MAV_FTP_OPCODE_RESETSESSION MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_RESETSESSION
	// ListDirectory. List files and directories in path from offset
	MAV_FTP_OPCODE_LISTDIRECTORY MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_LISTDIRECTORY
	// OpenFileRO: Opens file at path for reading, returns session
	MAV_FTP_OPCODE_OPENFILERO MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_OPENFILERO
	// ReadFile: Reads size bytes from offset in session
	MAV_FTP_OPCODE_READFILE MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_READFILE
	// CreateFile: Creates file at path for writing, returns session
	MAV_FTP_OPCODE_CREATEFILE MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_CREATEFILE
	// WriteFile: Writes size bytes to offset in session
	MAV_FTP_OPCODE_WRITEFILE MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_WRITEFILE
	// RemoveFile: Remove file at path
	MAV_FTP_OPCODE_REMOVEFILE MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_REMOVEFILE
	// CreateDirectory: Creates directory at path
	MAV_FTP_OPCODE_CREATEDIRECTORY MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_CREATEDIRECTORY
	// RemoveDirectory: Removes directory at path. The directory must be empty.
	MAV_FTP_OPCODE_REMOVEDIRECTORY MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_REMOVEDIRECTORY
	// OpenFileWO: Opens file at path for writing, returns session
	MAV_FTP_OPCODE_OPENFILEWO MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_OPENFILEWO
	// TruncateFile: Truncate file at path to offset length
	MAV_FTP_OPCODE_TRUNCATEFILE MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_TRUNCATEFILE
	// Rename: Rename path1 to path2
	MAV_FTP_OPCODE_RENAME MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_RENAME
	// CalcFileCRC32: Calculate CRC32 for file at path
	MAV_FTP_OPCODE_CALCFILECRC MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_CALCFILECRC
	// BurstReadFile: Burst download session file
	MAV_FTP_OPCODE_BURSTREADFILE MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_BURSTREADFILE
	// ACK: ACK response
	MAV_FTP_OPCODE_ACK MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_ACK
	// NAK: NAK response
	MAV_FTP_OPCODE_NAK MAV_FTP_OPCODE = common.MAV_FTP_OPCODE_NAK
)
