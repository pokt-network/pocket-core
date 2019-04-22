package db

import (
	"fmt"
	"github.com/pokt-network/pocket-core/util"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/service"
)

// "peersRefresh" updates the peerList and dispatchPeerList from the database every x time.
func peersRefresh() {
	for {
		var items []node.Node
		db := DB()
		db.Lock()
		output, err := DB().getAll()
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			db.Unlock()
			logs.NewLog(err.Error(), logs.PanicLevel, logs.JSONLogFormat)
		}
		// unmarshal the output from the database call
		err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &items)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			db.Unlock()
			logs.NewLog(err.Error(), logs.PanicLevel, logs.JSONLogFormat)
		}
		pl := node.PeerList()
		pl.Set(items)
		pl.CopyToDP()
		db.Unlock()
		// every x minutes
		time.Sleep(time.Duration(config.GlobalConfig().PRefresh) * time.Second)
	}
}

// "PeersRefresh" is a helper function that runs peersRefresh in a go routine
func PeersRefresh() {
	if config.GlobalConfig().Dispatch {
		go peersRefresh()
	}
}

// "checkPeers" checks each service node's liveness.
func checkPeers() {
	for {
		pl := node.PeerList()
		db := DB()
		dp := node.DispatchPeers()
		for _, p := range pl.M {
			p := p.(node.Node)
			if !isAlive(p) {
				// try again
				time.Sleep(5 * time.Second)
				if !isAlive(p) {
					fmt.Println("\n" + p.IP + " failed a liveness check from dispatcher at " + p.IP + ":" + p.RelayPort + "\n")
					pl.Remove(p)
					dp.Delete(p)
					db.Remove(p)
					service.HandleReport(&service.Report{
						IP:      p.IP,
						Message: " failed a livenss check from dispatcher at " + p.IP + ":" + p.RelayPort + "\n"})
				}
			}
		}
		time.Sleep(time.Duration(config.GlobalConfig().PRefresh) * time.Second)
	}
}

// "isAlive" checks a node and returns the status of that check.
func isAlive(n node.Node) bool { // TODO handle scenarios where the error is on the dispatch node side
	if resp, err := check(n); err != nil || resp == nil || resp.StatusCode < 200 {
		logs.NewLog(n.GID+" - "+n.IP+" failed liveness check: "+resp.Status, logs.WaringLevel, logs.JSONLogFormat)
		return false
	}
	return true
}

// "check" tests a node by doing an HTTP GET to API.
func check(n node.Node) (*http.Response, error) {
	u, err := util.URLProto(n.IP + ":" + n.RelayPort + "/v1/")
	if err != nil {
		logs.NewLog(n.GID+" - "+n.IP+" lLiveness check error: "+err.Error(), logs.WaringLevel, logs.JSONLogFormat)
		return nil, err
	}
	return http.Get(u)
}

// "CheckPeers" is a helper function to checks each service node's liveness. Runs checkPeers() in a go routine.
func CheckPeers() {
	if config.GlobalConfig().Dispatch {
		go checkPeers()
	}
}
