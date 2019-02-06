package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/types"
)

type Whitelist types.Set

var (
	SNWL   *Whitelist
	DevWL  *Whitelist
	wlOnce sync.Once
)

// "WhiteListInit()" initializes both whitelist structures.
func WhiteListInit() {
	wlOnce.Do(func() {
		SNWL = (*Whitelist)(types.NewSet())
		DevWL = (*Whitelist)(types.NewSet())
	})
}

// "SWL" returns service node white list.
func SWL() *Whitelist {
	if SNWL == nil { // just in case
		WhiteListInit()
	}
	return SNWL
}

// "DWL" returns developer white list.
func DWL() *Whitelist {
	if DevWL == nil { // just in case
		WhiteListInit()
	}
	return DevWL
}

// "Contains" returns if within whitelist.
func (w *Whitelist) Contains(s string) bool {
	return (*types.Set)(w).Contains(s)
}

// "Remove" removes item from whitelist.
func (w *Whitelist) Remove(s string) {
	(*types.Set)(w).Remove(s)
}

// "Add" appends item to whitelist.
func (w *Whitelist) Add(s string) {
	(*types.Set)(w).Add(s)
}

// "AddMulti" appends multiple items to whitelist
func (w *Whitelist) AddMulti(list []string) {
	w.Mux.Lock()
	defer w.Mux.Unlock()
	for _, v := range list {
		w.M[v] = struct{}{}
	}
}

// "Count" returns the length of the whitelist.
func (w *Whitelist) Count() int {
	return (*types.Set)(w).Count()
}

// "SWLFile" builds the service white list from a file.
func SWLFile() error {
	return SWL().wlFile(config.GlobalConfig().SNWL)
}

// "DWLFile" builds the develoeprs white list from a file.
func DWLFile() error {
	return DWL().wlFile(config.GlobalConfig().DWL)
}

// "wlFile" builds a whitelist structure from a file.
func (w *Whitelist) wlFile(filePath string) error {
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

// "EnsureWL" cross checks the whitelist for
func EnsureWL(whiteList *Whitelist, query string) bool {
	if !whiteList.Contains(query) {
		os.Stderr.WriteString("Node: " + query + " rejected because it is not within whitelist\n")
		fmt.Println("NOTE: if you are a developer, just add a service_whitelist.json file to " + config.GlobalConfig().DD)
		fmt.Println("An example of this file is under the docs directory")
		return false
	}
	return true
}
