package session

import (
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/net/session"
	"io/ioutil"
	"net"
	"testing"
)

func TestConnection(t *testing.T) {
	const port = "3333"
	const host = "localhost"
	go session.ServeAndListen(port,host)										// start the server
	conn, err := net.Dial(_const.SESSION_CONN_TYPE, host+":"+port)				// establish a connection
	if err != nil {
		t.Errorf("Unable to establish " + _const.SESSION_CONN_TYPE + " connection on port " + host+":"+port)
	}
	defer conn.Close()
	conn.Write([]byte("testing"))												// send message 'testing'
	response, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Errorf("Unable to read data from  " + _const.SESSION_CONN_TYPE + " connection on port " + host+":"+port)
	}
	t.Log(string(response))
}
