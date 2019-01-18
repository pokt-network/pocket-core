package service

import (
	"encoding/json"
	"os"

	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/plugin/rpc"
)

// "Relay" is a JSON structure that specifies information to complete reads and writes to other blockchains
// TODO convert to blockchain structure (see node.Blockchain)
type Relay struct {
	Blockchain string `json:"blockchain"`
	NetworkID  string `json:"netid"`
	Version    string `json:"version"`
	Data       string `json:"data"`
	DevID      string `json:"devid"`
}

// "RouteRelay" routes the relay to the specified hosted chain
func RouteRelay(relay Relay) (string, error) {
	if node.EnsureWL(node.GetDWL(), relay.DevID) {
		port := node.GetChainPort(node.Blockchain{Name: relay.Blockchain, NetID: relay.NetworkID, Version: relay.Version})
		if port == "" {
			logs.NewLog("Not a supported blockchain", logs.ErrorLevel, logs.JSONLogFormat)
			return "Error: not a supported blockchain", nil // TODO custom error here
		}
		return rpc.ExecuteRequest([]byte(relay.Data), port)
	}
	return "Invalid credentials", nil
}

// DISCLAIMER: The code below is for the centralized dispatcher of Pocket core mvp, may be removed for production
type Report struct {
	GID     string `json:"gid"`
	Message string `json:"message"`
}

// NOTE: This is for the centralized dispatcher of Pocket core mvp, may be removed for production
func HandleReport(report *Report) (string, error) {
	f, err := os.OpenFile(_const.REPORTFILENAME, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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
