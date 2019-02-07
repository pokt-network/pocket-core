package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"sync"
)

// A structure that specifies a non-native blockchain client running on a port.
type HostedChain struct {
	Blockchain `json:"blockchain"` // blockchain structure
	Port       string              `json:"port"`   // port that the client is running on
	Medium     string              `json:"medium"` // http, ws, tcp, etc.
}

var (
	// A globally acessed structure that holds the HostedChains structures.
	chains map[Blockchain]HostedChain
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

func ChainsSlice() []Blockchain {
	cs := make([]Blockchain, 0)
	for k := range Chains() {
		cs = append(cs, k)
	}
	return cs
}

// "ExportChains" converts chains into json.
func ExportChains() ([]byte, error) {
	mux.Lock()
	defer mux.Unlock()
	return json.Marshal(Chains())
}

// "UnmarshalChains" converts json into chains.
func UnmarshalChains(b []byte) error {
	h := Chains()
	data := make([]HostedChain, 0)
	mux.Lock()
	defer mux.Unlock()
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	// convert slice into map for quick access
	// this is o(n) but alternative is a more complicated chains.json file
	for _, hc := range data {
		h[hc.Blockchain] = hc
	}
	return nil
}

// "CFile" reads a file into chains.
func CFIle(filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return UnmarshalChains(file)
}

// "ChainPort" returns the port of a blockchain client.
func ChainPort(b Blockchain) string {
	mux.Lock()
	defer mux.Unlock()
	return Chains()[b].Port
}

// "TestChains" tests for hosted blockchain clients.
func TestChains() bool {
	hc := Chains()
	mux.Lock()
	defer mux.Unlock()
	for _, c := range hc {
		if err := pingPort(c.Port); err != nil {
			fmt.Println(c.Name, " client is not detected on port ", c.Port)
			return false
		}
	}
	return true
}

// "pingPort" attempts to connect to the specific port hosting the chain.
func pingPort(port string) error {
	_, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil
	}
	return errors.New("port: " + port + " is not in use")
}
