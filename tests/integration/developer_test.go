package integration

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/util"
)

/*
Integration Testing Assumptions:
1) Dispatcher is hosting a testrpc instance that is labeled as (Blockchain: 0 | NetworkID: 0 | Version: 0) in chains.json file
2) Dispatcher has white listed DEVID1 (Dev) and GID1 (SN)
3) Dispatcher is running on DispIP:DisRPort
 */

func TestRelay(t *testing.T) {
	filepath, err := filepath.Abs("fixtures" + _const.FILESEPARATOR + "relay.json")
	if err != nil {
		t.Fatalf(err.Error())
	}
	b, err := ioutil.ReadFile(filepath)
	fmt.Println(string(b))
	if err != nil {
		t.Fatalf(err.Error())
	}
	resp, err := util.RPCRequ("http://"+config.GlobalConfig().DisIP+":"+config.GlobalConfig().DisRPort+"/v1/relay", b, util.POST)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(resp)
}
