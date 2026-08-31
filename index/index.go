package index

import "github.com/nikhilbhatia08/EuphoriaDB/record"

type Index interface {
	BeforeFirst(searchKey interface{}) error

	Next() (bool, error)

	GetDataRID() (*record.Id, error)

	Insert(value interface{}, dataRID *record.Id) error

	Delete(value interface{}, dataRID *record.Id) error

	Close() error
}
