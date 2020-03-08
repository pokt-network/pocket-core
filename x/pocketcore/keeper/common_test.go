package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	apps "github.com/pokt-network/pocket-core/x/apps"
	appsKeeper "github.com/pokt-network/pocket-core/x/apps/keeper"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodesKeeper "github.com/pokt-network/pocket-core/x/nodes/keeper"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/crypto/keys"
	"github.com/pokt-network/posmint/store"
	storeTypes "github.com/pokt-network/posmint/store/types"
	posminttypes "github.com/pokt-network/posmint/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/pokt-network/posmint/x/params"
	"github.com/pokt-network/posmint/x/supply"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmStore "github.com/tendermint/tendermint/store"
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

type simulateRelayKeys struct {
	private crypto.PrivateKey
	client  crypto.PrivateKey
}

func NewTestKeybase() keys.Keybase {
	return keys.NewInMemory()
}

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
func newContext(t *testing.T, isCheckTx bool) sdk.Context {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	nodesKey := sdk.NewKVStoreKey(nodesTypes.StoreKey)
	appsKey := sdk.NewKVStoreKey(appsTypes.StoreKey)
	pocketKey := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(nodesKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(appsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(pocketKey, sdk.StoreTypeIAVL, db)
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
			Hash: types.Hash([]byte("fake")),
		},
	})
	return ctx
}

// nolint: deadcode unused
func createTestInput(t *testing.T, isCheckTx bool) (sdk.Ctx, []nodesTypes.Validator, []appsTypes.Application, []auth.BaseAccount, Keeper, map[string]*sdk.KVStoreKey) {
	initPower := int64(100000000000)
	nAccs := int64(5)

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	nodesKey := sdk.NewKVStoreKey(nodesTypes.StoreKey)
	appsKey := sdk.NewKVStoreKey(appsTypes.StoreKey)
	pocketKey := sdk.NewKVStoreKey(types.StoreKey)

	keys := make(map[string]*sdk.KVStoreKey)
	keys["params"] = keyParams
	keys["pos"] = nodesKey
	keys["application"] = appsKey

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(nodesKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(appsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(pocketKey, sdk.StoreTypeIAVL, db)
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
		Height: 976,
		Time:   time.Time{},
		LastBlockId: abci.BlockID{
			Hash: types.Hash([]byte("fake")),
		},
	})
	cdc := makeTestCodec()

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		appsTypes.StakedPoolName:  {supply.Burner, supply.Staking, supply.Minter},
		nodesTypes.StakedPoolName: {supply.Burner, supply.Staking},
		nodesTypes.DAOPoolName:    {supply.Burner, supply.Staking},
	}

	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}
	valTokens := sdk.TokensFromConsensusPower(initPower)

	ethereum, err := types.NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "geth",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}

	hb := types.HostedBlockchains{
		M: map[string]types.HostedBlockchain{ethereum: {
			Hash: ethereum,
			URL:  "https://www.google.com",
		}},
	}

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, modAccAddrs)
	sk := supply.NewKeeper(cdc, keySupply, ak, bk, maccPerms)
	nk := nodesKeeper.NewKeeper(cdc, nodesKey, ak, bk, sk, pk.Subspace(nodesTypes.DefaultParamspace), nodesTypes.ModuleName)
	appk := appsKeeper.NewKeeper(cdc, appsKey, bk, nk, sk, pk.Subspace(appsTypes.DefaultParamspace), appsTypes.ModuleName)
	appk.SetApplication(ctx, getTestApplication())
	keeper := NewPocketCoreKeeper(pocketKey, cdc, nk, appk, hb, pk.Subspace(types.DefaultParamspace), "test")
	kb := NewTestKeybase()
	_, err = kb.Create("test")
	assert.Nil(t, err)
	_, err = kb.GetCoinbase()
	assert.Nil(t, err)
	keeper.Keybase = kb
	moduleManager := module.NewManager(
		auth.NewAppModule(ak),
		bank.NewAppModule(bk, ak),
		supply.NewAppModule(sk, ak),
		nodes.NewAppModule(nk, ak, sk),
		apps.NewAppModule(appk, sk, nk),
	)
	genesisState := ModuleBasics.DefaultGenesis()
	moduleManager.InitGenesis(ctx, genesisState)
	initialCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, valTokens))
	accs := createTestAccs(ctx, int(nAccs), initialCoins, &ak)
	ap := createTestApps(ctx, int(nAccs), sdk.NewInt(10000000), appk, sk)
	vals := createTestValidators(ctx, int(nAccs), sdk.NewInt(10000000), sdk.ZeroInt(), &nk, sk, kb)
	appk.SetParams(ctx, appsTypes.DefaultParams())
	nk.SetParams(ctx, nodesTypes.DefaultParams())
	defaultPocketParams := types.DefaultParams()
	defaultPocketParams.SupportedBlockchains = []string{getTestSupportedBlockchain()}
	keeper.SetParams(ctx, defaultPocketParams)
	return ctx, vals, ap, accs, keeper, keys
}

var (
	testApp            appsTypes.Application
	testAppPrivateKey  crypto.PrivateKey
	testSupportedChain string
)

