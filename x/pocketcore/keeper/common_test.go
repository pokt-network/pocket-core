package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	types2 "github.com/pokt-network/pocket-core/codec/types"
	"math"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/crypto/keys"
	"github.com/pokt-network/pocket-core/store"
	storeTypes "github.com/pokt-network/pocket-core/store/types"
	pocketTypes "github.com/pokt-network/pocket-core/types"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/types/module"
	apps "github.com/pokt-network/pocket-core/x/apps"
	appsKeeper "github.com/pokt-network/pocket-core/x/apps/keeper"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/auth"
	"github.com/pokt-network/pocket-core/x/gov"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodesKeeper "github.com/pokt-network/pocket-core/x/nodes/keeper"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/privval"
	tmStore "github.com/tendermint/tendermint/store"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

var (
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		gov.AppModuleBasic{},
	)
)

func TestMain(m *testing.M) {
	m.Run()
	err := os.RemoveAll("data")
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}

type simulateRelayKeys struct {
	private crypto.PrivateKey
	client  crypto.PrivateKey
}

func NewTestKeybase() keys.Keybase {
	return keys.NewInMemory()
}

// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.NewCodec(types2.NewInterfaceRegistry())
	auth.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	crypto.RegisterAmino(cdc.AminoCodec().Amino)
	return cdc
}

// : deadcode unused
func createTestInput(t *testing.T, isCheckTx bool) (sdk.Ctx, []nodesTypes.Validator, []appsTypes.Application, []auth.BaseAccount, Keeper, map[string]*sdk.KVStoreKey, keys.Keybase) {
	sdk.VbCCache = sdk.NewCache(1)
	initPower := int64(100000000000)
	nAccs := int64(5)
	kb := NewTestKeybase()
	_, err := kb.Create("test")
	assert.Nil(t, err)
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	addr := tmtypes.Address(cb.GetAddress())
	pk, err := kb.ExportPrivateKeyObject(cb.GetAddress(), "test")
	assert.Nil(t, err)
	types.InitPVKeyFile(privval.FilePVKey{
		Address: addr,
		PubKey:  cb.PublicKey,
		PrivKey: pk,
	})
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.ParamsKey
	tkeyParams := sdk.ParamsTKey
	nodesKey := sdk.NewKVStoreKey(nodesTypes.StoreKey)
	appsKey := sdk.NewKVStoreKey(appsTypes.StoreKey)
	pocketKey := sdk.NewKVStoreKey(types.StoreKey)

	keys := make(map[string]*sdk.KVStoreKey)
	keys["params"] = keyParams
	keys["pos"] = nodesKey
	keys["application"] = appsKey

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db, false, 5000000)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(nodesKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(appsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(pocketKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err = ms.LoadLatestVersion()
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
		Height: 976,
		Time:   time.Time{},
		LastBlockId: abci.BlockID{
			Hash: types.Hash([]byte("fake")),
		},
	})
	cdc := makeTestCodec()

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		appsTypes.StakedPoolName:  {auth.Burner, auth.Staking, auth.Minter},
		nodesTypes.StakedPoolName: {auth.Burner, auth.Staking},
		govTypes.DAOAccountName:   {auth.Burner, auth.Staking},
	}

	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[auth.NewModuleAddress(acc).String()] = true
	}
	valTokens := sdk.TokensFromConsensusPower(initPower)

	ethereum := hex.EncodeToString([]byte{01})

	hb := types.HostedBlockchains{
		M: map[string]types.HostedBlockchain{ethereum: {
			ID:  ethereum,
			URL: "https://www.google.com:443",
		}},
	}
	types.InitConfig(&hb, log.NewTMLogger(os.Stdout), sdk.DefaultTestingPocketConfig())
	authSubspace := sdk.NewSubspace(auth.DefaultParamspace)
	nodesSubspace := sdk.NewSubspace(nodesTypes.DefaultParamspace)
	appSubspace := sdk.NewSubspace(appsTypes.DefaultParamspace)
	pocketSubspace := sdk.NewSubspace(types.DefaultParamspace)
	ak := auth.NewKeeper(cdc, keyAcc, authSubspace, maccPerms)
	nk := nodesKeeper.NewKeeper(cdc, nodesKey, ak, nodesSubspace, nodesTypes.ModuleName)
	appk := appsKeeper.NewKeeper(cdc, appsKey, nk, ak, nil, appSubspace, appsTypes.ModuleName)
	appk.SetApplication(ctx, getTestApplication())
	keeper := NewKeeper(pocketKey, cdc, ak, nk, appk, &hb, pocketSubspace)
	appk.PocketKeeper = keeper
	assert.Nil(t, err)
	moduleManager := module.NewManager(
		auth.NewAppModule(ak),
		nodes.NewAppModule(nk),
		apps.NewAppModule(appk),
	)
	genesisState := ModuleBasics.DefaultGenesis()
	moduleManager.InitGenesis(ctx, genesisState)
	initialCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, valTokens))
	accs := createTestAccs(ctx, int(nAccs), initialCoins, &ak)
	ap := createTestApps(ctx, int(nAccs), sdk.NewIntFromBigInt(new(big.Int).SetUint64(math.MaxUint64)), appk, ak)
	vals := createTestValidators(ctx, int(nAccs), sdk.ZeroInt(), &nk, ak, kb)
	appk.SetParams(ctx, appsTypes.DefaultParams())
	nk.SetParams(ctx, nodesTypes.DefaultParams())
	defaultPocketParams := types.DefaultParams()
	defaultPocketParams.SupportedBlockchains = []string{getTestSupportedBlockchain()}
	keeper.SetParams(ctx, defaultPocketParams)
	return ctx, vals, ap, accs, keeper, keys, kb
}

