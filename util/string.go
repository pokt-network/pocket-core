package util

import (
	"bytes"
	"encoding/json"
)

// "StringToPrettyJSON" takes a string object and returns a json.Indent byte array
func StringToPrettyJSON(s string) ([]byte, error) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(s), "", "\t")
	return prettyJSON.Bytes(), err
}
