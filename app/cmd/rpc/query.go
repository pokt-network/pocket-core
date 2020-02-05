package rpc

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/app"
	appTypes "github.com/pokt-network/pocket-core/x/apps/types"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	"net/http"
	"strconv"
	"strings"
)

func Version(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	WriteResponse(w, APIVersion, r.URL.Path, r.Host)
}

type heightParams struct {
	Height int64 `json:"height"`
}

type hashParams struct {
	Hash string `json:"hash"`
}

type heightAddrParams struct {
	Height  int64  `json:"height"`
	Address string `json:"address"`
}

type heightAndStakingStatusParams struct {
	Height        int64  `json:"height"`
	StakingStatus string `json:"staking_status"`
}

func Block(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryBlock(params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteJSONResponse(w, string(res), r.URL.Path, r.Host)
}

func Tx(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = hashParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryTx(params.Hash)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
	}
	s, er := json.MarshalIndent(res, "", "  ")
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteResponse(w, string(s), r.URL.Path, r.Host)
}

func Height(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res, err := app.QueryHeight()
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteResponse(w, strconv.Itoa(int(res)), r.URL.Path, r.Host)
}

func Balance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightAddrParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryBalance(params.Address, params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	s, err := app.Codec().MarshalJSON(res)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteResponse(w, string(s), r.URL.Path, r.Host)
}

func Account(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightAddrParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryAccount(params.Address, params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	s, err := app.Codec().MarshalJSON(res)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteResponse(w, string(s), r.URL.Path, r.Host)
}

func Nodes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightAndStakingStatusParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	var res nodeTypes.Validators
	var err error
	switch strings.ToLower(params.StakingStatus) {
	case "":
		// no status passed
		res, err = app.QueryAllNodes(params.Height)
	case "staked":
		// staked nodes
		res, err = app.QueryStakedNodes(params.Height)
	case "unstaked":
		// unstaked nodes
		res, err = app.QueryUnstakedNodes(params.Height)
	case "unstaking":
		// unstaking nodes
		res, err = app.QueryUnstakingNodes(params.Height)
	default:
		panic("invalid staking status, can be staked, unstaked, unstaking, or empty")
	}
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := res.JSON()
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
	}
	WriteResponse(w, string(j), r.URL.Path, r.Host)
}

func Node(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightAddrParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryNode(params.Address, params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := res.MarshalJSON()
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteResponse(w, string(j), r.URL.Path, r.Host)
}

func NodeParams(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryNodeParams(params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := app.Codec().MarshalJSON(res)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteResponse(w, string(j), r.URL.Path, r.Host)
}

func NodeProofs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightAddrParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryProofs(params.Address, params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := app.Codec().MarshalJSON(res)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteResponse(w, string(j), r.URL.Path, r.Host)
}

type queryNodeProof struct {
	Address      string `json:"address"`
	Blockchain   string `json:"blockchain"`
	AppPubKey    string `json:"app_pubkey"`
	SBlockHeight int64  `json:"session_block_height"`
	Height       int64  `json:"height"`
}

func NodeProof(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = queryNodeProof{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryProof(params.Blockchain, params.AppPubKey, params.Address, params.SBlockHeight, params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := app.Codec().MarshalJSON(res)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteResponse(w, string(j), r.URL.Path, r.Host)
}

func Apps(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightAndStakingStatusParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	var res appTypes.Applications
	var err error
	switch strings.ToLower(params.StakingStatus) {
	case "":
		// no status passed
		res, err = app.QueryAllApps(params.Height)
	case "staked":
		// staked nodes
		res, err = app.QueryStakedApps(params.Height)
	case "unstaked":
		// unstaked nodes
		res, err = app.QueryUnstakedApps(params.Height)
	case "unstaking":
		// unstaking nodes
		res, err = app.QueryUnstakingApps(params.Height)
	default:
		panic("invalid staking status, can be staked, unstaked, unstaking, or empty")
	}
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := res.JSON()
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
	}
	WriteResponse(w, string(j), r.URL.Path, r.Host)
}

func App(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightAddrParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryApp(params.Address, params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := res.MarshalJSON()
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteResponse(w, string(j), r.URL.Path, r.Host)
}

func AppParams(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryAppParams(params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := app.Codec().MarshalJSON(res)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteResponse(w, string(j), r.URL.Path, r.Host)
}

func PocketParams(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryPocketParams(params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := app.Codec().MarshalJSON(res)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteResponse(w, string(j), r.URL.Path, r.Host)
}

func SupportedChains(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryPocketSupportedBlockchains(params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := app.Codec().MarshalJSON(res)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteResponse(w, string(j), r.URL.Path, r.Host)
}

type querySupplyResponse struct {
	NodeStaked    int64 `json:"node_staked"`
	AppStaked     int64 `json:"app_staked"`
	Dao           int64 `json:"dao"`
	TotalStaked   int64 `json:"total_staked"`
	TotalUnstaked int64 `json:"total_unstaked"`
	Total         int64 `json:"total"`
}

func Supply(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	nodesStake, nodesUnstaked, err := app.QueryTotalNodeCoins(params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	appsStaked, _, err := app.QueryTotalAppCoins(params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	dao, err := app.QueryDaoBalance(params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	totalStaked := nodesStake.Add(appsStaked).Add(dao)
	totalUnstaked := nodesUnstaked
	total := totalStaked.Add(totalUnstaked)
	res, err := json.MarshalIndent(&querySupplyResponse{
		NodeStaked:    nodesStake.Int64(),
		AppStaked:     appsStaked.Int64(),
		Dao:           dao.Int64(),
		TotalStaked:   totalStaked.Int64(),
		TotalUnstaked: totalUnstaked.Int64(),
		Total:         total.Int64(),
	}, "", "  ")
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteJSONResponse(w, string(res), r.URL.Path, r.Host)
}
