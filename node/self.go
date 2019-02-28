package node

import (
	"sync"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/util"
)

var (
	self     *Node
	selfOnce sync.Once
)

func Self() (*Node, error) {
	var err error
	selfOnce.Do(func() {
		ip, err := util.IP()
		if err != nil {
			logs.NewLog(err.Error(), logs.FatalLevel, logs.JSONLogFormat)
			return
		}
		self = &Node{GID: config.GlobalConfig().GID, RelayPort: config.GlobalConfig().RRPCPort,
			IP: ip, ClientPort: config.GlobalConfig().CRPCPort, Blockchains: ChainsSlice(),
			ClientID: _const.CLIENTID, CliVersion: _const.VERSION}
	})
	if err != nil {
		return nil, err
	}
	return self, nil
}
