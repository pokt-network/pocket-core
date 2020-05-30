package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/app/cmd/rpc"
)

var (
	SendRawTxPath,
	GetNodePath,
	GetACLPath,
	GetUpgradePath,
	GetDAOOwnerPath,
	GetHeightPath,
	GetAccountPath,
	GetAppPath,
	GetTxPath,
	GetBlockPath,
	GetSupportedChainsPath,
	GetBalancePath,
	GetAccountTxsPath,
	GetNodeParamsPath,
	GetNodesPath,
	GetAppsPath,
	GetAppParamsPath,
	GetPocketParamsPath,
	GetNodeReceiptPath,
	GetNodeReceiptsPath,
	GetNodeClaimsPath,
	GetNodeClaimPath,
	GetBlockTxsPath,
	GetSupplyPath,
	GetAllParamsPath,
	GetParamPath string
)

func init() {
	routes := rpc.GetRoutes()
	for _, route := range routes {
		switch route.Name {
		case "SendRawTx":
			SendRawTxPath = route.Path
		case "QueryNode":
			GetNodePath = route.Path
		case "QueryACL":
			GetACLPath = route.Path
		case "QueryUpgrade":
			GetUpgradePath = route.Path
		case "QueryDAO":
			GetDAOOwnerPath = route.Path
		case "QueryHeight":
			GetHeightPath = route.Path
		case "QueryAccount":
			GetAccountPath = route.Path
		case "QueryApp":
			GetAppPath = route.Path
		case "QueryTX":
			GetTxPath = route.Path
		case "QueryBlock":
			GetBlockPath = route.Path
		case "QuerySupportedChains":
			GetSupportedChainsPath = route.Path
		case "QueryBalance":
			GetBalancePath = route.Path
		case "QueryAccountTxs":
			GetAccountTxsPath = route.Path
		case "QueryNodeParams":
			GetNodeParamsPath = route.Path
		case "QueryNodes":
			GetNodesPath = route.Path
		case "QueryApps":
			GetAppsPath = route.Path
		case "QueryAppParams":
			GetAppParamsPath = route.Path
		case "QueryPocketParams":
			GetPocketParamsPath = route.Path
		case "QueryNodeReceipt":
			GetNodeReceiptPath = route.Path
		case "QueryNodeReceipts":
			GetNodeReceiptsPath = route.Path
		case "QueryBlockTxs":
			GetBlockTxsPath = route.Path
		case "QuerySupply":
			GetSupplyPath = route.Path
		case "QueryNodeClaim":
			GetNodeClaimPath = route.Path
		case "QueryNodeClaims":
			GetNodeClaimsPath = route.Path
		case "QueryAllParams":
			GetAllParamsPath = route.Path
		case "QueryParam":
			GetParamPath = route.Path
		default:
			continue
		}
	}
}

func QueryRPC(path string, jsonArgs []byte) (string, error) {
	//cliURL := app.GlobalConfig.PocketConfig.RemoteCLIURL + ":" + app.GlobalConfig.PocketConfig.RPCPort + path
	cliURL := app.GlobalConfig.PocketConfig.RemoteCLIURL + path
	fmt.Println(cliURL)
	req, err := http.NewRequest("POST", cliURL, bytes.NewBuffer(jsonArgs))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	res, err := strconv.Unquote(string(bz))
	if err == nil {
		bz = []byte(res)
	}
	if resp.StatusCode == http.StatusOK {
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, bz, "", "    ")
		if err == nil {
			return prettyJSON.String(), nil
		}
		return string(bz), nil
	}
	return "", fmt.Errorf("the http status code was not okay: %d, and the status was: %s, with a response of %v", resp.StatusCode, resp.Status, string(bz))
}