var (
	testApp            appsTypes.Application
	testAppPrivateKey  crypto.PrivateKey
	testSupportedChain string
)

func getTestSupportedBlockchain() string {
	if testSupportedChain == "" {
		testSupportedChain = hex.EncodeToString([]byte{01})
	}
	return testSupportedChain
}

func getTestApplicationPrivateKey() crypto.PrivateKey {
	if testAppPrivateKey == nil {
		testAppPrivateKey = getRandomPrivateKey()
	}
	return testAppPrivateKey
}

func getTestApplication() appsTypes.Application {
	if testApp.Address == nil {
		pk := getTestApplicationPrivateKey().PublicKey()
		testApp = appsTypes.Application{
			Address:                 sdk.Address(pk.Address()),
			PublicKey:               pk,
			Jailed:                  false,
			Status:                  2,
			Chains:                  []string{getTestSupportedBlockchain()},
			StakedTokens:            sdk.NewInt(10000000),
			MaxRelays:               sdk.NewInt(10000000),
			UnstakingCompletionTime: time.Time{},
		}
	}
	return testApp
}

// : unparam deadcode unused
func createTestAccs(ctx sdk.Ctx, numAccs int, initialCoins sdk.Coins, ak *auth.Keeper) (accs []auth.BaseAccount) {
	for i := 0; i < numAccs; i++ {
		privKey := crypto.Ed25519PrivateKey{}.GenPrivateKey()
		pubKey := privKey.PublicKey()
		addr := sdk.Address(pubKey.Address())
		acc := auth.NewBaseAccountWithAddress(addr)
		acc.Coins = initialCoins
		acc.PubKey = pubKey
		ak.SetAccount(ctx, &acc)
		accs = append(accs, acc)
	}
	return
}

