package transactions

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
)

type SetIntRecord struct {
	TxNum  int
	offset int
	value  int
	block  *filemgr.BlockId
}

func NewSetIntRecord(page *filemgr.Page) (*SetIntRecord, error) {
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
	value := page.GetInt(valuePos)

	return &SetIntRecord{TxNum: txNum, offset: offset, value: value, block: block}, nil
}

func (si *SetIntRecord) Op() LogRecordType {
	return SetInt
}

func (si *SetIntRecord) String() string {
	return fmt.Sprintf("<SETINT %d %s %d %d>", si.TxNum, si.block, si.offset, si.value)
}

func (si *SetIntRecord) Undo(tx *Transaction) error {
	if err := tx.Pin(si.block); err != nil {
		return err
	}
	defer tx.Unpin(si.block)

	return tx.SetInt(si.block, si.offset, si.value, false)
}

func (si *SetIntRecord) TxNumber() int {
	return si.TxNum
}

func WriteSetIntToLog(logManager *log.LogManager, txNum int, block *filemgr.BlockId, offset, val int) (int, error) {
	intSize := getIntSize()

	txNumPos := intSize
	fileNamePos := txNumPos + intSize
	fileName := block.Filename()

	blockNumPos := fileNamePos + filemgr.MaxLength(len(block.File))
	blockNum := block.BlockNum()

	offsetPos := blockNumPos + intSize
	valuePos := offsetPos + intSize
	recordLen := valuePos + intSize

	recordBytes := make([]byte, recordLen)
	page := filemgr.NewPageFromBytes(recordBytes)

	page.SetInt(0, int(SetInt))
	page.SetInt(txNumPos, txNum)
	if err := page.SetString(fileNamePos, fileName); err != nil {
		return -1, err
	}
	page.SetInt(blockNumPos, blockNum)
	page.SetInt(offsetPos, offset)
	page.SetInt(valuePos, val)

	return logManager.Append(recordBytes)
}
