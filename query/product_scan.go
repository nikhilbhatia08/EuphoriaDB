package query

import (
	"fmt"
	"time"

	"github.com/nikhilbhatia08/EuphoriaDB/scan"
	"github.com/nikhilbhatia08/EuphoriaDB/record"
)

var _ scan.UpdateScan = (*ProductScan)(nil)

type ProductScan struct {
	s1 scan.TableScan
	s2 scan.TableScan
}

func NewProductScan(scan1 scan.TableScan, scan2 scan.TableScan) *ProductScan {
	return &ProductScan{
		s1: scan1,
		s2: scan2,
	}
}

func (ps *ProductScan) BeforeFirst() error {
	if err := ps.s1.BeforeFirst(); err != nil {
		return err
	}
	if _, err := ps.s1.Next(); err != nil {
		return err
	}
	return ps.s2.BeforeFirst()
}

func (ps *ProductScan) Next() (bool, error) {
	hasNextS2, err := ps.s2.Next()
	if err != nil {
		return false, err
	}
	if hasNextS2 {
		return true, nil
	}

	if err := ps.s2.BeforeFirst(); err != nil {
		return false, err
	}

	hasNextS2, err = ps.s2.Next()
	if err != nil || !hasNextS2 {
		return false, err
	}
	hasNextS1, err := ps.s1.Next()
	if err != nil || !hasNextS1 {
		return false, err
	}

	return true, nil
}

func (ps *ProductScan) HasField(fieldName string) bool {
	return ps.s1.HasField(fieldName) || ps.s2.HasField(fieldName)
}

func (ps *ProductScan) GetInt(fieldName string) (int, error) {
	if ps.s1.HasField(fieldName) {
		return ps.s1.GetInt(fieldName)
	}
	return ps.s2.GetInt(fieldName)
}

func (ps *ProductScan) GetString(fieldName string) (string, error) {
	if ps.s1.HasField(fieldName) {
		return ps.s1.GetString(fieldName)
	}
	return ps.s2.GetString(fieldName)
}

func (ps *ProductScan) GetBool(fieldName string) (bool, error) {
	if ps.s1.HasField(fieldName) {
		return ps.s1.GetBool(fieldName)
	}
	return ps.s2.GetBool(fieldName)
}

func (ps *ProductScan) GetDate(fieldName string) (time.Time, error) {
	if ps.s1.HasField(fieldName) {
		return ps.s1.GetDate(fieldName)
	}
	return ps.s2.GetDate(fieldName)
}

func (ps *ProductScan) GetVal(fieldName string) (interface{}, error) {
	if ps.s1.HasField(fieldName) {
		return ps.s1.GetVal(fieldName)
	}
	return ps.s2.GetVal(fieldName)
}

func (ps *ProductScan) SetInt(fieldName string, val int) error {
	if ps.s1.HasField(fieldName) {
		updateScan, ok := ps.s1.(scan.UpdateScan)
		if !ok {
			return fmt.Errorf("update not supported %T", ps.s1)
		}
		return updateScan.SetInt(fieldName, val)
	}
	updateScan, ok := ps.s2.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("update not supported %T", ps.s2)
	}
	return updateScan.SetInt(fieldName, val)
}

func (ps *ProductScan) SetString(fieldName string, val string) error {
	if ps.s1.HasField(fieldName) {
		updateScan, ok := ps.s1.(scan.UpdateScan)
		if !ok {
			return fmt.Errorf("update not supported %T", ps.s1)
		}
		return updateScan.SetString(fieldName, val)
	}
	updateScan, ok := ps.s2.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("update not supported %T", ps.s2)
	}
	return updateScan.SetString(fieldName, val)
}

func (ps *ProductScan) SetBool(fieldName string, val bool) error {
	if ps.s1.HasField(fieldName) {
		updateScan, ok := ps.s1.(scan.UpdateScan)
		if !ok {
			return fmt.Errorf("update not supported %T", ps.s1)
		}
		return updateScan.SetBool(fieldName, val)
	}
	updateScan, ok := ps.s2.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("update not supported %T", ps.s2)
	}
	return updateScan.SetBool(fieldName, val)
}

func (ps *ProductScan) SetDate(fieldName string, val time.Time) error {
	if ps.s1.HasField(fieldName) {
		updateScan, ok := ps.s1.(scan.UpdateScan)
		if !ok {
			return fmt.Errorf("update not supported %T", ps.s1)
		}
		return updateScan.SetDate(fieldName, val)
	}
	updateScan, ok := ps.s2.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("update not supported %T", ps.s2)
	}
	return updateScan.SetDate(fieldName, val)
}

func (ps *ProductScan) SetVal(fieldName string, val interface{}) error {
	if ps.s1.HasField(fieldName) {
		updateScan, ok := ps.s1.(scan.UpdateScan)
		if !ok {
			return fmt.Errorf("update not supported %T", ps.s1)
		}
		return updateScan.SetVal(fieldName, val)
	}
	updateScan, ok := ps.s2.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("update not supported %T", ps.s2)
	}
	return updateScan.SetVal(fieldName, val)
}

func (ps *ProductScan) Insert() error {
	return fmt.Errorf("insert not supported on ProductScan")
}

func (ps *ProductScan) Delete() error {
	return fmt.Errorf("delete not supported on ProductScan")
}

func (ps *ProductScan) GetRecordID() *record.Id {
	panic("GetRecordID not supported on ProductScan")
}

func (ps *ProductScan) MoveToRecordID(rid *record.Id) error {
	return fmt.Errorf("MoveToRecordID not supported on ProductScan")
}

func (ps *ProductScan) Close() {
	ps.s1.Close()
	ps.s2.Close()
}
