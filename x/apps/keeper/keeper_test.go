package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/store"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/types/module"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/auth"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodeskeeper "github.com/pokt-network/pocket-core/x/nodes/keeper"
	nodestypes "github.com/pokt-network/pocket-core/x/nodes/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"testing"
)

func TestKeeper_Codespace(t *testing.T) {
	_, _, keeper := createTestInput(t, true)
	if got := keeper.Codespace(); got != "apps" {
		t.Errorf("Codespace() = %v, want %v", got, "apps")
	}
}

func TestKeepers_NewKeeper(t *testing.T) {
	tests := []struct {
		name     string
		hasError bool
		msg      string
	}{
		{
			name:     "create a keeper",
			hasError: false,
		},
		{
			name:     "errors if no GetModuleAddress is nill",
			msg:      fmt.Sprintf("%s module account has not been set", types.StakedPoolName),
			hasError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initPower := int64(100000000000)
			nAccs := int64(4)

			keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
			keyParams := sdk.ParamsKey
			tkeyParams := sdk.ParamsTKey
			nodesKey := sdk.NewKVStoreKey(nodestypes.StoreKey)
			appsKey := sdk.NewKVStoreKey(types.StoreKey)

			db := dbm.NewMemDB()
			ms := store.NewCommitMultiStore(db, false, 5000000)
			ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
			ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
			ms.MountStoreWithDB(nodesKey, sdk.StoreTypeIAVL, db)
			ms.MountStoreWithDB(appsKey, sdk.StoreTypeIAVL, db)
			ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
			err := ms.LoadLatestVersion()
			if err != nil {
				t.FailNow()
			}

			ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain"}, true, log.NewNopLogger()).WithAppVersion("0.0.0")
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
				nodestypes.StakedPoolName: {auth.Burner, auth.Staking},
				govTypes.DAOAccountName:   {auth.Burner, auth.Staking},
			}
			if !tt.hasError {
				maccPerms[types.StakedPoolName] = []string{auth.Burner, auth.Staking, auth.Minter}
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

			_ = createTestAccs(ctx, int(nAccs), initialCoins, &ak)

			if tt.hasError {
				return
			}
			_ = NewKeeper(cdc, appsKey, nk, ak, MockPocketKeeper{}, appSubspace, "apps")
		})
	}
}
