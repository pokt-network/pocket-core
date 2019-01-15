package node

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/config"
	"io/ioutil"
	"sync"
)

type Whitelist struct{
	list map[string]struct{}
	mux sync.Mutex
}

var(
	dispatchWL Whitelist
	developerWL Whitelist
	wlOnce sync.Once
)

func WhiteListInit(){
	wlOnce.Do(func(){
		dispatchWL.list = make(map[string]struct{})
		developerWL.list = make(map[string]struct{})
	})
}

func GetDispatchWhitelist() Whitelist{
	if dispatchWL.list == nil {		// just in case
		WhiteListInit()
	}
	return dispatchWL
}

func GetDeveloperWhiteList() Whitelist{
	if developerWL.list == nil {	// just in case
		WhiteListInit()
	}
	return developerWL
}

func (w Whitelist) Contains(s string) bool{
	w.mux.Lock()
	defer w.mux.Unlock()
	_,ok := w.list[s]; return ok
}

func (w Whitelist) Delete(s string){
	w.mux.Lock()
	defer w.mux.Unlock()
	delete(w.list,s)
}

func (w Whitelist) Add (s string){
	w.mux.Lock()
	defer w.mux.Unlock()
	w.list[s]= struct{}{}
}

func (w Whitelist) AddMulti(list []string) {
	for _, v := range list {
		w.Add(v)
	}
}

func (w Whitelist) Size() int {
	w.mux.Lock()
	defer w.mux.Unlock()
	return len(w.list)
}

func DispatchWLFromFile() error{
	return GetDispatchWhitelist().whiteListFromFile(config.GetConfigInstance().DeveloperWL)
}

func DeveloperWLFromFile() error{
	return GetDeveloperWhiteList().whiteListFromFile(config.GetConfigInstance().ServiceNodeWL)
}

func (w Whitelist) whiteListFromFile(filePath string) error {
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	var data []string
	err = json.Unmarshal(f, &data)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	w.AddMulti(data)
	fmt.Println(w)
	return nil
}
