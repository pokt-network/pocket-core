// Pocket Core: This is the starting point of the CLI.
package main

import "fmt"

const (
	clientIdentifier = "pocket_core"
	version = "0.0.1"
)

func main() {
	welcome()
}

func welcome(){
	fmt.Println("Client:\t\t", clientIdentifier)
	fmt.Println("Version:\t", version)

}
