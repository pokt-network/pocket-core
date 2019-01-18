// This package is for global utility functions and structures.
package util

import (
	"reflect"
	"strings"
)

// "StructTags" converts the structure tags into a raw string variable.
func StructTags(b interface{}) []string {
	val := reflect.ValueOf(b)
	size := val.Type().NumField()
	result := make([]string, size)
	for i := 0; i < size; i++ {
		t := val.Type().Field(i)
		fieldName := t.Name

		if jsonTag := t.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
				fieldName = jsonTag[:commaIdx]
			}
			result[i] = fieldName
		}
	}
	return result
}
