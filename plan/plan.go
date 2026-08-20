package plan

import (
	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/scan"
)

type Plan interface {
	Open() (scan.TableScan, error)

	BlocksAccessed() int

	RecordsOutput() int

	DistinctValues(fieldName string) int

	Schema() *record.Schema
}
