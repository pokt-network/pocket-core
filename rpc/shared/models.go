// This package is shared between the different RPC packages
package shared


// Define all shared API modles


/*
"Example" is a basic JSON structure serving as a placeholder.
 */
type Example struct {
	Title string `json:"title"`
	Test  string `json:"test"`
	Data  string `json:"data"`
}

type Information struct{
	Endpoint 	string `json:"endpoint"`
	Method 		string `json:"method"`
	Parameters 	[]string `json:"params"`
	Returns		string `json:"returns"`
	Example		string `json:"example"`
}
