package query

import (
	"fmt"
	"time"

	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/scan"
)

var _ scan.UpdateScan = (*SelectScan)(nil)

type SelectScan struct {
	s         scan.TableScan
	predicate *Predicate
}

func (ss *SelectScan) BeforeFirst() error {
	return ss.s.BeforeFirst()
}

func (ss *SelectScan) Next() (bool, error) {
	for {
		next, err := ss.s.Next()
		if err != nil {
			return false, fmt.Errorf("error checking next: %w", err)
		}
		if !next {
			break
		}

		if ss.predicate == nil {
			return true, nil
		}

		if satisfied, err := ss.predicate.IsStatisfied(ss.s); err == nil {
			return satisfied, nil
		} else {
			return false, fmt.Errorf("error checking predicate: %w", err)
		}
	}

	return false, nil
}

func (ss *SelectScan) GetInt(fieldName string) (int, error) {
	return ss.s.GetInt(fieldName)
}

func (ss *SelectScan) GetString(fieldName string) (string, error) {
	return ss.s.GetString(fieldName)
}

func (ss *SelectScan) GetBool(fieldName string) (bool, error) {
	return ss.s.GetBool(fieldName)
}

func (ss *SelectScan) GetDate(fieldName string) (time.Time, error) {
	return ss.s.GetDate(fieldName)
}

func (ss *SelectScan) GetVal(fieldName string) (any, error) {
	return ss.s.GetVal(fieldName)
}

func (ss *SelectScan) HasField(fieldName string) bool {
	return ss.s.HasField(fieldName)
}

func (ss *SelectScan) Close() {
	ss.s.Close()
}

func (ss *SelectScan) SetInt(fieldName string, value int) error {
	updateScan, ok := ss.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}

	return updateScan.SetInt(fieldName, value)
}

func (ss *SelectScan) SetString(fieldName string, value string) error {
	updateScan, ok := ss.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}

	return updateScan.SetString(fieldName, value)
}

func (ss *SelectScan) SetBool(fieldName string, value bool) error {
	updateScan, ok := ss.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}

	return updateScan.SetBool(fieldName, value)
}

func (ss *SelectScan) SetDate(fieldName string, value time.Time) error {
	updateScan, ok := ss.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}

	return updateScan.SetDate(fieldName, value)
}

func (ss *SelectScan) SetVal(fieldName string, value any) error {
	updateScan, ok := ss.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}

	return updateScan.SetVal(fieldName, value)
}

func (ss *SelectScan) Delete() error {
	updateScan, ok := ss.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}

	return updateScan.Delete()
}

func (ss *SelectScan) Insert() error {
	updateScan, ok := ss.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}

	return updateScan.Insert()
}

func (ss *SelectScan) GetRecordID() *record.Id {
	updateScan, ok := ss.s.(scan.UpdateScan)
	if !ok {
		panic(fmt.Sprintf("error update not supported: %T", ss.s))
	}

	return updateScan.GetRecordID()
}

func (ss *SelectScan) MoveToRecordID(rid *record.Id) error {
	updateScan, ok := ss.s.(scan.UpdateScan)
	if !ok {
		return fmt.Errorf("error update not supported")
	}
	return updateScan.MoveToRecordID(rid)
}