func createTestValidators(ctx sdk.Ctx, numAccs int, valCoins sdk.BigInt, nk *nodesKeeper.Keeper, ak auth.Keeper, kb keys.Keybase) (accs nodesTypes.Validators) {
	ethereum := hex.EncodeToString([]byte{01})
	for i := 0; i < numAccs-1; i++ {
		privKey := crypto.Ed25519PrivateKey{}.GenPrivateKey()
		pubKey := privKey.PublicKey()
		addr := sdk.Address(pubKey.Address())
		privKey2 := crypto.Ed25519PrivateKey{}.GenPrivateKey()
		pubKey2 := privKey2.PublicKey()
		addr2 := sdk.Address(pubKey2.Address())
		val := nodesTypes.NewValidator(addr, pubKey, []string{ethereum}, "https://www.google.com:443", valCoins, addr2)
		// set the vals from the data
		nk.SetValidator(ctx, val)
		nk.SetStakedValidatorByChains(ctx, val)
		// ensure there's a signing info entry for the val (used in slashing)
		_, found := nk.GetValidatorSigningInfo(ctx, val.GetAddress())
		if !found {
			signingInfo := nodesTypes.ValidatorSigningInfo{
				Address:     val.GetAddress(),
				StartHeight: ctx.BlockHeight(),
				JailedUntil: time.Unix(0, 0),
			}
			nk.SetValidatorSigningInfo(ctx, val.GetAddress(), signingInfo)
		}
		accs = append(accs, val)
	}
	// add self node to it
	kp, er := kb.GetCoinbase()
	if er != nil {
		panic(er)
	}
	val := nodesTypes.NewValidator(sdk.Address(kp.GetAddress()), kp.PublicKey, []string{ethereum}, "https://www.google.com:443", valCoins, kp.GetAddress())
	// set the vals from the data
	nk.SetValidator(ctx, val)
	nk.SetStakedValidatorByChains(ctx, val)
	// ensure there's a signing info entry for the val (used in slashing)
	_, found := nk.GetValidatorSigningInfo(ctx, val.GetAddress())
	if !found {
		signingInfo := nodesTypes.ValidatorSigningInfo{
			Address:     val.GetAddress(),
			StartHeight: ctx.BlockHeight(),
			JailedUntil: time.Unix(0, 0),
		}
		nk.SetValidatorSigningInfo(ctx, val.GetAddress(), signingInfo)
	}
	accs = append(accs, val)
	// end self node logic
	stakedTokens := sdk.NewInt(int64(numAccs)).Mul(valCoins)
	// take the staked amount and create the corresponding coins object
	stakedCoins := sdk.NewCoins(sdk.NewCoin(nk.StakeDenom(ctx), stakedTokens))
	// check if the staked pool accounts exists
	stakedPool := nk.GetStakedPool(ctx)
	// if the stakedPool is nil
	if stakedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", nodesTypes.StakedPoolName))
	}
	// add coins if not provided on genesis (there's an option to provide the coins in genesis)
	if stakedPool.GetCoins().IsZero() {
		if err := stakedPool.SetCoins(stakedCoins); err != nil {
			panic(err)
		}
		ak.SetModuleAccount(ctx, stakedPool)
	} else {
		// if it is provided in the genesis file then ensure the two are equal
		if !stakedPool.GetCoins().IsEqual(stakedCoins) {
			panic(fmt.Sprintf("%s module account total does not equal the amount in each validator account", nodesTypes.StakedPoolName))
		}
	}
	return
}

func createTestApps(ctx sdk.Ctx, numAccs int, valCoins sdk.BigInt, ak appsKeeper.Keeper, sk auth.Keeper) (accs appsTypes.Applications) {
	ethereum := hex.EncodeToString([]byte{01})
	for i := 0; i < numAccs; i++ {
		privKey := crypto.Ed25519PrivateKey{}.GenPrivateKey()
		pubKey := privKey.PublicKey()
		addr := sdk.Address(pubKey.Address())
		app := appsTypes.NewApplication(addr, pubKey, []string{ethereum}, valCoins)
		// set the vals from the data
		// calculate relays
		app.MaxRelays = ak.CalculateAppRelays(ctx, app)
		ak.SetApplication(ctx, app)
		ak.SetStakedApplication(ctx, app)
		accs = append(accs, app)
	}
	stakedTokens := sdk.NewInt(int64(numAccs)).Mul(valCoins)
	// take the staked amount and create the corresponding coins object
	stakedCoins := sdk.NewCoins(sdk.NewCoin(ak.StakeDenom(ctx), stakedTokens))
	// check if the staked pool accounts exists
	stakedPool := ak.GetStakedPool(ctx)
	// if the stakedPool is nil
	if stakedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", appsTypes.StakedPoolName))
	}
	// add coins if not provided on genesis (there's an option to provide the coins in genesis)
	if stakedPool.GetCoins().IsZero() {
		if err := stakedPool.SetCoins(stakedCoins); err != nil {
			panic(err)
		}
		sk.SetModuleAccount(ctx, stakedPool)
	} else {
		// if it is provided in the genesis file then ensure the two are equal
		if !stakedPool.GetCoins().IsEqual(stakedCoins) {
			panic(fmt.Sprintf("%s module account total does not equal the amount in each app account", appsTypes.StakedPoolName))
		}
	}
	return
}

