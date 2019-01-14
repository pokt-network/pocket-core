package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"sync"
)

type HostedChain struct {
	Blockchain	   `json:"blockchain"`
	Port    string `json:"port"`
	Medium  string `json:"medium"`
}

var (
	hostedChains *[]HostedChain
	once sync.Once
	mux sync.Mutex
)

func UnmarshalChains(b []byte) error{
	h := GetHostedChains()
	mux.Lock()
	defer mux.Unlock()
	if err := json.Unmarshal(b, &h); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetHostedChains() *[]HostedChain{
	once.Do(func(){
		h := make([]HostedChain,0); hostedChains = &h
	})
	return hostedChains
}

func ExportHostedChains() ([]byte, error){
	mux.Lock()
	defer mux.Unlock()
	b, err := json.Marshal(GetHostedChains())
	return b, err
}

func HostedChainsFile(filepath string) error {
	file, err := ioutil.ReadFile(filepath)			// read the file from the specified path
	if err != nil {									// if error
		fmt.Println(err)
		return err
	}
	return UnmarshalChains(file) 					// call manPeers.Json on the byte[]
}

func GetHostedChainPort(name string, netid string, version string) string {	// TODO optimize this currently O(n)
// TODO this can be compared with a hosted chain structure so no need to do individual comparison!
	for _,chain := range *GetHostedChains() {
		if name == chain.Name && netid == chain.NetID && version == chain.Version {
			return chain.Port
		}
	}
	return ""
}

func TestForHostedChains() bool{
	hc := GetHostedChains()
	for _,chain := range *hc {
		if err := pingPort(chain.Port); err!= nil {
			fmt.Println(chain.Name, " client is not detected on port ", chain.Port)
			return false
		}
	}
	return true
}

func pingPort(port string) error{
	_, err := net.Listen("tcp", ":" + port)
	if err != nil {
		return nil
	}
	return errors.New("port: "+port+" is not in use")
}
