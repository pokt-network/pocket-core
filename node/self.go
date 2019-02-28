package node

import (
	"sync"
	
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/crypto"
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

func gidSetup() (string, error) {
	hashString, err := crypto.NewSHA1Hash()
	if err != nil {
		return "", err
	}
	return config.GlobalConfig().GID + ":" + hashString, nil
}

func Self() (*Node, error) {
	var err error
	selfOnce.Do(func() {
		ip, err := ipSetup()
		if err != nil {
			ExitGracefully("unable to obtain public ip " + err.Error())
		}
		gid, err := gidSetup()
		if err != nil {
			ExitGracefully("unable to generate GID " + err.Error())
		}
		self = &Node{GID: gid, RelayPort: config.GlobalConfig().RRPCPort,
			IP: ip, ClientPort: config.GlobalConfig().CRPCPort, Blockchains: ChainsSlice(),
			ClientID: _const.CLIENTID, CliVersion: _const.VERSION}
	})
	if err != nil {
		return nil, err
	}
	return self, nil
}
