// This package deals with all things networking related.
package session

import (
	"encoding/gob"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/message"
	"net"
)

/*
"NewPeer" returns a pointer to an empty peer structure.
 */
func NewPeer() *Peer {
	return &Peer{} // return an empty peer pointer
}

/*
"NewPeerFromConn" returns a new peer connection from an already established connection.
 */
func NewPeerFromConn(conn net.Conn) *Peer {
	peer := &Peer{}   															// init a peer pointer
	peer.Conn = conn  															// save connection to this peer instance
	go peer.Receive() 															// run receive to listen for incoming messages
	return peer       															// return the peer
}

/*
"CreateConnection" establishes a connection as a peer acting as a client (to a peer acting as a server).
 */
func (peer *Peer) CreateConnection(port string, host string) {
	conn, err := net.Dial(_const.SESSION_CONN_TYPE, host+":"+port) 				// establish a connection
	if err != nil { 															// handle connection error
		logs.NewLog("Unable to establish "+_const.SESSION_CONN_TYPE+" connection on port "+host+":"+port,
			logs.PanicLevel, logs.JSONLogFormat)
	}
	peer.Conn = conn  															// save the connection to this peer instance
	go peer.Receive() 															// run receive to listen for incoming messages
}

/*
"Listen" Listens on a specific port for incoming connection
 */
func Listen(port string, host string) {
	l, err := net.Listen(_const.SESSION_CONN_TYPE, host+":"+port)				// listen on port & host
	if err != nil {																// handle server creation error
		logs.NewLog("Unable to create a new "+_const.SESSION_CONN_TYPE+" server on port:"+port, logs.PanicLevel, logs.JSONLogFormat)
		logs.NewLog("ERROR: "+err.Error(), logs.PanicLevel, logs.JSONLogFormat)
	}
	defer l.Close()																// close the server after serve and listen finishes
	logs.NewLog("Listening on port :"+port, logs.InfoLevel, logs.JSONLogFormat) // log the new connection
	for {																		// for the duration of incoming requests
		conn, err := l.Accept()													// accept the connection
		if err != nil {															// handle request accept err
			logs.NewLog("Unable to accept the "+_const.SESSION_CONN_TYPE+" Conn on port:"+port, logs.PanicLevel, logs.JSONLogFormat)
			logs.NewLog("ERROR: "+err.Error(), logs.PanicLevel, logs.JSONLogFormat)
		}
		peer:=NewPeerFromConn(conn)												// create a new peer from connection
		RegisterSessionPeerConnection(*peer)									// register the peer to the global list
	}
}

/*
"Send" sends a message structure through the stream.
 */
func (peer *Peer) Send(message message.Message) {
	peer.Lock()																	// lock the peer for encoding
	encoder := gob.NewEncoder(peer.Conn)										// create a new gob encoder to the connection stream
	err := encoder.Encode(message)												// encode the structure into the stream
	peer.Unlock()																// unlock the peer
	if err != nil {																// handle any errors
		logs.NewLog("Unable to encode the message "+err.Error(), logs.PanicLevel, logs.JSONLogFormat)
	}
}

/*
"Receive" listens for messages within the stream.
 */
func (peer *Peer) Receive() {													// TODO curious if there is a more efficient way (blocking) to do this
	dec := gob.NewDecoder(peer.Conn)											// create a gob decoder object
	for {
		m := message.Message{}													// create an empty message
		dec.Decode(&m)															// decode the message from the stream
		if m != (message.Message{}) {											// if the message isn't empty
			HandleMessage(&m, peer)												// handle the message
		}
	}
}

/*
"CloseConnection" ends the persistent connection
 */
func (peer *Peer) CloseConnection() {
	peer.Lock()																	// lock the peer
	defer peer.Unlock()															// once complete unlock the peer
	peer.Conn.Close()															// close the connection
}
