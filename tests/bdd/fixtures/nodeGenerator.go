package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/common"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	chains = [...]string{"eth", "btc", "ltc", "dash", "aion", "eos"}
)

func RandomChains() []common.Blockchain {
	var res []common.Blockchain
	rand.Shuffle(len(chains), func(i, j int) { chains[i], chains[j] = chains[j], chains[i] })
	c := chains[rand.Intn(len(chains)-1):]
	res = make([]common.Blockchain, len(c))
	for i, chain := range c {
		res[i] = common.Blockchain{Name: chain, NetID: strconv.Itoa(rand.Intn(4)), Version: strconv.Itoa(rand.Intn(4))}
	}
	return res
}

func CreateNode(i int) common.NodeWorldState {
	hasher := sha256.New()
	hasher.Write([]byte("node" + strconv.Itoa(i)))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return common.NodeWorldState{Enode: "enode://" + hash + "@" + strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255)) + ":30303?discport=30301",
		Stake: rand.Intn(255), Active: rand.Intn(2) != 0, IsVal: rand.Intn(2) != 0, Chains: RandomChains()}
}

func CreateNodePool(amount int) []common.NodeWorldState {
	var nodePool []common.NodeWorldState
	for i := 0; i < amount; i++ {
		nodePool = append(nodePool, CreateNode(i))
	}
	return nodePool
}
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
func main() {
	prefix := []string{"xsmall", "small", "medium", "large"}
	sizes := []int{25, 500, 5000, 50000}
	for i := 0; i < 3; i++ { // don't run the large one for now
		absPath, _ := filepath.Abs("tests/bdd/fixtures/" + prefix[i] + "nodepool.json")
		nodePool := CreateNodePool(sizes[i])
		b, err := json.MarshalIndent(nodePool, "", "    ")
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(absPath)
		f, err := os.Create(absPath)
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = f.Write(b)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
