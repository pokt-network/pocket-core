// This package is node related code.
package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sync"

	"github.com/pokt-network/pocket-core/util"
)

// A structure that specifies a non-native blockchain.
type Blockchain struct {
	Name    string `json:"name"`
	NetID   string `json:"netid"`
	Version string `json:"version"`
}

// A structure that specifies a non-native blockchain client running on a port.
type HostedChain struct {
	Blockchain `json:"blockchain"` // blockchain structure
	Port       string              `json:"port"`   // port that the client is running on
	Medium     string              `json:"medium"` // http, ws, tcp, etc.
}

var (
	chains map[Blockchain]HostedChain // A structure to hold the hosted chains of the client.
	once   sync.Once
	mux    sync.Mutex
)

// "Chains" is the singleton accessor for chains.
func Chains() map[Blockchain]HostedChain {
	once.Do(func() {
		chains = make(map[Blockchain]HostedChain)
	})
	return chains
}

// "ChainsSlice" converts the chains structure into a slice of type Blockchain.
func ChainsSlice() []Blockchain {
	cs := make([]Blockchain, 0)
	for k := range Chains() {
		cs = append(cs, k)
	}
	return cs
}

// "jsonToChains" converts json into chains structure.
func jsonToChains(b []byte) error {
	h := Chains()
	data := make([]HostedChain, 0)
	mux.Lock()
	defer mux.Unlock()
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	for _, hc := range data {
		h[hc.Blockchain] = hc
	}
	return nil
}

// "CFile" reads a file into chains.
func CFile(filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return jsonToChains(file)
}

// "ChainPort" returns the port of a blockchain client.
func ChainPort(b Blockchain) string {
	mux.Lock()
	defer mux.Unlock()
	return Chains()[b].Port
}

// "pingPort" attempts to connect to the specific port hosting the chain.
func pingPort(port string) error {
	_, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil
	}
	return errors.New("port: " + port + " is not in use")
}

// "TestChains" tests for hosted blockchain clients.
func TestChains() {
	hc := Chains()
	mux.Lock()
	defer mux.Unlock()
	for _, c := range hc {
		if err := pingPort(c.Port); err != nil {
			fmt.Fprint(os.Stderr, c.Name+" client is not detected on port "+c.Port)
			util.ExitGracefully(c.Name + " client isn't detected")
		}
		fmt.Println(c.Name + " V:" + c.Version + " NetID:" + c.NetID + " client is active and ready for service on port " + c.Port)
	}
}
