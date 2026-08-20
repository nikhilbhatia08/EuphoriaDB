package scan

import "time"

type TableScan interface {
	BeforeFirst() error

	Next() (bool, error)

	GetInt(fieldName string) (int, error)

	GetString(fieldName string) (string, error)

	GetBool(fieldName string) (bool, error)

	GetDate(fieldName string) (time.Time, error)

	HasField(fieldName string) bool

	GetVal(fieldName string) (any, error)

	Close()
}
