// This package is network code relating to pocket 'sessions'
package sessio

import (
	"encoding/gob"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/message"
	"net"
	"sync"
)

// "sessionConnection.go" specifies the connection obj structure, methods, and functions.

// The peer structure represents a persistent connection between two nodes within a session
type Connection struct {
	Conn       net.Conn 														// the persistent connection
	sync.Mutex          														// lock for sending closing
	rLock      sync.Mutex														// the lock for for decoding messages
	Peer       interface{}														// peer object
}

/***********************************************************************************************************************
Connection Constructor
*/

/*
"NewConnection" creates a new Connection object from a net.Conn object
 */
func NewConnection(conn net.Conn) *Connection {
	connection := &Connection{} 												// init a connection pointer
	connection.Conn = conn      												// save connection to this instance
	go connection.Receive()     												// run receive to listen for messages
	return connection           												// return the connection
}

/***********************************************************************************************************************
Connection Methods
*/

/*
"CloseConnection" closes the net.Conn connection from the Connection object
 */
func (connection *Connection) CloseConnection() {
	connection.Lock()         													// sPoolLock the connection
	defer connection.Unlock() 													// o complete unlock the connection
	connection.Conn.Close()   													// close the connection
}

/*
"Send" sends a message via the Connection object
 */
func (connection *Connection) Send(message message.Message, registrants ...interface{}) error {
	connection.Lock() 															// sPoolLock the connection for encoding
	for _, r := range registrants {
		gob.Register(r)
	}
	encoder := gob.NewEncoder(connection.Conn) 									// new gob encoder to the stream
	err := encoder.Encode(message)             									// encode the structure into the stream
	connection.Unlock()                        									// unlock the connection
	if err != nil {                            									// handle any errors
		logs.NewLog("Unable to encode the message "+
			err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		return err																// return the error
	}
	return nil																	// return no error
}

/*
"receive" receives a message via the Connection object
 */
func (connection *Connection) Receive() {
	dec := gob.NewDecoder(connection.Conn) 										// create a gob decoder object
	for {
		connection.rLock.Lock()
		m := message.Message{}        											// create an empty message
		dec.Decode(&m)                											// decode the message from the stream
		if m != (message.Message{}) { 											// if the message isn't empty
		logs.NewLog("New message received",
			logs.InfoLevel, logs.JSONLogFormat)
			HandleMessage(&m) 													// handle the message
		}
		connection.rLock.Unlock()
	}
}

/*
"Listen" creates a server for the initial connection of the Connection object
 */
// TODO eventually derive port and host (need scheme to allow multiple sessions)
func (connection *Connection) Listen(port string, host string) error {
	l, err := net.Listen(_const.SESSCONNTYPE, host+":"+port) // listen on port & host
	if err != nil {                                               				// handle server creation error
		logs.NewLog("Unable to create a new "+_const.SESSCONNTYPE+
			" server on port:"+port, logs.ErrorLevel, logs.JSONLogFormat)
		logs.NewLog("ERROR: "+err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		return err
	}
	defer l.Close()                                                             // close the server after finish
	logs.NewLog("Listening on port :"+port, logs.InfoLevel, logs.JSONLogFormat) // log the new connection
	for {                                                                       // for the duration of incoming requests
		conn, err := l.Accept() 												// accept the connection
		if err != nil {         												// handle request accept err
			logs.NewLog("Unable to accept the "+_const.SESSCONNTYPE+
				" Conn on port:"+port, logs.ErrorLevel, logs.JSONLogFormat)
			logs.NewLog("ERROR: "+err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
			return err															// return the error
		}
		connection.Conn = conn 													// create a new connection from net.Conn
		go connection.Receive()													// receive on a new thread
	}
}

/*
"Dial" instantiates a connection with a specific host and port for the Connection object
 */
// TODO eventually derive port and host from connection.SessionPeer
func (connection *Connection) Dial(port string, host string) error {
	conn, err := net.Dial(_const.SESSCONNTYPE, host+":"+port) // establish a connection
	if err != nil {                                                				// handle connection error
		logs.NewLog("Unable to establish "+_const.SESSCONNTYPE+" connection on port "+host+":"+port,
			logs.ErrorLevel, logs.JSONLogFormat)
		return err																// return the error
	}
	connection.Conn = conn  													// save net.Conn to this object
	go connection.Receive() 													// receive on a new thread
	return nil																	// return no error
}
