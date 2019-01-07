// This package is all message related code
package message

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/logs"
	"net"
	"sync"
)

/***********************************************************************************************************************
Connection types
 */
type MSGTYPE int

var (
	serverMUX sync.Mutex
	clientMUX sync.Mutex
)

const (
	BLOCKCHAIN MSGTYPE = iota + 1
	RELAY
)

/***********************************************************************************************************************
Exported Calls
 */
func SendMessage(msgtype MSGTYPE, message Message, ip string, registrants ...interface{}) error {
	var port string
	for _,r := range registrants {
		gob.Register(r)
	}
	switch msgtype {
	case BLOCKCHAIN:
		port = _const.BCMSGPORT
	case RELAY:
		port = _const.RMSGPORT
	}
	addr, err := net.ResolveUDPAddr(_const.MSGCONNTYPE, ip+":"+port)
	if err != nil {
		return err
	}
	return dial(message, addr, registrants)
}

func RunMessageServers() {
	go runBCMSGServer()
	go runRelayMSGServer()
}

/***********************************************************************************************************************
Server
 */

func runRelayMSGServer() { // TODO handle error
	if err := listen(_const.RMSGPORT); err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
	}
}

func runBCMSGServer() { // TODO handle error
	if err := listen(_const.BCMSGPORT); err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
	}
}

func listen(port string) error {
	serverAddr, err := net.ResolveUDPAddr(_const.MSGCONNTYPE, _const.MSGHOST+":"+port)
	if err != nil {
		return err
	}
	serverConn, err := net.ListenUDP(_const.MSGCONNTYPE, serverAddr)
	if err != nil {
		return err
	}
	defer serverConn.Close()
	buf := make([]byte, 1024) // create a gob decoder object
	for {
		if err = receive(serverConn, buf); err != nil {
			return err
		}
	}
	return nil
}

func receive(conn *net.UDPConn, buf []byte) error {
	serverMUX.Lock()
	defer serverMUX.Unlock()
	n, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		return err
	}
	m := new(Message)
	if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(m); err != nil {
		return err
	}
	HandleMessage(m, addr)
	return nil
}

/***********************************************************************************************************************
Client
 */

func dial(message Message, ip *net.UDPAddr, registrants ...interface{}) error {
	localAddr, err := net.ResolveUDPAddr("udp", "localhost"+":"+"0")
	if err != nil {
		fmt.Println(err.Error())
	}
	conn, err := net.DialUDP("udp", localAddr, ip)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()
	return send(conn, message, registrants)
}

func send(conn *net.UDPConn, message Message, registrants ...interface{}) error {
	clientMUX.Lock()
	defer clientMUX.Unlock()
	var buf bytes.Buffer
	for r := range registrants {
		gob.Register(r)
	}
	if err := gob.NewEncoder(&buf).Encode(message); err != nil {
		fmt.Println(err.Error())
	}
	_, err := conn.Write(buf.Bytes())
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}
