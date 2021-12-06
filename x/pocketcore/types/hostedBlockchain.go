package types

import (
	sdk "github.com/pokt-network/pocket-core/types"
	"sync"
)

// HostedBlockchain" - An object that represents a local hosted non-native blockchain
type HostedBlockchain struct {
	ID        string    `json:"id"`         // network identifier of the hosted blockchain
	URL       string    `json:"url"`        // url of the hosted blockchain
	BasicAuth BasicAuth `json:"basic_auth"` // basic http auth optinal
}

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// HostedBlockchains" - An object that represents the local hosted non-native blockchains
type HostedBlockchains struct {
	M map[string]HostedBlockchain // M[addr] -> addr, url
	L sync.Mutex
}

// "Contains" - Checks to see if the hosted chain is within the HostedBlockchains object
func (c *HostedBlockchains) Contains(id string) bool {
	c.L.Lock()
	defer c.L.Unlock()
	// quick map check
	_, found := c.M[id]
	return found
}

// "GetChainURL" - Returns the url or error of the hosted blockchain using the hex network identifier
func (c *HostedBlockchains) GetChain(id string) (chain HostedBlockchain, err sdk.Error) {
	c.L.Lock()
	defer c.L.Unlock()
	// map check
	res, found := c.M[id]
	if !found {
		return HostedBlockchain{}, NewErrorChainNotHostedError(ModuleName)
	}
	return res, nil
}

// "GetChainURL" - Returns the url or error of the hosted blockchain using the hex network identifier
func (c *HostedBlockchains) GetChainURL(id string) (url string, err sdk.Error) {
	chain, err := c.GetChain(id)
	if err != nil {
		return "", err
	}
	return chain.URL, nil
}

// "Validate" - Validates the hosted blockchain object
func (c *HostedBlockchains) Validate() error {
	c.L.Lock()
	defer c.L.Unlock()
	// loop through all of the chains
	for _, chain := range c.M {
		// validate not empty
		if chain.ID == "" || chain.URL == "" {
			return NewInvalidHostedChainError(ModuleName)
		}
		// validate the merkleHash
		if err := NetworkIdentifierVerification(chain.ID); err != nil {
			return err
		}
	}
	return nil
}
