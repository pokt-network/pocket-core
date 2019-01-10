package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/pokt-network/pocket-core/logs"
)

type PeerList struct {
	List map[string]Node
	sync.Mutex
}

var (
	o     sync.Once
	pList *PeerList
)

/*
"GetPeerList" returns the global list of peers
*/
func GetPeerList() *PeerList {
	o.Do(func() {
		pList = &PeerList{}
		pList.List = make(map[string]Node)
	})
	return pList
}

/*
"AddPeer" adds a peer object to the global peerlist
 */
func (pList *PeerList) AddPeer(node Node) {
	pList.Lock()
	defer pList.Unlock()
	if !pList.contains(node.GID) {
		logs.NewLog("Added new peer: "+node.GID, logs.InfoLevel, logs.JSONLogFormat)
		pList.List[node.GID] = node
	}
}

/*
"RemovePeer" removes a peer object from the global list
 */
func (pList *PeerList) RemovePeer(node Node) {
	pList.Lock()         // lock the list
	defer pList.Unlock() // after the function completes unlock the list

	logs.NewLog("Removed peer: "+node.GID, logs.InfoLevel, logs.JSONLogFormat)
	delete(pList.List, node.GID)
}

func (pList *PeerList) contains(GID string) bool {
	_, ok := pList.List[GID]
	return ok
}

/*
"Contains" returns true if node is within peerlist
*/
func (pList *PeerList) Contains(GID string) bool {
	pList.Lock()
	defer pList.Unlock()
	return pList.contains(GID)
}

/*
"Count" returns the count of peers within the list
*/
func (pList *PeerList) Count() int {
	pList.Lock()
	defer pList.Unlock()
	return len(pList.List)
}

/*
"Print" prints the peerlist to the CLI
*/
func (pList *PeerList) Print() {
	fmt.Println(pList.List)
}

// NOTE Centralized Dispatch for MVP Only
func (pList *PeerList) AddPeersToDispatchStructure() {
	pList.Lock()
	defer pList.Unlock()
	for _, peer := range pList.List {
		NewDispatchPeer(peer)
	}
}

/***********************************************************************************************************************
peerList Functions
*/

/*
"ManualPeersFile" adds peers from a peers.json to the peerlist
*/
func ManualPeersFile(filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return manualPeersJSON(file)
}

/*
"manualPeersJSON" adds peers from a json []byte to the peerlist
<<<<<<< HEAD:node/peerList.go
 */
func manualPeersJSON(b []byte) error {
	var data []Node // create an empty structure to hold the data temporarily
	if err := json.Unmarshal(b, &data); err != nil { // unmarshal the byte array into the struct
		return err
	}
	for _, n := range data {
		pList := GetPeerList()
		pList.AddPeer(n)
	}
	return nil
}

/*
"GetPeerCount" returns the number of peers
*/
func GetPeerCount() int {
	pList := GetPeerList()
	pList.Lock()
	defer pList.Unlock()
	return len(pList.List)
}
