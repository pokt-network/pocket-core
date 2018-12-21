// This package deals with all things networking related.
package sessio

import (
	"github.com/pokt-network/pocket-core/node"
	"net"
	"sync"
)
type Role int

const (
	VALIDATOR Role = iota+1
	SERVICER
	DISPATCHER // TODO
)
// The peer structure represents a persistent connection between two nodes within a session
type Connection struct {
	Conn       net.Conn		`json:"conn"` 					// the persistent connection between the two
	sync.Mutex 				`json:"mutex"`                	// the sPoolLock for sending messages and closing the connection
	node.Node  				`json:"node"`                 	// the peer that is connected
	Role       Role 		`json:"Role"`            		// Role the node plays within the session
}

/*
This is the session structure.
 */
type Session struct {
	DevID    string                `json:"DevID"` 			// "DevID" is the developer's ID that identifies the sessio
	ConnList map[string]Connection `json:"connList"`		// "ConnList" is the list of peer connections
	sync.Mutex
}

/*
This holds a list of list that are active (needs to confirm using liveness check).
 */
type sessionPool struct {
	list map[string]Session // "list" is the local list of ongoing list.
}


