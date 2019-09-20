package types

import (
	"encoding/hex"
	"encoding/json"
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

func (c *HostedBlockchains) GetChainFromBytes(chainHash []byte) HostedBlockchain {
	if c == nil || len(c.M) == 0 {
		return HostedBlockchain{}
	}
	h := hex.EncodeToString(chainHash)
	res := (*List)(c).Get(h)
	if res == nil {
		return HostedBlockchain{}
	}
	return res.(HostedBlockchain)
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
