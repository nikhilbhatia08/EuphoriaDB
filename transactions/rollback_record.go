package transactions

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
)

type RollbackRecord struct {
	LogRecord
	TxNum int
}

func NewRollbackRecord(txNum int) (*RollbackRecord, error) {
	return &RollbackRecord{
		TxNum: txNum,
	}, nil
}

func (rr *RollbackRecord) String() string {
	return fmt.Sprintf("<ROLLBACK %d>", rr.TxNum)
}

func WriteRollbackRecordToLog(logmgr *log.LogManager, txNum int) (int, error) {
	intSize := getIntSize()
	record := make([]byte, 2*intSize)

	page := filemgr.NewPageFromBytes(record)
	page.SetInt(0, int(Rollback))
	page.SetInt(intSize, txNum)

	return logmgr.Append(record)
}
