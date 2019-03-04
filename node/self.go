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

func ipSetup() (string, error) {
	var err error
	ip := config.GlobalConfig().IP
	if ip == _const.DEFAULTIP {
		ip, err = util.IP()
		if err != nil {
			logs.NewLog(err.Error(), logs.FatalLevel, logs.JSONLogFormat)
			return "", err
		}
	}
	return ip, nil
}

func Self() (*Node, error) {
	var err error
	selfOnce.Do(func() {
		ip, err := ipSetup()
		if err != nil {
			ExitGracefully("unable to obtain public ip " + err.Error())
		}
		if err != nil {
			ExitGracefully("unable to generate GID " + err.Error())
		}
		self = &Node{GID: config.GlobalConfig().GID, RelayPort: config.GlobalConfig().Port, // notice this change
			IP: ip, ClientPort: config.GlobalConfig().CRPCPort, Blockchains: ChainsSlice(),
			ClientID: _const.CLIENTID, CliVersion: _const.VERSION}
	})
	if err != nil {
		return nil, err
	}
	return self, nil
}
