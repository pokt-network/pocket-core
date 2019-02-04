package node

import (
	"sync"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/util"
)

var (
	self     *Node
	selfOnce sync.Once
)

func Self() *Node {
	selfOnce.Do(func() {
		ip, err := util.IP()
		if err != nil {
			logs.NewLog(err.Error(), logs.FatalLevel, logs.JSONLogFormat)
		}
		self = &Node{GID: config.Get().GID, RelayPort: config.Get().RRPCPort, IP: ip}
	})
	return self
}
