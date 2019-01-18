package node

import (
	"sync"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/util"
)

var (
	self     *Node
	selfOnce sync.Once
)

func GetSelf() *Node {
	selfOnce.Do(func() {
		ip, err := util.GetIPAdress()
		if err != nil {
			// TODO handle error
		}
		self = &Node{GID: config.GetInstance().GID, RelayPort: config.GetInstance().RRPCPort, IP: ip}
	})
	return self
}
