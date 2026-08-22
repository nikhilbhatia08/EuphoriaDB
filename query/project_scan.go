package query

import (
	"fmt"
	"time"

	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/scan"
)

var _ scan.UpdateScan = (*ProjectScan)(nil)

type ProjectScan struct {
	s scan.TableScan
	fieldList []string
}

func NewProjectScan(scan scan.TableScan, fieldList []string) *ProjectScan {
	return &ProjectScan{
		s: scan,
		fieldList: fieldList,
	}
}

func (ps *ProjectScan) BeforeFirst() error {
	return ps.BeforeFirst()
}

func (ps *ProjectScan) Next() (bool, error) {
	return ps.s.Next()
}

func (ps *ProjectScan) HasField(fieldName string) bool {
	for _, field := range ps.fieldList {
		if field == fieldName {
			return true
		}
	}

	return false
}

func (ps *ProjectScan) GetInt(fieldName string) (int, error) {
	if ps.HasField(fieldName) {
		return ps.s.GetInt(fieldName)
	}

	return -1, fmt.Errorf("field not found: %s.", fieldName)
}

func (ps *ProjectScan) GetString(fieldName string) (string, error) {
	if ps.HasField(fieldName) {
		return ps.s.GetString(fieldName)
	}

	return "", fmt.Errorf("field not found: %s.", fieldName)
}

func (ps *ProjectScan) GetBool(fieldName string) (bool, error) {
	if ps.HasField(fieldName) {
		return ps.s.GetBool(fieldName)
	}

	return false, fmt.Errorf("field not found: %s.", fieldName)
}

func (ps *ProjectScan) GetDate(fieldName string) (time.Time, error) {
	if ps.HasField(fieldName) {
		return ps.s.GetDate(fieldName)
	}

	return time.Time{}, fmt.Errorf("field not found: %s.", fieldName)
}

func (ps *ProjectScan) GetVal(fieldName string) (any, error) {
	if ps.HasField(fieldName) {
		return ps.s.GetVal(fieldName)
	}

	return nil, fmt.Errorf("field not found: %s.", fieldName)
}

func (ps *ProjectScan) SetInt(fieldName string, value int) error {
	if !ps.HasField(fieldName) {
		return fmt.Errorf("field not found: %s.", fieldName)
	}

	updateScan, ok := ps.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}

	return updateScan.SetInt(fieldName, value)
}

func (ps *ProjectScan) SetString(fieldName string, value string) error {
	if !ps.HasField(fieldName) {
		return fmt.Errorf("field not found: %s.", fieldName)
	}

	updateScan, ok := ps.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}

	return updateScan.SetString(fieldName, value)
}

func (ps *ProjectScan) SetBool(fieldName string, value bool) error {
	if !ps.HasField(fieldName) {
		return fmt.Errorf("field not found: %s.", fieldName)
	}

	updateScan, ok := ps.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}

	return updateScan.SetBool(fieldName, value)
}

func (ps *ProjectScan) SetDate(fieldName string, value time.Time) error {
	if !ps.HasField(fieldName) {
		return fmt.Errorf("field not found: %s.", fieldName)
	}

	updateScan, ok := ps.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}

	return updateScan.SetDate(fieldName, value)
}

func (ps *ProjectScan) SetVal(fieldName string, value any) error {
	if !ps.HasField(fieldName) {
		return fmt.Errorf("field not found: %s.", fieldName)
	}

	updateScan, ok := ps.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}

	return updateScan.SetVal(fieldName, value)
}

func (ps *ProjectScan) Delete() error {
	updateScan, ok := ps.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}

	return updateScan.Delete()
}

func (ps *ProjectScan) Insert() error {
	updateScan, ok := ps.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}

	return updateScan.Insert()
}

func (ps *ProjectScan) GetRecordID() *record.Id {
	updateScan, ok := ps.s.(scan.UpdateScan)
	if !ok {
		panic(fmt.Sprintf("error update not supported: %v", ps.s))
	}

	return updateScan.GetRecordID()
}

func (ps *ProjectScan) MoveToRecordID(rid *record.Id) error {
	updateScan, ok := ps.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}
	return updateScan.MoveToRecordID(rid)
}

func (ps *ProjectScan) Close() {
	ps.s.Close()
}


