package shared

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/logs"
)

// The "Route" structure defines the generalization of an api route.
type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

// "Routes" is a slice that holds all of the routes within one structure.
type Routes []Route

// "WriteRoutes" handles the localhost:<relay-port>/routes call.
func WriteRoutes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, routes Routes) {
	var paths []string
	for _, v := range routes {
		if v.Method != "GET" {
			paths = append(paths, v.Path)
		}
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	j, err := json.MarshalIndent(paths, "", "    ")
	if err != nil {
		logs.NewLog("Unable to marshal WriteRoutes to JSON", logs.ErrorLevel, logs.JSONLogFormat)
	}
	WriteRawJSONResponse(w, j, r.URL.Path, r.Host)
}
