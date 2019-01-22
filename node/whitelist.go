package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/types"
)

type Whitelist types.Set

var (
	SNWL   *Whitelist
	DWL    *Whitelist
	wlOnce sync.Once
)

// "WhiteListInit()" initializes both whitelist structures.
func WhiteListInit() {
	wlOnce.Do(func() {
		SNWL = (*Whitelist)(types.NewSet())
		DWL = (*Whitelist)(types.NewSet())
	})
}

// "GetSWL" returns service node white list.
func GetSWL() *Whitelist {
	if SNWL == nil { // just in case
		WhiteListInit()
	}
	return SNWL
}

// "GetDWL" returns developer white list.
func GetDWL() *Whitelist {
	if DWL == nil { // just in case
		WhiteListInit()
	}
	return DWL
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
	return GetSWL().wlFile(config.GetInstance().DWL)
}

// "DWLFile" builds the develoeprs white list from a file.
func DWLFile() error {
	return GetDWL().wlFile(config.GetInstance().SNWL)
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
		fmt.Println("Node: ", query, "rejected because it is not within whitelist")
		return false
	}
	return true
}
