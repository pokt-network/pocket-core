package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pokt-network/pocket-core/util"
	"io/ioutil"
	"os"
	"strings"
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

// "Clear" removes all items from the whitelist
func (w *Whitelist) Clear() {
	(*types.Set)(w).Clear()
}

func (w *Whitelist) ToSlice() []string {
	var res []string
	for entry := range w.M {
		if entry != "" {
			res = append(res, entry.(string))
		}
	}
	return res
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
	w.Clear()
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
	return nil
}

func UpdateWhiteList() error {
	res, arr, err := GetWhiteList()
	if err != nil {
		return err
	}
	// update
	dwl := DWL()
	for _, s := range arr {
		dwl.Clear()
		dwl.Add(s)
	}
	// write devwl
	err = ioutil.WriteFile(config.GlobalConfig().DWL, []byte(res), 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetWhiteList() (string, []string, error) {
	res, err := getWhiteList()
	if err != nil {
		return "", nil, err
	}
	var arr []string
	err = json.Unmarshal([]byte(res), &arr)
	if err != nil {
		return "", nil, err
	}
	return res, arr, nil
}

func getWhiteList() (string, error) {
	if !config.GlobalConfig().Dispatch {
		url := "http://" + config.GlobalConfig().DisIP + ":" + config.GlobalConfig().DisRPort + "/v1/whitelist"
		pl, err := Self()
		if err != nil {
			return "", err
		}
		return util.StructRPCReq(url, pl, util.POST)
	}
	return "", errors.New("dispatch node can't execute this call")
}

// "EnsureWL" cross checks the whitelist for
func EnsureWL(whiteList *Whitelist, query string) bool {
	if index := strings.IndexByte(query, ':'); index > 0 { // delimited by ':'
		query = query[:index]
	}
	if !whiteList.Contains(query) {
		os.Stderr.WriteString("Node: " + query + " rejected because it is not within whitelist. Code: 1\n")
		return false
	}
	return true
	os.Stderr.WriteString("Node: " + query + " rejected because it is not within whitelist. Code: 2\n")
	return false
}
