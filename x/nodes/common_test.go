package nodes

import (
	"github.com/pokt-network/pocket-core/x/nodes/keeper"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/store"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/gov"
	govTypes "github.com/pokt-network/posmint/x/gov/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"math/rand"
	"testing"
)

// nolint: deadcode unused
var (
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		gov.AppModuleBasic{},
	)
)

// nolint: deadcode unused
// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

func GetTestTendermintClient() client.Client {
	var tmNodeURI string
	var defaultTMURI = "tcp://localhost:26657"

	if tmNodeURI == "" {
		return client.NewHTTP(defaultTMURI, "/websocket")
	}
	return client.NewHTTP(tmNodeURI, "/websocket")
}

// nolint: deadcode unused
func createTestInput(t *testing.T, isCheckTx bool) (sdk.Context, []auth.Account, keeper.Keeper) {
	initPower := int64(100000000000)
	nAccs := int64(4)

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyPOS := sdk.NewKVStoreKey(types.ModuleName)
	keyParams := sdk.ParamsKey
	tkeyParams := sdk.ParamsTKey

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyPOS, sdk.StoreTypeIAVL, db)
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
		auth.FeeCollectorName:   nil,
		types.StakedPoolName:    {auth.Burner, auth.Staking, auth.Minter},
		govTypes.DAOAccountName: {auth.Burner, auth.Staking, auth.Minter},
	}
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[auth.NewModuleAddress(acc).String()] = true
	}
	valTokens := sdk.TokensFromConsensusPower(initPower)

	accSubspace := sdk.NewSubspace(auth.DefaultParamspace)
	posSubspace := sdk.NewSubspace(types.DefaultParamspace)

	ak := auth.NewKeeper(cdc, keyAcc, accSubspace, maccPerms)
	moduleManager := module.NewManager(
		auth.NewAppModule(ak),
	)

	genesisState := ModuleBasics.DefaultGenesis()
	moduleManager.InitGenesis(ctx, genesisState)

	initialCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, valTokens))
	accs := createTestAccs(ctx, int(nAccs), initialCoins, &ak)

	keeper := keeper.NewKeeper(cdc, keyPOS, ak, posSubspace, sdk.CodespaceType("pos"))

	params := types.DefaultParams()
	keeper.SetParams(ctx, params)
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

//func addMintedCoinsToModule(t *testing.T, ctx sdk.Ctx, k *keeper.Keeper, module string) {
//	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), sdk.NewInt(100000000000)))
//	mintErr := k.supplyKeeper.MintCoins(ctx, module, coins.Add(coins))
//	if mintErr != nil {
//		t.Fail()
//	}
//}
//
//func sendFromModuleToAccount(t *testing.T, ctx sdk.Ctx, k *keeper.Keeper, module string, address sdk.Address, amount sdk.Int) {
//	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
//	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, module, sdk.Address(address), coins)
//	if err != nil {
//		t.Fail()
//	}
//}

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

func getGenesisStateForTest(ctx sdk.Ctx, keeper keeper.Keeper, defaultparams bool) types.GenesisState {
	keeper.SetPreviousProposer(ctx, sdk.GetAddress(getRandomPubKey()))
	var prm = types.DefaultParams()

	if !defaultparams {
		prm = keeper.GetParams(ctx)
	}
	prevStateTotalPower := keeper.PrevStateValidatorsPower(ctx)
	validators := keeper.GetAllValidators(ctx)
	var prevStateValidatorPowers []types.PrevStatePowerMapping
	keeper.IterateAndExecuteOverPrevStateValsByPower(ctx, func(addr sdk.Address, power int64) (stop bool) {
		prevStateValidatorPowers = append(prevStateValidatorPowers, types.PrevStatePowerMapping{Address: addr, Power: power})
		return false
	})
	signingInfos := make(map[string]types.ValidatorSigningInfo)
	missedBlocks := make(map[string][]types.MissedBlock)
	keeper.IterateAndExecuteOverValSigningInfo(ctx, func(address sdk.Address, info types.ValidatorSigningInfo) (stop bool) {
		addrstring := address.String()
		signingInfos[addrstring] = info
		localMissedBlocks := []types.MissedBlock{}

		keeper.IterateAndExecuteOverMissedArray(ctx, address, func(index int64, missed bool) (stop bool) {
			localMissedBlocks = append(localMissedBlocks, types.MissedBlock{index, missed})
			return false
		})
		missedBlocks[addrstring] = localMissedBlocks

		return false
	})
	prevProposer := keeper.GetPreviousProposer(ctx)

	return types.GenesisState{
		Params:                   prm,
		PrevStateTotalPower:      prevStateTotalPower,
		PrevStateValidatorPowers: prevStateValidatorPowers,
		Validators:               validators,
		Exported:                 true,
		SigningInfos:             signingInfos,
		MissedBlocks:             missedBlocks,
		PreviousProposer:         prevProposer,
	}

}
