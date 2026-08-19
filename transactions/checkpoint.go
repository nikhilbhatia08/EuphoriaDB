package transactions

import (
	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
)

type CheckpointRecord struct {
	LogRecord
}

func NewCheckpointRecord() (*CheckpointRecord, error) {
	return &CheckpointRecord{}, nil
}

func (ck *CheckpointRecord) Op() LogRecordType {
	return Checkpoint
}

func (ck *CheckpointRecord) String() string {
	return "<CHECKPOINT>"
}

func WriteCheckpointRecordToLog(logManager *log.LogManager) (int, error) {
	record := make([]byte, 4)

	page := filemgr.NewPageFromBytes(record)
	page.SetInt(0, int(Checkpoint))

	return logManager.Append(record)
}
