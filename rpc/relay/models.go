// This package contains files for the Relay API
package relay

/*
"models.go" defines API models in this file.
 */

/*
"Example" is a basic JSON structure serving as a placeholder.
 */
type Example struct {
	Title string `json:"title"`
	Test  string `json:"test"`
	Data  string `json:"data"`
}

type Dispatch struct {
	DevID string `json:"devid"`
}
