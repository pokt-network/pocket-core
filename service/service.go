package service

import (
	"errors"
	
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
	hc := node.ChainToHosted(node.Blockchain{Name: relay.Blockchain, NetID: relay.NetworkID, Version: relay.Version})
	port := hc.Port
	host := hc.Host
	if port == "" || host == "" {
		logs.NewLog("Not a supported blockchain", logs.ErrorLevel, logs.JSONLogFormat)
		return "This blockchain is not supported by this node", errors.New("not a supported blockchain")
	}
	return rpc.ExecuteRequest([]byte(relay.Data), host, port)
}
