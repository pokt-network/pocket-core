package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/exported"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodeskeeper "github.com/pokt-network/pocket-core/x/nodes/keeper"
	nodestypes "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/supply"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"math/rand"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/store"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/pokt-network/posmint/x/params"

	sdk "github.com/pokt-network/posmint/types"
)

// nolint: deadcode unused
var (
	multiPerm    = "multiple permissions account"
	randomPerm   = "random permission"
	holder       = "holder"
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
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
	params.RegisterCodec(cdc)
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
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
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
		nodestypes.DAOPoolName:    {supply.Burner, supply.Staking},
	}

	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}
	valTokens := sdk.TokensFromConsensusPower(initPower)

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, modAccAddrs)
	sk := supply.NewKeeper(cdc, keySupply, ak, bk, maccPerms)
	nk := nodeskeeper.NewKeeper(cdc, nodesKey, ak, bk, sk, pk.Subspace(nodestypes.DefaultParamspace), "pos")

	moduleManager := module.NewManager(
		auth.NewAppModule(ak),
		bank.NewAppModule(bk, ak),
		supply.NewAppModule(sk, ak),
		nodes.NewAppModule(nk, ak, sk),
	)

	genesisState := ModuleBasics.DefaultGenesis()
	moduleManager.InitGenesis(ctx, genesisState)

	appSubspace := pk.Subspace(DefaultParamspace)
	initialCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, valTokens))
	accs := createTestAccs(ctx, int(nAccs), initialCoins, &ak)

	keeper := NewKeeper(cdc, keySupply, bk, nk, sk, appSubspace, "apps")

	p := types.DefaultParams()
	keeper.SetParams(ctx, p)
	return ctx, accs, keeper
}

// nolint: unparam deadcode unused
func createTestAccs(ctx sdk.Context, numAccs int, initialCoins sdk.Coins, ak *auth.AccountKeeper) (accs []auth.Account) {
	for i := 0; i < numAccs; i++ {
		privKey := crypto.GenerateEd25519PrivKey()
		pubKey := privKey.PublicKey()
		addr := sdk.Address(pubKey.Address())
		acc := auth.NewBaseAccountWithAddress(addr)
		acc.Coins = initialCoins
		acc.PubKey = pubKey
		acc.AccountNumber = uint64(i)
		ak.SetAccount(ctx, &acc)
		accs = append(accs, &acc)
	}
	return
}

func addMintedCoinsToModule(t *testing.T, ctx sdk.Context, k *Keeper, module string) {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), sdk.NewInt(100000000000)))
	mintErr := k.supplyKeeper.MintCoins(ctx, module, coins.Add(coins))
	if mintErr != nil {
		t.Fail()
	}
}

func sendFromModuleToAccount(t *testing.T, ctx sdk.Context, k *Keeper, module string, address sdk.Address, amount sdk.Int) {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, module, sdk.Address(address), coins)
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

type AppHooks struct {
	mock.Mock
}

func (ah *AppHooks) BeforeApplicationRegistered(ctx sdk.Context, valAddr sdk.Address) {
	ah.Called(ctx, valAddr)
}
func (ah *AppHooks) AfterApplicationRegistered(ctx sdk.Context, valAddr sdk.Address) {
	ah.Called(ctx, valAddr)
}
func (ah *AppHooks) BeforeApplicationRemoved(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ah.Called(ctx, consAddr, valAddr)
}
func (ah *AppHooks) AfterApplicationRemoved(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ah.Called(ctx, consAddr, valAddr)
}
func (ah *AppHooks) BeforeApplicationStaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ah.Called(ctx, consAddr, valAddr)
}
func (ah *AppHooks) AfterApplicationStaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ah.Called(ctx, consAddr, valAddr)
}
func (ah *AppHooks) BeforeApplicationBeginUnstaking(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ah.Called(ctx, consAddr, valAddr)
}
func (ah *AppHooks) AfterApplicationBeginUnstaking(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ah.Called(ctx, consAddr, valAddr)
}
func (ah *AppHooks) BeforeApplicationUnstaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ah.Called(ctx, consAddr, valAddr)
}
func (ah *AppHooks) AfterApplicationUnstaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ah.Called(ctx, consAddr, valAddr)
}
func (ah *AppHooks) BeforeApplicationSlashed(ctx sdk.Context, valAddr sdk.Address, fraction sdk.Dec) {
	ah.Called(ctx, valAddr, fraction)
}
func (ah *AppHooks) AfterApplicationSlashed(ctx sdk.Context, valAddr sdk.Address, fraction sdk.Dec) {
	ah.Called(ctx, valAddr, fraction)
}

// InvariantRegistry is an autogenerated mock type for the InvariantRegistry type
type InvariantRegistry struct {
	mock.Mock
}

// RegisterRoute provides a mock function with given fields: moduleName, route, invar
func (_m *InvariantRegistry) RegisterRoute(moduleName string, route string, invar sdk.Invariant) {
	_m.Called(moduleName, route, invar)
}
