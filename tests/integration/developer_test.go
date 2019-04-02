package integration

import (
	"flag"
	"github.com/pokt-network/pocket-core/config"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/util"
)

const assumptions = "Integration Testing Assumptions:\n" +
	"1) Dispatcher is hosting a testrpc instance that is labeled as (Blockchain: ETH | NetworkID: 4) in chains.json file\n" +
	"2) Dispatcher has white listed DEVID1 (Dev) and GID1 (SN)\n" +
	"3) Dispatcher is running on DispIP:DisRPort\n" +
	"4) Dispatcher has valid aws credentials for DB test"

const (
	relay    = "relay"
	report   = "report"
	dispatch = "dispatch"
)

var dispatchU *string

func init() {
	dispatchU = flag.String("url", config.GlobalConfig().DisIP+":"+config.GlobalConfig().DisRPort, "the host:port for the test dispatch node")
	flag.Parse()
	*dispatchU = *dispatchU + "/v1/"
}
func requestFromFile(urlSuffix string) (string, error) {
	dispatchU, err := util.URLProto(*dispatchU)
	if err != nil {
		return "", err
	}
	fp, err := filepath.Abs("fixtures" + _const.FILESEPARATOR + urlSuffix + ".json")
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return "", err
	}
	return util.RPCRequ(dispatchU+urlSuffix, b, util.POST)
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
