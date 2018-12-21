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

/*
This is the session structure.
 */
type Session struct {
	DevID    string                	`json:"devid"` 			// "DevID" is the developer's ID that identifies the sessio
	ConnList map[string]Connection 	`json:"connList"`		// "ConnList" is the List of peer connections [GID] Connection
	sync.Mutex						`json:"mutex"`
}

/*
This holds a List of List that are active (needs to confirm using liveness check).
 */
type sessionPool struct {
	List map[string]Session // "List" is the local List of ongoing List.
	sync.Mutex              // for thread safety
}

type SessionPeer struct {
	Role Role				`json:"role"`
	node.Node				`json:"node"`
}

// The peer structure represents a persistent connection between two nodes within a session
type Connection struct {
	Conn       net.Conn		`json:"conn"` 					// the persistent connection between the two
	sync.Mutex 				`json:"mutex"`                	// the sPoolLock for sending messages and closing the connection
	SessionPeer
}

type NewSessionPayload struct {
	DevID string			`json:"devid"`
	Peers []SessionPeer		`json:"peers"`
}


