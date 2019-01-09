// This package is shared between the different RPC packages
package shared

import (
	"github.com/pokt-network/pocket-core/util"
	"net/http"
)

// "info.go" specifies methods pertaining to APIInfo endpoints

/*
"CreateInfoStruct" generates the specific APIReference structure dynamically.
*/
// TODO adapt this for slice structures
func CreateInfoStruct(r *http.Request, method string, model interface{}, returns string) APIReference {
	params := util.StructTagsToString(model)
	return APIReference{"localhost:port" + r.URL.String(), method,
		params, returns,
		createExampleString(r.URL.String(), params)}
}

/*
"createExampleString" creates the APIReference example string shown to the devs.
*/
func createExampleString(url string, params []string) string {
	var data string
	data = "'{"
	for index, s := range params {
		if index == len(params)-1 { // last iteration
			data += "\"" + s + "\"" + ": \"foo\"" + "}'"
		} else {
			data += "\"" + s + "\"" + ": \"foo\", "
		}
	}
	return "curl --data " + data + " localhost:port" + url
}
