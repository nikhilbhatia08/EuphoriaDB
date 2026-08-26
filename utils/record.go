package utils

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/types"
)

func ValidateStringLength(schema *record.Schema, fieldName string, value interface{}) error {
	if schema.FieldType(fieldName) == types.Varchar {
		fieldLength := schema.Length(fieldName)
		if strValue, ok := value.(string); ok {
			if len(strValue) > fieldLength {
				return fmt.Errorf("field %s exceeds the character limit VARCHAR(%d) specified.", fieldName, fieldLength)
			}
		} else {
			return fmt.Errorf("field %s is not of type string.", fieldName)
		}
	}

	return nil
}