package keeper

import (
	"github.com/pokt-network/pocket-core/codec/types"
	"math/rand"
	"testing"

	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/store"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/types/module"
	"github.com/pokt-network/pocket-core/x/auth"
	"github.com/pokt-network/pocket-core/x/auth/keeper"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

var (
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
	)
)

// nolint: deadcode unused
// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.NewCodec(types.NewInterfaceRegistry())
	auth.RegisterCodec(cdc)
	govTypes.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	crypto.RegisterAmino(cdc.AminoCodec().Amino)
	return cdc
}

func getRandomPubKey() crypto.Ed25519PublicKey {
	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	return pub
}

func getRandomValidatorAddress() sdk.Address {
	return sdk.Address(getRandomPubKey().Address())
}

// nolint: deadcode unused
func createTestKeeperAndContext(t *testing.T, isCheckTx bool) (sdk.Context, Keeper) {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db, false, 5000000)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(sdk.ParamsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(sdk.ParamsTKey, sdk.StoreTypeTransient, db)
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
	cdc := makeTestCodec()
	maccPerms := map[string][]string{
		auth.FeeCollectorName:   nil,
		govTypes.DAOAccountName: {"burner", "staking", "minter"},
		"FAKE":                  {"burner", "staking", "minter"},
	}
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[auth.NewModuleAddress(acc).String()] = true
	}
	akSubspace := sdk.NewSubspace(auth.DefaultParamspace)
	ak := keeper.NewKeeper(cdc, keyAcc, akSubspace, maccPerms)
	ak.GetModuleAccount(ctx, "FAKE")
	pk := NewKeeper(cdc, sdk.ParamsKey, sdk.ParamsTKey, govTypes.DefaultParamspace, ak, akSubspace)
	moduleManager := module.NewManager(
		auth.NewAppModule(ak),
	)
	genesisState := ModuleBasics.DefaultGenesis()
	moduleManager.InitGenesis(ctx, genesisState)
	params := govTypes.DefaultParams()
	pk.SetParams(ctx, params)
	gs := govTypes.DefaultGenesisState()
	acl := createTestACL()
	gs.Params.ACL = acl
	pk.InitGenesis(ctx, gs)
	return ctx, pk
}

var testACL govTypes.ACL

func createTestACL() govTypes.ACL {
	if testACL == nil {
		acl := govTypes.ACL(make([]govTypes.ACLPair, 0))
		acl.SetOwner("auth/MaxMemoCharacters", getRandomValidatorAddress())
		acl.SetOwner("auth/TxSigLimit", getRandomValidatorAddress())
		acl.SetOwner("auth/FeeMultipliers", getRandomValidatorAddress())
		acl.SetOwner("gov/daoOwner", getRandomValidatorAddress())
		acl.SetOwner("gov/acl", getRandomValidatorAddress())
		acl.SetOwner("gov/upgrade", getRandomValidatorAddress())
		testACL = acl
	}
	return testACL
}

// Checks wether or not a Events slice contains an event that equals the values of event
func ContainsEvent(events sdk.Events, event abci.Event) bool {
	stringEvents := sdk.StringifyEvents(events.ToABCIEvents())
	stringEventStr := sdk.StringEvents{sdk.StringifyEvent(event)}.String()
	for _, item := range stringEvents {
		itemStr := sdk.StringEvents{item}.String()
		if itemStr == stringEventStr {
			return true
		}
	}
	return false
}
