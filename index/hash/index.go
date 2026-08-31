package hash

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/index"
	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/scan"
	"github.com/nikhilbhatia08/EuphoriaDB/table"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
	"github.com/nikhilbhatia08/EuphoriaDB/types"
	"github.com/nikhilbhatia08/EuphoriaDB/utils"
)

var _ index.Index = (*HashIndex)(nil)

type HashIndex struct {
	tx         *transactions.Transaction
	indexName  string
	layout     *record.Layout
	searchKey  interface{}
	tableScan  scan.UpdateScan
	numBuckets int
}

func NewHashIndex(tx *transactions.Transaction, indexName string, layout *record.Layout) *HashIndex {
	return &HashIndex{
		numBuckets: 100,
		tx:         tx,
		indexName:  indexName,
		layout:     layout,
	}
}

func (hi *HashIndex) BeforeFirst(searchKey interface{}) error {
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

func (hi *HashIndex) Next() (bool, error) {
	for {
		next, err := hi.tableScan.Next()
		if err != nil {
			return false, err
		}
		if !next {
			return false, nil
		}

		value, err := hi.tableScan.GetVal("dataval")
		if err != nil {
			return false, fmt.Errorf("error fetching value from the index.")
		}

		if types.CompareSupportedTypes(hi.searchKey, value, types.EQ) {
			return true, nil
		}
	}
}

func (hi *HashIndex) GetDataRID() (*record.Id, error) {
	blkNum, err := hi.tableScan.GetInt("block")
	if err != nil {
		return nil, err
	}

	id, err := hi.tableScan.GetInt("id")
	if err != nil {
		return nil, err
	}

	return record.NewID(blkNum, id), nil
}

func (hi *HashIndex) Insert(value interface{}, dataRID *record.Id) error {
	if err := hi.BeforeFirst(value); err != nil {
		return err
	}

	if err := hi.tableScan.Insert(); err != nil {
		return err
	}

	if err := hi.tableScan.SetInt("block", dataRID.BlockNumber()); err != nil {
		return err
	}
	if err := hi.tableScan.SetInt("id", dataRID.Slot()); err != nil {
		return err
	}
	if err := hi.tableScan.SetVal("dataval", value); err != nil {
		return err
	}

	return nil
}

func (hi *HashIndex) Delete(value interface{}, dataRID *record.Id) error {
	if err := hi.BeforeFirst(value); err != nil {
		return err
	}

	for {
		next, err := hi.Next()
		if err != nil {
			return err
		}
		if !next {
			break
		}

		rid, err := hi.GetDataRID()
		if err != nil {
			return err
		}

		if rid == dataRID {
			if err := hi.tableScan.Delete(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (hi *HashIndex) Close() error {
	if hi.tableScan != nil {
		hi.tableScan.Close()
		hi.tableScan = nil
	}
	return nil
}
