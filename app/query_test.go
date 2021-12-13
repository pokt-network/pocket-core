// nolint
package app

import (
	"encoding/hex"
	"encoding/json"
	"github.com/pokt-network/pocket-core/codec"
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/crypto/keys"
	sdk "github.com/pokt-network/pocket-core/types"
	apps "github.com/pokt-network/pocket-core/x/apps"
	types3 "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/gov"
	"github.com/pokt-network/pocket-core/x/nodes"
	types2 "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/node"
	tmTypes "github.com/tendermint/tendermint/types"
	"gopkg.in/h2non/gock.v1"
)

func TestQueryBlock(t *testing.T) {

	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query block amino", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query block proto", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, _, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			height := int64(1)
			<-evtChan // Wait for block
			got, err := PCA.QueryBlock(&height)
			assert.Nil(t, err)
			assert.NotNil(t, got)

			cleanup()
			stopCli()
		})
	}
}

func TestQueryChainHeight(t *testing.T) {

	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query height amino", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query height proto", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, _, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			got, err := PCA.QueryHeight()
			assert.Nil(t, err)
			assert.Equal(t, int64(1), got) // should not be 0 due to empty blocks

			cleanup()
			stopCli()
		})
	}
}

func TestQueryTx(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query tx from proto account with proto cdec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			time.Sleep(time.Second * 2)
			cb, err := kb.GetCoinbase()
			assert.Nil(t, err)
			kp, err := kb.Create("test")
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var tx *sdk.TxResponse
			<-evtChan // Wait for block
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			tx, err = nodes.Send(memCodec(), memCli, kb, cb.GetAddress(), kp.GetAddress(), "test", sdk.NewInt(1000), tc.upgrades.codecUpgrade.upgradeMod)
			assert.Nil(t, err)
			assert.NotNil(t, tx)

			<-evtChan // Wait for tx
			got, err := PCA.QueryTx(tx.TxHash, false)
			assert.Nil(t, err)
			balance, err := PCA.QueryBalance(kp.GetAddress().String(), PCA.BaseApp.LastBlockHeight())
			assert.Nil(t, err)
			assert.Equal(t, int64(1000), balance.Int64())
			assert.NotNil(t, got)

			cleanup()
			stopCli()
		})
	}
}

func TestQueryAminoTx(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query tx amino", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			time.Sleep(time.Second * 2)
			cb, err := kb.GetCoinbase()
			assert.Nil(t, err)
			kp, err := kb.Create("test")
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var tx *sdk.TxResponse
			<-evtChan // Wait for block
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			tx, err = nodes.Send(memCodec(), memCli, kb, cb.GetAddress(), kp.GetAddress(), "test", sdk.NewInt(1000), tc.upgrades.codecUpgrade.upgradeMod)
			assert.Nil(t, err)
			assert.NotNil(t, tx)

			<-evtChan // Wait for tx
			got, err := PCA.QueryTx(tx.TxHash, false)
			assert.Nil(t, err)
			validator, err := PCA.QueryBalance(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.True(t, validator.Equal(sdk.NewInt(1000)))
			assert.NotNil(t, got)

			cleanup()
			stopCli()
		})
	}
}

func TestQueryValidators(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query validators proto", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 1}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			gen, _ := twoValTwoNodeGenesisState()
			_, _, cleanup := tc.memoryNodeFn(t, gen)
			time.Sleep(2 * time.Second)
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			got, err := PCA.QueryNodes(PCA.LastBlockHeight(), types2.QueryValidatorsParams{Page: 1, Limit: 1})
			assert.Nil(t, err)
			res := got.Result.([]types2.Validator)
			assert.Equal(t, 1, len(res))
			got, err = PCA.QueryNodes(0, types2.QueryValidatorsParams{Page: 2, Limit: 1})
			assert.Nil(t, err)
			res = got.Result.([]types2.Validator)
			assert.Equal(t, 1, len(res))
			got, err = PCA.QueryNodes(0, types2.QueryValidatorsParams{Page: 1, Limit: 1000})
			assert.Nil(t, err)
			res = got.Result.([]types2.Validator)
			assert.Equal(t, 2, len(res))
			cleanup()
			stopCli()
		})
	}
}
func TestQueryApps(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query apps from amino account with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query apps from proto account with proto cdec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform necessary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			time.Sleep(time.Second * 2)
			kp, err := kb.GetCoinbase()
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var tx *sdk.TxResponse
			var chains = []string{"0001"}

			<-evtChan // Wait for block
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			tx, err = apps.StakeTx(memCodec(), memCli, kb, chains, sdk.NewInt(1000000), kp, "test", tc.codecUpgrade.upgradeMod)
			assert.Nil(t, err)
			assert.NotNil(t, tx)

			<-evtChan // Wait for tx
			got, err := PCA.QueryApps(PCA.LastBlockHeight(), types3.QueryApplicationsWithOpts{
				Page:  1,
				Limit: 1,
			})
			assert.Nil(t, err)
			slice, ok := takeArg(got.Result, reflect.Slice)
			if !ok {
				t.Fatalf("couldn't convert arg to slice")
			}
			assert.Equal(t, 1, slice.Len())
			got, err = PCA.QueryApps(PCA.LastBlockHeight(), types3.QueryApplicationsWithOpts{
				Page:  2,
				Limit: 1,
			})
			assert.Nil(t, err)
			slice, ok = takeArg(got.Result, reflect.Slice)
			if !ok {
				t.Fatalf("couldn't convert arg to slice")
			}
			assert.Equal(t, 1, slice.Len())
			got, err = PCA.QueryApps(PCA.LastBlockHeight(), types3.QueryApplicationsWithOpts{
				Page:  1,
				Limit: 2,
			})
			assert.Nil(t, err)
			slice, ok = takeArg(got.Result, reflect.Slice)
			if !ok {
				t.Fatalf("couldn't convert arg to slice")
			}
			assert.Equal(t, 2, slice.Len())

			stopCli()
			cleanup()
		})
	}
}

