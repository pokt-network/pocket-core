package keeper

import (
	"fmt"
	cdcTypes "github.com/pokt-network/pocket-core/codec/types"
	"github.com/pokt-network/pocket-core/crypto"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/store"
	"github.com/pokt-network/pocket-core/x/auth/types"

	sdk "github.com/pokt-network/pocket-core/types"
)

// nolint: deadcode unused
var (
	multiPerm  = "multiple permissions account"
	randomPerm = "random permission"
	holder     = "holder"
)

// nolint: deadcode unused
// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.NewCodec(cdcTypes.NewInterfaceRegistry())
	types.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	crypto.RegisterAmino(cdc.AminoCodec().Amino)

	return cdc
}

// nolint: deadcode unused
func createTestInput(t *testing.T, isCheckTx bool, initPower int64, nAccs int64) (sdk.Context, Keeper) {
	keyAcc := sdk.NewKVStoreKey(types.StoreKey)
	keyParams := sdk.ParamsKey
	tkeyParams := sdk.ParamsTKey
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db, false, 5000000)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "supply-chain"}, isCheckTx, log.NewNopLogger()).WithAppVersion("0.0.0")
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)
	cdc := makeTestCodec()
	maccPerms := map[string][]string{
		holder:       nil,
		types.Minter: {types.Minter},
		types.Burner: {types.Burner},
		multiPerm:    {types.Minter, types.Burner, types.Staking},
		randomPerm:   {"random"},
	}
	keeper := NewKeeper(cdc, keyAcc, sdk.NewSubspace(types.StoreKey), maccPerms)
	valTokens := sdk.TokensFromConsensusPower(initPower)
	initialCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, valTokens))
	createTestAccs(ctx, int(nAccs), initialCoins, &keeper)
	totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, valTokens.MulRaw(nAccs)))
	keeper.SetSupply(ctx, types.NewSupply(totalSupply))
	keeper.SetParams(ctx, types.DefaultParams())
	return ctx, keeper
}

// nolint: unparam deadcode unused
func createTestAccs(ctx sdk.Ctx, numAccs int, initialCoins sdk.Coins, ak *Keeper) (accs types.Accounts) {
	var err error
	for i := 0; i < numAccs; i++ {
		privKey := crypto.Secp256k1PrivateKey{}.GenPrivateKey()
		pubKey := privKey.PubKey()
		addr := sdk.Address(pubKey.Address())
		acc := types.NewBaseAccountWithAddress(addr)
		acc.Coins = initialCoins
		acc.PubKey, err = crypto.PubKeyToPublicKey(pubKey)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ak.SetAccount(ctx, &acc)
	}
	return
}
