package hash

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/index"
	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/scan"
	"github.com/nikhilbhatia08/EuphoriaDB/table"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
	"github.com/nikhilbhatia08/EuphoriaDB/utils"
)

var _ index.Index = (*HasHIndex)(nil)

type HasHIndex struct {
	tx *transactions.Transaction
	indexName string 
	layout *record.Layout
	searchKey interface{}
	tableScan scan.TableScan
	numBuckets int
}

func NewHashIndex(tx *transactions.Transaction, indexName string, layout *record.Layout) *HasHIndex {
	return &HasHIndex{
		numBuckets: 100,
		tx: tx,
		indexName: indexName,
		layout: layout,
	}
}

func (hi *HasHIndex) BeforeFirst(searchKey interface{}) error {
	hi.Close()
	hi.searchKey = searchKey
	hash, err := utils.HashValue(searchKey)
	if err != nil {
		return err
	}

	bucket := hash % uint32(hi.numBuckets)
	tableName := fmt.Sprintf("%s-%d", hi.indexName, bucket)
	hi.tableScan, err = table.NewTableScan(hi.tx, tableName, hi.layout)
	if err != nil {
		return err
	}

	return nil
}


func (hi *HasHIndex) Close() {
	if hi.tableScan != nil {
		hi.tableScan.Close()
		hi.tableScan = nil
	}
}