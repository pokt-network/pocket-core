// This package is all message related code
package message

import "fmt"

// TODO call upon message handlers based on payload index

func HandleMessage(message *Message){
	fmt.Println(message.Payload.Data)
}

