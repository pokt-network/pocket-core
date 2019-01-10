package node

import (
	"github.com/pokt-network/pocket-core/config"
	"sync"
)

// This file holds a singleton node structure that holds all of the information pertaining to self
var (
	self *Node
	selfOnce sync.Once
)

func GetSelf() *Node{
	selfOnce.Do(func(){
		self = &Node{GID:config.GetConfigInstance().GID}
	})
	return self
}
