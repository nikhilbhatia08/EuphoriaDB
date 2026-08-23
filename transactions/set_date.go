package transactions

import (
	"fmt"
	"time"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
)

type SetDateRecord struct {
	TxNum  int
	offset int
	value  time.Time
	block  *filemgr.BlockId
}

func NewSetDateRecord(page *filemgr.Page) (*SetDateRecord, error) {
	intSize := getIntSize()
	txNumPos := intSize
	txNum := page.GetInt(txNumPos)

	fileNamePos := txNumPos + intSize
	fileName, err := page.GetString(fileNamePos)
	if err != nil {
		return nil, err
	}

	blockNumPos := fileNamePos + filemgr.MaxLength(len(fileName))
	blockNum := page.GetInt(blockNumPos)
	block := &filemgr.BlockId{File: fileName, Blknum: int(blockNum)}

	offsetPos := blockNumPos + intSize
	offset := page.GetInt(offsetPos)

	valuePos := offsetPos + intSize
	val := page.GetDate(valuePos)

	return &SetDateRecord{TxNum: txNum, offset: offset, value: val, block: block}, nil
}

func (sd *SetDateRecord) Op() LogRecordType {
	return SetDate
}

func (sd *SetDateRecord) String() string {
	return fmt.Sprintf("<SETDATE %d %s %d %s>", sd.TxNum, sd.block, sd.offset, sd.value.String())
}

func (sd *SetDateRecord) Undo(tx *Transaction) error {
	if err := tx.Pin(sd.block); err != nil {
		return err
	}
	defer tx.Unpin(sd.block)

	return tx.SetDate(sd.block, sd.offset, sd.value, false)
}

func (sd *SetDateRecord) TxNumber() int {
	return sd.TxNum
}

func WriteSetDateToLog(logManager *log.LogManager, txNum int, block *filemgr.BlockId, offset int, val time.Time) (int, error) {
	intSize := getIntSize()

	txNumPos := intSize
	fileNamePos := txNumPos + intSize
	fileName := block.Filename()

	blockNumPos := fileNamePos + filemgr.MaxLength(len(fileName))
	blockNum := block.Blknum

	offsetPos := blockNumPos + intSize
	valuePos := offsetPos + intSize
	recordLen := valuePos + 8

	recordBytes := make([]byte, recordLen)
	page := filemgr.NewPageFromBytes(recordBytes)

	page.SetInt(0, int(SetDate))
	page.SetInt(txNumPos, txNum)
	if err := page.SetString(fileNamePos, fileName); err != nil {
		return -1, err
	}
	page.SetInt(blockNumPos, blockNum)
	page.SetInt(offsetPos, offset)
	page.SetDate(valuePos, val)

	return logManager.Append(recordBytes)
}
