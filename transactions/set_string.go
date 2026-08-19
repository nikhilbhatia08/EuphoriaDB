package transactions

import (
	"fmt"
	"runtime"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
)

type SetStringRecord struct {
	LogRecord
	TxNum  int
	Offset int
	Value  string
	Block  *filemgr.BlockId
}

func NewSetStringRecord(page *filemgr.Page) (*SetStringRecord, error) {
	intSize := getIntSize()

	txPos := intSize
	txNum := page.GetInt(txPos)

	fileNamePos := txPos + intSize
	fileName, err := page.GetString(fileNamePos)
	if err != nil {
		return nil, fmt.Errorf("error fetching file name. err: %w", err)
	}

	blockNumPos := fileNamePos + intSize
	blockNum := page.GetInt(blockNumPos)
	block := filemgr.NewBlockId(fileName, blockNum)

	offsetPos := blockNumPos + intSize
	offset := page.GetInt(offsetPos)

	valuePos := offsetPos + intSize
	value, err := page.GetString(valuePos)
	if err != nil {
		return nil, err
	}

	return &SetStringRecord{
		TxNum:  txNum,
		Offset: offset,
		Value:  value,
		Block:  block,
	}, nil
}

func (str *SetStringRecord) Op() LogRecordType {
	return SetString
}

func (str *SetStringRecord) TxNumber() int {
	return str.TxNum
}

func (str *SetStringRecord) String() string {
	return fmt.Sprintf("<SETSTRING %d %v %v %s>", str.TxNum, str.Block, str.Offset, str.Value)
}

func (str *SetStringRecord) Undo(tx *Transaction) error {
	if err := tx.Pin(str.Block); err != nil {
		return err
	}
	defer tx.Unpin(str.Block)

	return tx.SetString(str.Block, str.Offset, str.Value, false)
}

func WriteSetStringToLog(logmgr *log.LogManager, txNum int, block *filemgr.BlockId, offset int, value string) (int, error) {
	intSize := getIntSize()

	tpos := intSize
	fpos := tpos + intSize
	bpos := fpos + filemgr.MaxLength(len(block.Filename()))
	opos := bpos + intSize
	vpos := opos + intSize

	reclen := vpos + filemgr.MaxLength(len(value))
	record := make([]byte, reclen)

	page := filemgr.NewPageFromBytes(record)
	page.SetInt(0, int(SetString))
	page.SetInt(tpos, txNum)
	page.SetString(fpos, block.Filename())
	page.SetInt(bpos, block.BlockNum())
	page.SetInt(opos, offset)
	page.SetString(vpos, value)

	return logmgr.Append(record)
}

func getIntSize() int {
	if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
		return 4
	}

	return 8
}
