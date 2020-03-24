package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodeskeeper "github.com/pokt-network/pocket-core/x/nodes/keeper"
	nodestypes "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/store"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/bank"
	govTypes "github.com/pokt-network/posmint/x/gov/types"
	"github.com/pokt-network/posmint/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"strings"
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
		name   string
		panics bool
		msg    string
	}{
		{
			name:   "create a keeper",
			panics: false,
		},
		{
			name:   "errors if no GetModuleAddress is nill",
			msg:    fmt.Sprintf("%s module account has not been set", types.StakedPoolName),
			panics: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			if err != nil {
				t.FailNow()
			}

			ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain"}, true, log.NewNopLogger())
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
				nodestypes.StakedPoolName: {supply.Burner, supply.Staking},
				govTypes.DAOAccountName:   {supply.Burner, supply.Staking},
			}
			if !tt.panics {
				maccPerms[types.StakedPoolName] = []string{supply.Burner, supply.Staking, supply.Minter}
			}

			modAccAddrs := make(map[string]bool)
			for acc := range maccPerms {
				modAccAddrs[supply.NewModuleAddress(acc).String()] = true
			}
			valTokens := sdk.TokensFromConsensusPower(initPower)

			accSubspace := sdk.NewSubspace(auth.DefaultParamspace)
			bankSubspace := sdk.NewSubspace(bank.DefaultParamspace)
			nodesSubspace := sdk.NewSubspace(nodestypes.DefaultParamspace)
			appSubspace := sdk.NewSubspace(DefaultParamspace)
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

			if tt.panics {
				defer func() {
					err := recover()
					if !strings.Contains(err.(string), tt.msg) {
						t.Errorf("SetHooks(): got %v want %v", err, tt.msg)
					}
				}()
			}
			_ = NewKeeper(cdc, keySupply, bk, nk, sk, appSubspace, "apps")
		})
	}
}
