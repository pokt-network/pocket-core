package rpc

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/app"
	appTypes "github.com/pokt-network/pocket-core/x/apps/types"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	"math/big"
	"net/http"
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

type heightAndValidatorsOptsParams struct {
	Height int64                           `json:"height"`
	Opts   nodeTypes.QueryValidatorsParams `json:"opts"`
}

type heightAndApplicationsOptsParams struct {
	Height int64                              `json:"height"`
	Opts   appTypes.QueryApplicationsWithOpts `json:"opts"`
}

type heightAndStakingStatusParams struct {
	Height        int64  `json:"height"`
	StakingStatus string `json:"staking_status"`
	Page          int    `json:"page,omitempty"`
	PerPage       int    `json:"per_page,omitempty"`
}

type paginatedAddressParams struct {
	Address  string `json:"address"`
	Page     int    `json:"page,omitempty"`
	PerPage  int    `json:"per_page,omitempty"`
	Received bool   `json:"received,omitempty"`
	Prove    bool   `json:"prove,omitempty"`
}

type paginatedHeightParams struct {
	Height  int64 `json:"height"`
	Page    int   `json:"page,omitempty"`
	PerPage int   `json:"per_page,omitempty"`
	Prove   bool  `json:"prove,omitempty"`
}

func Block(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryBlock(&params.Height)
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
	WriteJSONResponse(w, string(s), r.URL.Path, r.Host)
}

func AccountTxs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = paginatedAddressParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	var res *core_types.ResultTxSearch
	var err error
	if params.Received == false {
		res, err = app.QueryAccountTxs(params.Address, params.Page, params.PerPage, params.Prove)
	} else {
		res, err = app.QueryRecipientTxs(params.Address, params.Page, params.PerPage, params.Prove)
	}

	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
	}
	s, er := json.MarshalIndent(res, "", "  ")
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteJSONResponse(w, string(s), r.URL.Path, r.Host)
}

func BlockTxs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = paginatedHeightParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryBlockTxs(params.Height, params.Page, params.PerPage, params.Prove)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
	}
	s, er := json.MarshalIndent(res, "", "  ")
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteJSONResponse(w, string(s), r.URL.Path, r.Host)
}

type queryHeightResponse struct {
	Height int64 `json:"height"`
}

func Height(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res, err := app.QueryHeight()
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	height, err := json.Marshal(&queryHeightResponse{Height: res})
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteJSONResponse(w, string(height), r.URL.Path, r.Host)
}

type queryBalanceResponse struct {
	Balance *big.Int `json:"balance"`
}

func Balance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightAddrParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	balance, err := app.QueryBalance(params.Address, params.Height)

	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	s, err := json.MarshalIndent(&queryBalanceResponse{Balance: balance.BigInt()}, "", "")
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteJSONResponse(w, string(s), r.URL.Path, r.Host)
}

func Account(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightAddrParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryAccount(params.Address, params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	s, err := json.Marshal(res)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteJSONResponse(w, string(s), r.URL.Path, r.Host)
}

func Nodes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightAndValidatorsOptsParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryNodes(params.Height, params.Opts)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := res.JSON()
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(j)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
	}
}

func Node(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightAddrParams{Height: 0}
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
	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

func NodeParams(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightParams{Height: 0}
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
	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

func NodeReceipts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightAddrParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryReceipts(params.Address, params.Height)
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

type queryNodeReceipts struct {
	Address      string `json:"address"`
	Blockchain   string `json:"blockchain"`
	AppPubKey    string `json:"app_pubkey"`
	SBlockHeight int64  `json:"session_block_height"`
	Height       int64  `json:"height"`
	ReceiptType  string `json:"receipt_type"`
}

func NodeReceipt(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = queryNodeReceipts{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryReceipt(params.Blockchain, params.AppPubKey, params.Address, params.ReceiptType, params.SBlockHeight, params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := app.Codec().MarshalJSON(res)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

func Apps(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightAndApplicationsOptsParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryApps(params.Height, params.Opts)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := res.JSON()
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(j)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
	}
}

func App(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightAddrParams{Height: 0}
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
	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

func AppParams(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightParams{Height: 0}
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
	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

func PocketParams(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightParams{Height: 0}
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
	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

func SupportedChains(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightParams{Height: 0}
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
	NodeStaked    string `json:"node_staked"`
	AppStaked     string `json:"app_staked"`
	Dao           string `json:"dao"`
	TotalStaked   string `json:"total_staked"`
	TotalUnstaked string `json:"total_unstaked"`
	Total         string `json:"total"`
}

func Supply(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	nodesStake, total, err := app.QueryTotalNodeCoins(params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	appsStaked, err := app.QueryTotalAppCoins(params.Height)
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
	totalUnstaked := total.Sub(totalStaked)
	res, err := json.MarshalIndent(&querySupplyResponse{
		NodeStaked:    nodesStake.String(),
		AppStaked:     appsStaked.String(),
		Dao:           dao.String(),
		TotalStaked:   totalStaked.BigInt().String(),
		TotalUnstaked: totalUnstaked.BigInt().String(),
		Total:         total.BigInt().String(),
	}, "", "  ")
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteJSONResponse(w, string(res), r.URL.Path, r.Host)
}

func DAOOwner(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryDaoOwner(0)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	s, err := json.Marshal(res)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteResponse(w, string(s), r.URL.Path, r.Host)
}

func Upgrade(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryUpgrade(0)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	s, err := json.Marshal(res)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteResponse(w, string(s), r.URL.Path, r.Host)
}

func ACL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = heightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryACL(params.Height)
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

func State(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res, err := app.ExportState()
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteRaw(w, string(res), r.URL.Path, r.Host)
}
