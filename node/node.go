// This package contains files related specifically to nodes.
package node

// models.go maintains the structures for the node package.

type Node struct {
	GID          string  		`json:"gid"` 					// node's global id (could be public address)
	RemoteIP     string  		`json:"remoteid"`				// holds the remote IP address
	LocalIP      string  		`json:"localip"`				// holds the local IP address
	RelayPort    string  		`json:"relayport"`				// specifies the port for relay API
	ClientPort   string  		`json:"clientport"`				// specifies the port for the client API
	ClientID     string  		`json:"clientid"`				// holds the identifier string for the client "pocket_core"
	CliVersion   string  		`json:"cliversion"`				// holds the version of the client
	Blockchains []Blockchain	`json:"blockchains"`			// holds the hosted blockchains
}

type Validator struct {
	// TODO add Validator specific data here
	Node
}

type Service struct {
	// TODO add Service specific data here
	// NOTE: May not be applicable considering all nodes on the pocket network are service nodes
	Node
}
