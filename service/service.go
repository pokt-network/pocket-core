// This package is for the service functionality of Pocket Core. In other words, this package is for executing nonnative relay requests from clients
package service

import (
	"encoding/json"
	"net/url"
	"os"

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
}

// "RouteRelay" routes the relay to the specified hosted chain
func RouteRelay(relay Relay) (string, error) {
	if node.EnsureDWL(node.DWL(), relay.DevID) {
		hc := node.ChainToHosted(node.Blockchain{Name: relay.Blockchain, NetID: relay.NetworkID})
		u, err := url.ParseRequestURI(hc.Host + ":" + hc.Port)
		if err != nil {
			return "", err
		}
		if hc.Path != "" {
			u.Path = hc.Path
		}
		return rpc.ExecuteRequest([]byte(relay.Data), u)
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
