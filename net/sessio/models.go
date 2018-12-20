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
	DISPATCHER
)
// The peer structure represents a persistent connection between two nodes within a session
type Connection struct {
	Conn net.Conn		// the persistent connection between the two
	sync.Mutex			// the lock for sending messages and closing the connection
	node.Node			// the peer that is connected
	role Role			// role the node plays within the session

}

/*
This is the session structure.
 */
type Session struct {
	devID string											// "devID" is the developer's ID that identifies the sessio
	connectionList map[string]Connection
}

/*
This holds a list of list that are active (needs to confirm using liveness check).
 */
type sessionPool struct {
	list map[string]Session // "list" is the local list of ongoing list.
}

var once sync.Once


