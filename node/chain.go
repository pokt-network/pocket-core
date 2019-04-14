// This package is node related code.
package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pokt-network/pocket-core/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
)

// A structure that specifies a non-native blockchain.
type Blockchain struct {
	Name  string `json:"name"`
	NetID string `json:"netid"`
}

// A structure that specifies a non-native blockchain client running on a port.
type HostedChain struct {
	Blockchain `json:"blockchain"`
	Port       string `json:"port"`
	Host       string `json:"host"`
	Path       string `json:"path"`   // url path for token based authentication
	Medium     string `json:"medium"` // http, ws, tcp, etc.
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

// "ChainToHosted" returns the hostedChain Object from a blockchain.
func ChainToHosted(b Blockchain) HostedChain {
	mux.Lock()
	defer mux.Unlock()
	return Chains()[b]
}

// "dialHC" attempts to connect to the specific host:port hosting the chain.
func dialHC(u *url.URL) error {
	resp, _ := http.Get(u.String())
	if resp.StatusCode >= 200 {
		return nil
	}
	return errors.New(strconv.Itoa(resp.StatusCode) + " : " + resp.Status)
}

// "TestChains" tests for hosted blockchain clients.
func TestChains() {
	hc := Chains()
	mux.Lock()
	defer mux.Unlock()
	for _, c := range hc {
		s, err := util.URLProto(c.Host + ":" + c.Port)
		if err != nil {
			ExitGracefully(err.Error())
		}
		u, err := url.ParseRequestURI(s)
		if err != nil {
			ExitGracefully(err.Error())
		}
		if c.Path != "" {
			u.Path = c.Path
		}
		if err := dialHC(u); err != nil {
			fmt.Fprint(os.Stderr, c.Name+" client is not detected @ "+u.String()+"\n")
			ExitGracefully(c.Name + " client isn't detected" + "\n")
		}
		fmt.Println(c.Name + " NetID:" + c.NetID + " client is active and ready for service @ " + u.String())
	}
}
