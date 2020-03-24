package pos

import (
	"github.com/pokt-network/pocket-core/x/apps/keeper"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodeskeeper "github.com/pokt-network/pocket-core/x/nodes/keeper"
	nodestypes "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/store"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/pokt-network/posmint/x/gov"
	govTypes "github.com/pokt-network/posmint/x/gov/types"
	"github.com/pokt-network/posmint/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"math/rand"
	"testing"

	sdk "github.com/pokt-network/posmint/types"
)

// nolint: deadcode unused
var (
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		gov.AppModuleBasic{},
		supply.AppModuleBasic{},
		nodes.AppModuleBasic{},
	)
)

// nolint: deadcode unused
// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.New()

	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	nodestypes.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

// nolint: deadcode unused
func createTestInput(t *testing.T, isCheckTx bool) (sdk.Ctx, keeper.Keeper, types.SupplyKeeper, types.PosKeeper) {
	initPower := int64(100000000000)
	nAccs := int64(4)

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.ParamsKey
	tkeyParams := sdk.ParamsTKey
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	nodesKey := sdk.NewKVStoreKey(nodestypes.StoreKey)
	appsKey := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(nodesKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(appsKey, sdk.StoreTypeIAVL, db)
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
	cdc := makeTestCodec()

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		types.StakedPoolName:      {supply.Burner, supply.Staking, supply.Minter},
		nodestypes.StakedPoolName: {supply.Burner, supply.Staking},
		govTypes.DAOAccountName:   {supply.Burner, supply.Staking},
	}

	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}
	valTokens := sdk.TokensFromConsensusPower(initPower)

	accSubspace := sdk.NewSubspace(auth.DefaultParamspace)
	bankSubspace := sdk.NewSubspace(bank.DefaultParamspace)
	nodesSubspace := sdk.NewSubspace(nodestypes.DefaultParamspace)
	appSubspace := sdk.NewSubspace(types.DefaultParamspace)
	ak := auth.NewAccountKeeper(cdc, keyAcc, accSubspace, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, bankSubspace, bank.DefaultCodespace, modAccAddrs)
	sk := supply.NewKeeper(cdc, keySupply, ak, bk, maccPerms)
	nk := nodeskeeper.NewKeeper(cdc, nodesKey, ak, bk, sk, nodesSubspace, "pos")
	moduleManager := module.NewManager(
		auth.NewAppModule(ak),
		bank.NewAppModule(bk, ak),
		supply.NewAppModule(sk, ak),
		nodes.NewAppModule(nk, ak, sk),
	)

	genesisState := ModuleBasics.DefaultGenesis()
	moduleManager.InitGenesis(ctx, genesisState)

	initialCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, valTokens))
	_ = createTestAccs(ctx, int(nAccs), initialCoins, &ak)

	keeper := keeper.NewKeeper(cdc, keySupply, bk, nk, sk, appSubspace, "apps")

	p := types.DefaultParams()
	keeper.SetParams(ctx, p)
	return ctx, keeper, sk, nk
}

// nolint: unparam deadcode unused
func createTestAccs(ctx sdk.Ctx, numAccs int, initialCoins sdk.Coins, ak *auth.AccountKeeper) (accs []auth.Account) {
	for i := 0; i < numAccs; i++ {
		privKey := crypto.GenerateEd25519PrivKey()
		pubKey := privKey.PublicKey()
		addr := sdk.Address(pubKey.Address())
		acc := auth.NewBaseAccountWithAddress(addr)
		acc.Coins = initialCoins
		acc.PubKey = pubKey
		ak.SetAccount(ctx, &acc)
		accs = append(accs, &acc)
	}
	return
}

func getRandomPubKey() crypto.Ed25519PublicKey {
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])
	return pub
}

func getApplication() types.Application {
	pub := getRandomPubKey()
	return types.Application{
		Address:      sdk.Address(pub.Address()),
		StakedTokens: sdk.NewInt(100000000000),
		PublicKey:    pub,
		Jailed:       false,
		Status:       sdk.Staked,
		MaxRelays:    sdk.NewInt(100000000000),
		Chains:       []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"},
	}
}
