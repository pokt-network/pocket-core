package node

type Node struct {
	GID         string       `json:"gid"`         // node's global id (could be public address)
	IP          string       `json:"ip"`          // holds the remote IP address
	RelayPort   string       `json:"relayport"`   // specifies the port for relay API
	ClientID    string       `json:"clientid"`    // holds the identifier string for the client "pocket_core"
	CliVersion  string       `json:"cliversion"`  // holds the version of the client
	Blockchains []Blockchain `json:"blockchains"` // holds the hosted blockchains
}

type Validator struct {
	// TODO add Validator specific data here
	Node
}
