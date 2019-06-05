// This package is for the service functionality of Pocket Core. In other words, this package is for executing nonnative relay requests from clients
package service

import (
	"encoding/json"
	"os"
	"strings"
	
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/plugin/rpc"
)

// "Relay" is a JSON structure that specifies information to complete reads and writes to other blockchains
type Relay struct {
	Blockchain string `json:"blockchain"`
	NetworkID  string `json:"netid"`
	Data       string `json:"data"`
	DevID      string `json:"devid"`
	Method     string `json:"method"`
	Path       string `json:"path"`
}

// "RouteRelay" routes the relay to the specified hosted chain
// This call handles REST and traditional JSON RPC
func RouteRelay(relay Relay) (string, error) {
	if node.EnsureDWL(node.DWL(), relay.DevID) {
		var url string
		hc := node.ChainToHosted(node.Blockchain{Name: relay.Blockchain, NetID: relay.NetworkID})
		url = hc.URL
		if relay.Path != "" {
			url = strings.TrimSuffix(url, "/")
			relay.Path = strings.TrimPrefix(relay.Path, "/")
			relay.Path = strings.TrimSuffix(relay.Path, "/")
			url += "/" + relay.Path
		}
		return rpc.ExecuteHTTPRequest([]byte(relay.Data), url, string(relay.Method))
	}
	return "Invalid credentials", nil
}

type Report struct {
	IP      string `json:"ip"`
	Message string `json:"message"`
}

// NOTE: This is for the centralized dispatcher of Pocket core mvp, may be removed for production
func HandleReport(report *Report) (string, error) {
	f, err := os.OpenFile(config.GlobalConfig().DD+_const.FILESEPARATOR+_const.REPORTFILENAMEPLACEHOLDER, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	text, err := json.Marshal(report)
	if err != nil {
		return "500 ERROR", err
	}
	if _, err = f.WriteString(string(text) + "\n"); err != nil {
		return "500 ERROR", err
	}
	return "Okay! The node has been successfully reported to our servers and will be reviewed! Thank you!", err
}
