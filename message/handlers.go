// This package is all message related code
package message

import "fmt"

// TODO call upon message handlers based on payload index
// TODO think about concurrency

/*
"HandleMessage" is the parent function that calls upon specific message handler functions
				based on the index
 */
func HandleMessage(message *Message){
	fmt.Println(message.Payload.Data)	// currently just prints out the message but should use a switch soon
}

