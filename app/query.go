package app

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"

	sdk "github.com/pokt-network/pocket-core/types"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/auth/exported"
	"github.com/pokt-network/pocket-core/x/auth/util"
	"github.com/pokt-network/pocket-core/x/gov/types"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
)

const (
	messageSenderQuery     = "tx.signer='%s'"
	transferRecipientQuery = "tx.recipient='%s'"
	txHeightQuery          = "tx.height=%d"
)

// zero for height = latest
func (app PocketCoreApp) QueryBlock(height *int64) (blockJSON []byte, err error) {
	tmClient := app.GetClient()
	defer func() { _ = tmClient.Stop() }()
	b, err := tmClient.Block(height)
	if err != nil {
		return nil, err
	}
	return Codec().MarshalJSONIndent(b, "", "  ")
}

func (app PocketCoreApp) QueryTx(hash string, prove bool) (res *core_types.ResultTx, err error) {
	tmClient := app.GetClient()
	defer func() { _ = tmClient.Stop() }()
	h, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	res, err = tmClient.Tx(h, prove)
	return
}

func (app PocketCoreApp) QueryAccountTxs(addr string, page, perPage int, prove bool, sort string) (res *core_types.ResultTxSearch, err error) {
	tmClient := app.GetClient()
	defer func() { _ = tmClient.Stop() }()
	_, err = hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(messageSenderQuery, addr)
	page, perPage = checkPagination(page, perPage)
	res, err = tmClient.TxSearch(query, prove, page, perPage, checkSort(sort))
	return
}
func (app PocketCoreApp) QueryRecipientTxs(addr string, page, perPage int, prove bool, sort string) (res *core_types.ResultTxSearch, err error) {
	tmClient := app.GetClient()
	defer func() { _ = tmClient.Stop() }()
	_, err = hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(transferRecipientQuery, addr)
	page, perPage = checkPagination(page, perPage)
	res, err = tmClient.TxSearch(query, prove, page, perPage, checkSort(sort))
	return
}

func (app PocketCoreApp) QueryBlockTxs(height int64, page, perPage int, prove bool, sort string) (res *core_types.ResultTxSearch, err error) {
	tmClient := app.GetClient()
	defer func() { _ = tmClient.Stop() }()
	query := fmt.Sprintf(txHeightQuery, height)
	page, perPage = checkPagination(page, perPage)
	res, err = tmClient.TxSearch(query, prove, page, perPage, checkSort(sort))
	return
}

func (app PocketCoreApp) QueryHeight() (res int64, err error) {
	tmClient := app.GetClient()
	defer func() { _ = tmClient.Stop() }()
	status, err := tmClient.Status()
	if err != nil {
		return -1, err
	}

	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}

func (app PocketCoreApp) QueryNodeStatus() (res *core_types.ResultStatus, err error) {
	tmClient := app.GetClient()
	defer func() { _ = tmClient.Stop() }()
	return tmClient.Status()
}

func (app PocketCoreApp) QueryBalance(addr string, height int64) (res sdk.BigInt, err error) {
	acc, err := app.QueryAccount(addr, height)
	if err != nil {
		return
	}
	if (*acc) == nil {
		return sdk.NewInt(0), nil
	}
	return (*acc).GetCoins().AmountOf(sdk.DefaultStakeDenom), nil
}

func (app PocketCoreApp) QueryAccount(addr string, height int64) (res *exported.Account, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	acc := app.accountKeeper.GetAccount(ctx, a)
	return &acc, nil
}

func (app PocketCoreApp) QueryNodes(height int64, opts nodesTypes.QueryValidatorsParams) (res Page, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	opts.Page, opts.Limit = checkPagination(opts.Page, opts.Limit)
	nodes := app.nodesKeeper.GetAllValidatorsWithOpts(ctx, opts)
	return paginate(opts.Page, opts.Limit, nodes, int(app.nodesKeeper.GetParams(ctx).MaxValidators))
}

func (app PocketCoreApp) QueryNode(addr string, height int64) (res nodesTypes.Validator, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return res, err
	}
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	res, found := app.nodesKeeper.GetValidator(ctx, a)
	if !found {
		err = fmt.Errorf("validator not found for %s", a.String())
	}
	return
}

func (app PocketCoreApp) QueryNodeParams(height int64) (res nodesTypes.Params, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	return app.nodesKeeper.GetParams(ctx), nil
}

func (app PocketCoreApp) QuerySigningInfo(height int64, addr string) (res nodesTypes.ValidatorSigningInfo, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return nodesTypes.ValidatorSigningInfo{}, err
	}
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	res, found := app.nodesKeeper.GetValidatorSigningInfo(ctx, a)
	if !found {
		err = fmt.Errorf("signing info not found for %s", a.String())
	}
	return
}

