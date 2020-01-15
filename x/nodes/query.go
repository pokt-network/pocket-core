package nodes

import (
	"encoding/hex"
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/util"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func QueryAccountBalance(cdc *codec.Codec, tmNode rpcclient.Client, addr sdk.Address, height int64) (sdk.Int, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryAccountBalanceParams{Address: addr}
	bz, err := cdc.MarshalBinaryBare(params)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	balanceBz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/account_balance", types.StoreKey), bz)
	var balance sdk.Int
	err = cdc.UnmarshalJSON(balanceBz, balance)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return balance, nil
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
	resKVs, _, err := cliCtx.QuerySubspace(types.StakedValidatorsKey, types.StoreKey)
	if err != nil {
		return types.Validators{}, err
	}
	var validators types.Validators
	for _, kv := range resKVs {
		validators = append(validators, types.MustUnmarshalValidator(cdc, kv.Value))
	}
	return validators, nil
}

func QueryUnstakedValidators(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (types.Validators, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	resKVs, _, err := cliCtx.QuerySubspace(types.UnstakedValidatorsKey, types.StoreKey)
	if err != nil {
		return types.Validators{}, err
	}
	var validators types.Validators
	for _, kv := range resKVs {
		validators = append(validators, types.MustUnmarshalValidator(cdc, kv.Value))
	}
	return validators, nil
}

func QueryUnstakingValidators(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (types.Validators, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	resKVs, _, err := cliCtx.QuerySubspace(types.UnstakingValidatorsKey, types.StoreKey)
	if err != nil {
		return types.Validators{}, err
	}
	var validators types.Validators
	for _, kv := range resKVs {
		validators = append(validators, types.MustUnmarshalValidator(cdc, kv.Value))
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
	return types.ValidatorSigningInfo{}, nil
}

func QuerySupply(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (stakedCoins sdk.Int, unstakedCoins sdk.Int, err error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	stakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/stakedPool", types.StoreKey), nil)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	unstakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/unstakedPool", types.StoreKey), nil)
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
	tx, err := (tmNode).Tx(res, false)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func QueryBlock(tmNode rpcclient.Client, height *int64) ([]byte, error) {
	res, err := (tmNode).Block(height)
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
