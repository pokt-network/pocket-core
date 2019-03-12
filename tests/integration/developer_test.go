package integration

import (
	"flag"
	"fmt"
	"github.com/pokt-network/pocket-core/config"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/util"
)

const assumptions = "Integration Testing Assumptions:\n" +
	"1) Dispatcher is hosting a testrpc instance that is labeled as (Blockchain: ETH | NetworkID: 4 | Version: 0) in chains.json file\n" +
	"2) Dispatcher has white listed DEVID1 (Dev) and GID1 (SN)\n" +
	"3) Dispatcher is running on DispIP:DisRPort\n" +
	"4) Dispatcher has valid aws credentials for DB test"

type PORT int

const (
	Relay PORT = iota
	Client

	relay    = "relay"
	report   = "report"
	dispatch = "dispatch"
)
var dispatchU, serviceU *string
func init(){
	dispatchU = flag.String("dispatchtesturl", config.GlobalConfig().DisIP, "the host:port for the test dispatch node")
	serviceU = flag.String("servicetesturl", config.GlobalConfig().DisIP, "the host:port for the test service node")
	flag.Parse()
}
func requestFromFile(urlSuffix string) (string, error) {
	const http = "http://"
	if !strings.Contains(*dispatchU, http) {
		*dispatchU = http + *dispatchU+"/v1/"
	}
	if !strings.Contains(*serviceU, http) {
		*serviceU = http + *serviceU+"/v1/"
	}
	fp, err := filepath.Abs("fixtures" + _const.FILESEPARATOR + urlSuffix + ".json")
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return "", err
	}
	switch urlSuffix {

	}
	if urlSuffix == relay {
		return util.RPCRequ(*serviceU+urlSuffix, b, util.POST)
	}
	return util.RPCRequ(*dispatchU+urlSuffix, b, util.POST)
}

func TestRelay(t *testing.T) {
	resp, err := requestFromFile(relay)
	if err != nil {
		t.Log(assumptions)
		t.Fatalf(err.Error())
	}
	t.Log(resp)
}

func TestReport(t *testing.T) {
	resp, err := requestFromFile(report)
	if err != nil {
		t.Log(assumptions)
		t.Fatalf(err.Error())
	}
	t.Log(resp)
}

func TestDispatch(t *testing.T) {
	resp, err := requestFromFile(dispatch)
	if err != nil {
		t.Log(assumptions)
		t.Fatalf(err.Error())
	}
	t.Log(resp)
}
