// This package is for node assignment to a developer.
package dispatch

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/node"
)

type Dispatch struct {
	DevID       string            `json:"devid"`
	Blockchains []node.Blockchain `json:"blockchains"`
}

type DispatchServe struct {
	Name    string   `json:"name"`
	Version string   `json:"version"`
	NetID   string   `json:"netid"`
	Ips     []string `json:"ips"`
}

// NOTE: this call has been augmented for the Pocket Core MVP Centralized Dispatcher
// "Serve" formats Dispatch PL for an API request.
func Serve(dispatch *Dispatch) ([]byte, error, int) {
	if node.EnsureWL(node.DWL(), dispatch.DevID) {
		var result []DispatchServe
		for _, bc := range dispatch.Blockchains {
			ips := make([]string, 0)
			nodes := node.DispatchPeers().PeersByChain(bc)
			for _, n := range nodes {
				ips = append(ips, n.IP+":"+n.RelayPort)
			}
			result = append(result, DispatchServe{Name: strings.ToUpper(bc.Name), Version: strings.ToUpper(bc.Version), NetID: strings.ToUpper(bc.NetID), Ips: ips})
		}
		res, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return nil, err, 500
			logs.NewLog("Couldn't convert node array to json array: "+err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		}
		return res, nil, 200
	}
	return []byte(""), errors.New("invalid Credentials"), 401
}
