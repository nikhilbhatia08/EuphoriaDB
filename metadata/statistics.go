package metadata

import (
	"fmt"
	"sync"

	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/table"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
)

type StatManager struct {
	tableMetadata *TableMetadata
	tableStats    map[string]*StatInfo
	numCalls      int

	mu sync.Mutex
}

func NewStatManager(tableMetadata *TableMetadata, tx *transactions.Transaction) (*StatManager, error) {
	sm := &StatManager{
		tableMetadata: tableMetadata,
	}
	if err := sm.refereshStatistics(tx); err != nil {
		return nil, fmt.Errorf("error to referesh statistics: %w", err)
	}

	return sm, nil
}

func (sm *StatManager) GetStatInfo(tableName string, layout *record.Layout, tx *transactions.Transaction) (*StatInfo, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.numCalls++
	if sm.numCalls > 100 {
		if err := sm.refereshStatistics(tx); err != nil {
			return nil, fmt.Errorf("failed to referesh statistics: %w", err)
		}
	}

	var statInfo *StatInfo
	statInfo, ok := sm.tableStats[tableName]
	if !ok {
		statInfo, err := sm.calcTableStats(tableName, layout, tx)
		if err != nil {
			return nil, fmt.Errorf("error calculating table %s stats: %w", tableName, err)
		}

		sm.tableStats[tableName] = statInfo
	}

	return statInfo, nil
}

func (sm *StatManager) refereshStatistics(tx *transactions.Transaction) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.numCalls = 0
	layout, err := sm.tableMetadata.GetLayout("table_catalog", tx)
	if err != nil {
		return err
	}

	tableCatalogScan, err := table.NewTableScan(tx, "table_catalog", layout)
	if err != nil {
		return fmt.Errorf("error scanning table_catalog: %w", err)
	}
	defer tableCatalogScan.Close()

	for {
		next, err := tableCatalogScan.Next()
		if err != nil {
			return fmt.Errorf("error checking next of table catalog scan: %w", err)
		}
		if !next {
			break
		}

		tableName, err := tableCatalogScan.GetString("tablename")
		if err != nil {
			return fmt.Errorf("error getting table name : %w", err)
		}

		tableLayout, err := sm.tableMetadata.GetLayout(tableName, tx)
		if err != nil {
			return fmt.Errorf("error getting table %s layout: %w", tableName, err)
		}

		statInfo, err := sm.calcTableStats(tableName, tableLayout, tx)
		if err != nil {
			return fmt.Errorf("error calculating table stats: %w", err)
		}
		sm.tableStats[tableName] = statInfo
	}

	return nil
}

func (sm *StatManager) calcTableStats(tableName string, layout *record.Layout, tx *transactions.Transaction) (*StatInfo, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	numRecords := 0
	numBlocks := 0

	tableScan, err := table.NewTableScan(tx, tableName, layout)
	if err != nil {
		return nil, fmt.Errorf("failed to scan %s: %w", tableName, err)
	}
	defer tableScan.Close()

	for {
		next, err := tableScan.Next()
		if err != nil {
			return nil, fmt.Errorf("error checking next record: %w", err)
		}
		if !next {
			break
		}

		numRecords++
		numBlocks = tableScan.GetRecordID().BlockNumber() + 1
	}

	return NewStatInfo(numBlocks, numRecords), nil
}
