// This package is network code relating to other nodes within the network.
package node

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/logs"
	"io/ioutil"
	"sync"
)

// "peerList.go" specifies the peerlist structure, methods, and functions

/***********************************************************************************************************************
Peerlist structure
 */
type PeerList struct {
	List map[string]Node
	sync.Mutex
}

/***********************************************************************************************************************
Peerlist instance
 */
var (
	o  sync.Once
	pList *PeerList
)

/***********************************************************************************************************************
Singleton getter
*/

/*
"GetPeerList" returns the global list of peers
 */
func GetPeerList() *PeerList {
	o.Do(func() {									// only do once (thread safety)
		pList = &PeerList{}							// init empty peerlist
		pList.List = make(map[string]Node) 			// make the map [GID]Node
	})
	return pList									// return the peerlist
}

/***********************************************************************************************************************
peerList Methods
*/

/*
"AddPeer" adds a peer object to the global peerlist
 */
func (pList *PeerList) AddPeer(node Node) {
	if !pList.Contains(node.GID) { 					// if node not within peerlist
		pList.Lock()								// lock the list
		defer pList.Unlock()						// after function completes unlock the list
		logs.NewLog("Added new peer: "+node.GID, logs.InfoLevel, logs.JSONLogFormat)
		pList.List[node.GID] = node					// add the node to the global map
	}
}

/*
"RemovePeer" removes a peer object from the global list
 */
func (pList *PeerList) RemovePeer(node Node) {
	pList.Lock()									// lock the list
	defer pList.Unlock()							// after the function completes unlock the list
	logs.NewLog("Removed peer: "+node.GID, logs.InfoLevel, logs.JSONLogFormat)
	delete(pList.List, node.GID)					// delete the item from the map
}

/*
"Contains" returns true if node is within peerlist
 */
func (pList *PeerList) Contains(GID string) bool {
	pList.Lock()									// lock the list
	defer pList.Unlock()							// after the function completes unlock the list
	_, ok := pList.List[GID]						// check if within the list
	return ok										// return the bool
}

/*
"Count" returns the count of peers within the list
 */
func (pList *PeerList) Count() int {
	pList.Lock()									// lock the list
	defer pList.Unlock()							// after the function completes unlock the list
	return len(pList.List)							// return the length of the list
}
/*
"Print" prints the peerlist to the CLI
 */
func (pList *PeerList) Print() {
	fmt.Println(pList.List)							// print the list to the console
}

/***********************************************************************************************************************
peerList Functions
*/

/*
"ManualPeersFile" adds peers from a peers.json to the peerlist
 */
func ManualPeersFile(filepath string) error {
	file, err := ioutil.ReadFile(filepath)			// read the file from the specified path
	if err != nil {									// if error
		return err
	}
	return manualPeersJSON(file) // call manPeers.Json on the byte[]
}
/*
"manualPeersJSON" adds peers from a json []byte to the peerlist
 */
func manualPeersJSON(b []byte) error{
	var data []Node									// create an empty structure to hold the data temporarily
	if err:=json.Unmarshal(b, &data); err != nil{	// unmarshal the byte array into the struct
		return err
	}
	for _, n := range data {						// copy struct into global peerlist
		pList := GetPeerList()
		pList.AddPeer(n)
	}
	return nil
}

/*
"GetPeerCount" returns the number of peers
 */
func GetPeerCount() int {
	pList := GetPeerList()							// get the peerlist
	pList.Lock()									// lock the list
	defer pList.Unlock()							// unlock once function completes
	return len(pList.List)							// return the length of the list
}
