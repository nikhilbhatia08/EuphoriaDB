package transactions

import (
	"errors"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
)

type LogRecordType int

const (
	Checkpoint LogRecordType = iota
	Start
	Commit
	Rollback
	SetInt
	SetString
	SetBool
	SetDate
)

type LogRecord interface {
	Op() LogRecordType
	Undo(*Transaction) error
	TxNumber() int
}

func CreateLogRecord(bytes []byte) (LogRecord, error) {
	page := filemgr.NewPageFromBytes(bytes)
	switch page.GetInt(0) {
	case int(Checkpoint):
		return NewCheckpointRecord()
	case int(Start):
		return NewStartRecord(page)
	case int(Commit):
		return NewCommitRecord(page)
	case int(Rollback):
		return NewRollbackRecord(page)
	case int(SetString):
		return NewSetStringRecord(page)
	case int(SetInt):
		return NewSetIntRecord(page)
	case int(SetBool):
		return NewSetBoolRecord(page)
	case int(SetDate):
		return NewSetDateRecord(page)
	default:
		return nil, errors.New("LogRecordType not present with such value")
	}
}
