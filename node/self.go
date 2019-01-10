package node

import (
	"fmt"
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/util"
	"sync"
)

// This file holds a singleton node structure that holds all of the information pertaining to self
var (
	self *Node
	selfOnce sync.Once
)

func GetSelf() *Node{
	selfOnce.Do(func(){
		ip, err := util.GetIPAdress()
		if err != nil {
			// TODO handle ip error
			fmt.Println(err.Error())
		}
		self = &Node{GID:config.GetConfigInstance().GID, RelayPort:config.GetConfigInstance().Relayrpcport, IP:ip}
	})
	return self
}
