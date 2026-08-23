package transactions

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
)

type RollbackRecord struct {
	TxNum int
}

func NewRollbackRecord(page *filemgr.Page) (*RollbackRecord, error) {
	txNumPos := getIntSize()
	txNum := page.GetInt(txNumPos)

	return &RollbackRecord{
		TxNum: txNum,
	}, nil
}

func (rr *RollbackRecord) String() string {
	return fmt.Sprintf("<ROLLBACK %d>", rr.TxNum)
}

func (rr *RollbackRecord) Op() LogRecordType {
	return Rollback
}

func (rr *RollbackRecord) TxNumber() int {
	return rr.TxNum
}

func (rr *RollbackRecord) Undo(tx *Transaction) error {
	// Rollback record has no undo action
	return nil
}

func WriteRollbackRecordToLog(logmgr *log.LogManager, txNum int) (int, error) {
	intSize := getIntSize()
	record := make([]byte, 2*intSize)

	page := filemgr.NewPageFromBytes(record)
	page.SetInt(0, int(Rollback))
	page.SetInt(intSize, txNum)

	return logmgr.Append(record)
}
