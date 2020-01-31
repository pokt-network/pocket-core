package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/supply"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/store"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/pokt-network/posmint/x/params"
	//"github.com/pokt-network/posmint/x/supply/internal/types"

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

// nolint: deadcode unused
func createTestInput(t *testing.T, isCheckTx bool) (sdk.Context, []auth.Account, Keeper) {
	initPower := int64(100000000000)
	nAccs := int64(4)

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
	cdc := makeTestCodec()

	maccPerms := map[string][]string{
		auth.FeeCollectorName: nil,
		types.StakedPoolName:  {supply.Burner, supply.Staking, supply.Minter},
		types.DAOPoolName:     {supply.Burner, supply.Staking, supply.Minter},
		//ADDED THIS CHECK
		types.ModuleName: {supply.Burner, supply.Staking, supply.Minter},
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

	moduleManager := module.NewManager(
		auth.NewAppModule(ak),
		bank.NewAppModule(bk, ak),
		supply.NewAppModule(sk, ak),
	)

	genesisState := ModuleBasics.DefaultGenesis()
	moduleManager.InitGenesis(ctx, genesisState)

	posSubSpace := pk.Subspace(DefaultParamspace)
	initialCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, valTokens))
	accs := createTestAccs(ctx, int(nAccs), initialCoins, &ak)

	keeper := NewKeeper(cdc, keySupply, ak, bk, sk, posSubSpace, sdk.CodespaceType("pos"))

	params := types.DefaultParams()
	keeper.SetParams(ctx, params)
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

func getRandomValidatorAddress() sdk.Address {
	return sdk.Address(getRandomPubKey().Address())
}

func getValidator() types.Validator {
	pub := getRandomPubKey()
	return types.Validator{
		Address:      sdk.Address(pub.Address()),
		StakedTokens: sdk.NewInt(100000000000),
		PublicKey:    pub,
		Jailed:       false,
		Status:       sdk.Staked,
		ServiceURL:   "google.com",
		Chains:       []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"},
	}
}

func getStakedValidator() types.Validator {
	return getValidator()
}

func getUnstakedValidator() types.Validator {
	v := getValidator()
	return v.UpdateStatus(sdk.Unstaked)
}

func getUnstakingValidator() types.Validator {
	v := getValidator()
	return v.UpdateStatus(sdk.Unstaking)
}
func modifyFn(i *int) func(index int64, Validator exported.ValidatorI) (stop bool) {
	return func(index int64, validator exported.ValidatorI) (stop bool) {
		val := validator.(types.Validator)
		val.StakedTokens = sdk.NewInt(100)
		if index == 1 {
			stop = true
		}
		*i++
		return
	}
}

type POSHooks struct {
	mock.Mock
}

func (ph *POSHooks) BeforeValidatorRegistered(ctx sdk.Context, valAddr sdk.Address) {
	ph.Called(ctx, valAddr)
}
func (ph *POSHooks) AfterValidatorRegistered(ctx sdk.Context, valAddr sdk.Address) {
	ph.Called(ctx, valAddr)
}
func (ph *POSHooks) BeforeValidatorRemoved(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ph.Called(ctx, consAddr, valAddr)
}
func (ph *POSHooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ph.Called(ctx, consAddr, valAddr)
}
func (ph *POSHooks) BeforeValidatorStaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ph.Called(ctx, consAddr, valAddr)
}
func (ph *POSHooks) AfterValidatorStaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ph.Called(ctx, consAddr, valAddr)
}
func (ph *POSHooks) BeforeValidatorBeginUnstaking(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ph.Called(ctx, consAddr, valAddr)
}
func (ph *POSHooks) AfterValidatorBeginUnstaking(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ph.Called(ctx, consAddr, valAddr)
}
func (ph *POSHooks) BeforeValidatorUnstaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ph.Called(ctx, consAddr, valAddr)
}
func (ph *POSHooks) AfterValidatorUnstaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	ph.Called(ctx, consAddr, valAddr)
}
func (ph *POSHooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.Address, fraction sdk.Dec) {
	ph.Called(ctx, valAddr, fraction)
}
func (ph *POSHooks) AfterValidatorSlashed(ctx sdk.Context, valAddr sdk.Address, fraction sdk.Dec) {
	ph.Called(ctx, valAddr, fraction)
}
