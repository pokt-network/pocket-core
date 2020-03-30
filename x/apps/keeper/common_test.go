package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/exported"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodeskeeper "github.com/pokt-network/pocket-core/x/nodes/keeper"
	nodestypes "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/types/module"
	govTypes "github.com/pokt-network/posmint/x/gov/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"math/rand"
	"testing"

	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/store"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/gov"
)

// nolint: deadcode unused
var (
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		nodes.AppModuleBasic{},
	)
)

// nolint: deadcode unused
// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	nodestypes.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

// nolint: deadcode unused
func createTestInput(t *testing.T, isCheckTx bool) (sdk.Context, []auth.Account, Keeper) {
	initPower := int64(100000000000)
	nAccs := int64(4)

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.ParamsKey
	tkeyParams := sdk.ParamsTKey
	nodesKey := sdk.NewKVStoreKey(nodestypes.StoreKey)
	appsKey := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(nodesKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(appsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain"}, isCheckTx, log.NewNopLogger()).WithAppVersion("0.0.0")
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
		types.StakedPoolName:      {auth.Burner, auth.Staking, auth.Minter},
		nodestypes.StakedPoolName: {auth.Burner, auth.Staking},
		govTypes.DAOAccountName:   {auth.Burner, auth.Staking},
	}

	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[auth.NewModuleAddress(acc).String()] = true
	}
	valTokens := sdk.TokensFromConsensusPower(initPower)
	accSubspace := sdk.NewSubspace(auth.DefaultParamspace)
	nodesSubspace := sdk.NewSubspace(nodestypes.DefaultParamspace)
	appSubspace := sdk.NewSubspace(DefaultParamspace)
	ak := auth.NewKeeper(cdc, keyAcc, accSubspace, maccPerms)
	nk := nodeskeeper.NewKeeper(cdc, nodesKey, ak, nodesSubspace, "pos")
	moduleManager := module.NewManager(
		auth.NewAppModule(ak),
		nodes.NewAppModule(nk),
	)
	genesisState := ModuleBasics.DefaultGenesis()
	moduleManager.InitGenesis(ctx, genesisState)
	initialCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, valTokens))
	accs := createTestAccs(ctx, int(nAccs), initialCoins, &ak)
	keeper := NewKeeper(cdc, appsKey, nk, ak, appSubspace, "apps")
	p := types.DefaultParams()
	keeper.SetParams(ctx, p)
	return ctx, accs, keeper
}

// nolint: unparam deadcode unused
func createTestAccs(ctx sdk.Ctx, numAccs int, initialCoins sdk.Coins, ak *auth.Keeper) (accs []auth.Account) {
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

func addMintedCoinsToModule(t *testing.T, ctx sdk.Ctx, k *Keeper, module string) {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), sdk.NewInt(100000000000)))
	mintErr := k.AccountsKeeper.MintCoins(ctx, module, coins.Add(coins))
	if mintErr != nil {
		t.Fail()
	}
}

func sendFromModuleToAccount(t *testing.T, ctx sdk.Ctx, k *Keeper, module string, address sdk.Address, amount sdk.Int) {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	err := k.AccountsKeeper.SendCoinsFromModuleToAccount(ctx, module, sdk.Address(address), coins)
	if err != nil {
		t.Fail()
	}
}

func getRandomPubKey() crypto.Ed25519PublicKey {
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])
	return pub
}

func getRandomApplicationAddress() sdk.Address {
	return sdk.Address(getRandomPubKey().Address())
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

func getStakedApplication() types.Application {
	return getApplication()
}

func getUnstakedApplication() types.Application {
	v := getApplication()
	return v.UpdateStatus(sdk.Unstaked)
}

func getUnstakingApplication() types.Application {
	v := getApplication()
	return v.UpdateStatus(sdk.Unstaking)
}

func modifyFn(i *int) func(index int64, application exported.ApplicationI) (stop bool) {
	return func(index int64, application exported.ApplicationI) (stop bool) {
		app := application.(types.Application)
		app.StakedTokens = sdk.NewInt(100)
		if index == 1 {
			stop = true
		}
		*i++
		return
	}
}

// InvariantRegistry is an autogenerated mock type for the InvariantRegistry type
type InvariantRegistry struct {
	mock.Mock
}

// RegisterRoute provides a mock function with given fields: moduleName, route, invar
func (_m *InvariantRegistry) RegisterRoute(moduleName string, route string, invar sdk.Invariant) {
	_m.Called(moduleName, route, invar)
}
