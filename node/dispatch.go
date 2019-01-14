package node

import (
  "fmt"
  "net/http"
  "sync"
)

// This file holds a global structure used for dispatching
// NOTE: currently this is developed for PC MVP Shift and may be phased out depending on performance and need

var (
  dispatchPeers map[Blockchain]map[string]Node // <BlockchainOBJ><GID><Node>
  m             sync.Mutex
  one           sync.Once
)

func getDispatchPeers() map[Blockchain]map[string]Node {
  one.Do(func() {
    dispatchPeers = make(map[Blockchain]map[string]Node)
  })
  return dispatchPeers
}

func NewDispatchPeer(newNode Node) {
  m.Lock()
  defer m.Unlock()
  dispatchPeers := getDispatchPeers()
  for _, blockchain := range newNode.Blockchains {
    nodes := dispatchPeers[blockchain] // type map[GID]Node
    if nodes == nil { // blockchain not within list
      dispatchPeers[blockchain] = map[string]Node{newNode.GID: newNode} // add new node to empty map
    } else { // blockchain is within list
      nodes[newNode.GID] = newNode      // add node to inner map
      dispatchPeers[blockchain] = nodes // update outer map
    }
  }
}

func getPeersByBlockchain(blockchain Blockchain) map[string]Node {
  return getDispatchPeers()[blockchain]
}

func GetPeersByBlockchain(blockchain Blockchain) map[string]Node {
  m.Lock()
  defer m.Unlock()
  return getPeersByBlockchain(blockchain)
}

func DeleteDispatchPeer(delNode Node) {
  m.Lock()
  defer m.Unlock()
  for _, blockchain := range delNode.Blockchains {
    delete(getPeersByBlockchain(blockchain), delNode.GID) // delete node from map via GID
  }
}

func PrintDispatchPeers() {
  m.Lock()
  defer m.Unlock()
  for blockchain, nodeMap := range getDispatchPeers() {
    fmt.Println(blockchain.Name, "Version:", blockchain.Version, "NetID:", blockchain.NetID)
    fmt.Println("  GID's:")
    for gid, _ := range nodeMap {
      fmt.Println("   ", gid)
    }
    fmt.Println("")
  }
}

func DispatchLivenessCheck(){
  pl := GetPeerList()
  for _,peer := range pl.List{
    if !isAlive(peer){
      if !isAlive(peer){
        pl.RemovePeer(peer)
        DeleteDispatchPeer(peer)
      }
    }
  }
}

func isAlive(peer Node) bool{ // TODO handle scenarios where the error is on the dispatch node side
  if resp,err:= liveCheck(peer); err!=nil || resp==nil || resp.StatusCode != http.StatusOK {
    fmt.Println("Failed Liveness Check")
    return false
  }
  fmt.Println("Successful liveness check")
  return true
}

func liveCheck(peer Node) (*http.Response, error){
  u := "http://"+peer.IP+":"+peer.RelayPort+"/v1/"
  return http.Get(u)
}
