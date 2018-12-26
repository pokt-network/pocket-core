// This package deals with all things networking related.
package peers

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/node"
	"io/ioutil"
	"log"
	"sync"
)

// "peers.go" specifies peer related code.

// TODO turn all panics into error correction (do research into this, next RC)
// TODO standard network errors (next RC)
// TODO TODO document and reorder message indexing (ongoing)
// TODO restructure packages (next RC)
// TODO add logging (next RC)

type PeerList struct {
	List map[string]node.Node
	sync.Mutex
}

var (
	once  sync.Once
	pList *PeerList
)
/***********************************************************************************************************************
Singleton getter
 */
func GetPeerList() *PeerList {
	if pList == nil {
		once.Do(func() {
			pList = &PeerList{}
			pList.List = make(map[string]node.Node) // only make the peerlist once
		})
	}
	return pList
}
/***********************************************************************************************************************
peerList Methods
 */
func (pList *PeerList) AddPeer(node node.Node) {
	if !pList.Contains(node.GID) { // if node not within peerlist
		pList.Lock()
		defer pList.Unlock()
		pList.List[node.GID] = node
	}
}

func (pList *PeerList) RemovePeer(node node.Node) {
	pList.Lock()
	defer pList.Unlock()
	delete(pList.List, node.GID)
}

func (pList *PeerList) Contains(GID string) bool{
	pList.Lock()
	defer pList.Unlock()
	_, ok := pList.List[GID]
	return ok
}

func (pList *PeerList)Count() int{
	pList.Lock()
	defer pList.Unlock()
	return len(pList.List)
}

func (pList *PeerList) Print() {
	fmt.Println(pList.List)
}


/***********************************************************************************************************************
peerList Functions
 */

func ManualPeersFile(filepath string){
	file, err := ioutil.ReadFile(filepath)
	if err!=nil {
		log.Fatalf(err.Error())
	}
	ManualPeersJSON(file)
}

func ManualPeersJSON(b []byte){
	var data [] node.Node
	json.Unmarshal(b,&data)
	for _,n:= range data{
		pList:= GetPeerList()
		pList.AddPeer(n)
	}
}

func GetPeerCount() int{
	pList := GetPeerList()
	pList.Lock()
	defer pList.Unlock()
	return len(pList.List)
}
