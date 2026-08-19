package record

import (
	"fmt"
	"time"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
	"github.com/nikhilbhatia08/EuphoriaDB/types"
)

const (
	EMPTY int = 0
	USED  int = 1
)

type RecordPage struct {
	tx     *transactions.Transaction
	block  *filemgr.BlockId
	layout *Layout
}

func NewRecordPage(tx *transactions.Transaction, block *filemgr.BlockId, layout *Layout) (*RecordPage, error) {
	if err := tx.Pin(block); err != nil {
		return nil, err
	}

	return &RecordPage{
		tx:     tx,
		block:  block,
		layout: layout,
	}, nil
}

func (rp *RecordPage) GetInt(slot int, fieldName string) (int, error) {
	fldPos := rp.offset(slot) + rp.layout.Offset(fieldName)
	return rp.tx.GetInt(rp.block, fldPos)
}

func (rp *RecordPage) GetString(slot int, fieldName string) (string, error) {
	fldPos := rp.offset(slot) + rp.layout.Offset(fieldName)
	return rp.tx.GetString(rp.block, fldPos)
}

func (rp *RecordPage) GetBool(slot int, fieldName string) (bool, error) {
	fldPos := rp.offset(slot) + rp.layout.Offset(fieldName)
	return rp.tx.GetBool(rp.block, fldPos)
}

func (rp *RecordPage) GetDate(slot int, fieldName string) (time.Time, error) {
	fldPos := rp.offset(slot) + rp.layout.Offset(fieldName)
	return rp.tx.GetDate(rp.block, fldPos)
}

func (rp *RecordPage) SetInt(slot int, fieldName string, value int) error {
	fldPos := rp.offset(slot) + rp.layout.Offset(fieldName)
	return rp.tx.SetInt(rp.block, fldPos, value, true)
}

func (rp *RecordPage) SetString(slot int, fieldName string, value string) error {
	fldPos := rp.offset(slot) + rp.layout.Offset(fieldName)
	return rp.tx.SetString(rp.block, fldPos, value, true)
}

func (rp *RecordPage) SetBool(slot int, fieldName string, value bool) error {
	fldPos := rp.offset(slot) + rp.layout.Offset(fieldName)
	return rp.tx.SetBool(rp.block, fldPos, value, true)
}

func (rp *RecordPage) SetDate(slot int, fieldName string, value time.Time) error {
	fldPos := rp.offset(slot) + rp.layout.Offset(fieldName)
	return rp.tx.SetDate(rp.block, fldPos, value, true)
}

func (rp *RecordPage) Delete(slot int) error {
	return rp.setFlag(slot, EMPTY)
}

func (rp *RecordPage) Block() *filemgr.BlockId {
	return rp.block
}

func (rp *RecordPage) NextAfter(slot int) (int, error) {
	return rp.searchAfter(slot, USED)
}

func (rp *RecordPage) InsertAfter(slot int) (int, error) {
	newSlot, err := rp.searchAfter(slot, EMPTY)
	if err != nil {
		return -1, err
	}

	if newSlot >= 0 {
		if err := rp.setFlag(slot, USED); err != nil {
			return -1, err
		}
	}

	return newSlot, nil
}

func (rp *RecordPage) searchAfter(slot int, flag int) (int, error) {
	slot++
	for rp.isValidSlot(slot) {
		flagValue, err := rp.tx.GetInt(rp.block, rp.offset(slot))
		if err != nil {
			return -1, err
		}
		if flag == flagValue {
			return slot, nil
		}

		slot++
	}

	return -1, nil
}

func (rp *RecordPage) Format() error {
	if rp.layout.SlotSize() > rp.tx.BlockSize() {
		return fmt.Errorf("record slot size (%d) exceeds block size (%d)", rp.layout.SlotSize(), rp.tx.BlockSize())
	}

	slot := 0
	for rp.isValidSlot(slot) {
		if err := rp.tx.SetInt(rp.block, rp.offset(slot), EMPTY, false); err != nil {
			return err
		}

		schema := rp.layout.Schema()

		for _, fieldName := range schema.Fields() {
			fldPos := rp.offset(slot) + rp.layout.Offset(fieldName)

			var err error
			switch schema.FieldType(fieldName) {
			case types.Integer:
				err = rp.tx.SetInt(rp.block, fldPos, EMPTY, false)
			case types.Varchar:
				err = rp.tx.SetString(rp.block, fldPos, "", false)
			case types.Boolean:
				err = rp.tx.SetBool(rp.block, fldPos, false, false)
			case types.Date:
				err = rp.tx.SetDate(rp.block, fldPos, time.Time{}, false)
			}

			if err != nil {
				return err
			}

			slot++
		}
	}

	return nil
}

func (rp *RecordPage) setFlag(slot int, flag int) error {
	return rp.tx.SetInt(rp.block, rp.offset(slot), flag, true)
}

func (rp *RecordPage) isValidSlot(slot int) bool {
	return slot >= 0 && rp.offset(slot+1) <= rp.tx.BlockSize()
}

func (rp *RecordPage) offset(slot int) int {
	return slot * rp.layout.SlotSize()
}
