// This package deals with all things networking related.
package session

import (
	"encoding/gob"
	"fmt"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/message"
	"net"
)
/* NOTE: 	The ideology behind this design is once the persistent connection has been established between the actors within
			the session, they can send and receive messages back and forth.

			This requires a few things:
				STEP 1
					1) If initially acting as a client... CreateConnection()
					2) If initially acting as a server... NewPeerFromConn()

				STEP 2
					1) Register the session peer connection in a data structure specific to the session
					2) Establish the role of the peer TODO this will be done by blockchain lookup

				STEP 3
					1) Handle messages accordingly
*/

func NewPeer() *Peer {
	return &Peer{}
}

func NewPeerFromConn(conn net.Conn) *Peer {
	fmt.Println("New peer created from connection")
	peer := &Peer{}
	peer.Conn = conn
	return peer
}

/*
"CreateConnection" sends message to server.
 */
func (peer *Peer) CreateConnection(port string, host string) {
	conn, err := net.Dial(_const.SESSION_CONN_TYPE, host+":"+port) // establish a connection
	if err != nil {
		logs.NewLog("Unable to establish "+_const.SESSION_CONN_TYPE+" connection on port "+host+":"+port,
			logs.PanicLevel, logs.JSONLogFormat)
	}
	peer.Conn = conn
}

func (peer *Peer) Send(message message.Message) {
	peer.Lock()
	encoder := gob.NewEncoder(peer.Conn)
	err := encoder.Encode(message)
	peer.Unlock()
	if err != nil {
		logs.NewLog("Unable to encode the message " + err.Error(),logs.PanicLevel, logs.JSONLogFormat) // TODO fix all panics to return error
	}
}

func (peer *Peer) Receive(){
	peer.Lock()
	dec := gob.NewDecoder(peer.Conn)
	m := &message.Message{}
	dec.Decode(m)
	peer.Unlock()
	message.HandleMessage(m)
}

func (peer *Peer) CloseConnection(){
	peer.Lock()
	defer peer.Unlock()
	peer.Conn.Close()
}

