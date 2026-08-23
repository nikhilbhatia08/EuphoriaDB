package metadata

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
)

type MetadataManager struct {
	tableMetadata *TableMetadata
	statManager   *StatManager
	viewManager   *ViewManager
}

func NewMetadataManager(isNew bool, tx *transactions.Transaction) (*MetadataManager, error) {
	tableMetadata, err := NewTableMetadata(isNew, tx)
	if err != nil {
		return nil, fmt.Errorf("error initializing table metadata: %w", err)
	}
	StatManager, err := NewStatManager(tableMetadata, tx)
	if err != nil {
		return nil, fmt.Errorf("error initializing stat manager: %w", err)
	}
	viewManager, err := NewViewManager(isNew, tableMetadata, tx)
	if err != nil {
		return nil, fmt.Errorf("error initializing view manager: %w", err)
	}

	return &MetadataManager{
		tableMetadata: tableMetadata,
		statManager:   StatManager,
		viewManager:   viewManager,
	}, nil
}

func (mm *MetadataManager) CreateTable(tableName string, schema *record.Schema, tx *transactions.Transaction) error {
	if err := mm.tableMetadata.CreateTable(tableName, schema, tx); err != nil {
		return fmt.Errorf("error creating table %s: %w", tableName, err)
	}

	return nil
}

func (mm *MetadataManager) CreateView(viewName string, viewDef string, tx *transactions.Transaction) error {
	return mm.CreateView(viewName, viewDef, tx)
}

func (mm *MetadataManager) GetViewDefinition(viewName string, tx *transactions.Transaction) (string, error) {
	return mm.viewManager.GetViewDefinition(viewName, tx)
}

func (mm *MetadataManager) GetLayout(tableName string, tx *transactions.Transaction) (*record.Layout, error) {
	return mm.tableMetadata.GetLayout(tableName, tx)
}

func (mm *MetadataManager) GetStatInfo(tableName string, layout *record.Layout, tx *transactions.Transaction) (*StatInfo, error) {
	return mm.statManager.GetStatInfo(tableName, layout, tx)
}