func getTestSupportedBlockchain() string {
	if testSupportedChain == "" {
		testSupportedChain, _ = types.NonNativeChain{
			Ticker:  "eth",
			Netid:   "4",
			Version: "v1.9.9",
			Client:  "geth",
			Inter:   "",
		}.HashString()
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

// nolint: unparam deadcode unused
func createTestAccs(ctx sdk.Ctx, numAccs int, initialCoins sdk.Coins, ak *auth.AccountKeeper) (accs []auth.BaseAccount) {
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

func createTestValidators(ctx sdk.Ctx, numAccs int, valCoins sdk.Int, daoCoins sdk.Int, nk *nodesKeeper.Keeper, sk supply.Keeper, kb keys.Keybase) (accs nodesTypes.Validators) {
	ethereum, err := types.NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "geth",
		Inter:   "",
	}.HashString()
	if err != nil {
		panic(err)
	}
	for i := 0; i < numAccs-1; i++ {
		privKey := crypto.Ed25519PrivateKey{}.GenPrivateKey()
		pubKey := privKey.PublicKey()
		addr := sdk.Address(pubKey.Address())
		val := nodesTypes.NewValidator(addr, pubKey, []string{ethereum}, "https://www.google.com", valCoins)
		// set the vals from the data
		nk.SetValidator(ctx, val)
		nk.SetStakedValidator(ctx, val)
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
		panic(err)
	}
	val := nodesTypes.NewValidator(sdk.Address(kp.GetAddress()), kp.PublicKey, []string{ethereum}, "https://www.google.com", valCoins)
	// set the vals from the data
	nk.SetValidator(ctx, val)
	nk.SetStakedValidator(ctx, val)
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
	// check if the dao pool account exists
	daoPool := nk.GetDAOPool(ctx)
	if daoPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", nodesTypes.DAOPoolName))
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
			panic(fmt.Sprintf("%s module account total does not equal the amount in each validator account", nodesTypes.StakedPoolName))
		}
	}
	// if the dao pool has zero tokens (not provided in genesis file)
	if daoPool.GetCoins().IsZero() {
		// ad the coins
		if err := daoPool.SetCoins(sdk.NewCoins(sdk.NewCoin(nk.StakeDenom(ctx), daoCoins))); err != nil {
			panic(err)
		}
	}
	return
}

