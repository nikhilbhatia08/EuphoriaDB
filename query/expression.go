package query

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/scan"
)

type Expression struct {
	value     any
	fieldName string
}

func NewFieldExpression(fieldName string) *Expression {
	return &Expression{value: nil, fieldName: fieldName}
}

func NewExpressionWithValue(value any) *Expression {
	return &Expression{value: value, fieldName: ""}
}

func (e *Expression) IsFieldName() bool {
	return e.fieldName != ""
}

func (e *Expression) AsConstant() any {
	return e.value
}

func (e *Expression) AsFieldName() string {
	return e.fieldName
}

func (e *Expression) Evaluate(scan scan.TableScan) (any, error) {
	if e.value != nil {
		return e.value, nil
	}

	return scan.GetVal(e.fieldName)
}

func (e *Expression) AppliesTo(schema *record.Schema) bool {
	if e.value != nil {
		return true
	}

	return schema.HasField(e.fieldName)
}

func (e *Expression) String() string {
	if e.value != nil {
		return fmt.Sprintf("%v", e.value)
	}

	return e.fieldName
}
