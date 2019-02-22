package relay

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "Routes" is a function that returns all of the routes of the API.
func Routes() shared.Routes {
	routes := shared.Routes{
		shared.Route{Name: "Version", Method: "GET", Path: "/v1", HandlerFunc: Version},
		shared.Route{Name: "WriteRoutes", Method: "GET", Path: "/v1/routes", HandlerFunc: WriteRoutes},
		shared.Route{Name: "Report", Method: "POST", Path: "/v1/report", HandlerFunc: Report},
		shared.Route{Name: "ReportInfo", Method: "GET", Path: "/v1/report", HandlerFunc: ReportInfo},
		shared.Route{Name: "Dispatch", Method: "POST", Path: "/v1/dispatch", HandlerFunc: Dispatch},
		shared.Route{Name: "DispatchInfo", Method: "GET", Path: "/v1/dispatch", HandlerFunc: DispatchInfo},
		shared.Route{Name: "Relay", Method: "POST", Path: "/v1/relay/", HandlerFunc: Relay},
		shared.Route{Name: "RelayInfo", Method: "GET", Path: "/v1/relay/", HandlerFunc: RelayInfo},
		shared.Route{Name: "Register", Method: "POST", Path: "/v1/register", HandlerFunc: Register},
		shared.Route{Name: "UnRegister", Method: "POST", Path: "/v1/unregister", HandlerFunc: UnRegister},
		shared.Route{Name: "RegisterInfo", Method: "GET", Path: "/v1/register", HandlerFunc: RegisterInfo},
		shared.Route{Name: "UnRegisterInfo", Method: "GET", Path: "/v1/unregister", HandlerFunc: UnRegisterInfo},
		shared.Route{Name: "Flags", Method: "GET", Path: "/v1/flags", HandlerFunc: Flags},
	}
	return routes
}

// "WriteRoutes" handles the localhost:<relay-port>/routes call.
func WriteRoutes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteRoutes(w, r, ps, Routes())
}
