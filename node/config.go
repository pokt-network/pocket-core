package node

import "github.com/pokt-network/pocket-core/config"

func Files() {
	// Map.json
	if err := ManualPeersFile(config.GetInstance().PFile); err != nil { // add Map from file
		// TODO handle error (note: if file doesn't exist this still should work)
	}
	// chains.json
	if err := CFIle(config.GetInstance().CFile); err != nil {
		// TODO handle error (note: if hosted chains file doesn't exist how to proceed?"
	}
	// whitelists for centralized dispatcher
	WhiteListInit()
	if err := SWLFile(); err != nil {
		// TODO handle error
	}
	if err := DWLFile(); err != nil {
		// TODO handle error
	}
}
