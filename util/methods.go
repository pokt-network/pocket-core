// This package is for global utility functions and structures.
package util

import (
	"encoding/hex"
	"reflect"
	"runtime"
	"strings"
)

// "methods.go" defines global utility functions

// "Caller" returns the caller of the function that called it :)
func Caller() (*runtime.Func, uintptr) {

	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return nil, fpcs[0] // proper error her would be better
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0])
	if fun == nil {
		return nil, fpcs[0]
	}

	// return its name
	return fun, fpcs[0]
}

/*
"BytesToHex" converts a byte array into a hex string
 */
func BytesToHex(h []byte) string {
	return hex.EncodeToString(h)
}

/*
"StructTagsToString" converts the structure tags into a raw string variable.
 */
func StructTagsToString(b interface{}) []string {
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
			result[i]=fieldName
		}
	}
	return result
}
