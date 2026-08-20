package scan

import (
	"time"

	"github.com/nikhilbhatia08/EuphoriaDB/record"
)

type UpdateScan interface {
	TableScan

	SetVal(fieldName string, val any) error

	SetInt(fieldName string, val int) error

	SetString(fieldName string, val string) error

	SetBool(fieldName string, val bool) error

	SetDate(fieldName string, val time.Time) error

	Insert() error

	Delete() error

	GetRecordID() *record.Id

	MoveToRecordID(rid *record.Id) error
}
