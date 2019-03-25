package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	chains = [...]string{"eth", "btc", "ltc", "dash", "aion", "eos"}
)

type NodeWorldState struct {
	Enode  string   `json:"enode"`
	Stake  int      `json:"stake"`
	Active bool     `json:"status"`
	IsVal  bool     `json:"isval"`
	Chains []string `json:"chains"`
}

func RandomChains() []string {
	rand.Shuffle(len(chains), func(i, j int) { chains[i], chains[j] = chains[j], chains[i] })
	return chains[rand.Intn(len(chains)-1):]
}

func CreateNode(i int) NodeWorldState {
	hasher := sha1.New()
	hasher.Write([]byte("node"+strconv.Itoa(i)))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return NodeWorldState{"enode://"+hash+"@" + strconv.Itoa(rand.Intn(100)) + "." + strconv.Itoa(rand.Intn(100)) + "." + strconv.Itoa(rand.Intn(100)) + "." + strconv.Itoa(rand.Intn(100)) + ":30303?discport=30301",
		rand.Intn(100), rand.Intn(2) != 0, rand.Intn(2) != 0, RandomChains()}
}

func CreateNodePool(amount int) []NodeWorldState {
	var nodePool []NodeWorldState
	for i := 0; i < amount; i++ {
		nodePool = append(nodePool, CreateNode(i))
	}
	return nodePool
}
func init(){
	rand.Seed(time.Now().UTC().UnixNano())
}
func main() {
	absPath, _ := filepath.Abs("tests/bdd/fixtures/nodepool.json")
	nodePool := CreateNodePool(500)
	b, err := json.MarshalIndent(nodePool, "", "    ")
	if err != nil{
		fmt.Println(err.Error())
	}
	fmt.Println(absPath)
	f, err:=os.Create(absPath)
	if err != nil{
		fmt.Println(err.Error())
	}
	_, err=f.Write(b)
	if err != nil{
		fmt.Println(err.Error())
	}
}
