package parse

import "github.com/nikhilbhatia08/EuphoriaDB/record"

type CreateTableData struct {
	tableName string
	schema    *record.Schema
}

func NewCreateTableData(tableName string, sch *record.Schema) *CreateTableData {
	return &CreateTableData{
		tableName: tableName,
		schema:    sch,
	}
}

func (ctd *CreateTableData) TableName() string {
	return ctd.tableName
}

func (ctd *CreateTableData) NewSchema() *record.Schema {
	return ctd.schema
}