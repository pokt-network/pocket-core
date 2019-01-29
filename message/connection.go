// This package defines p2p messaging between nodes.
package message

import (
	"bytes"
	"encoding/gob"
	"net"
	"sync"

	"github.com/pokt-network/pocket-core/logs"
)

type MSGTYPE int

const (
	BLOCKCHAIN MSGTYPE = iota + 1
	RELAY
)

var (
	server sync.Mutex
	client sync.Mutex
)

// "SendMessage" sends a message struct over the wire.
func SendMessage(msgtype MSGTYPE, message Message, ip string, registrants ...interface{}) error {
	var port string
	for _, r := range registrants {
		gob.Register(r)
	}
	switch msgtype {
	case BLOCKCHAIN:
		port = BCMSGPORT
	case RELAY:
		port = RMSGPORT
	}
	addr, err := net.ResolveUDPAddr(UDP, ip+":"+port)
	if err != nil {
		return err
	}
	return dial(message, addr, registrants)
}

// "StartServers" starts the blockchain and relay servers.
func StartServers() {
	go bmsgServer()
	go rmsgServer()
}

// "rmsgServer" runs the relay messaging server.
func rmsgServer() {
	if err := listen(RMSGPORT); err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
	}
}

// "runBCMSGserver" runs the blockchain messaging server.
func bmsgServer() {
	if err := listen(BCMSGPORT); err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
	}
}

// "listen" starts a server on a specific port.
func listen(port string) error {
	// get the local udp address
	serverAddr, err := net.ResolveUDPAddr(UDP, MSGHOST+":"+port)
	if err != nil {
		return err
	}
	// start listening locally
	conn, err := net.ListenUDP(UDP, serverAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		if err = receive(conn, buf); err != nil {
			logs.NewLog("Receive errored out "+err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		}
	}
	return nil
}

// "receive" handles incoming messages accordingly.
func receive(conn *net.UDPConn, buf []byte) error {
	server.Lock()
	defer server.Unlock()
	n, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		return err
	}
	m := new(Message)
	// use gob to decode the message struct
	if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(m); err != nil {
		return err
	}
	HandleMSG(m, addr)
	return nil
}

// "dial" initializes a connection.
func dial(message Message, ip *net.UDPAddr, registrants ...interface{}) error {
	// port shouldn't matter because of constant messaging ports throughtout
	localAddr, err := net.ResolveUDPAddr(UDP, MSGHOST+":0")
	if err != nil {
		return err
	}
	conn, err := net.DialUDP(UDP, localAddr, ip)
	if err != nil {
		return err
	}
	defer conn.Close()
	return send(conn, message, registrants)
}

// "send" sends a message over the udp connection
func send(conn *net.UDPConn, message Message, registrants ...interface{}) error {
	client.Lock()
	defer client.Unlock()
	var buf bytes.Buffer
	for r := range registrants {
		gob.Register(r)
	}
	// use gob to encode a new message structure
	if err := gob.NewEncoder(&buf).Encode(message); err != nil {
		return err
	}
	_, err := conn.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}
