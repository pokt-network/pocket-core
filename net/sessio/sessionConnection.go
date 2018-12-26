// This package deals with all things networking related.
package sessio

import (
	"encoding/gob"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/message"
	"net"
	"sync"
)

// The peer structure represents a persistent connection between two nodes within a session
type Connection struct {
	Conn       net.Conn								// the persistent connection between the two
	sync.Mutex 				              			// the sPoolLock for sending messages and closing the connection
	rLock sync.Mutex
	Peer interface{}
}

/***********************************************************************************************************************
Connection Constructor
 */
func NewConnection(conn net.Conn) *Connection {
	connection := &Connection{} 												// init a connection pointer
	connection.Conn = conn      												// save connection to this connection instance
	go connection.Receive()     												// run receive to listen for incoming messages
	return connection           												// return the connection
}

/***********************************************************************************************************************
Connection Methods
 */

func (connection *Connection) CloseConnection() {
	connection.Lock()															// sPoolLock the connection
	defer connection.Unlock()													// o complete unlock the connection
	connection.Conn.Close()														// close the connection
}

func (connection *Connection) Send(message message.Message, registrants ...interface{}) {
	connection.Lock()															// sPoolLock the connection for encoding
	for _, r := range registrants {
		gob.Register(r)
	}
	encoder := gob.NewEncoder(connection.Conn)									// create a new gob encoder to the connection stream
	err := encoder.Encode(message)												// encode the structure into the stream
	connection.Unlock()															// unlock the connection
	if err != nil {																// handle any errors
		panic(err.Error())
		logs.NewLog("Unable to encode the message "+err.Error(), logs.PanicLevel, logs.JSONLogFormat)
	}
}

func (connection *Connection) Receive() {
	dec := gob.NewDecoder(connection.Conn)										// create a gob decoder object
	for {
		connection.rLock.Lock()
		m := message.Message{}													// create an empty message
		dec.Decode(&m)															// decode the message from the stream
		if m != (message.Message{}) {											// if the message isn't empty
			HandleMessage(&m)													// handle the message
		}

		connection.rLock.Unlock()
	}
}

func (connection *Connection) Listen(port string, host string) {	// TODO eventually derive port and host (need scheme to allow multiple sessions)
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
		connection.Conn = conn   												// create a new connection from connection
		go connection.Receive()
	}
}

func (connection *Connection) Dial(port string, host string) { // TODO eventually derive port and host from connection.SessionPeer
	conn, err := net.Dial(_const.SESSION_CONN_TYPE, host+":"+port) 				// establish a connection
	if err != nil { 															// handle connection error
		logs.NewLog("Unable to establish "+_const.SESSION_CONN_TYPE+" connection on port "+host+":"+port,
			logs.PanicLevel, logs.JSONLogFormat)
	}
	connection.Conn = conn            // save the connection to this connection instance
	go connection.Receive()           // run receive to listen for incoming messages
}
