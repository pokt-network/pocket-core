package types

import (
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/store"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/pokt-network/posmint/x/params"
	"github.com/pokt-network/posmint/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"testing"
	"time"
)

// nolint: deadcode unused
var (
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
		supply.AppModuleBasic{},
	)
)

// nolint: deadcode unused
// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.New()

	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	params.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

var (
	testSupportedChain string
)

func getTestSupportedBlockchain() string {
	if testSupportedChain == "" {
		testSupportedChain, _ = NonNativeChain{
			Ticker:  "eth",
			Netid:   "4",
			Version: "v1.9.9",
			Client:  "geth",
			Inter:   "",
		}.HashString()
	}
	return testSupportedChain
}

// nolint: deadcode unused
func newContext(t *testing.T, isCheckTx bool) sdk.Context {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain"}, isCheckTx, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)
	ctx = ctx.WithBlockHeader(abci.Header{
		Height: 1,
		Time:   time.Time{},
		LastBlockId: abci.BlockID{
			Hash: hash([]byte("fake")),
		},
	})
	return ctx
}

func GetRandomPrivateKey() crypto.Ed25519PrivateKey {
	return crypto.Ed25519PrivateKey{}.GenPrivateKey().(crypto.Ed25519PrivateKey)
}

func getRandomPubKey() crypto.Ed25519PublicKey {
	pk := GetRandomPrivateKey()
	return pk.PublicKey().(crypto.Ed25519PublicKey)
}

func getRandomValidatorAddress() sdk.Address {
	return sdk.Address(getRandomPubKey().Address())
}