func getRandomPrivateKey() crypto.Ed25519PrivateKey {
	return crypto.Ed25519PrivateKey{}.GenPrivateKey().(crypto.Ed25519PrivateKey)
}

func getRandomPubKey() crypto.Ed25519PublicKey {
	pk := crypto.Ed25519PrivateKey{}.GenPrivateKey()
	return pk.PublicKey().(crypto.Ed25519PublicKey)
}

func getRandomValidatorAddress() sdk.Address {
	return sdk.Address(getRandomPubKey().Address())
}

func simulateRelays(t *testing.T, k Keeper, ctx *sdk.Ctx, maxRelays int) (npk crypto.PublicKey, validHeader types.SessionHeader, keys simulateRelayKeys) {
	npk = getRandomPubKey()
	ethereum := hex.EncodeToString([]byte{01})
	clientKey := getRandomPrivateKey()
	validHeader = types.SessionHeader{
		ApplicationPubKey:  getTestApplication().PublicKey.RawString(),
		Chain:              ethereum,
		SessionBlockHeight: 1,
	}
	logger := log.NewNopLogger()
	types.InitConfig(&types.HostedBlockchains{
		M: make(map[string]types.HostedBlockchain),
	}, logger, sdk.DefaultTestingPocketConfig())

	// NOTE Add a minimum of 5 proofs to memInvoice to be able to create a merkle tree
	for j := 0; j < maxRelays; j++ {
		proof := createProof(getTestApplicationPrivateKey(), clientKey, npk, ethereum, j)
		types.SetProof(validHeader, types.RelayEvidence, proof, sdk.NewInt(100000))
	}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", k.storeKey).Return((*ctx).KVStore(k.storeKey))
	mockCtx.On("PrevCtx", validHeader.SessionBlockHeight).Return(*ctx, nil)
	mockCtx.On("Logger").Return((*ctx).Logger())
	keys = simulateRelayKeys{getTestApplicationPrivateKey(), clientKey}
	return
}
func createProof(private, client crypto.PrivateKey, npk crypto.PublicKey, chain string, entropy int) types.Proof {
	aat := types.AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: private.PublicKey().RawString(),
		ClientPublicKey:      client.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	sig, err := private.Sign(aat.Hash())
	if err != nil {
		panic(err)
	}
	aat.ApplicationSignature = hex.EncodeToString(sig)
	proof := types.RelayProof{
		Entropy:            int64(entropy + 1),
		RequestHash:        aat.HashString(), // fake
		SessionBlockHeight: 1,
		ServicerPubKey:     npk.RawString(),
		Blockchain:         chain,
		Token:              aat,
		Signature:          "",
	}
	clientSig, er := client.Sign(proof.Hash())
	if er != nil {
		panic(er)
	}
	proof.Signature = hex.EncodeToString(clientSig)
	return proof
}

// Ctx is an autogenerated mock type for the Ctx type
type Ctx struct {
	mock.Mock
}

