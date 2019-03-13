// This package is for the service functionality of Pocket Core. In other words, this package is for executing nonnative relay requests from clients
package service

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/plugin/rpc"
)

// "Relay" is a JSON structure that specifies information to complete reads and writes to other blockchains
type Relay struct {
	Blockchain string `json:"blockchain"`
	NetworkID  string `json:"netid"`
	Version    string `json:"version"`
	Data       string `json:"data"`
	DevID      string `json:"devid"`
}

// "RouteRelay" routes the relay to the specified hosted chain
func RouteRelay(relay Relay) (string, error) {
	if node.EnsureWL(node.DWL(), relay.DevID) {
		hc := node.ChainToHosted(node.Blockchain{Name: relay.Blockchain, NetID: relay.NetworkID, Version: relay.Version})
		port := hc.Port
		host := hc.Host
		if port == "" || host == "" {
			logs.NewLog("Not a supported blockchain", logs.ErrorLevel, logs.JSONLogFormat)
			return "This blockchain is not supported by this node", errors.New("not a supported blockchain")
		}
		return rpc.ExecuteRequest([]byte(relay.Data), host, port)
	}
	return "Invalid credentials", nil
}

type Report struct {
	GID     string `json:"gid"`
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
