package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/core"
	"golang.org/x/crypto/sha3"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	chains = [...]string{"eth", "btc", "ltc", "dash", "aion", "eos"}
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func randomchains() []core.Blockchain {
	var res []core.Blockchain
	rand.Shuffle(len(chains), func(i, j int) { chains[i], chains[j] = chains[j], chains[i] })
	c := chains[rand.Intn(len(chains)-1):]
	res = make([]core.Blockchain, len(c))
	for i, chain := range c {
		res[i] = core.Blockchain{Name: chain, NetID: strconv.Itoa(1), Version: strconv.Itoa(1)}
	}
	return res
}

<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
func createNode(nodeNumber int) common.NodeWorldState {
	hasher := sha256.New()
<<<<<<< 834f04a866176f01faca62715e36541236df2511
	hasher.Write([]byte("node" + strconv.Itoa(i)))
=======
=======
func createNode(nodeNumber int) core.NodeWorldState {
	hasher := sha3.New256()
>>>>>>> servicing milestone addition
	hasher.Write([]byte(node + strconv.Itoa(nodeNumber)))
>>>>>>> cleanup encapsulation
	hash := hex.EncodeToString(hasher.Sum(nil))
<<<<<<< 102119ed0565fdc940c3fef58bc5e981e94649ee
	return common.NodeWorldState{Enode: "enode://" + hash + "@" + strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255)) + ":30303?discport=30301",
		Stake: rand.Intn(255), Active: rand.Intn(2) != 0, IsVal: rand.Intn(2) != 0, Chains: []common.Blockchain{common.Blockchain{Name: "eth", NetID: strconv.Itoa(1), Version: strconv.Itoa(1)}}}
=======
	return core.NodeWorldState{Enode: enodePrefix + hash + enodeDelimiter + strconv.Itoa(rand.Intn(255)) + dot +
		strconv.Itoa(rand.Intn(255)) + dot + strconv.Itoa(rand.Intn(255)) + dot + strconv.Itoa(rand.Intn(255)) + enodeDisport,
		Stake: rand.Intn(255), Active: rand.Intn(2) != 0, IsVal: rand.Intn(2) != 0,
		Chains: []core.Blockchain{{Name: eth, NetID: strconv.Itoa(1), Version: strconv.Itoa(1)}}}
>>>>>>> servicing milestone addition
}

func CreateNodePool(amount int) []core.NodeWorldState {
	var nodePool []core.NodeWorldState
	for i := 0; i < amount; i++ {
		nodePool = append(nodePool, createNode(i))
	}
	return nodePool
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