func (app PocketCoreApp) QueryTotalNodeCoins(height int64) (stakedTokens sdk.BigInt, totalTokens sdk.BigInt, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	stakedTokens = app.nodesKeeper.GetStakedTokens(ctx)
	totalTokens = app.nodesKeeper.TotalTokens(ctx)
	return
}

func (app PocketCoreApp) QueryDaoBalance(height int64) (res sdk.BigInt, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	return app.govKeeper.GetDAOTokens(ctx), nil
}

func (app PocketCoreApp) QueryDaoOwner(height int64) (res sdk.Address, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	return app.govKeeper.GetDAOOwner(ctx), nil
}

func (app PocketCoreApp) QueryUpgrade(height int64) (res types.Upgrade, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	return app.govKeeper.GetUpgrade(ctx), nil
}

func (app PocketCoreApp) QueryACL(height int64) (res types.ACL, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	return app.govKeeper.GetACL(ctx), nil
}

type AllParamsReturn struct {
	AppParams    []SingleParamReturn `json:"app_params"`
	NodeParams   []SingleParamReturn `json:"node_params"`
	PocketParams []SingleParamReturn `json:"pocket_params"`
	GovParams    []SingleParamReturn `json:"gov_params"`
	AuthParams   []SingleParamReturn `json:"auth_params"`
}

type SingleParamReturn struct {
	Key   string `json:"param_key"`
	Value string `json:"param_value"`
}

func (app PocketCoreApp) QueryAllParams(height int64) (r AllParamsReturn, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	//get all the parameters from gov module
	allmap := app.govKeeper.GetAllParamNameValue(ctx)

	//transform for easy handling
	for k, v := range allmap {
		sub, _ := types.SplitACLKey(k)
		s, err2 := strconv.Unquote(v)
		if err2 != nil {
			//ignoring this error as content is a json object
			s = v
		}
		switch sub {
		case "pos":
			r.NodeParams = append(r.NodeParams, SingleParamReturn{
				Key:   k,
				Value: s,
			})
		case "application":
			r.AppParams = append(r.AppParams, SingleParamReturn{
				Key:   k,
				Value: s,
			})
		case "pocketcore":
			r.PocketParams = append(r.PocketParams, SingleParamReturn{
				Key:   k,
				Value: s,
			})
		case "gov":
			r.GovParams = append(r.GovParams, SingleParamReturn{
				Key:   k,
				Value: s,
			})
		case "auth":
			r.AuthParams = append(r.AuthParams, SingleParamReturn{
				Key:   k,
				Value: s,
			})
		default:
		}
	}

	return r, nil
}

func (app PocketCoreApp) QueryParam(height int64, paramkey string) (r SingleParamReturn, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	//get all the parameters from gov module
	allmap := app.govKeeper.GetAllParamNameValue(ctx)

	if val, ok := allmap[paramkey]; ok {
		r.Key = paramkey
		s, err2 := strconv.Unquote(val)
		if err2 != nil {
			//ignoring this error as content is a json object
			r.Value = val
			return r, err
		}
		r.Value = s
	}
	return
}

func (app PocketCoreApp) QueryApps(height int64, opts appsTypes.QueryApplicationsWithOpts) (res Page, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	opts.Page, opts.Limit = checkPagination(opts.Page, opts.Limit)
	applications := app.appsKeeper.GetAllApplicationsWithOpts(ctx, opts)
	return paginate(opts.Page, opts.Limit, applications, int(app.appsKeeper.GetParams(ctx).MaxApplications))
}

func (app PocketCoreApp) QueryApp(addr string, height int64) (res appsTypes.Application, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return res, err
	}
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	res, found := app.appsKeeper.GetApplication(ctx, a)
	if !found {
		err = appsTypes.ErrNoApplicationFound(appsTypes.ModuleName)
		return
	}
	return
}

func (app PocketCoreApp) QueryTotalAppCoins(height int64) (staked sdk.BigInt, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	return app.appsKeeper.GetStakedTokens(ctx), nil
}

func (app PocketCoreApp) QueryAppParams(height int64) (res appsTypes.Params, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	return app.appsKeeper.GetParams(ctx), nil
}

func (app PocketCoreApp) QueryPocketSupportedBlockchains(height int64) (res []string, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	sb := app.pocketKeeper.SupportedBlockchains(ctx)
	return sb, nil
}

