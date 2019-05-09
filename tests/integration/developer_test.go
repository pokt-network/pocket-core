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
	"1) Service node is hosting a testrpc instance that is labeled as (blockchain: ETH | netid: 4) in chains.json file\n" +
	"2) Dispatcher node has white listed DEVID1 (Dev) and GID1 (SN)\n" +
	"3) Dispatcher node is running on DispIP:DisRPort\n" +
	"4) Dispatcher node has valid aws credentials for DB test"

const (
	relay     = "relay"
	report    = "report"
	dispatch  = "dispatch"
	urlstring = "disIP:disrPort"
)

var dispatchU, serviceU *string

func init() {
	dispatchU = flag.String("dispatchurl", urlstring, "the host:port for the test dispatch node")
	serviceU = flag.String("serviceurl", urlstring, "the host:port for the test service node")
	config.Init()
	if *dispatchU == urlstring {
		*dispatchU = config.GlobalConfig().DisIP + ":" + config.GlobalConfig().DisRPort
	}
	if *serviceU == urlstring {
		*serviceU = config.GlobalConfig().DisIP + ":" + config.GlobalConfig().DisRPort
	}
	*dispatchU = *dispatchU + "/v1/"
	*serviceU = *serviceU + "/v1/"
}
func requestFromFile(urlSuffix string) (string, error) {
	dispatchU, err := util.URLProto(*dispatchU)
	serviceU, err := util.URLProto(*serviceU)
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

	if urlSuffix == relay {
		return util.RPCRequ(serviceU+urlSuffix, b, util.POST)
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