// GetPrevBlockHash provides a mock function with given fields: height
func (_m *Ctx) GetPrevBlockHash(height int64) ([]byte, error) {
	ret := _m.Called(height)
	var r0 []byte
	if rf, ok := ret.Get(0).(func(int64) []byte); ok {
		r0 = rf(height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}
	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(height)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

func (_m *Ctx) IsPrevCtx() bool {
	return true
}

// BlockGasMeter provides a mock function with given fields:
func (_m *Ctx) BlockGasMeter() storeTypes.GasMeter {
	ret := _m.Called()

	var r0 storeTypes.GasMeter
	if rf, ok := ret.Get(0).(func() storeTypes.GasMeter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(storeTypes.GasMeter)
		}
	}

	return r0
}

// BlockHeader provides a mock function with given fields:
func (_m *Ctx) BlockHeader() abcitypes.Header {
	ret := _m.Called()

	var r0 abcitypes.Header
	if rf, ok := ret.Get(0).(func() abcitypes.Header); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(abcitypes.Header)
	}

	return r0
}

// BlockHeight provides a mock function with given fields:
func (_m *Ctx) BlockHeight() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// BlockStore provides a mock function with given fields:
func (_m *Ctx) BlockStore() *tmStore.BlockStore {
	ret := _m.Called()

	var r0 *tmStore.BlockStore
	if rf, ok := ret.Get(0).(func() *tmStore.BlockStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tmStore.BlockStore)
		}
	}

	return r0
}

// BlockTime provides a mock function with given fields:
func (_m *Ctx) BlockTime() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// CacheContext provides a mock function with given fields:
func (_m *Ctx) CacheContext() (pocketTypes.Context, func()) {
	ret := _m.Called()

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func() pocketTypes.Context); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	var r1 func()
	if rf, ok := ret.Get(1).(func() func()); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(func())
		}
	}

	return r0, r1
}

// ChainID provides a mock function with given fields:
func (_m *Ctx) ChainID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ConsensusParams provides a mock function with given fields:
func (_m *Ctx) ConsensusParams() *abcitypes.ConsensusParams {
	ret := _m.Called()

	var r0 *abcitypes.ConsensusParams
	if rf, ok := ret.Get(0).(func() *abcitypes.ConsensusParams); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcitypes.ConsensusParams)
		}
	}

	return r0
}

// Context provides a mock function with given fields:
func (_m *Ctx) Context() context.Context {
	ret := _m.Called()

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// EventManager provides a mock function with given fields:
func (_m *Ctx) EventManager() *pocketTypes.EventManager {
	ret := _m.Called()

	var r0 *pocketTypes.EventManager
	if rf, ok := ret.Get(0).(func() *pocketTypes.EventManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pocketTypes.EventManager)
		}
	}

	return r0
}

// GasMeter provides a mock function with given fields:
func (_m *Ctx) GasMeter() storeTypes.GasMeter {
	ret := _m.Called()

	var r0 storeTypes.GasMeter
	if rf, ok := ret.Get(0).(func() storeTypes.GasMeter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(storeTypes.GasMeter)
		}
	}

	return r0
}

// IsCheckTx provides a mock function with given fields:
func (_m *Ctx) IsCheckTx() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsZero provides a mock function with given fields:
func (_m *Ctx) IsZero() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsZero provides a mock function with given fields:
func (_m *Ctx) IsAfterUpgradeHeight() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsZero provides a mock function with given fields:
func (_m *Ctx) IsOnUpgradeHeight() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// KVStore provides a mock function with given fields: key
func (_m *Ctx) KVStore(key storeTypes.StoreKey) storeTypes.KVStore {
	ret := _m.Called(key)

	var r0 storeTypes.KVStore
	if rf, ok := ret.Get(0).(func(storeTypes.StoreKey) storeTypes.KVStore); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(storeTypes.KVStore)
		}
	}

	return r0
}

// Logger provides a mock function with given fields:
func (_m *Ctx) Logger() log.Logger {
	ret := _m.Called()

	var r0 log.Logger
	if rf, ok := ret.Get(0).(func() log.Logger); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(log.Logger)
		}
	}

	return r0
}

