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

/*
"Relay" is a JSON structure that specifies information to complete reads and writes to other blockchains
 */
type Relay struct {
	Blockchain string 	`json:"blockchain"`
	NetworkID string	`json:"netid"`
	Data string 		`json:"data"`
}
