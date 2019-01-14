package node

import (
  "testing"
  "time"
  
  "github.com/pokt-network/pocket-core/config"
  "github.com/pokt-network/pocket-core/node"
  "github.com/pokt-network/pocket-core/rpc"
)

func TestDispatchLiveness(t *testing.T) {
  // start API servers
  go rpc.StartRelayRPC(config.GetConfigInstance().Relayrpcport)
  time.Sleep(time.Second)
  // get peer list
  pl:=node.GetPeerList()
  // create arbitrary nodes
  self:=node.Node{GID:"self",IP:"localhost", RelayPort:config.GetConfigInstance().Relayrpcport}
  dead:=node.Node{GID:"deadNode",IP:"0.0.0.0",RelayPort:"0"}
  // add self to peerlist
  pl.AddPeer(self)
  // add dead node to peerlist
  pl.AddPeer(dead)
  // check liveness port of self
  node.DispatchLivenessCheck()
  // ensure that dead node is deleted and live node is still within list
  if !pl.Contains(self.GID) || pl.Contains(dead.GID){
    t.Fatalf("Peerlist result not correct, expected: " + self.GID +" only, and not: " + dead.GID)
  }
  pl.Print()
}
