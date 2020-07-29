package gov

import (
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth/util"
	"github.com/pokt-network/pocket-core/x/gov/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

func QueryACL(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (acl types.ACL, err error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryACLParams{}
	bz, err := cdc.MarshalBinaryBare(params)
	if err != nil {
		return nil, err
	}
	balanceBz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryACL), bz)
	err = cdc.UnmarshalJSON(balanceBz, &acl)
	if err != nil {
		return nil, err
	}
	return acl, nil
}

func QueryDAOOwner(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (daoOwner sdk.Address, err error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	daoPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryDAOOwner), nil)
	if err != nil {
		return nil, err
	}
	if err := cdc.UnmarshalJSON(daoPoolBytes, &daoOwner); err != nil {
		return nil, err
	}
	return daoOwner, err
}

func QueryDAO(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (daoCoins sdk.Int, err error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryDAOParams{}
	bz, err := cdc.MarshalBinaryBare(params)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	daoPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryDAO), bz)
	if err != nil {
		return sdk.Int{}, err
	}
	var daoPool sdk.Int
	if err := cdc.UnmarshalJSON(daoPoolBytes, &daoPool); err != nil {
		return sdk.Int{}, err
	}
	return daoPool, err
}

func QueryUpgrade(cdc *codec.Codec, tmNode rpcclient.Client, height int64) (upgrade types.Upgrade, err error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryUpgradeParams{}
	bz, err := cdc.MarshalBinaryBare(params)
	if err != nil {
		return upgrade, err
	}
	upgradeBz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryUpgrade), bz)
	if err != nil {
		return upgrade, err
	}
	var u types.Upgrade
	if err := cdc.UnmarshalJSON(upgradeBz, &u); err != nil {
		return upgrade, err
	}
	return u, err
}