func (app PocketCoreApp) QueryClaim(address, appPubkey, chain, evidenceType string, sessionBlockHeight int64, height int64) (res *pocketTypes.MsgClaim, err error) {
	a, err := sdk.AddressFromHex(address)
	if err != nil {
		return nil, err
	}
	header := pocketTypes.SessionHeader{
		ApplicationPubKey:  appPubkey,
		Chain:              chain,
		SessionBlockHeight: sessionBlockHeight,
	}
	err = header.ValidateHeader()
	if err != nil {
		return nil, err
	}
	et, err := pocketTypes.EvidenceTypeFromString(evidenceType)
	if err != nil {
		return nil, err
	}
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	claim, found := app.pocketKeeper.GetClaim(ctx, a, header, et)
	if !found {
		return nil, pocketTypes.NewClaimNotFoundError(pocketTypes.ModuleName)
	}
	return &claim, nil
}

func (app PocketCoreApp) QueryClaims(address string, height int64, page, perPage int) (res Page, err error) {
	var a sdk.Address
	var claims []pocketTypes.MsgClaim
	ctx, err := app.NewContext(height)
	page, perPage = checkPagination(page, perPage)
	if err != nil {
		return
	}
	if address != "" {
		a, err = sdk.AddressFromHex(address)
		if err != nil {
			return Page{}, err
		}
		claims, err = app.pocketKeeper.GetClaims(ctx, a)
		if err != nil {
			return Page{}, err
		}
	} else {
		claims = app.pocketKeeper.GetAllClaims(ctx)
	}
	p, err := paginate(page, perPage, claims, 10000)
	if err != nil {
		return Page{}, err
	}
	return p, nil
}

func (app PocketCoreApp) QueryPocketParams(height int64) (res pocketTypes.Params, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	p := app.pocketKeeper.GetParams(ctx)
	return p, nil
}

func (app PocketCoreApp) HandleChallenge(c pocketTypes.ChallengeProofInvalidData) (res *pocketTypes.ChallengeResponse, err error) {
	ctx, err := app.NewContext(app.LastBlockHeight())
	if err != nil {
		return nil, err
	}
	return app.pocketKeeper.HandleChallenge(ctx, c)
}

func (app PocketCoreApp) HandleDispatch(header pocketTypes.SessionHeader) (res *pocketTypes.DispatchResponse, err error) {
	ctx, err := app.NewContext(app.LastBlockHeight())
	if err != nil {
		return nil, err
	}
	return app.pocketKeeper.HandleDispatch(ctx, header)
}

func (app PocketCoreApp) HandleRelay(r pocketTypes.Relay) (res *pocketTypes.RelayResponse, dispatch *pocketTypes.DispatchResponse, err error) {
	ctx, err := app.NewContext(app.LastBlockHeight())
	if err != nil {
		return nil, nil, err
	}

	status, err := app.pocketKeeper.TmNode.Status()
	if err != nil {
		return nil, nil, fmt.Errorf("pocket node is unable to retrieve status from tendermint node, cannot service in this state")
	}
	if status.SyncInfo.CatchingUp {
		return nil, nil, fmt.Errorf("pocket node is currently syncing to the blockchain, cannot service in this state")
	}
	res, err = app.pocketKeeper.HandleRelay(ctx, r)
	var err1 error
	if err != nil && pocketTypes.ErrorWarrantsDispatch(err) {
		dispatch, err1 = app.HandleDispatch(r.Proof.SessionHeader())
		if err1 != nil {
			return
		}
	}
	return
}

func checkPagination(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 30
	}
	return page, limit
}
func checkSort(s string) string {
	switch s {
	case "asc":
		return s
	case "desc":
		return s
	default:
		return "desc"
	}
}

func paginate(page, limit int, items interface{}, max int) (res Page, error error) {
	slice, success := takeArg(items, reflect.Slice)
	if !success {
		return Page{}, fmt.Errorf("invalid argument, non slice input to paginate")
	}
	l := slice.Len()
	start, end := util.Paginate(l, page, limit, max)
	if start == -1 && end == -1 {
		return Page{}, nil
	}
	if start < 0 || end < 0 {
		return Page{}, fmt.Errorf("invalid bounds error: start %d finish %d", start, end)
	} else {
		items = slice.Slice(start, end).Interface()
	}
	totalPages := int(math.Ceil(float64(l) / float64(end-start)))
	if totalPages < 1 {
		totalPages = 1
	}
	return Page{Result: items, Total: totalPages, Page: page}, nil
}

func takeArg(arg interface{}, kind reflect.Kind) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)
	if val.Kind() == kind {
		ok = true
	}
	return
}

type Page struct {
	Result interface{} `json:"result"`
	Total  int         `json:"total_pages"`
	Page   int         `json:"page"`
}

// Marshals struct into JSON
func (p Page) JSON() (out []byte, err error) {
	// each element should be a JSON
	return json.Marshal(p)
}

// String returns a human readable string representation of a validator page
func (p Page) String() string {
	return fmt.Sprintf("Total:\t\t%d\nPage:\t\t%d\nResult:\t\t\n====\n%v\n====\n", p.Total, p.Page, p.Result)
}
