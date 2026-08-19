package record

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/types"
)

type Layout struct {
	schema *Schema
	offsets map[string]int
	slotSize int
}

func NewLayout(schema *Schema) *Layout {
	layout := &Layout{
		schema: schema,
		offsets: map[string]int{},
	}

	pos := types.GetIntSize()
	for _, fieldName := range schema.Fields() {
		layout.offsets[fieldName] = pos
		pos += layout.lengthInBytes(fieldName)
	}
	layout.slotSize = pos

	return layout
}

func NewLayoutWithInfo(schema *Schema, offsets map[string]int, slotSize int) *Layout {
	return  &Layout{
		schema: schema,
		offsets: offsets,
		slotSize: slotSize,
	}
}

func (ly *Layout) Schema() *Schema {
	return ly.schema
}

func (ly *Layout) Offsets() map[string]int {
	return ly.offsets
}

func (ly *Layout) Offset(fieldName string) int {
	offset, _ := ly.offsets[fieldName]
	return offset
}

func (ly *Layout) SlotSize() int {
	return ly.slotSize
}

func (l *Layout) lengthInBytes(fieldName string) int {
	fieldType := l.schema.FieldType(fieldName)

	switch fieldType {
	case types.Integer:
		return types.GetIntSize()
	case types.Boolean:
		return 1
	case types.Date:
		return 8
	case types.Varchar:
		return filemgr.MaxLength(l.schema.Length(fieldName))
	default:
		panic(fmt.Sprintf("Unknown field type: %d", fieldType))
	}
}