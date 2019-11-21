package types

import (
	"encoding/json"
	sdk "github.com/pokt-network/posmint/types"
	"io/ioutil"
	"sync"
)

type HostedBlockchain struct {
	Hash string
	URL  string
}

type HostedBlockchains List

var (
	hostedChains *HostedBlockchains // [HostedBlockchain Hash] -> Hosted HostedBlockchain
	chainOnce    sync.Once
)

func GetHostedChains() *HostedBlockchains {
	chainOnce.Do(func() {
		hostedChains = (*HostedBlockchains)(NewList())
	})
	return hostedChains
}

func (c *HostedBlockchains) AddChain(chain HostedBlockchain) {
	(*List)(c).Add(chain.Hash, chain)
}

func (c *HostedBlockchains) RemoveChain(chain HostedBlockchain) {
	(*List)(c).Remove(chain.Hash)
}

func (c *HostedBlockchains) Len() int {
	return (*List)(c).Count()
}

func (c *HostedBlockchains) ContainsFromObject(chain HostedBlockchain) bool {
	return (*List)(c).Contains(chain.Hash)
}

func (c *HostedBlockchains) ContainsFromBytes(chainHash []byte) bool {
	h := string(chainHash) // todo  -> use amino
	return (*List)(c).Contains(h)
}

func (c *HostedBlockchains) ContainsFromString(chainHash string) bool {
	return (*List)(c).Contains(chainHash)
}

func (c *HostedBlockchains) Clear() {
	(*List)(c).Clear()
}

func (c *HostedBlockchains) GetChain(hexString string) (HostedBlockchain, sdk.Error) {
	res := (*List)(c).Get(hexString)
	if res == nil {
		return HostedBlockchain{}, NewErrorChainNotHostedError()
	}
	return res.(HostedBlockchain), nil
}

func (c *HostedBlockchains) GetChainURL(hexString string) (url string, err sdk.Error) {
	res := (*List)(c).Get(hexString)
	if res == nil {
		return "", NewErrorChainNotHostedError()
	}
	return res.(HostedBlockchain).URL, nil
}

func (c *HostedBlockchains) Validate() error {
	c.Mux.Lock()
	defer c.Mux.Unlock()
	for _, chain := range c.M {
		if chain.(HostedBlockchain).Hash == "" || chain.(HostedBlockchain).URL == "" {
			return InvalidHostedChainError
		}
	}
	return nil
}

// "HostedChainsFromFile" reads a file into the hosted chains object.
func HostedChainsFromFile(filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	data := make([]HostedBlockchain, 0)
	// unmarshal json file into a slice of type HostedBlockchain
	if err := json.Unmarshal(file, &data); err != nil {
		return err
	}
	// retrieve the global chain object
	h := GetHostedChains()
	// add each chain one by one
	for _, chain := range data {
		h.AddChain(chain)
	}
	return nil
}
