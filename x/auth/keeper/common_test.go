package keeper

// DONTCOVER

import (
	"github.com/pokt-network/pocket-core/codec"
	cdcTypes "github.com/pokt-network/pocket-core/codec/types"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/store/rootmulti"
	sdk "github.com/pokt-network/pocket-core/types"
	authTypes "github.com/pokt-network/pocket-core/x/auth/types"
	govKeeper "github.com/pokt-network/pocket-core/x/gov/keeper"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type testInput struct {
	cdc    *codec.Codec
	ctx    sdk.Context
	Keeper Keeper
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()

	cdc := codec.NewCodec(cdcTypes.NewInterfaceRegistry())
	authTypes.RegisterCodec(cdc)
	crypto.RegisterAmino(cdc.AminoCodec().Amino)

	authCapKey := sdk.NewKVStoreKey("auth")
	keyParams := sdk.ParamsKey

	ms := rootmulti.NewStore(db, false, 5000000)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	_ = ms.LoadLatestVersion()
	akSubspace := sdk.NewSubspace(authTypes.DefaultCodespace)
	ak := NewKeeper(
		cdc, authCapKey, akSubspace, nil,
	)
	govKeeper.NewKeeper(cdc, sdk.ParamsKey, govTypes.DefaultCodespace, ak, akSubspace)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	ak.SetParams(ctx, authTypes.DefaultParams())
	return testInput{Keeper: ak, cdc: cdc, ctx: ctx}
}
