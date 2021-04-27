package rpc

import (
	"encoding/hex"
	"encoding/json"
	sdk "github.com/pokt-network/pocket-core/types"
	"math/big"
	"net/http"

	types2 "github.com/pokt-network/pocket-core/x/auth/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/types"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/app"
	appTypes "github.com/pokt-network/pocket-core/x/apps/types"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
)

func Version(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	WriteResponse(w, APIVersion, r.URL.Path, r.Host)
}

type HeightParams struct {
	Height int64 `json:"height"`
}
type HeightAndKeyParams struct {
	Height int64  `json:"height"`
	Key    string `json:"key"`
}

type HashAndProveParams struct {
	Hash  string `json:"hash"`
	Prove bool   `json:"prove"`
}

type HeightAndAddrParams struct {
	Height  int64  `json:"height"`
	Address string `json:"address"`
}

type HeightAndValidatorOptsParams struct {
	Height int64                           `json:"height"`
	Opts   nodeTypes.QueryValidatorsParams `json:"opts"`
}

type HeightAndApplicaitonOptsParams struct {
	Height int64                              `json:"height"`
	Opts   appTypes.QueryApplicationsWithOpts `json:"opts"`
}

type PaginateAddrParams struct {
	Address  string `json:"address"`
	Page     int    `json:"page,omitempty"`
	PerPage  int    `json:"per_page,omitempty"`
	Received bool   `json:"received,omitempty"`
	Prove    bool   `json:"prove,omitempty"`
	Sort     string `json:"order,omitempty"`
}

type PaginatedHeightParams struct {
	Height  int64  `json:"height"`
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
	Prove   bool   `json:"prove,omitempty"`
	Sort    string `json:"order,omitempty"`
}

type PaginatedHeightAndAddrParams struct {
	Height  int64  `json:"height"`
	Addr    string `json:"address"`
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
}

func Block(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = HeightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.PCA.QueryBlock(&params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteJSONResponse(w, string(res), r.URL.Path, r.Host)
}

func Tx(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = HashAndProveParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.PCA.QueryTx(params.Hash, params.Prove)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	rpcResponse := ResultTxToRPC(res)
	s, er := json.MarshalIndent(rpcResponse, "", "  ")
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteJSONResponse(w, string(s), r.URL.Path, r.Host)
}

// Result of searching for txs
type RPCResultTxSearch struct {
	Txs        []*RPCResultTx `json:"txs"`
	TotalCount int            `json:"total_count"`
}

// Result of querying for a tx
type RPCResultTx struct {
	Hash     bytes.HexBytes       `json:"hash"`
	Height   int64                `json:"height"`
	Index    uint32               `json:"index"`
	TxResult RPCResponseDeliverTx `json:"tx_result"`
	Tx       types.Tx             `json:"tx"`
	Proof    types.TxProof        `json:"proof,omitempty"`
	StdTx    RPCStdTx             `json:"stdTx"`
}

type RPCResponseDeliverTx struct {
	Code        uint32        `json:"code,omitempty"`
	Data        []byte        `json:"data,omitempty"`
	Log         string        `json:"log,omitempty"`
	Info        string        `json:"info,omitempty"`
	Events      []abci.Event  `json:"events,omitempty"`
	Codespace   string        `json:"codespace,omitempty"`
	Signer      types.Address `json:"signer,omitempty"`
	Recipient   types.Address `json:"recipient,omitempty"`
	MessageType string        `json:"message_type,omitempty"`
}

type RPCStdTx types2.StdTx

type rPCStdTx struct {
	Msg       json.RawMessage `json:"msg" yaml:"msg"`
	Fee       sdk.Coins       `json:"fee" yaml:"fee"`
	Signature RPCStdSignature `json:"signature" yaml:"signature"`
	Memo      string          `json:"memo" yaml:"memo"`
	Entropy   int64           `json:"entropy" yaml:"entropy"`
}

type RPCStdSignature struct {
	PublicKey string `json:"pub_key"`
	Signature string `json:"signature"`
}

func (r RPCStdTx) MarshalJSON() ([]byte, error) {
	msgBz := (types2.StdTx)(r).Msg.GetSignBytes()
	sig := RPCStdSignature{
		PublicKey: r.Signature.RawString(),
		Signature: hex.EncodeToString(r.Signature.Signature),
	}
	return json.Marshal(rPCStdTx{
		Msg:       msgBz,
		Fee:       r.Fee,
		Signature: sig,
		Memo:      r.Memo,
		Entropy:   r.Entropy,
	})
}

func ResultTxSearchToRPC(res *core_types.ResultTxSearch) RPCResultTxSearch {
	if res == nil {
		return RPCResultTxSearch{}
	}
	rpcTxSearch := RPCResultTxSearch{
		Txs:        make([]*RPCResultTx, 0, res.TotalCount),
		TotalCount: res.TotalCount,
	}
	for _, result := range res.Txs {
		rpcTxSearch.Txs = append(rpcTxSearch.Txs, ResultTxToRPC(result))
	}
	return rpcTxSearch
}

func ResultTxToRPC(res *core_types.ResultTx) *RPCResultTx {
	if res == nil {
		return nil
	}
	tx := app.UnmarshalTx(res.Tx, res.Height)
	if app.GlobalConfig.PocketConfig.DisableTxEvents {
		res.TxResult.Events = nil
	}
	rpcDeliverTx := RPCResponseDeliverTx{
		Code:        res.TxResult.Code,
		Data:        res.TxResult.Data,
		Log:         res.TxResult.Log,
		Info:        res.TxResult.Info,
		Events:      res.TxResult.Events,
		Codespace:   res.TxResult.Codespace,
		Signer:      res.TxResult.Signer,
		Recipient:   res.TxResult.Recipient,
		MessageType: res.TxResult.MessageType,
	}
	rpcStdTx := RPCStdTx(tx)
	r := &RPCResultTx{
		Hash:     res.Hash,
		Height:   res.Height,
		Index:    res.Index,
		TxResult: rpcDeliverTx,
		Tx:       res.Tx,
		Proof:    res.Proof,
		StdTx:    rpcStdTx,
	}
	return r
}

func AccountTxs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = PaginateAddrParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	var res *core_types.ResultTxSearch
	var err error
	if !params.Received {
		res, err = app.PCA.QueryAccountTxs(params.Address, params.Page, params.PerPage, params.Prove, params.Sort)
	} else {
		res, err = app.PCA.QueryRecipientTxs(params.Address, params.Page, params.PerPage, params.Prove, params.Sort)
	}
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	rpcResponse := ResultTxSearchToRPC(res)
	s, er := json.MarshalIndent(rpcResponse, "", "  ")
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteJSONResponse(w, string(s), r.URL.Path, r.Host)
}

