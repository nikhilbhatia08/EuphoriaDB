package metadata

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/table"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
	"github.com/nikhilbhatia08/EuphoriaDB/types"
)

const MAX_NAME_LENGTH = 16

type TableMetadata struct {
	tableCatalogLayout *record.Layout
	fieldCatalogLayout *record.Layout
}

func NewTableMetadata(isNew bool, tx *transactions.Transaction) (*TableMetadata, error) {
	tableMetadata := &TableMetadata{}

	tableCatalogSchema := record.NewSchema()
	tableCatalogSchema.AddStringField("tablename", MAX_NAME_LENGTH)
	tableCatalogSchema.AddIntField("slotsize")
	tableCatalogLayout := record.NewLayout(tableCatalogSchema)
	tableMetadata.tableCatalogLayout = tableCatalogLayout

	fieldCatalogSchema := record.NewSchema()
	fieldCatalogSchema.AddStringField("tablename", MAX_NAME_LENGTH)
	fieldCatalogSchema.AddStringField("fieldname", MAX_NAME_LENGTH)
	fieldCatalogSchema.AddIntField("offset")
	fieldCatalogSchema.AddIntField("type")
	fieldCatalogSchema.AddIntField("length")
	fieldCatalogLayout := record.NewLayout(fieldCatalogSchema)
	tableMetadata.fieldCatalogLayout = fieldCatalogLayout

	if isNew {
		if err := tableMetadata.CreateTable("table_catalog", tableCatalogSchema, tx); err != nil {
			return nil, fmt.Errorf("error creating table_catalog table: %w", err)
		}
		if err := tableMetadata.CreateTable("field_catalog", fieldCatalogSchema, tx); err != nil {
			return nil, fmt.Errorf("error creationg field_catalog table: %w", err)
		}
	}

	return tableMetadata, nil
}

func (tbm *TableMetadata) CreateTable(tableName string, schema *record.Schema, tx *transactions.Transaction) error {
	layout := record.NewLayout(schema)

	tableScanCatalog, err := table.NewTableScan(tx, "table_catalog", tbm.tableCatalogLayout)
	if err != nil {
		return fmt.Errorf("error scanning table catalog: %w", err)
	}
	defer tableScanCatalog.Close()

	if err := tableScanCatalog.Insert(); err != nil {
		return fmt.Errorf("error inserting: %w", err)
	}

	tableScanCatalog.SetString("tablename", tableName)
	tableScanCatalog.SetInt("slotsize", layout.SlotSize())

	fieldCatalogTableScan, err := table.NewTableScan(tx, "field_catalog", tbm.fieldCatalogLayout)
	if err != nil {
		return fmt.Errorf("error scanning field catalog: %w", err)
	}
	defer fieldCatalogTableScan.Close()

	for _, fieldName := range layout.Schema().Fields() {
		fieldCatalogTableScan.Insert()
		fieldCatalogTableScan.SetString("tablename", tableName)
		fieldCatalogTableScan.SetString("fieldname", fieldName)
		fieldCatalogTableScan.SetInt("offset", layout.Offset(fieldName))
		fieldCatalogTableScan.SetInt("length", schema.Length(fieldName))
		fieldCatalogTableScan.SetInt("type", int(schema.FieldType(fieldName)))
	}

	return nil
}

func (tbm *TableMetadata) GetLayout(tableName string, tx *transactions.Transaction) (*record.Layout, error) {
	size := -1
	tableCatalogScan, err := table.NewTableScan(tx, "table_catalog", tbm.tableCatalogLayout)
	if err != nil {
		return nil, fmt.Errorf("error fetching table catalog: %w", err)
	}
	defer tableCatalogScan.Close()

	for {
		next, err := tableCatalogScan.Next()
		if err != nil {
			return nil, err
		}
		if !next {
			break
		}

		scannedTableName, err := tableCatalogScan.GetString("tablename")
		if err != nil {
			return nil, fmt.Errorf("error fetching column: %w", err)
		}
		if tableName == scannedTableName {
			size, err = tableCatalogScan.GetInt("slotsize")
			if err != nil {
				return nil, err
			}

			break
		}
	}

	if size < 0 {
		return nil, fmt.Errorf("table %s not found", tableName)
	}

	schema := record.NewSchema()
	offsets := map[string]int{}
	fieldCatalogScan, err := table.NewTableScan(tx, "field_catalog", tbm.fieldCatalogLayout)
	if err != nil {
		return nil, fmt.Errorf("error scanning field_catalog table: %w", err)
	}
	defer fieldCatalogScan.Close()

	for {
		next, err := fieldCatalogScan.Next()
		if err != nil {
			return nil, err
		}
		if !next {
			break
		}

		scannedTableName, err := fieldCatalogScan.GetString("tablename")
		if err != nil {
			return nil, fmt.Errorf("error fetching column: %w", err)
		}
		if scannedTableName == tableName {
			fieldName, err := fieldCatalogScan.GetString("fieldname")
			if err != nil {
				return nil, err
			}
			offset, err := fieldCatalogScan.GetInt("offset")
			if err != nil {
				return nil, err
			}
			length, err := fieldCatalogScan.GetInt("length")
			if err != nil {
				return nil, err
			}
			fieldType, err := fieldCatalogScan.GetInt("type")
			if err != nil {
				return nil, err
			}

			offsets[fieldName] = offset
			schema.AddField(fieldName, types.Type(fieldType), length)
		}
	}

	return record.NewLayout(schema), nil
}
