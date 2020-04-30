package types

import (
	sdk "github.com/pokt-network/posmint/types"
	"sync"
)

// HostedBlockchain" - An object that represents a local hosted non-native blockchain
type HostedBlockchain struct {
	ID  string `json:"id"`  // network identifier of the hosted blockchain
	URL string `json:"url"` // url of the hosted blockchain
}

// HostedBlockchains" - An object that represents the local hosted non-native blockchains
type HostedBlockchains struct {
	M map[string]HostedBlockchain // m[addr] -> addr, url
	l sync.Mutex
	o sync.Once
}

// "Contains" - Checks to see if the hosted chain is within the HostedBlockchains object
func (c *HostedBlockchains) Contains(id string) bool {
	c.l.Lock()
	defer c.l.Unlock()
	// quick map check
	_, found := c.M[id]
	return found
}

// "GetChainURL" - Returns the url or error of the hosted blockchain using the hex network identifier
func (c *HostedBlockchains) GetChainURL(id string) (url string, err sdk.Error) {
	c.l.Lock()
	defer c.l.Unlock()
	// map check
	res, found := c.M[id]
	if !found {
		return "", NewErrorChainNotHostedError(ModuleName)
	}
	return res.URL, nil
}

// "Validate" - Validates the hosted blockchain object
func (c *HostedBlockchains) Validate() error {
	c.l.Lock()
	defer c.l.Unlock()
	// loop through all of the chains
	for _, chain := range c.M {
		// validate not empty
		if chain.ID == "" || chain.URL == "" {
			return NewInvalidHostedChainError(ModuleName)
		}
		// validate the hash
		if err := NetworkIdentifierVerification(chain.ID); err != nil {
			return err
		}
	}
	return nil
}
