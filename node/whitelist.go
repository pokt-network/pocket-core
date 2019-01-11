package node

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/config"
	"io/ioutil"
	"sync"
)

type Whitelist map[string]struct{} // gid and 0 space (essentially a set)

var(
	dispatchWL Whitelist
	dispatchWLOnce sync.Once
	dispatchWLMux sync.Mutex
)

func GetDispatchWhitelist() Whitelist{
	dispatchWLOnce.Do(func() {
		dispatchWL = make(map[string]struct{})
	})
	return dispatchWL
}

func (w Whitelist) Contains(s string) bool{
	dispatchWLMux.Lock()
	defer dispatchWLMux.Unlock()
	_,ok := w[s]; return ok
}

func (w Whitelist) Delete(s string){
	dispatchWLMux.Lock()
	defer dispatchWLMux.Unlock()
	delete(w,s)
}

func (w Whitelist) Add (s string){
	dispatchWLMux.Lock()
	defer dispatchWLMux.Unlock()
	w[s]= struct{}{}
}

func (w Whitelist) AddMulti(list []string) {
	for _, v := range list {
		w.Add(v)
	}
}

func (w Whitelist) Size() int {
	dispatchWLMux.Lock()
	defer dispatchWLMux.Unlock()
	return len(w)
}

func WhitelistFromFile() error{
	plan, err := ioutil.ReadFile(config.GetConfigInstance().Whitelist)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	var data []string
	err = json.Unmarshal(plan, &data)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	GetDispatchWhitelist().AddMulti(data)
	fmt.Println(GetDispatchWhitelist())
	return nil
}
