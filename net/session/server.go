package session

import (
	"fmt"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/logs"
	"net"
)

/*
"Listen" Listens on a specific port
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
			logs.NewLog("Unable to accept the "+_const.SESSION_CONN_TYPE+" Conn on port:"+port, logs.PanicLevel, logs.JSONLogFormat) // TODO fix all panic levels with consecutive logs
			logs.NewLog("ERROR: "+err.Error(), logs.PanicLevel, logs.JSONLogFormat)
		}
		go HandleNewConn(conn) // creates new client instance
	}
}

func HandleNewConn(conn net.Conn){
	peer:=NewPeerFromConn(conn)
	AddPeerToSessionPeersList(*peer) // TODO error handling
	fmt.Println("Added new peer to the sessionPeerList @ " +peer.Conn.RemoteAddr().String())
}
