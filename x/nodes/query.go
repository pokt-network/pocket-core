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

func QueryAccountBalance(cdc *codec.Codec, tmNode rpcclient.Client, addr sdk.Address, height int64) (sdk.Int, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryAccountParams{Address: addr}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	path := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryAccountBalance)
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
	path := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryAccount)
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

func QueryValidators(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (types.Validators, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	resKVs, _, err := cliCtx.QuerySubspace(types.AllValidatorsKey, types.StoreKey)
	if err != nil {
		return types.Validators{}, err
	}
	var validators types.Validators
	for _, kv := range resKVs {
		validators = append(validators, types.MustUnmarshalValidator(cdc, kv.Value))
	}
	return validators, nil
}

func QueryStakedValidators(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (types.Validators, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryStakedValidatorsParams{
		Page:  1,
		Limit: 10000,
	}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryStakedValidators), bz)
	if err != nil {
		return types.Validators{}, err
	}
	validators := types.Validators{}
	err = cdc.UnmarshalJSON(res, &validators)
	if err != nil {
		return validators, err
	}
	return validators, nil
}

func QueryUnstakedValidators(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (types.Validators, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryUnstakedValidatorsParams{
		Page:  1,
		Limit: 10000,
	}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryUnstakedValidators), bz)
	if err != nil {
		return types.Validators{}, err
	}
	validators := types.Validators{}
	err = cdc.UnmarshalJSON(res, &validators)
	if err != nil {
		return validators, err
	}
	return validators, nil
}

func QueryUnstakingValidators(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (types.Validators, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryUnstakingValidatorsParams{
		Page:  1,
		Limit: 10000,
	}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryUnstakingValidators), bz)
	if err != nil {
		return types.Validators{}, err
	}
	validators := types.Validators{}
	err = cdc.UnmarshalJSON(res, &validators)
	if err != nil {
		return validators, err
	}
	return validators, nil
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

func QuerySupply(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (stakedCoins sdk.Int, unstakedCoins sdk.Int, err error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	stakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryStakedPool), nil)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	unstakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryUnstakedPool), nil)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	var stakedPool types.StakingPool
	if err := cdc.UnmarshalJSON(stakedPoolBytes, &stakedPool); err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	var unstakedPool types.StakingPool
	if err := cdc.UnmarshalJSON(unstakedPoolBytes, &unstakedPool); err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	return stakedPool.Tokens, unstakedPool.Tokens, nil
}

func QueryDAO(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (daoCoins sdk.Int, err error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	daoPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/dao", types.StoreKey), nil)
	if err != nil {
		return sdk.Int{}, err
	}
	var daoPool types.DAOPool
	if err := cdc.UnmarshalJSON(daoPoolBytes, &daoPool); err != nil {
		return sdk.Int{}, err
	}
	return daoPool.Tokens, err
}

func QueryPOSParams(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (types.Params, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	route := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryParameters)
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
