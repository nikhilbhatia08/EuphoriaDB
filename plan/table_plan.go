package plan

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/metadata"
	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/scan"
	"github.com/nikhilbhatia08/EuphoriaDB/table"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
)

type TablePlan struct {
	tx        *transactions.Transaction
	tableName string
	layout    *record.Layout
	statInfo  *metadata.StatInfo
}

func NewTablePlan(tx *transactions.Transaction, tableName string, metadataMgr *metadata.MetadataManager) (*TablePlan, error) {
	tablePlan := &TablePlan{
		tx:        tx,
		tableName: tableName,
	}

	layout, err := metadataMgr.GetLayout(tableName, tx)
	if err != nil {
		return nil, fmt.Errorf("error creating table plan: %w", err)
	}

	statInfo, err := metadataMgr.GetStatInfo(tableName, layout, tx)
	if err != nil {
		return nil, fmt.Errorf("error creating table plan: %w", err)
	}

	tablePlan.layout = layout
	tablePlan.statInfo = statInfo
	return tablePlan, nil
}

func (tp *TablePlan) Open() (scan.TableScan, error) {
	return table.NewTableScan(tp.tx, tp.tableName, tp.layout)
}

func (tp *TablePlan) BlocksAccessed() int {
	return tp.statInfo.NumBlocks()
}
