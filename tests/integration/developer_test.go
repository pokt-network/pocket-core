package integration

import (
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
4) Dispatcher has valid aws credentials for DB test
 */
var dispatchURL = "http://" + config.GlobalConfig().DisIP + ":"

type PORT int

const (
	Relay PORT = iota
	Client
)

func requestFromFile(urlSuffix string, port PORT) (string, error) {
	switch port {
	case Relay:
		dispatchURL = dispatchURL + config.GlobalConfig().DisRPort + "/v1/"
	case Client:
		dispatchURL = dispatchURL + config.GlobalConfig().DisCPort + "/v1/"
	}
	fp, err := filepath.Abs("fixtures" + _const.FILESEPARATOR + urlSuffix + ".json")
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return "", err
	}
	return util.RPCRequ(dispatchURL+urlSuffix, b, util.POST)
}

func TestRelay(t *testing.T) {
	resp, err := requestFromFile("relay", Relay)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(resp)
}

func TestReport(t *testing.T) {
	resp, err := requestFromFile("report", Relay)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(resp)
}

func TestDispatch(t *testing.T) {
	resp, err := requestFromFile("dispatch", Relay)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(resp)
}
