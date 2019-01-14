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
  time.Sleep(time.Second*2)
  // add self to dispatch peers and peerlist
  node.GetPeerList().AddPeer(node.Node{GID:"self",IP:"localhost", RelayPort:config.GetConfigInstance().Relayrpcport})
  // check liveness port of self
  node.DispatchLivenessCheck()
}
