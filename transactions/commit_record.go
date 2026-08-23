package transactions

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
)

type CommitRecord struct {
	TxNum int
}

func NewCommitRecord(page *filemgr.Page) (*CommitRecord, error) {
	txNumPos := getIntSize()
	txNum := page.GetInt(txNumPos)

	return &CommitRecord{
		TxNum: txNum,
	}, nil
}

func (cr *CommitRecord) Op() LogRecordType {
	return Commit
}

func (cr *CommitRecord) String() string {
	return fmt.Sprintf("<COMMIT %d>", cr.TxNum)
}

func (cr *CommitRecord) TxNumber() int {
	return cr.TxNum
}

func (cr *CommitRecord) Undo(tx *Transaction) error {
	// Commit records have no undo action
	return nil
}

func WriteCommitRecordToLog(logMgr *log.LogManager, txNum int) (int, error) {
	intSize := getIntSize()
	record := make([]byte, 2*intSize)

	page := filemgr.NewPageFromBytes(record)
	page.SetInt(0, int(Commit))
	page.SetInt(intSize, txNum)

	return logMgr.Append(record)
}
