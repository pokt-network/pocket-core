package db

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/const"
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
		fmt.Println("retrieved from db", items)
		pl := node.PeerList()
		pl.Set(items)
		pl.CopyToDP()
		db.Unlock()
		// every x minutes
		time.Sleep(_const.DBREFRESH * time.Second)
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
				if !isAlive(p) {
					fmt.Println("\n" + p.GID + " failed a livenss check from dispatcher at " + p.IP + ":" + p.ClientPort + "\n")
					pl.Remove(p)
					dp.Delete(p)
					db.Remove(p)
					service.HandleReport(&service.Report{
						GID:     p.GID,
						Message: " failed a livenss check from dispatcher at " + p.IP + ":" + p.ClientPort + "\n"})
				}
			}
		}
		time.Sleep(_const.DBREFRESH * time.Second)
	}
}

// "isAlive" checks a node and returns the status of that check.
func isAlive(n node.Node) bool { // TODO handle scenarios where the error is on the dispatch node side
	if resp, err := check(n); err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

// "check" tests a node by doing an HTTP GET to API.
func check(n node.Node) (*http.Response, error) {
	return http.Get("http://" + n.IP + ":" + n.RelayPort + "/v1/")
}

// "CheckPeers" is a helper function to checks each service node's liveness. Runs checkPeers() in a go routine.
func CheckPeers() {
	go checkPeers()
}
