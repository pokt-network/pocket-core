// This package is for global utility functions and structures.
package util

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// "methods.go" defines global utility functions

// "Caller" returns the caller of the function that called it
func Caller() (*runtime.Frame){
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return nil
	}
	info, err:=runtime.CallersFrames(fpcs).Next()
	if(err){
		return nil
	}
	return &info
}

/*
"BytesToHex" converts a byte array into a hex string
 */
func BytesToHex(h []byte) string {
	return hex.EncodeToString(h)
	return hex.EncodeToString(h)[2:]	// remove the 0X
}

/*
"ByteToUInt16" converts a byte array into an unsigned integer
 */
func BytesToUInt16(h []byte) uint16 {
	return binary.BigEndian.Uint16(h)
}

/*
"ByteToUInt32" converts a byte array into an unsigned integer
 */
func BytesToUInt32(h []byte) uint32 {
	return binary.BigEndian.Uint32(h)
}

/*
"ByteToUInt64" converts a byte array into an unsigned integer
*/
func BytesToUInt64(h []byte) uint64 {
	return binary.BigEndian.Uint64(h)
}

/*
"ArrayToString" converts array into comma separated String
 */
func ArrayToString(a interface{}, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
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
