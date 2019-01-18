package message

import (
	"testing"
	"time"

	"github.com/pokt-network/pocket-core/message"
)

func TestMessage(t *testing.T) {
	// run the servers
	message.StartServers()
	time.Sleep(time.Second * 2)
	// create a new message
	m := message.NewMessage(message.Payload{Data: "Hello World"})
	// create a new UDP addr object to send
	// send the message
	message.SendMessage(message.RELAY, m, message.MSGHOST)
	time.Sleep(time.Second * 2)
}
