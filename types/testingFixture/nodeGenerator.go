package main

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	"io/ioutil"
	"math/rand"
	"time"
)

const (
	numberOfNodes = 50
)

// writes json nodepool for testing
func main() {
	types.RegisterPOKT()
	var result []types.Node
	fmt.Println()
	for i := 0; i < numberOfNodes; i++ {
		result = append(result, generateAliveNode())
	}
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = ioutil.WriteFile("types/testingFixture/randomPool.json", output, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func generateAliveNode() (node types.Node) {
	randomSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randomSource)
	balance := types.NewPOKT(random.Int63())
	stakeAmount := types.NewPOKT(random.Int63())
	_, pubKey := crypto.NewKeypair()
	node = types.Node{
		Account: types.Account{
			Address:     nil, // todo
			PubKey:      pubKey,
			Balance:     balance,
			StakeAmount: stakeAmount,
		},
		URL:             nil, // todo
		SupportedChains: nil, // todo
		IsAlive:         true,
	}

	return
}
