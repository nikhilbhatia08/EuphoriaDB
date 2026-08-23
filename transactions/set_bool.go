package transactions

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
)

type SetBoolRecord struct {
	TxNum  int
	offset int
	value  bool
	block  *filemgr.BlockId
}

func NewSetBoolRecord(page *filemgr.Page) (*SetBoolRecord, error) {
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
	block := &filemgr.BlockId{File: fileName, Blknum: blockNum}

	offsetPos := blockNumPos + intSize
	offset := page.GetInt(offsetPos)

	valuePos := offsetPos + intSize
	val := page.GetBool(valuePos)

	return &SetBoolRecord{TxNum: txNum, offset: offset, value: val, block: block}, nil
}

func (sr *SetBoolRecord) Op() LogRecordType {
	return SetBool
}

func (sr *SetBoolRecord) String() string {
	return fmt.Sprintf("<SETBOOL %d %s %d %t>", sr.TxNum, sr.block, sr.offset, sr.value)
}

func (sr *SetBoolRecord) Undo(tx *Transaction) error {
	if err := tx.Pin(sr.block); err != nil {
		return err
	}
	defer tx.Unpin(sr.block)

	return tx.SetBool(sr.block, sr.offset, sr.value, false)
}

func (sr *SetBoolRecord) TxNumber() int {
	return sr.TxNum
}

func WriteSetBoolToLog(logManager *log.LogManager, txNum int, block *filemgr.BlockId, offset int, val bool) (int, error) {
	intSize := getIntSize()

	txNumPos := intSize
	fileNamePos := txNumPos + intSize
	fileName := block.Filename()

	blockNumPos := fileNamePos + filemgr.MaxLength(len(fileName))
	blockNum := block.Blknum

	offsetPos := blockNumPos + intSize
	valuePos := offsetPos + intSize

	recordLen := valuePos + 1

	recordBytes := make([]byte, recordLen)
	page := filemgr.NewPageFromBytes(recordBytes)

	page.SetInt(0, int(SetBool))
	page.SetInt(txNumPos, txNum)
	if err := page.SetString(fileNamePos, fileName); err != nil {
		return -1, err
	}
	page.SetInt(blockNumPos, blockNum)
	page.SetInt(offsetPos, offset)
	page.SetBool(valuePos, val)

	return logManager.Append(recordBytes)
}
