package core

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pokt-network/pocket-core/types"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

type Chain struct {
	Hash string
	URL  string
}

type HostedChains types.List

var (
	hostedChains *HostedChains // [Hex Chain Hash] -> Hosted Chain
	chainOnce    sync.Once
)

func GetHostedChains() *HostedChains {
	chainOnce.Do(func() {
		hostedChains = (*HostedChains)(types.NewList())
	})
	return hostedChains
}

func (c *HostedChains) AddChain(chain Chain) {
	(*types.List)(c).Add(chain.Hash, chain)
}

func (c *HostedChains) RemoveChain(chain Chain) {
	(*types.List)(c).Remove(chain.Hash)
}

func (c *HostedChains) Len() int {
	return (*types.List)(c).Count()
}

func (c *HostedChains) ContainsFromObject(chain Chain) bool {
	return (*types.List)(c).Contains(chain.Hash)
}

func (c *HostedChains) ContainsFromBytes(chainHash []byte) bool {
	h := hex.EncodeToString(chainHash)
	return (*types.List)(c).Contains(h)
}

func (c *HostedChains) Clear() {
	(*types.List)(c).Clear()
}

func (c *HostedChains) GetChainFromHex(chainHash string) Chain {
	if c == nil || len(c.M) == 0 {
		return Chain{}
	}
	return (*types.List)(c).Get(chainHash).(Chain)
}

func (c *HostedChains) GetChainFromBytes(chainHash []byte) Chain {
	if c == nil || len(c.M) == 0 {
		return Chain{}
	}
	h := hex.EncodeToString(chainHash)
	return (*types.List)(c).Get(h).(Chain)
}

// "jsonToChains" converts json into chains structure.
func jsonToChains(b []byte) error {
	data := make([]Chain, 0)
	// unmarshal json file into a slice of type Chain
	if err := json.Unmarshal(b, &data); err != nil {
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

// "HostedChainsFromFile" reads a file into the hosted chains object.
func HostedChainsFromFile(filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return jsonToChains(file)
}

// "HostedChainsToFile" converts the hosted chains object into a json file
func HostedChainsToFile(filepath string) error {
	var chainsSlice []Chain
	hc := GetHostedChains()
	for _, chain := range hc.M {
		chainsSlice = append(chainsSlice, chain.(Chain))
	}
	res, err := json.MarshalIndent(chainsSlice, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath, res, 0644)
	if err != nil {
		return err
	}
	return nil
}

// "httpTest" attempts to connect to the specific host:port hosting the chain.
func httpTest(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(strconv.Itoa(resp.StatusCode) + " : " + resp.Status)
	}
	return nil
}

// "TestChains" tests for hosted blockchain clients.
func TestChains() error {
	hc := GetHostedChains()
	hc.Mux.Lock()
	defer hc.Mux.Unlock()
	for _, c := range hc.M {
		if err := httpTest(c.(Chain).URL); err != nil {
			return UnreachableBlockchainErr(c.(Chain).Hash, c.(Chain).URL)
		}
	}
	return nil
}
