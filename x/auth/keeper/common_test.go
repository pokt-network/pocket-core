package keeper

// DONTCOVER

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/store"
	sdk "github.com/pokt-network/pocket-core/types"
	authTypes "github.com/pokt-network/pocket-core/x/auth/types"
	govKeeper "github.com/pokt-network/pocket-core/x/gov/keeper"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
)

type testInput struct {
	cdc    *codec.Codec
	ctx    sdk.Context
	Keeper Keeper
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()

	cdc := codec.New()
	authTypes.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	authCapKey := sdk.NewKVStoreKey("auth")
	keyParams := sdk.ParamsKey
	tkeyParams := sdk.ParamsTKey

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()
	akSubspace := sdk.NewSubspace(authTypes.DefaultCodespace)
	ak := NewKeeper(
		cdc, authCapKey, akSubspace, nil,
	)
	govKeeper.NewKeeper(cdc, sdk.ParamsKey, sdk.ParamsTKey, govTypes.DefaultCodespace, ak, akSubspace)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	ak.SetParams(ctx, authTypes.DefaultParams())
	return testInput{Keeper: ak, cdc: cdc, ctx: ctx}
}
