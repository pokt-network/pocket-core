package types

import (
	sdk "github.com/pokt-network/posmint/types"
	"sync"
)

type HostedBlockchain struct {
	Hash string `json:"addr"`
	URL  string `json:"url"`
}

type HostedBlockchains struct {
	M map[string]HostedBlockchain // m[addr] -> addr, url
	l sync.Mutex
	o sync.Once
}

var (
	globalHostedChains *HostedBlockchains // [HostedBlockchain HashString] -> Hosted HostedBlockchain
	chainOnce          sync.Once
)

func GetHostedChains() *HostedBlockchains { // todo getHostedChains never called!
	chainOnce.Do(func() {
		globalHostedChains = &HostedBlockchains{
			M: make(map[string]HostedBlockchain),
		}
	})
	return globalHostedChains
}

func (c *HostedBlockchains) Add(chain HostedBlockchain) {
	c.l.Lock()
	defer c.l.Unlock()
	c.M[chain.Hash] = chain
}

func (c *HostedBlockchains) Delete(chain HostedBlockchain) {
	c.l.Lock()
	defer c.l.Unlock()
	delete(c.M, chain.Hash)
}

func (c *HostedBlockchains) Len() int {
	c.l.Lock()
	defer c.l.Unlock()
	return len(c.M)
}

func (c *HostedBlockchains) ContainsFromString(chainHash string) bool {
	c.l.Lock()
	defer c.l.Unlock()
	_, found := c.M[chainHash]
	return found
}

func (c *HostedBlockchains) Clear() {
	c.l.Lock()
	defer c.l.Unlock()
	c.M = make(map[string]HostedBlockchain)
}

func (c *HostedBlockchains) GetChain(hexChain string) (HostedBlockchain, sdk.Error) {
	c.l.Lock()
	defer c.l.Unlock()
	res := c.M[hexChain]
	if res.Hash == "" {
		return HostedBlockchain{}, NewErrorChainNotHostedError(ModuleName)
	}
	return res, nil
}

func (c *HostedBlockchains) GetChainURL(hexChain string) (url string, err sdk.Error) {
	c.l.Lock()
	defer c.l.Unlock()
	res := c.M[hexChain]
	if res.Hash == "" {
		return "", NewErrorChainNotHostedError(ModuleName)
	}
	return res.URL, nil
}

func (c *HostedBlockchains) Validate() error {
	c.l.Lock()
	defer c.l.Unlock()
	for _, chain := range c.M {
		if chain.Hash == "" || chain.URL == "" {
			return NewInvalidHostedChainError(ModuleName)
		}
		if err := HashVerification(chain.Hash); err != nil {
			return err
		}
	}
	return nil
}