func TestQueryValidator(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query validator amino", memoryNodeFn: NewInMemoryTendermintNodeAmino},
		{name: "query validator proto", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			cb, err := kb.GetCoinbase()
			if err != nil {
				t.Fatal(err)
			}
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			got, err := PCA.QueryNode(cb.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.Equal(t, cb.GetAddress(), got.Address)
			assert.False(t, got.Jailed)
			assert.True(t, got.StakedTokens.Equal(sdk.NewInt(1000000000000000)))

			cleanup()
			stopCli()
		})
	}
}

func TestQueryDaoBalance(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query dao balance from amino account with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query dao balance from proto account with proto cdec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, _, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			got, err := gov.QueryDAO(memCodec(), memCli, PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.Equal(t, big.NewInt(1000), got.BigInt())

			cleanup()
			stopCli()
		})
	}
}

func TestQueryACL(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query dao balance from amino account with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query dao balance from proto account with proto cdec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, _, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			got, err := gov.QueryACL(memCodec(), memCli, PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.Equal(t, got, testACL)

			cleanup()
			stopCli()
		})
	}
}

func TestQueryDaoOwner(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query dao owner from amino account with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query dao owner from proto account with proto cdec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, _, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			kb := getInMemoryKeybase()
			cb, err := kb.GetCoinbase()
			if err != nil {
				t.Fatal(err)
			}
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			got, err := gov.QueryDAOOwner(memCodec(), memCli, PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.Equal(t, got.String(), cb.GetAddress().String())

			cleanup()
			stopCli()
		})
	}
}

func TestQueryUpgrade(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query upgrade with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query upgrade with proto codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, _, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			var err error
			got, err := gov.QueryUpgrade(memCodec(), memCli, PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.Equal(t, got.UpgradeHeight(), int64(10000))

			cleanup()
			stopCli()
		})
	}
}

func TestQuerySupply(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query supply amino", memoryNodeFn: NewInMemoryTendermintNodeAmino},
		{name: "query supply proto", memoryNodeFn: NewInMemoryTendermintNodeProto},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, _, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			gotStaked, total, err := PCA.QueryTotalNodeCoins(PCA.LastBlockHeight())
			//fmt.Println(err)
			assert.Nil(t, err)
			//fmt.Println(gotStaked, total)
			assert.True(t, gotStaked.Equal(sdk.NewInt(1000000000000000)))
			assert.True(t, total.Equal(sdk.NewInt(1000002010001000)))

			cleanup()
			stopCli()
		})
	}
}

func TestQueryPOSParams(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query POS params amino", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query POS params proto", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, _, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			got, err := PCA.QueryNodeParams(PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.Equal(t, int64(5000), got.MaxValidators)
			assert.Equal(t, int64(1000000), got.StakeMinimum)
			assert.Equal(t, int64(10), got.DAOAllocation)
			assert.Equal(t, sdk.DefaultStakeDenom, got.StakeDenom)

			cleanup()
			stopCli()
		})
	}
}

func TestAccountBalance(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query account balance from amino", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query account balance from proto", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			cb, err := kb.GetCoinbase()
			assert.Nil(t, err)
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			got, err := PCA.QueryBalance(cb.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, got, sdk.NewInt(1000000000))

			cleanup()
			stopCli()
		})
	}
}

func TestQuerySigningInfo(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query signign info amino ", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query signing info proto", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			cb, err := kb.GetCoinbase()
			assert.Nil(t, err)
			cbAddr := cb.GetAddress()
			assert.Nil(t, err)
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			got, err := PCA.QuerySigningInfo(0, cbAddr.String())
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, got.Address.String(), cbAddr.String())

			cleanup()
			stopCli()
		})
	}
}

