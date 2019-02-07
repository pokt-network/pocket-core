package client

import (
	"flag"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

// "Flags" handles the localhost:<client-port>/v1/client/id call.
func Flags(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag.CommandLine.SetOutput(w)
	flag.PrintDefaults()
	flag.CommandLine.SetOutput(os.Stderr)
}
