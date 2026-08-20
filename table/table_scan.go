package table

import (
	"fmt"
	"time"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/scan"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
	"github.com/nikhilbhatia08/EuphoriaDB/types"
)

var _ scan.UpdateScan = (*TableScan)(nil)

type TableScan struct {
	tx          *transactions.Transaction
	layout      *record.Layout
	recordPage  *record.RecordPage
	filename    string
	currentSlot int
}

func NewTableScan(tx *transactions.Transaction, tableName string, layout *record.Layout) (*TableScan, error) {
	fileName := fmt.Sprintf("%s.tbl", tableName)
	tableScan := &TableScan{
		tx:          tx,
		layout:      layout,
		filename:    fileName,
		currentSlot: -1,
	}

	size, err := tx.Size(fileName)
	if err != nil {
		return nil, fmt.Errorf("error getting file size: %w", err)
	}

	if size == 0 {
		if err := tableScan.moveToNewBlock(); err != nil {
			return nil, fmt.Errorf("error moving to new block: %w", err)
		}
	} else {
		if err := tableScan.moveToBlock(0); err != nil {
			return nil, fmt.Errorf("error moving to block: %w", err)
		}
	}

	return tableScan, nil
}

func (sc *TableScan) BeforeFirst() error {
	return sc.moveToBlock(0)
}

func (sc *TableScan) Next() (bool, error) {
	slot, err := sc.recordPage.NextAfter(sc.currentSlot)
	if err != nil {
		atLastBlock, err := sc.atLastBlock()
		if err != nil {
			return false, err
		}
		if atLastBlock {
			return false, nil
		}

		if err := sc.moveToBlock(sc.recordPage.Block().Blknum + 1); err != nil {
			return false, err
		}
		slot, err = sc.recordPage.NextAfter(sc.currentSlot)
		if err != nil {
			return false, err
		}
	}
	sc.currentSlot = slot

	return true, nil
}

func (sc *TableScan) GetInt(fieldName string) (int, error) {
	return sc.recordPage.GetInt(sc.currentSlot, fieldName)
}

func (sc *TableScan) GetString(fieldName string) (string, error) {
	return sc.recordPage.GetString(sc.currentSlot, fieldName)
}

func (sc *TableScan) GetDate(fieldName string) (time.Time, error) {
	return sc.recordPage.GetDate(sc.currentSlot, fieldName)
}

func (sc *TableScan) GetBool(fieldName string) (bool, error) {
	return sc.recordPage.GetBool(sc.currentSlot, fieldName)
}

func (sc *TableScan) GetVal(fieldName string) (any, error) {
	fieldType := sc.layout.Schema().FieldType(fieldName)

	switch fieldType {
	case types.Integer:
		return sc.GetInt(fieldName)
	case types.Varchar:
		return sc.GetString(fieldName)
	case types.Boolean:
		return sc.GetBool(fieldName)
	case types.Date:
		return sc.GetDate(fieldName)
	default:
		return nil, fmt.Errorf("field type does not exist: %v", fieldType)
	}
}

func (sc *TableScan) SetInt(fieldName string, value int) error {
	return sc.recordPage.SetInt(sc.currentSlot, fieldName, value)
}

func (sc *TableScan) SetString(fieldName string, value string) error {
	return sc.recordPage.SetString(sc.currentSlot, fieldName, value)
}

func (sc *TableScan) SetBool(fieldName string, value bool) error {
	return sc.recordPage.SetBool(sc.currentSlot, fieldName, value)
}

func (sc *TableScan) SetDate(fieldName string, value time.Time) error {
	return sc.recordPage.SetDate(sc.currentSlot, fieldName, value)
}

func (sc *TableScan) SetVal(fieldName string, value any) error {
	fieldType := sc.layout.Schema().FieldType(fieldName)

	switch fieldType {
	case types.Integer:
		if v, ok := value.(int); ok {
			return sc.SetInt(fieldName, v)
		}
	case types.Varchar:
		if v, ok := value.(string); ok {
			return sc.SetString(fieldName, v)
		}
	case types.Boolean:
		if v, ok := value.(bool); ok {
			return sc.SetBool(fieldName, v)
		}
	case types.Date:
		if v, ok := value.(time.Time); ok {
			return sc.SetDate(fieldName, v)
		}
	}

	return fmt.Errorf("field type does not exist: %v", fieldType)
}

func (sc *TableScan) HasField(fieldName string) bool {
	return sc.layout.Schema().HasField(fieldName)
}

func (sc *TableScan) Insert() error {
	if sc.layout.SlotSize() > sc.tx.BlockSize() {
		return fmt.Errorf("record slot size %d is greater than block size %d", sc.layout.SlotSize(), sc.tx.BlockSize())
	}

	for {
		slot, err := sc.recordPage.InsertAfter(sc.currentSlot)
		if err == nil {
			sc.currentSlot = slot
			return nil
		}

		atLastBlock, err := sc.atLastBlock()
		if err != nil {
			return err
		}

		if atLastBlock {
			if err := sc.moveToNewBlock(); err != nil {
				return fmt.Errorf(" error move to new block: %w", err)
			}
		} else {
			if err := sc.moveToBlock(sc.recordPage.Block().Blknum + 1); err != nil {
				return fmt.Errorf("error moving to next block: %w", err)
			}
		}
	}
}

func (sc *TableScan) Delete() error {
	return sc.recordPage.Delete(sc.currentSlot)
}

func (sc *TableScan) GetRecordID() *record.Id {
	return record.NewID(sc.recordPage.Block().Blknum, sc.currentSlot)
}

func (sc *TableScan) MoveToRecordID(rid *record.Id) error {
	sc.Close()

	blk := &filemgr.BlockId{
		File:   sc.filename,
		Blknum: rid.BlockNumber(),
	}

	page, err := record.NewRecordPage(sc.tx, blk, sc.layout)
	if err != nil {
		return fmt.Errorf("create new page: %w", err)
	}

	sc.recordPage = page
	sc.currentSlot = rid.Slot()
	return nil
}

func (sc *TableScan) Close() {
	if sc.recordPage != nil {
		sc.tx.Unpin(sc.recordPage.Block())
	}
}

func (sc *TableScan) moveToBlock(blockNum int) error {
	sc.Close()

	blk := filemgr.NewBlockId(sc.filename, blockNum)

	page, err := record.NewRecordPage(sc.tx, blk, sc.layout)
	if err != nil {
		return fmt.Errorf("error creating new record page: %w", err)
	}

	sc.recordPage = page
	sc.currentSlot = -1
	return nil
}

func (sc *TableScan) moveToNewBlock() error {
	sc.Close()

	blk, err := sc.tx.Append(sc.filename)
	if err != nil {
		return fmt.Errorf("cannot append block: %w", err)
	}

	page, err := record.NewRecordPage(sc.tx, blk, sc.layout)
	if err != nil {
		return fmt.Errorf("record page: %w", err)
	}

	if err := page.Format(); err != nil {
		return fmt.Errorf("record page format err: %w", err)
	}

	sc.recordPage = page
	sc.currentSlot = -1
	return nil
}

func (sc *TableScan) atLastBlock() (bool, error) {
	fileSize, err := sc.tx.Size(sc.filename)
	if err != nil {
		return false, fmt.Errorf("error getting file size: %w", err)
	}

	return sc.recordPage.Block().BlockNum() == fileSize-1, nil
}