// MinGasPrices provides a mock function with given fields:
func (_m *Ctx) MinGasPrices() pocketTypes.DecCoins {
	ret := _m.Called()

	var r0 pocketTypes.DecCoins
	if rf, ok := ret.Get(0).(func() pocketTypes.DecCoins); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pocketTypes.DecCoins)
		}
	}

	return r0
}

// MultiStore provides a mock function with given fields:
func (_m *Ctx) MultiStore() storeTypes.MultiStore {
	ret := _m.Called()

	var r0 storeTypes.MultiStore
	if rf, ok := ret.Get(0).(func() storeTypes.MultiStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(storeTypes.MultiStore)
		}
	}

	return r0
}

// MustGetPrevCtx provides a mock function with given fields: height
func (_m *Ctx) MustGetPrevCtx(height int64) pocketTypes.Context {
	ret := _m.Called(height)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(int64) pocketTypes.Context); ok {
		r0 = rf(height)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// PrevCtx provides a mock function with given fields: height
func (_m *Ctx) PrevCtx(height int64) (pocketTypes.Context, error) {
	ret := _m.Called(height)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(int64) pocketTypes.Context); ok {
		r0 = rf(height)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransientStore provides a mock function with given fields: key
func (_m *Ctx) TransientStore(key storeTypes.StoreKey) storeTypes.KVStore {
	ret := _m.Called(key)

	var r0 storeTypes.KVStore
	if rf, ok := ret.Get(0).(func(storeTypes.StoreKey) storeTypes.KVStore); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(storeTypes.KVStore)
		}
	}

	return r0
}

// TxBytes provides a mock function with given fields:
func (_m *Ctx) TxBytes() []byte {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// Value provides a mock function with given fields: key
func (_m *Ctx) Value(key interface{}) interface{} {
	ret := _m.Called(key)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(interface{}) interface{}); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// VoteInfos provides a mock function with given fields:
func (_m *Ctx) VoteInfos() []abcitypes.VoteInfo {
	ret := _m.Called()

	var r0 []abcitypes.VoteInfo
	if rf, ok := ret.Get(0).(func() []abcitypes.VoteInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]abcitypes.VoteInfo)
		}
	}

	return r0
}

