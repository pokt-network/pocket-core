package session

import (
	"github.com/pokt-network/pocket-core/message"
	"github.com/pokt-network/pocket-core/net/sessio"
	"testing"
	"time"
)

func TestConnection(t *testing.T) { // broken test
	// establish a new listener
	const port = "3333"
	const host = "localhost"
	go sessio.Listen(port,host)
	// create a new peer
	connection1 := sessio.NewConnection()
	connection1.CreateConnection(port,host)
	// wait for peer registration
	for sessio.GetSessionConnList()[connection1.Conn.LocalAddr().String()].GID==""{}
	// get connection2 from global peer list
	connection2 := sessio.GetSessionConnList()[connection1.Conn.LocalAddr().String()]
	// create a new message to send
	peer1Payload := message.Payload{0,"I am peer 1: "+ connection1.Conn.LocalAddr().String()}
	peer2Payload := message.Payload{0,"I am peer 2: "+ connection2.Conn.LocalAddr().String()}
	// send several messages back and forth
	connection1.Send(message.NewMessage(peer1Payload))
	connection2.Send(message.NewMessage(peer2Payload))
	connection1.Send(message.NewMessage(peer1Payload))
	connection2.Send(message.NewMessage(peer2Payload))
	// don't close the connection without waiting for a second to allow for threads to complete their process
	time.Sleep(1 * time.Second)
	// close the connections
	connection1.CloseConnection()
	connection2.CloseConnection()
}