func createTestApps(ctx sdk.Ctx, numAccs int, valCoins sdk.Int, ak appsKeeper.Keeper, sk supply.Keeper) (accs appsTypes.Applications) {
	ethereum, err := types.NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "geth",
		Inter:   "",
	}.HashString()
	if err != nil {
		panic(err)
	}
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
func simulateRelays(t *testing.T, k Keeper, ctx *sdk.Ctx, maxRelays int) (npk crypto.PublicKey, evidenceMap *types.EvidenceMap, validHeader types.SessionHeader, keys simulateRelayKeys, receipt types.Receipt) {
	npk = getRandomPubKey()
	evidenceMap = types.GetEvidenceMap()

	ethereum, err := types.NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "geth",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	clientKey := getRandomPrivateKey()

	validHeader = types.SessionHeader{
		ApplicationPubKey:  getTestApplication().PublicKey.RawString(),
		Chain:              ethereum,
		SessionBlockHeight: 1,
	}
	receipt = types.Receipt{
		SessionHeader:   validHeader,
		ServicerAddress: sdk.Address(npk.Address()).String(),
		Total:           2000,
		EvidenceType:    types.RelayEvidence,
	}

	// NOTE Add a minimum of 5 proofs to memInvoice to be able to create a merkle tree
	for j := 0; j < maxRelays; j++ {
		proof := createProof(getTestApplicationPrivateKey(), clientKey, npk, ethereum, j)
		evidenceMap.AddToEvidence(validHeader, proof)
	}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", k.storeKey).Return((*ctx).KVStore(k.storeKey))
	mockCtx.On("PrevCtx", validHeader.SessionBlockHeight).Return(*ctx, nil)
	mockCtx.On("Logger").Return((*ctx).Logger())
	k.SetReceipts(mockCtx, []types.Receipt{receipt})
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
func (_m *Ctx) CacheContext() (posminttypes.Context, func()) {
	ret := _m.Called()

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func() posminttypes.Context); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
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
func (_m *Ctx) EventManager() *posminttypes.EventManager {
	ret := _m.Called()

	var r0 *posminttypes.EventManager
	if rf, ok := ret.Get(0).(func() *posminttypes.EventManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*posminttypes.EventManager)
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
func (_m *Ctx) MinGasPrices() posminttypes.DecCoins {
	ret := _m.Called()

	var r0 posminttypes.DecCoins
	if rf, ok := ret.Get(0).(func() posminttypes.DecCoins); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(posminttypes.DecCoins)
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
func (_m *Ctx) MustGetPrevCtx(height int64) posminttypes.Context {
	ret := _m.Called(height)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(int64) posminttypes.Context); ok {
		r0 = rf(height)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// PrevCtx provides a mock function with given fields: height
func (_m *Ctx) PrevCtx(height int64) (posminttypes.Context, error) {
	ret := _m.Called(height)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(int64) posminttypes.Context); ok {
		r0 = rf(height)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
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
func (_m *Ctx) WithBlockGasMeter(meter storeTypes.GasMeter) posminttypes.Context {
	ret := _m.Called(meter)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(storeTypes.GasMeter) posminttypes.Context); ok {
		r0 = rf(meter)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithBlockHeader provides a mock function with given fields: header
func (_m *Ctx) WithBlockHeader(header abcitypes.Header) posminttypes.Context {
	ret := _m.Called(header)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(abcitypes.Header) posminttypes.Context); ok {
		r0 = rf(header)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithBlockHeight provides a mock function with given fields: height
func (_m *Ctx) WithBlockHeight(height int64) posminttypes.Context {
	ret := _m.Called(height)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(int64) posminttypes.Context); ok {
		r0 = rf(height)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithBlockStore provides a mock function with given fields: bs
func (_m *Ctx) WithBlockStore(bs *tmStore.BlockStore) posminttypes.Context {
	ret := _m.Called(bs)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(*tmStore.BlockStore) posminttypes.Context); ok {
		r0 = rf(bs)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithBlockTime provides a mock function with given fields: newTime
func (_m *Ctx) WithBlockTime(newTime time.Time) posminttypes.Context {
	ret := _m.Called(newTime)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(time.Time) posminttypes.Context); ok {
		r0 = rf(newTime)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithChainID provides a mock function with given fields: chainID
func (_m *Ctx) WithChainID(chainID string) posminttypes.Context {
	ret := _m.Called(chainID)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(string) posminttypes.Context); ok {
		r0 = rf(chainID)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithConsensusParams provides a mock function with given fields: params
func (_m *Ctx) WithConsensusParams(params *abcitypes.ConsensusParams) posminttypes.Context {
	ret := _m.Called(params)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(*abcitypes.ConsensusParams) posminttypes.Context); ok {
		r0 = rf(params)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithContext provides a mock function with given fields: ctx
func (_m *Ctx) WithContext(ctx context.Context) posminttypes.Context {
	ret := _m.Called(ctx)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(context.Context) posminttypes.Context); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithEventManager provides a mock function with given fields: em
func (_m *Ctx) WithEventManager(em *posminttypes.EventManager) posminttypes.Context {
	ret := _m.Called(em)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(*posminttypes.EventManager) posminttypes.Context); ok {
		r0 = rf(em)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithGasMeter provides a mock function with given fields: meter
func (_m *Ctx) WithGasMeter(meter storeTypes.GasMeter) posminttypes.Context {
	ret := _m.Called(meter)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(storeTypes.GasMeter) posminttypes.Context); ok {
		r0 = rf(meter)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithIsCheckTx provides a mock function with given fields: isCheckTx
func (_m *Ctx) WithIsCheckTx(isCheckTx bool) posminttypes.Context {
	ret := _m.Called(isCheckTx)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(bool) posminttypes.Context); ok {
		r0 = rf(isCheckTx)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithLogger provides a mock function with given fields: logger
func (_m *Ctx) WithLogger(logger log.Logger) posminttypes.Context {
	ret := _m.Called(logger)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(log.Logger) posminttypes.Context); ok {
		r0 = rf(logger)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithMinGasPrices provides a mock function with given fields: gasPrices
func (_m *Ctx) WithMinGasPrices(gasPrices posminttypes.DecCoins) posminttypes.Context {
	ret := _m.Called(gasPrices)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(posminttypes.DecCoins) posminttypes.Context); ok {
		r0 = rf(gasPrices)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithMultiStore provides a mock function with given fields: ms
func (_m *Ctx) WithMultiStore(ms storeTypes.MultiStore) posminttypes.Context {
	ret := _m.Called(ms)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(storeTypes.MultiStore) posminttypes.Context); ok {
		r0 = rf(ms)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithProposer provides a mock function with given fields: addr
func (_m *Ctx) WithProposer(addr posminttypes.Address) posminttypes.Context {
	ret := _m.Called(addr)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(posminttypes.Address) posminttypes.Context); ok {
		r0 = rf(addr)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithTxBytes provides a mock function with given fields: txBytes
func (_m *Ctx) WithTxBytes(txBytes []byte) posminttypes.Context {
	ret := _m.Called(txBytes)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func([]byte) posminttypes.Context); ok {
		r0 = rf(txBytes)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithValue provides a mock function with given fields: key, value
func (_m *Ctx) WithValue(key interface{}, value interface{}) posminttypes.Context {
	ret := _m.Called(key, value)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func(interface{}, interface{}) posminttypes.Context); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}

// WithVoteInfos provides a mock function with given fields: voteInfo
func (_m *Ctx) WithVoteInfos(voteInfo []abcitypes.VoteInfo) posminttypes.Context {
	ret := _m.Called(voteInfo)

	var r0 posminttypes.Context
	if rf, ok := ret.Get(0).(func([]abcitypes.VoteInfo) posminttypes.Context); ok {
		r0 = rf(voteInfo)
	} else {
		r0 = ret.Get(0).(posminttypes.Context)
	}

	return r0
}
