package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/types"
	"io/ioutil"
	"path/filepath"
)

const (
	filep = "tests/bdd/fixtures/chains.json"
)

var (
	ethChain    = types.Blockchain{Name: "eth", NetID: "1", Version: "1"}
	btcChain    = types.Blockchain{Name: "btc", NetID: "1", Version: "1"}
	etcChain    = types.Blockchain{Name: "etc", NetID: "1", Version: "1"}
	bchChain    = types.Blockchain{Name: "bch", NetID: "1", Version: "1"}
	blockchains = []types.Blockchain{ethChain, btcChain, etcChain, bchChain}
)

func main() {
	fmt.Println(createChainsFile())
}

// Creating a chains.json file requires the Pocket Protocol Format
// Currently there is only 2 fields needed for a valid 'chain' in
// a chains.json file:
// 	URL and HASH
// A hash can be calculated by converting the chain into bytes
// using flatbuffers (see common/fbs/blockchain.fbs) and then taking
// the Pocket Protocol Hash of the bytes

func createChainsFile() error {
	var chains []types.Chain
	for _, chain := range blockchains {
		ch, err := types.GenerateChainHash(chain)
		if err != nil {
			return err
		}
		chains = append(chains, types.Chain{Hash: hex.EncodeToString(ch), URL: "test-url"})
	}
	absPath, err := filepath.Abs(filep)
	if err != nil {
		return err
	}
	res, err := json.MarshalIndent(chains, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(absPath, res, 0644)
	if err != nil {
		return err
	}
	return nil
}
