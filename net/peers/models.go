package peers

import (
	"github.com/pokt-network/pocket-core/node"
	"sync"
)

type PeerList struct {
	List map[string]node.Node
	sync.Mutex
}
