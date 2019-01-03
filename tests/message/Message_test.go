package message

import (
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/message"
	"testing"
	"time"
)

func TestMessage(t *testing.T) {
	// run the servers
	message.RunMessageServers()
	time.Sleep(time.Second*2)
	// create a new message
	m := message.NewMessage(message.Payload{Data: "Hello World"})
	// create a new UDP addr object to send
	// send the message
	message.SendMessage(message.RELAY, m, _const.MSGHOST)
	time.Sleep(time.Second*2)
}
