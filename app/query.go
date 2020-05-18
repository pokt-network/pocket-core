package app

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"

	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/exported"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/pokt-network/posmint/x/gov/types"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
)

const (
	messageSenderQuery     = "message.sender='%s'"
	transferRecipientQuery = "transfer.recipient='%s'"
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
	return codec.Cdc.MarshalJSONIndent(b, "", "  ")
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

func (app PocketCoreApp) QueryAccountTxs(addr string, page, perPage int, prove bool) (res *core_types.ResultTxSearch, err error) {
	tmClient := app.GetClient()
	defer func() { _ = tmClient.Stop() }()
	_, err = hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(messageSenderQuery, addr)
	page, perPage = checkPagination(page, perPage)
	res, err = tmClient.TxSearch(query, prove, page, perPage)
	return
}
func (app PocketCoreApp) QueryRecipientTxs(addr string, page, perPage int, prove bool) (res *core_types.ResultTxSearch, err error) {
	tmClient := app.GetClient()
	defer func() { _ = tmClient.Stop() }()
	_, err = hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(transferRecipientQuery, addr)
	page, perPage = checkPagination(page, perPage)
	res, err = tmClient.TxSearch(query, prove, page, perPage)
	return
}

func (app PocketCoreApp) QueryBlockTxs(height int64, page, perPage int, prove bool) (res *core_types.ResultTxSearch, err error) {
	tmClient := app.GetClient()
	defer func() { _ = tmClient.Stop() }()
	query := fmt.Sprintf(txHeightQuery, height)
	page, perPage = checkPagination(page, perPage)
	res, err = tmClient.TxSearch(query, prove, page, perPage)
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

func (app PocketCoreApp) QueryBalance(addr string, height int64) (res sdk.Int, err error) {
	acc, err := app.QueryAccount(addr, height)
	if err != nil {
		return
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

func (app PocketCoreApp) QueryTotalNodeCoins(height int64) (stakedTokens sdk.Int, totalTokens sdk.Int, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	stakedTokens = app.nodesKeeper.GetStakedTokens(ctx)
	totalTokens = app.nodesKeeper.TotalTokens(ctx)
	return
}

func (app PocketCoreApp) QueryDaoBalance(height int64) (res sdk.Int, err error) {
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

func (app PocketCoreApp) QueryTotalAppCoins(height int64) (staked sdk.Int, err error) {
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

func (app PocketCoreApp) QueryReceipts(addr string, height int64, page, perPage int) (res Page, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return
	}
	page, perPage = checkPagination(page, perPage)
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	r, err := app.pocketKeeper.GetReceipts(ctx, a)
	if err != nil {
		return
	}
	return paginate(page, perPage, r, 1000)
}

func (app PocketCoreApp) QueryReceipt(blockchain, appPubKey, addr, receiptType string, sessionblockHeight, height int64) (res *pocketTypes.Receipt, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	h := pocketTypes.SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              blockchain,
		SessionBlockHeight: sessionblockHeight,
	}
	var et pocketTypes.EvidenceType
	switch strings.ToLower(receiptType) {
	case "relay":
		et = pocketTypes.RelayEvidence
	case "challenge":
		et = pocketTypes.ChallengeEvidence
	default:
		return nil, sdk.ErrInternal("type in the receipt query is not recognized: (relay or challenge)")
	}
	app.pocketKeeper.GetReceipt(ctx, a, h, et)
	return
}

func (app PocketCoreApp) QueryPocketSupportedBlockchains(height int64) (res []string, err error) {
	ctx, err := app.NewContext(height)
	if err != nil {
		return
	}
	sb := app.pocketKeeper.SupportedBlockchains(ctx)
	return sb, nil
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

func (app PocketCoreApp) HandleRelay(r pocketTypes.Relay) (res *pocketTypes.RelayResponse, err error) {
	ctx, err := app.NewContext(app.LastBlockHeight())
	if err != nil {
		return nil, err
	}
	return app.pocketKeeper.HandleRelay(ctx, r)
}

func checkPagination(page, limit int) (int, int) {
	if page < 0 {
		page = 1
	}
	if limit < 0 {
		limit = 30
	}
	return page, limit
}

func paginate(page, limit int, items interface{}, MaxValidators int) (res Page, error error) {
	slice, success := takeArg(items, reflect.Slice)
	if !success {
		return Page{}, fmt.Errorf("invalid argument, non slice input to paginate")
	}
	l := slice.Len()
	start, end := util.Paginate(l, page, limit, MaxValidators)
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
