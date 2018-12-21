// This package deals with all things networking related.
package sessio

import (
	"encoding/gob"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/message"
	"net"
)

/*
"NewPeerFromConn" returns a new connection connection from an already established connection.
 */
func NewPeerFromConn(conn net.Conn) *Connection {
	connection := &Connection{} // init a connection pointer
	connection.Conn = conn      // save connection to this connection instance
	go connection.Receive()     // run receive to listen for incoming messages
	return connection           // return the connection
}

/*
"CreateConnection" establishes a connection as a connection acting as a client (to a connection acting as a server).
 */
func (connection *Connection) CreateConnection(port string, host string, session Session) {
	conn, err := net.Dial(_const.SESSION_CONN_TYPE, host+":"+port) 				// establish a connection
	if err != nil { 															// handle connection error
		logs.NewLog("Unable to establish "+_const.SESSION_CONN_TYPE+" connection on port "+host+":"+port,
			logs.PanicLevel, logs.JSONLogFormat)
	}
	connection.Conn = conn  													// save the connection to this connection instance
	go connection.Receive() 													// run receive to listen for incoming messages
	session.RegisterSessionConn(*connection) // TODO consider returning the connection and register it from caller
}

/*
"Listen" Listens on a specific port for incoming connection
 */
func Listen(port string, host string, session Session) {
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
		connection:=NewPeerFromConn(conn) // create a new connection from connection
		session.RegisterSessionConn(*connection)  // TODO consider returning the connection and register it from caller
	}
}

/*
"Send" sends a message structure through the stream.
 */
func (connection *Connection) Send(message message.Message) {
	connection.Lock()															// sPoolLock the connection for encoding
	encoder := gob.NewEncoder(connection.Conn)									// create a new gob encoder to the connection stream
	err := encoder.Encode(message)												// encode the structure into the stream
	connection.Unlock()															// unlock the connection
	if err != nil {																// handle any errors
		logs.NewLog("Unable to encode the message "+err.Error(), logs.PanicLevel, logs.JSONLogFormat)
	}
}

/*
"Receive" listens for messages within the stream. TODO may need to lock up for the decoding
 */
func (connection *Connection) Receive() { // TODO curious if there is a more efficient way (blocking) to do this
	dec := gob.NewDecoder(connection.Conn)										// create a gob decoder object
	for {
		m := message.Message{}													// create an empty message
		dec.Decode(&m)															// decode the message from the stream
		if m != (message.Message{}) {											// if the message isn't empty
			HandleMessage(&m, connection)										// handle the message
		}
	}
}

/*
"CloseConnection" ends the persistent connection
 */
func (connection *Connection) CloseConnection() {
	connection.Lock()															// sPoolLock the connection
	defer connection.Unlock()													// o complete unlock the connection
	connection.Conn.Close()														// close the connection
}

func (connection *Connection) SetRole(role Role){
	connection.Role = role
}

func (connection *Connection) GetRole() Role{
	return connection.Role
}
