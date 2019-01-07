// This package contains files related specifically to nodes.
package node

// models.go maintains the structures for the node package.

type Node struct {
	GID          string   // node's global id (could be public address)
	RemoteIP     string   // holds the remote IP address
	LocalIP      string   // holds the local IP address
	RelayPort    string   // specifies the port for relay API
	ClientPort   string   // specifies the port for the client API
	ClientID     string   // holds the identifier string for the client "pocket_core"
	CliVersion   string   // holds the version of the client
	HostedChains []string // holds the hosted blockchains
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