func BlockTxs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = PaginatedHeightParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	res, err := app.PCA.QueryBlockTxs(params.Height, params.Page, params.PerPage, params.Prove, params.Sort)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
	}
	rpcResponse := ResultTxSearchToRPC(res)
	s, er := json.MarshalIndent(rpcResponse, "", "  ")
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
	res, err := app.PCA.QueryHeight()
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
	var params = HeightAndAddrParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	balance, err := app.PCA.QueryBalance(params.Address, params.Height)
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
	var params = HeightAndAddrParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	res, err := app.PCA.QueryAccount(params.Address, params.Height)
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
	var params = HeightAndValidatorOptsParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	if params.Opts.Page == 0 {
		params.Opts.Page = 1
	}
	if params.Opts.Limit == 0 {
		params.Opts.Limit = 1000
	}
	res, err := app.PCA.QueryNodes(params.Height, params.Opts)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := res.JSON()
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(j)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
	}
}

func Node(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = HeightAndAddrParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	res, err := app.PCA.QueryNode(params.Address, params.Height)
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
	var params = HeightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.PCA.QueryNodeParams(params.Height)
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

type QueryNodeReceiptParam struct {
	Address      string `json:"address"`
	Blockchain   string `json:"blockchain"`
	AppPubKey    string `json:"app_pubkey"`
	SBlockHeight int64  `json:"session_block_height"`
	Height       int64  `json:"height"`
	ReceiptType  string `json:"receipt_type"`
}

func NodeClaim(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = QueryNodeReceiptParam{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	res, err := app.PCA.QueryClaim(params.Address, params.AppPubKey, params.Blockchain, params.ReceiptType, params.SBlockHeight, params.Height)
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

func NodeClaims(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = PaginatedHeightAndAddrParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	res, err := app.PCA.QueryClaims(params.Addr, params.Height, params.Page, params.PerPage)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := res.JSON()
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

func Apps(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = HeightAndApplicaitonOptsParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	if params.Opts.Page == 0 {
		params.Opts.Page = 1
	}
	if params.Opts.Limit == 0 {
		params.Opts.Limit = 1000
	}
	res, err := app.PCA.QueryApps(params.Height, params.Opts)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := res.JSON()
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(j)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
	}
}

func App(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = HeightAndAddrParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	res, err := app.PCA.QueryApp(params.Address, params.Height)
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
	var params = HeightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	res, err := app.PCA.QueryAppParams(params.Height)
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
	var params = HeightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	res, err := app.PCA.QueryPocketParams(params.Height)
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
	var params = HeightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	res, err := app.PCA.QueryPocketSupportedBlockchains(params.Height)
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
	var params = HeightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	nodesStake, total, err := app.PCA.QueryTotalNodeCoins(params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	appsStaked, err := app.PCA.QueryTotalAppCoins(params.Height)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	dao, err := app.PCA.QueryDaoBalance(params.Height)
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
	var params = HeightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	res, err := app.PCA.QueryDaoOwner(0)
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
	var params = HeightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	res, err := app.PCA.QueryUpgrade(0)
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
	var params = HeightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	res, err := app.PCA.QueryACL(params.Height)
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

func AllParams(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = HeightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	res, err := app.PCA.QueryAllParams(params.Height)
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
func Param(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = HeightAndKeyParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	res, err := app.PCA.QueryParam(params.Height, params.Key)
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

func State(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = HeightParams{Height: 0}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if params.Height == 0 {
		params.Height = app.PCA.BaseApp.LastBlockHeight()
	}
	res, err := app.PCA.ExportState(params.Height, "")
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteRaw(w, res, r.URL.Path, r.Host)
}
