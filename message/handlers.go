package message

import (
	"fmt"
	"net"
)

func HandleMessage(message *Message, addr *net.UDPAddr){
	switch message.Payload.ID {
	case 0: // simple print message (testing purposes)
		PrintMessage(message, addr)
	}
}

/*
Prints the payload of a message to the CLI (payload index 0)
*/
func PrintMessage(message *Message, addr *net.UDPAddr) {
	fmt.Println(message.Payload.Data, "from " + addr.IP.String() + ":", addr.Port)	// prints the payload to the CLI
}
