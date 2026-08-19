package record

import "github.com/nikhilbhatia08/EuphoriaDB/types"

type Schema struct {
	fields []string
	info   map[string]*types.FieldInfo
}

func NewSchema() *Schema {
	return &Schema{
		fields: []string{},
		info:   map[string]*types.FieldInfo{},
	}
}

func (sc *Schema) AddField(fieldName string, fieldType types.Type, length int) {
	sc.fields = append(sc.fields, fieldName)
	sc.info[fieldName] = &types.FieldInfo{
		Type:   fieldType,
		Length: length,
	}
}

func (sc *Schema) AddIntField(fieldName string) {
	sc.AddField(fieldName, types.Integer, 0)
}

func (sc *Schema) AddStringField(fieldName string, length int) {
	sc.AddField(fieldName, types.Varchar, length)
}

func (sc *Schema) AddDateField(fieldName string) {
	sc.AddField(fieldName, types.Date, 0)
}

func (sc *Schema) AddBoolField(fieldName string) {
	sc.AddField(fieldName, types.Boolean, 0)
}

func (sc *Schema) Add(fieldName string, otherSchema *Schema) {
	fieldType := otherSchema.FieldType(fieldName)
	fieldLength := otherSchema.Length(fieldName)
	sc.AddField(fieldName, fieldType, fieldLength)
}

func (sc *Schema) AddAll(otherSchema *Schema) {
	for _, fieldName := range otherSchema.fields {
		sc.Add(fieldName, otherSchema)
	}
}

func (sc *Schema) Fields() []string {
	return sc.fields
}

func (sc *Schema) HasField(fieldName string) bool {
	if _, ok := sc.info[fieldName]; ok {
		return true
	}

	return false
}

func (sc *Schema) FieldType(fieldName string) types.Type {
	return sc.info[fieldName].Type
}

func (sc *Schema) Length(fieldName string) int {
	return sc.info[fieldName].Length
}
