package session

import (
	"fmt"
	"github.com/pokt-network/pocket-core/net/sessio"
	"testing"
)

func TestSessionMessage(t *testing.T) {
	const LPORT = "3333"										// port for listener
	const LHOST = "localhost"									// host for listener
	const SHOST = "localhost"									// host for sender
	// STEP 1: CREATE DUMMY SESSION MESSAGE TO SEND
	dSess := sessio.NewDummySession("dummy-dev-id") 		// get new dummy session (prefilled)
	message := sessio.NewSessionMessage(dSess)      			// create new message object with dSess as payload
	// STEP 2: START LISTENING ON PORT FOR SESSION AS REC
	go dSess.Listen(LPORT, LHOST)
	// STEP 3: ESTABLISH A CONNECTION AS SENDER
	dSess.Dial(LPORT,LHOST, sessio.Connection{})
	for len(dSess.ConnList) == 0{}
	// STEP 4: SEND THE MESSAGE OVER THE WIRE
	fmt.Println(dSess.ConnList)
	sender := dSess.GetConnectionByIP("127.0.0.1:"+LPORT)
	sender.Send(message)
	// STEP 5: CHECK PEERLIST FOR EACH NODE
}