func TestQueryPocketSupportedBlockchains(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query supported blockchains amino with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query supported blockchains from proto with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, _, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			var err error
			got, err := PCA.QueryPocketSupportedBlockchains(PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Contains(t, got, sdk.PlaceholderHash)

			cleanup()
			stopCli()
		})
	}
}

func TestQueryPocketParams(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query pocket params amino ", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query pocket params proto", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, _, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			got, err := PCA.QueryPocketParams(PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, int64(5), got.SessionNodeCount)
			assert.Equal(t, int64(3), got.ClaimSubmissionWindow)
			assert.Equal(t, int64(100), got.ClaimExpiration)
			assert.Contains(t, got.SupportedBlockchains, sdk.PlaceholderHash)

			cleanup()
			stopCli()
		})
	}
}

func TestQueryAccount(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query account amino", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query account proto", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			acc := getUnstakedAccount(kb)
			assert.NotNil(t, acc)
			<-evtChan // Wait for block
			got, err := PCA.QueryAccount(acc.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, acc.GetAddress(), (*got).GetAddress())

			cleanup()
			stopCli()
		})
	}
}

func TestQueryStakedApp(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query query staked app amino", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query query staked app proto", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			time.Sleep(2 * time.Second)
			kp, err := kb.GetCoinbase()
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var tx *sdk.TxResponse
			var chains = []string{"0001"}
			<-evtChan // Wait for block
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			tx, err = apps.StakeTx(memCodec(), memCli, kb, chains, sdk.NewInt(1000000), kp, "test", tc.upgrades.codecUpgrade.upgradeMod)
			assert.Nil(t, err)
			assert.NotNil(t, tx)

			<-evtChan // Wait for  tx
			got, err := PCA.QueryApp(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, sdk.Staked, got.Status)
			assert.Equal(t, false, got.Jailed)

			cleanup()
			stopCli()
		})
	}
}

func TestRelayGenerator(t *testing.T) {
	const appPrivKey = "70906c8e250352e811a6ca994b674c4da1c6ba4be1e0b3edeadaf59979236c96a25e182d490e9722e72ba90eb21fe0124d03bcb75d2bf6f45b2a1d2b1dc92fac"
	const nodePublicKey = "a25e182d490e9722e72ba90eb21fe0124d03bcb75d2bf6f45b2a1d2b1dc92fac"
	const sessionBlockheight = 1
	const query = `{"jsonrpc":"2.0","method":"net_version","params":[],"id":67}`
	const supportedBlockchain = "0001"
	apkBz, err := hex.DecodeString(appPrivKey)
	if err != nil {
		panic(err)
	}
	var appPrivateKey crypto.Ed25519PrivateKey
	copy(appPrivateKey[:], apkBz)
	aat := types.AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPrivateKey.PublicKey().RawString(),
		ClientPublicKey:      appPrivateKey.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	sig, err := appPrivateKey.Sign(aat.Hash())
	if err != nil {
		panic(err)
	}
	aat.ApplicationSignature = hex.EncodeToString(sig)
	payload := types.Payload{
		Data: query,
	}
	// setup relay
	relay := types.Relay{
		Payload: payload,
		Proof: types.RelayProof{
			Entropy:            int64(rand.Int()),
			SessionBlockHeight: sessionBlockheight,
			ServicerPubKey:     nodePublicKey,
			Blockchain:         supportedBlockchain,
			Token:              aat,
			Signature:          "",
		},
	}
	relay.Proof.RequestHash = relay.RequestHashString()
	sig, err = appPrivateKey.Sign(relay.Proof.Hash())
	if err != nil {
		panic(err)
	}
	relay.Proof.Signature = hex.EncodeToString(sig)
	_, err = json.MarshalIndent(relay, "", "  ")
	if err != nil {
		panic(err)
	}
}

