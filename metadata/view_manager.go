package metadata

import (
	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/table"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
)

const (
	maxViewDefinitionLength = 100
	viewNameField           = "view_name"
	viewDefinitionField     = "view_definition"
	viewCatalogTable        = "view_catalog"
)

type ViewManager struct {
	tableMetadataManager *TableMetadata
}

func NewViewManager(isNew bool, tableMetadataManager *TableMetadata, tx *transactions.Transaction) (*ViewManager, error) {
	viewManager := &ViewManager{tableMetadataManager: tableMetadataManager}

	if isNew {
		schema := record.NewSchema()
		schema.AddStringField(viewNameField, MAX_NAME_LENGTH)
		schema.AddStringField(viewDefinitionField, maxViewDefinitionLength)
		if err := viewManager.tableMetadataManager.CreateTable(viewCatalogTable, schema, tx); err != nil {
			return nil, err
		}
	}

	return viewManager, nil
}

func (vm *ViewManager) CreateView(viewName, viewDefinition string, tx *transactions.Transaction) error {
	layout, err := vm.tableMetadataManager.GetLayout(viewCatalogTable, tx)
	if err != nil {
		return err
	}

	viewCatalogTableScan, err := table.NewTableScan(tx, viewCatalogTable, layout)
	if err != nil {
		return err
	}
	defer viewCatalogTableScan.Close()

	if err := viewCatalogTableScan.Insert(); err != nil {
		return err
	}
	if err := viewCatalogTableScan.SetString(viewNameField, viewName); err != nil {
		return err
	}
	return viewCatalogTableScan.SetString(viewDefinitionField, viewDefinition)
}

func (vm *ViewManager) GetViewDefinition(viewName string, tx *transactions.Transaction) (string, error) {
	layout, err := vm.tableMetadataManager.GetLayout(viewCatalogTable, tx)
	if err != nil {
		return "", err
	}

	viewCatalogTableScan, err := table.NewTableScan(tx, viewCatalogTable, layout)
	if err != nil {
		return "", err
	}
	defer viewCatalogTableScan.Close()

	for {
		hasNext, err := viewCatalogTableScan.Next()
		if err != nil {
			return "", err
		}
		if !hasNext {
			break
		}

		name, err := viewCatalogTableScan.GetString(viewNameField)
		if err != nil {
			return "", err
		}

		if name == viewName {
			definition, err := viewCatalogTableScan.GetString(viewDefinitionField)
			if err != nil {
				return "", err
			}

			return definition, nil
		}
	}

	return "", nil
}
