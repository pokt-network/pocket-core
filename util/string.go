package util

import (
	"bytes"
	"encoding/json"
)

func StringToPrettyJSON(s string) ([]byte, error) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(s), "", "\t")
	return prettyJSON.Bytes(), err
}
