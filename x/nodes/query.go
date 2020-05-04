package nodes

import (
	"encoding/hex"
	"fmt"

	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

const (
	defaultPage            = 1
	defaultPerPage         = 30
	customQuery            = "custom/%s/%s"
	messageSenderQuery     = "message.sender='%s'"
	transferRecipientQuery = "transfer.recipient='%s'"
	txHeightQuery          = "tx.height=%d"
)

func checkPagination(page, limit int) (int, int) {
	if page < 0 {
		page = 1
	}
	if limit < 0 {
		limit = 10000
	}
	return page, limit
}

func QueryAccountBalance(cdc *codec.Codec, tmNode rpcclient.Client, addr sdk.Address, height int64) (sdk.Int, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryAccountParams{Address: addr}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	path := fmt.Sprintf(customQuery, types.StoreKey, types.QueryAccountBalance)
	balanceBz, _, err := cliCtx.QueryWithData(path, bz)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	var balance sdk.Int
	err = cdc.UnmarshalJSON(balanceBz, &balance)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return balance, nil
}

func QueryAccount(cdc *codec.Codec, tmNode rpcclient.Client, addr sdk.Address, height int64) (*auth.BaseAccount, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryAccountParams{Address: addr}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf(customQuery, types.StoreKey, types.QueryAccount)
	balanceBz, _, err := cliCtx.QueryWithData(path, bz)
	if err != nil {
		return nil, err
	}
	var acc auth.BaseAccount
	err = cdc.UnmarshalJSON(balanceBz, &acc)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func QueryValidator(cdc *codec.Codec, tmNode rpcclient.Client, addr sdk.Address, height int64) (types.Validator, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	res, _, err := cliCtx.QueryStore(types.KeyForValByAllVals(addr), types.StoreKey)
	if err != nil {
		return types.Validator{}, err
	}
	if len(res) == 0 {
		return types.Validator{}, fmt.Errorf("no validator found with address %s", addr)
	}
	return types.MustUnmarshalValidator(cdc, res), nil
}

func QueryValidators(cdc *codec.Codec, tmNode rpcclient.Client, height int64, opts types.QueryValidatorsParams) (types.ValidatorsPage, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	opts.Page, opts.Limit = checkPagination(opts.Page, opts.Limit)
	bz, err := cdc.MarshalJSON(opts)
	if err != nil {
		return types.ValidatorsPage{}, err
	}
	validatorsPage := types.ValidatorsPage{}
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf(customQuery, types.StoreKey, types.QueryValidators), bz)
	if err != nil {
		return validatorsPage, err
	}

	err = cdc.UnmarshalJSON(res, &validatorsPage)
	if err != nil {
		return validatorsPage, err
	}

	return validatorsPage, nil
}

func QuerySigningInfo(cdc *codec.Codec, tmNode rpcclient.Client, height int64, consAddr sdk.Address) (types.ValidatorSigningInfo, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	key := types.GetValidatorSigningInfoKey(consAddr)
	res, _, err := cliCtx.QueryStore(key, types.StoreKey)
	if err != nil {
		return types.ValidatorSigningInfo{}, err
	}
	if len(res) == 0 {
		return types.ValidatorSigningInfo{}, fmt.Errorf("validator %s not found in slashing store", consAddr)
	}
	var signingInfo types.ValidatorSigningInfo
	cdc.MustUnmarshalBinaryLengthPrefixed(res, &signingInfo)
	return signingInfo, nil
}

func QuerySupply(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (stakedCoins sdk.Int, totalTokens sdk.Int, err error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	stakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf(customQuery, types.StoreKey, types.QueryStakedPool), nil)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	if err := cdc.UnmarshalJSON(stakedPoolBytes, &stakedCoins); err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	totalTokensBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf(customQuery, types.StoreKey, types.QueryTotalSupply), nil)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	if err := cdc.UnmarshalJSON(totalTokensBytes, &totalTokens); err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	return
}

func QueryPOSParams(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (types.Params, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	route := fmt.Sprintf(customQuery, types.StoreKey, types.QueryParameters)
	bz, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return types.Params{}, err
	}
	var params types.Params
	cdc.MustUnmarshalJSON(bz, &params)
	return params, nil
}

func QueryTransaction(tmNode rpcclient.Client, hash string) (*ctypes.ResultTx, error) {
	res, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	tx, err := tmNode.Tx(res, false)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func validatePageAndPerPage(page, perPage int) (resPage, resPerPage int) {
	resPage = defaultPage
	resPerPage = defaultPerPage
	if page > 0 {
		resPage = page
	}
	if perPage > 0 {
		resPerPage = perPage
	}
	return resPage, resPerPage
}

func QueryAccountTransactions(tmNode rpcclient.Client, addr string, page, perPage int, recipient bool, prove bool) (*ctypes.ResultTxSearch, error) {
	_, err := hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	var query string
	if recipient == true {
		query = fmt.Sprintf(transferRecipientQuery, addr)
	} else {
		query = fmt.Sprintf(messageSenderQuery, addr)
	}
	page, perPage = validatePageAndPerPage(page, perPage)
	result, err := tmNode.TxSearch(query, prove, page, perPage)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func QueryBlockTransactions(tmNode rpcclient.Client, height int64, page, perPage int, prove bool) (*ctypes.ResultTxSearch, error) {
	query := fmt.Sprintf(txHeightQuery, height)
	page, perPage = validatePageAndPerPage(page, perPage)
	result, err := tmNode.TxSearch(query, prove, page, perPage)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func QueryBlock(tmNode rpcclient.Client, height *int64) ([]byte, error) {
	res, err := tmNode.Block(height)
	if err != nil {
		return nil, err
	}

	return codec.Cdc.MarshalJSONIndent(res, "", "  ")
}

// get the current blockchain height
func QueryChainHeight(tmNode rpcclient.Client) (int64, error) {
	client := (tmNode)
	status, err := client.Status()
	if err != nil {
		return -1, err
	}

	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}

func QueryNodeStatus(tmNode rpcclient.Client) (*ctypes.ResultStatus, error) {
	res, err := (tmNode).Status()
	if err != nil {
		return nil, nil
	}
	return res, nil
}