func TestQueryRelay(t *testing.T) {
	const headerKey = "foo"
	const headerVal = "bar"

	expectedRequest := `"jsonrpc":"2.0","method":"web3_sha3","params":["0x68656c6c6f20776f726c64"],"id":64`
	expectedResponse := "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad"
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query relay amino", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query relay proto", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sdk.VbCCache = sdk.NewCache(1)
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			time.Sleep(time.Second * 2)
			genBz, _, validators, app := fiveValidatorsOneAppGenesis()
			// setup relay endpoint
			gock.New(sdk.PlaceholderURL).
				Post("").
				BodyString(expectedRequest).
				MatchHeader(headerKey, headerVal).
				Reply(200).
				BodyString(expectedResponse)
			_, kb, cleanup := tc.memoryNodeFn(t, genBz)
			appPrivateKey, err := kb.ExportPrivateKeyObject(app.Address, "test")
			assert.Nil(t, err)
			// setup AAT
			aat := types.AAT{
				Version:              "0.0.1",
				ApplicationPublicKey: appPrivateKey.PublicKey().RawString(),
				ClientPublicKey:      appPrivateKey.PublicKey().RawString(),
				ApplicationSignature: "",
			}
			sig, err := appPrivateKey.Sign(aat.Hash())
			if err != nil {
				panic(err)
			}
			aat.ApplicationSignature = hex.EncodeToString(sig)
			payload := types.Payload{
				Data:    expectedRequest,
				Headers: map[string]string{headerKey: headerVal},
			}
			// setup relay
			relay := types.Relay{
				Payload: payload,
				Meta:    types.RelayMeta{BlockHeight: 5}, // todo race condition here
				Proof: types.RelayProof{
					Entropy:            32598345349034509,
					SessionBlockHeight: 1,
					ServicerPubKey:     validators[0].PublicKey.RawString(),
					Blockchain:         sdk.PlaceholderHash,
					Token:              aat,
					Signature:          "",
				},
			}
			relay.Proof.RequestHash = relay.RequestHashString()
			sig, err = appPrivateKey.Sign(relay.Proof.Hash())
			if err != nil {
				panic(err)
			}
			relay.Proof.Signature = hex.EncodeToString(sig)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			res, _, err := PCA.HandleRelay(relay)
			assert.Nil(t, err, err)
			assert.Equal(t, expectedResponse, res.Response)
			gock.New(sdk.PlaceholderURL).
				Post("").
				BodyString(expectedRequest).
				Reply(200).
				BodyString(expectedResponse)
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			select {
			case <-evtChan:
				inv, err := types.GetEvidence(types.SessionHeader{
					ApplicationPubKey:  aat.ApplicationPublicKey,
					Chain:              relay.Proof.Blockchain,
					SessionBlockHeight: relay.Proof.SessionBlockHeight,
				}, types.RelayEvidence, sdk.NewInt(10000))
				assert.Nil(t, err)
				assert.NotNil(t, inv)
				assert.Equal(t, inv.NumOfProofs, int64(1))
				cleanup()
				stopCli()
				gock.Off()
				return
			}
		})
	}
}

func TestQueryDispatch(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query dispatch amino", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query dispatch proto", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sdk.VbCCache = sdk.NewCache(1)
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			genBz, _, validators, app := fiveValidatorsOneAppGenesis()
			_, kb, cleanup := tc.memoryNodeFn(t, genBz)
			appPrivateKey, err := kb.ExportPrivateKeyObject(app.Address, "test")
			assert.Nil(t, err)
			// Setup HandleDispatch Request
			key := types.SessionHeader{
				ApplicationPubKey:  appPrivateKey.PublicKey().RawString(),
				Chain:              sdk.PlaceholderHash,
				SessionBlockHeight: 1,
			}
			// setup the query
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			res, err := PCA.HandleDispatch(key)
			assert.Nil(t, err)
			for _, val := range validators {
				assert.Contains(t, res.Session.SessionNodes, val)
			}
			cleanup()
			stopCli()
		})
	}
}

func TestQueryAllParams(t *testing.T) {

	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query all params amino ", memoryNodeFn: NewInMemoryTendermintNodeAmino},
		{name: "query all params proto", memoryNodeFn: NewInMemoryTendermintNodeProto},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resetTestACL()
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, _, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			res, err := PCA.QueryAllParams(PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.NotNil(t, res)

			assert.NotZero(t, len(res.AppParams))
			cleanup()
		})
	}
}
func TestQueryParam(t *testing.T) {

	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query param amino ", memoryNodeFn: NewInMemoryTendermintNodeAmino},
		{name: "query param proto ", memoryNodeFn: NewInMemoryTendermintNodeProto},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resetTestACL()
			_, _, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			res, err := PCA.QueryParam(0, "pocketcore/SupportedBlockchains")
			assert.Nil(t, err)
			assert.NotNil(t, res)

			assert.NotNil(t, res.Value)
			cleanup()
		})
	}
}

func TestQueryAccountBalance(t *testing.T) {

	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query staked app amino with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		// {name: "query staked app from amino with proto codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
		{name: "query staked app params from proto with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			acc := getUnstakedAccount(kb)
			assert.NotNil(t, acc)
			<-evtChan // Wait for block
			got, err := PCA.QueryBalance(acc.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, sdk.NewInt(1000000000), got)
			cleanup()
			stopCli()
		})
	}
}

func TestQueryNonExistingAccountBalance(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "query non existing account balance amino with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "query staked app from amino with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, _, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			got, err := PCA.QueryBalance("802fddec29f99cae7a601cf648eafced1c062d39", PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, sdk.NewInt(0), got)
			cleanup()
			stopCli()
		})
	}
}
