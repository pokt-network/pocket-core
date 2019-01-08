package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type HostedChain struct {
	Name    string `json:"name"`
	NetID   string `json:"netid"`
	Version string `json:"version"`
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
	for _,chain := range *GetHostedChains() {
		if name == chain.Name && netid == chain.NetID && version == chain.Version {
			return chain.Port
		}
	}
	return ""
}

func TestForHostedChains(){
	// TODO runtime test that checks for the hosted chain
}