// WithBlockGasMeter provides a mock function with given fields: meter
func (_m *Ctx) WithBlockGasMeter(meter storeTypes.GasMeter) pocketTypes.Context {
	ret := _m.Called(meter)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(storeTypes.GasMeter) pocketTypes.Context); ok {
		r0 = rf(meter)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithBlockHeader provides a mock function with given fields: header
func (_m *Ctx) WithBlockHeader(header abcitypes.Header) pocketTypes.Context {
	ret := _m.Called(header)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(abcitypes.Header) pocketTypes.Context); ok {
		r0 = rf(header)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithBlockHeight provides a mock function with given fields: height
func (_m *Ctx) WithBlockHeight(height int64) pocketTypes.Context {
	ret := _m.Called(height)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(int64) pocketTypes.Context); ok {
		r0 = rf(height)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithBlockStore provides a mock function with given fields: bs
func (_m *Ctx) WithBlockStore(bs *tmStore.BlockStore) pocketTypes.Context {
	ret := _m.Called(bs)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(*tmStore.BlockStore) pocketTypes.Context); ok {
		r0 = rf(bs)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithBlockTime provides a mock function with given fields: newTime
func (_m *Ctx) WithBlockTime(newTime time.Time) pocketTypes.Context {
	ret := _m.Called(newTime)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(time.Time) pocketTypes.Context); ok {
		r0 = rf(newTime)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithChainID provides a mock function with given fields: chainID
func (_m *Ctx) WithChainID(chainID string) pocketTypes.Context {
	ret := _m.Called(chainID)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(string) pocketTypes.Context); ok {
		r0 = rf(chainID)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithConsensusParams provides a mock function with given fields: params
func (_m *Ctx) WithConsensusParams(params *abcitypes.ConsensusParams) pocketTypes.Context {
	ret := _m.Called(params)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(*abcitypes.ConsensusParams) pocketTypes.Context); ok {
		r0 = rf(params)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithContext provides a mock function with given fields: ctx
func (_m *Ctx) WithContext(ctx context.Context) pocketTypes.Context {
	ret := _m.Called(ctx)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(context.Context) pocketTypes.Context); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithEventManager provides a mock function with given fields: em
func (_m *Ctx) WithEventManager(em *pocketTypes.EventManager) pocketTypes.Context {
	ret := _m.Called(em)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(*pocketTypes.EventManager) pocketTypes.Context); ok {
		r0 = rf(em)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithGasMeter provides a mock function with given fields: meter
func (_m *Ctx) WithGasMeter(meter storeTypes.GasMeter) pocketTypes.Context {
	ret := _m.Called(meter)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(storeTypes.GasMeter) pocketTypes.Context); ok {
		r0 = rf(meter)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithIsCheckTx provides a mock function with given fields: isCheckTx
func (_m *Ctx) WithIsCheckTx(isCheckTx bool) pocketTypes.Context {
	ret := _m.Called(isCheckTx)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(bool) pocketTypes.Context); ok {
		r0 = rf(isCheckTx)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithLogger provides a mock function with given fields: logger
func (_m *Ctx) WithLogger(logger log.Logger) pocketTypes.Context {
	ret := _m.Called(logger)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(log.Logger) pocketTypes.Context); ok {
		r0 = rf(logger)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithMinGasPrices provides a mock function with given fields: gasPrices
func (_m *Ctx) WithMinGasPrices(gasPrices pocketTypes.DecCoins) pocketTypes.Context {
	ret := _m.Called(gasPrices)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(pocketTypes.DecCoins) pocketTypes.Context); ok {
		r0 = rf(gasPrices)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithMultiStore provides a mock function with given fields: ms
func (_m *Ctx) WithMultiStore(ms storeTypes.MultiStore) pocketTypes.Context {
	ret := _m.Called(ms)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(storeTypes.MultiStore) pocketTypes.Context); ok {
		r0 = rf(ms)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithProposer provides a mock function with given fields: addr
func (_m *Ctx) WithProposer(addr pocketTypes.Address) pocketTypes.Context {
	ret := _m.Called(addr)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(pocketTypes.Address) pocketTypes.Context); ok {
		r0 = rf(addr)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithTxBytes provides a mock function with given fields: txBytes
func (_m *Ctx) WithTxBytes(txBytes []byte) pocketTypes.Context {
	ret := _m.Called(txBytes)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func([]byte) pocketTypes.Context); ok {
		r0 = rf(txBytes)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithValue provides a mock function with given fields: key, value
func (_m *Ctx) WithValue(key interface{}, value interface{}) pocketTypes.Context {
	ret := _m.Called(key, value)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func(interface{}, interface{}) pocketTypes.Context); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

// WithValue provides a mock function with given fields: key, value
func (_m *Ctx) AppVersion() string {
	return ""
}

// WithVoteInfos provides a mock function with given fields: voteInfo
func (_m *Ctx) WithVoteInfos(voteInfo []abcitypes.VoteInfo) pocketTypes.Context {
	ret := _m.Called(voteInfo)

	var r0 pocketTypes.Context
	if rf, ok := ret.Get(0).(func([]abcitypes.VoteInfo) pocketTypes.Context); ok {
		r0 = rf(voteInfo)
	} else {
		r0 = ret.Get(0).(pocketTypes.Context)
	}

	return r0
}

func (_m *Ctx) ClearGlobalCache() {
	_m.Called()
}

func (_m *Ctx) BlockHash(cdc *codec.Codec, _ int64) ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).([]byte)
	}

	var r1 error
	if rf1, ok := ret.Get(1).(func() error); ok {
		r1 = rf1()
	} else {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
