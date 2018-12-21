// This package deals with all things networking related.
package sessio

import (
	"encoding/gob"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/message"
	"net"
)

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

func (connection *Connection) Send(message message.Message) { // TODO register type upon caller
	connection.Lock()															// sPoolLock the connection for encoding
	gob.Register(NewSessionPayload{})
	encoder := gob.NewEncoder(connection.Conn)									// create a new gob encoder to the connection stream
	err := encoder.Encode(message)												// encode the structure into the stream
	connection.Unlock()															// unlock the connection
	if err != nil {																// handle any errors
		panic(err.Error())
		logs.NewLog("Unable to encode the message "+err.Error(), logs.PanicLevel, logs.JSONLogFormat)
	}
}

func (connection *Connection) Receive() {		 								// TODO consider locking for decoding messages (needs its own lock 'rLock')
	dec := gob.NewDecoder(connection.Conn)										// create a gob decoder object
	for {
		m := message.Message{}													// create an empty message
		dec.Decode(&m)															// decode the message from the stream
		if m != (message.Message{}) {											// if the message isn't empty
			HandleMessage(&m)													// handle the message
		}
	}
}
