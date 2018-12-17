// This package is for 'session' networking development
package session

import (
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/logs"
	"net"
)

/*
"ServeAndListen" creates a new 'persistent conn' server.
 */
func ServeAndListen(port string, host string) { // TODO get this from flags
	l, err := net.Listen(_const.SESSION_CONN_TYPE, host+":"+port)				// listen on port & host
	if err != nil {																// handle server creation error
		logs.NewLog("Unable to create a new "+_const.SESSION_CONN_TYPE+" server on port:"+port, logs.PanicLevel, logs.JSONLogFormat)
		logs.NewLog("ERROR: "+err.Error(), logs.PanicLevel, logs.JSONLogFormat)
	}
	defer l.Close()																// close the server after serve and listen finishes // TODO close from the caller
	logs.NewLog("Listening on port :"+port, logs.InfoLevel, logs.JSONLogFormat) // log the new connection
	for {																		// for the duration of incoming requests
		conn, err := l.Accept()													// accept the connection
		if err != nil {															// handle request acception err
			logs.NewLog("Unable to accept the "+_const.SESSION_CONN_TYPE+" conn on port:"+port, logs.PanicLevel, logs.JSONLogFormat)
			logs.NewLog("ERROR: "+err.Error(), logs.PanicLevel, logs.JSONLogFormat)
		}
		go handleRequest(conn) // TODO handle diff types of requests here
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		logs.NewLog("Unable to read the incoming message", logs.PanicLevel, logs.JSONLogFormat)
		logs.NewLog("ERROR: "+err.Error(), logs.PanicLevel, logs.JSONLogFormat)
	}
	conn.Write([]byte("ACK"))
	conn.Close()
}
