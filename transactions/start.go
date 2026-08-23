package transactions

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
)

type StartRecord struct {
	TxNum int
}

func NewStartRecord(page *filemgr.Page) (*StartRecord, error) {
	txPos := getIntSize()
	txNum := page.GetInt(txPos)

	return &StartRecord{
		TxNum: txNum,
	}, nil
}

func (sr *StartRecord) Op() LogRecordType {
	return Start
}

func (sr *StartRecord) String() string {
	return fmt.Sprintf("<START %d>", sr.TxNum)
}

func (sr *StartRecord) TxNumber() int {
	return sr.TxNum
}

func (sr *StartRecord) Undo(tx *Transaction) error {
	// Start record has no undo
	return nil
}

func WriteStartToLog(logmgr *log.LogManager, txNum int) (int, error) {
	intSize := getIntSize()
	record := make([]byte, 2*intSize)
	page := filemgr.NewPageFromBytes(record)
	page.SetInt(0, int(Start))
	page.SetInt(intSize, txNum)

	return logmgr.Append(record)
}
