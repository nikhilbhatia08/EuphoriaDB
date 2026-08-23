package transactions

import (
	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
)

type CheckpointRecord struct {
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
	intSize := getIntSize()
	record := make([]byte, intSize)

	page := filemgr.NewPageFromBytes(record)
	page.SetInt(0, int(Checkpoint))

	return logManager.Append(record)
}

func (ck *CheckpointRecord) TxNumber() int {
	return -1
}

func (ck *CheckpointRecord) Undo(tx *Transaction) error {
	// Checkpoint has no undo
	return nil
}
