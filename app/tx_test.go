// nolint
package app

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/crypto/keys"
	sdk "github.com/pokt-network/pocket-core/types"
	apps "github.com/pokt-network/pocket-core/x/apps"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/auth/types"
	"github.com/pokt-network/pocket-core/x/gov"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
	rand2 "github.com/tendermint/tendermint/libs/rand"
	"github.com/tendermint/tendermint/node"
	tmTypes "github.com/tendermint/tendermint/types"
)

func TestMain(m *testing.M) {
	pocketTypes.CleanPocketNodes()
	sdk.InitCtxCache(1)
	m.Run()
}

func TestUnstakeApp(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "unstake an amino app with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "unstake a proto app with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}}, // todo: FULL PROTO SCENARIO
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			time.Sleep(1 * time.Second)
			kp, err := kb.GetCoinbase()
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var tx *sdk.TxResponse
			var chains = []string{"0001"}
			<-evtChan // Wait for block
			memCli, _, evtChan := subscribeTo(t, tmTypes.EventTx)
			tx, err = apps.StakeTx(memCodec(), memCli, kb, chains, sdk.NewInt(1000000), kp, "test", tc.codecUpgrade.upgradeMod)
			assert.Nil(t, err)
			assert.NotNil(t, tx)

			<-evtChan // Wait for tx
			got, err := PCA.QueryApps(PCA.LastBlockHeight(), appsTypes.QueryApplicationsWithOpts{
				Page:  1,
				Limit: 1})
			assert.Nil(t, err)
			res := got.Result.(appsTypes.Applications)
			assert.Equal(t, 1, len(res))
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			_, _ = apps.UnstakeTx(memCodec(), memCli, kb, kp.GetAddress(), "test", tc.codecUpgrade.upgradeMod)

			<-evtChan // Wait for tx
			got, err = PCA.QueryApps(PCA.LastBlockHeight(), appsTypes.QueryApplicationsWithOpts{
				Page:          1,
				Limit:         1,
				StakingStatus: 1,
			})
			assert.Nil(t, err)
			res = got.Result.(appsTypes.Applications)
			assert.Equal(t, 1, len(res))
			got, err = PCA.QueryApps(PCA.LastBlockHeight(), appsTypes.QueryApplicationsWithOpts{
				Page:          1,
				Limit:         1,
				StakingStatus: 2,
			})
			assert.Nil(t, err)
			res = got.Result.(appsTypes.Applications)
			assert.Equal(t, 1, len(res)) // default genesis application

			cleanup()
			stopCli()
		})
	}
}
func TestStakeApp(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "stake app with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "stake a proto app with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}}, // TODO FULL PROTO SCENARIO
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}

			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			time.Sleep(1 * time.Second)
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
			got, err := PCA.QueryApps(PCA.LastBlockHeight(), appsTypes.QueryApplicationsWithOpts{
				Page:  1,
				Limit: 2,
			})
			assert.Nil(t, err)
			res := got.Result.(appsTypes.Applications)
			assert.Equal(t, 2, len(res))

			stopCli()
			cleanup()
		})
	}
}
func TestEditStakeApp(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "editStake a proto application with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			var newChains = []string{"2121"}
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			time.Sleep(1 * time.Second)
			kp, err := kb.GetCoinbase()
			assert.Nil(t, err)
			kps, err := kb.List()
			assert.Nil(t, err)
			for _, k := range kps {
				if !k.GetAddress().Equals(kp.GetAddress()) {
					kp = k
					break
				}
			}
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var tx *sdk.TxResponse
			<-evtChan // Wait for block
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			balance, err := PCA.QueryBalance(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			n, err := PCA.QueryApp(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			var newBalance = balance.Sub(sdk.NewInt(100000)).Add(n.StakedTokens)
			tx, err = apps.StakeTx(memCodec(), memCli, kb, newChains, newBalance, kp, "test", tc.upgrades.codecUpgrade.upgradeMod)
			assert.Nil(t, err)
			assert.NotNil(t, tx)
			<-evtChan // Wait for tx
			appUpdated, err := PCA.QueryApp(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			// assert not the same as the old node
			assert.NotEqual(t, appUpdated, n)
			// assert chains and stake updated
			assert.Equal(t, newChains, appUpdated.Chains)
			// assert chains and stake updated
			assert.Equal(t, newBalance, appUpdated.StakedTokens)
			cleanup()
			stopCli()
		})
	}
}

func TestUnstakeNode(t *testing.T) {
	tt := []struct {
		name           string
		memoryNodeFn   func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		outputIsSigner bool
		*upgrades
	}{
		{name: "unstake a proto node with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
		{name: "unstake an amino node with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "unstake a node with new msg before 8 upgrade", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}, eight0Upgrade: upgrade{height: 999999}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}

			// 8.0 release
			isAfter8 := false
			if tc.eight0Upgrade.height != 0 {
				codec.UpgradeFeatureMap[codec.NonCustodialUpdateKey] = tc.eight0Upgrade.height
				isAfter8 = true
			}

			var chains = []string{"0001"}
			gen, _ := twoValTwoNodeGenesisState()
			_, kb, cleanup := tc.memoryNodeFn(t, gen)
			time.Sleep(1 * time.Second)
			kp, err := kb.GetCoinbase()
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var tx *sdk.TxResponse
			<-evtChan // Wait for block
			memCli, _, evtChan := subscribeTo(t, tmTypes.EventTx)
			_, err = PCA.QueryBalance(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			signer := kp.GetAddress()
			if tc.outputIsSigner {
				list, err := kb.List()
				assert.Nil(t, err)
				signer = list[2].GetAddress()
			}
			tx, err = nodes.UnstakeTx(memCodec(), memCli, kb, kp.GetAddress(), signer, "test", tc.upgrades.codecUpgrade.upgradeMod, isAfter8)
			assert.Nil(t, err)
			assert.NotNil(t, tx)
			if isAfter8 {
				assert.Equal(t, 2, int(tx.Code))
				cleanup()

			} else {
				<-evtChan // Wait for tx
				_, _, evtChan = subscribeTo(t, tmTypes.EventNewBlockHeader)
				for {
					select {
					case res := <-evtChan:
						if len(res.Events["begin_unstake.module"]) == 1 {
							got, err := PCA.QueryNodes(PCA.LastBlockHeight(), nodeTypes.QueryValidatorsParams{StakingStatus: 1, JailedStatus: 0, Blockchain: "", Page: 1, Limit: 1}) // unstaking
							assert.Nil(t, err)
							res := got.Result.([]nodeTypes.Validator)
							assert.Equal(t, 1, len(res))
							got, err = PCA.QueryNodes(PCA.LastBlockHeight(), nodeTypes.QueryValidatorsParams{StakingStatus: 2, JailedStatus: 0, Blockchain: "", Page: 1, Limit: 1}) // staked
							assert.Nil(t, err)
							res = got.Result.([]nodeTypes.Validator)
							assert.Equal(t, 1, len(res))
							memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlockHeader)
							header := <-evtChan // Wait for header
							if len(header.Events["unstake.module"]) == 1 {
								got, err := PCA.QueryNodes(PCA.LastBlockHeight(), nodeTypes.QueryValidatorsParams{StakingStatus: 0, JailedStatus: 0, Blockchain: "", Page: 1, Limit: 1})
								assert.Nil(t, err)
								res := got.Result.([]nodeTypes.Validator)
								assert.Equal(t, 1, len(res))
								vals := got.Result.([]nodeTypes.Validator)
								addr := vals[0].Address
								balance, err := PCA.QueryBalance(addr.String(), PCA.LastBlockHeight())
								assert.Nil(t, err)
								assert.NotEqual(t, balance, sdk.ZeroInt())
								tx, err = nodes.StakeTx(memCodec(), memCli, kb, chains, "https://myPocketNode.com:8080", sdk.NewInt(10000000), kp, signer, "test", tc.upgrades.codecUpgrade.upgradeMod, false, signer)
								assert.Nil(t, err)
								assert.NotNil(t, tx)
								assert.Equal(t, tx.Code, uint32(0x0))
								cleanup()
								stopCli()

							}
							return
						}
					default:
						continue
					}
				}
			}
		})
	}

}
func TestStakeNode(t *testing.T) {
	tt := []struct {
		name           string
		memoryNodeFn   func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		outputIsSigner bool
		*upgrades
	}{
		{name: "stake node with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "stake a proto node with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
		{name: "stake a proto node with proto codec bad signer", memoryNodeFn: NewInMemoryTendermintNodeProto, outputIsSigner: true, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
		{name: "stake non-custodial before 8 upgrade ", memoryNodeFn: NewInMemoryTendermintNodeProto, outputIsSigner: false, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}, eight0Upgrade: upgrade{height: 999999}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			// 8.0 release
			isAfter8 := false
			if tc.eight0Upgrade.height != 0 {
				codec.UpgradeFeatureMap[codec.NonCustodialUpdateKey] = tc.eight0Upgrade.height
				isAfter8 = true
			}
			gen, vals := twoValTwoNodeGenesisState()
			_, kb, cleanup := tc.memoryNodeFn(t, gen)
			time.Sleep(1 * time.Second)
			kp, err := kb.GetCoinbase()
			signer := kp.GetAddress()
			if tc.outputIsSigner {
				for _, val := range vals {
					if val.Address.String() != signer.String() {
						signer = val.Address
						break
					}
				}
			}
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var tx *sdk.TxResponse
			var chains = []string{"0001"}
			<-evtChan // Wait for block
			memCli, stopCli, _ := subscribeTo(t, tmTypes.EventTx)
			tx, err = nodes.StakeTx(memCodec(), memCli, kb, chains, "https://myPocketNode.com:8080", sdk.NewInt(10000000), kp, signer, "test", tc.upgrades.codecUpgrade.upgradeMod, isAfter8, signer)
			assert.Nil(t, err)
			assert.NotNil(t, tx)
			if isAfter8 == false && tc.outputIsSigner {
				assert.Equal(t, 4, int(tx.Code))
			} else if isAfter8 == true {
				assert.Equal(t, 2, int(tx.Code))
			} else {
				assert.Equal(t, 0, int(tx.Code))
			}
			cleanup()
			stopCli()

		})
	}
}
func TestEditStakeNode(t *testing.T) {
	tt := []struct {
		name           string
		memoryNodeFn   func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		outputIsSigner bool
		*upgrades
	}{
		{name: "editStake after proto upgrade", memoryNodeFn: NewInMemoryTendermintNodeProto, outputIsSigner: false, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			codec.TestMode = 0
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			// 8.0 upgrade
			isAfter8 := false
			if tc.eight0Upgrade.height != 0 {
				codec.TestMode = -3
				codec.UpgradeFeatureMap[codec.NonCustodialUpdateKey] = tc.eight0Upgrade.height
				isAfter8 = true
			} else {
				codec.TestMode = -2
			}
			var newChains = []string{"2121"}
			var newServiceURL = "https://newServiceUrl.com:8081"
			gen, vals := twoValTwoNodeGenesisState()
			_, kb, cleanup := tc.memoryNodeFn(t, gen)
			time.Sleep(1 * time.Second)
			kp, err := kb.GetCoinbase()
			assert.Nil(t, err)
			signer := kp.GetAddress()
			assert.Nil(t, err)
			if tc.outputIsSigner {
				for _, val := range vals {
					if val.Address.String() == signer.String() {
						signer = val.OutputAddress
					}
				}
			}
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var tx *sdk.TxResponse
			<-evtChan // Wait for block
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			balance, err := PCA.QueryBalance(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			n, err := PCA.QueryNode(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			var newBalance = balance.Sub(sdk.NewInt(100000)).Add(n.StakedTokens)
			fmt.Println(signer.String())
			tx, err = nodes.StakeTx(memCodec(), memCli, kb, newChains, newServiceURL, newBalance, kp, signer, "test", tc.upgrades.codecUpgrade.upgradeMod, isAfter8, signer)
			assert.Nil(t, err)
			assert.NotNil(t, tx)
			<-evtChan // Wait for tx
			nodeUpdated, err := PCA.QueryNode(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			// assert not the same as the old node
			assert.NotEqual(t, nodeUpdated, n)
			// assert chains, serviceurl, and stake updated
			assert.Equal(t, newChains, nodeUpdated.Chains)
			// assert chains, serviceurl, and stake updated
			assert.Equal(t, newServiceURL, nodeUpdated.ServiceURL)
			// assert chains, serviceurl, and stake updated
			assert.Equal(t, newBalance, nodeUpdated.StakedTokens)
			cleanup()
			stopCli()
		})
	}
}

func TestEditStakeNodeOutput(t *testing.T) {
	tt := []struct {
		name           string
		memoryNodeFn   func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		outputIsSigner bool
		*upgrades
	}{
		{name: "editStake output update flow", memoryNodeFn: NewInMemoryTendermintNodeProto, outputIsSigner: true, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			codec.TestMode = 0
			if tc.upgrades != nil { // NOTE: Use to perform necessary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			codec.TestMode = -2
			var newChains = []string{"2121"}
			var newServiceURL = "https://newServiceUrl.com:8081"
			gen, _ := twoValTwoNodeGenesisState()
			_, kb, cleanup := tc.memoryNodeFn(t, gen)
			time.Sleep(1 * time.Second)
			kp, err := kb.GetCoinbase()
			assert.Nil(t, err)
			signer := kp.GetAddress()
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var tx *sdk.TxResponse
			<-evtChan // Wait for block
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			balance, err := PCA.QueryBalance(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			n, err := PCA.QueryNode(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			var newBalance = balance.Sub(sdk.NewInt(100000)).Add(n.StakedTokens)
			tx, err = nodes.StakeTx(memCodec(), memCli, kb, newChains, newServiceURL, newBalance, kp, signer, "test", tc.upgrades.codecUpgrade.upgradeMod, false, signer)
			assert.Nil(t, err)
			assert.NotNil(t, tx)
			<-evtChan // Wait for tx
			nodeUpdated, err := PCA.QueryNode(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			// assert not the same as the old node
			assert.NotEqual(t, nodeUpdated, n)
			// assert chains, serviceurl, and stake updated
			assert.Equal(t, newChains, nodeUpdated.Chains)
			// assert chains, serviceurl, and stake updated
			assert.Equal(t, newServiceURL, nodeUpdated.ServiceURL)
			// assert chains, serviceurl, and stake updated
			assert.Equal(t, newBalance, nodeUpdated.StakedTokens)

			// 8.0 upgrade
			codec.TestMode = -3
			codec.UpgradeFeatureMap[codec.NonCustodialUpdateKey] = tc.eight0Upgrade.height
			isAfter8 := true
			newBalance = nodeUpdated.StakedTokens.Add(sdk.NewInt(1000))
			tx, err = nodes.StakeTx(memCodec(), memCli, kb, newChains, newServiceURL, newBalance, kp, signer, "test", tc.upgrades.codecUpgrade.upgradeMod, isAfter8, signer)
			assert.Nil(t, err)
			assert.NotNil(t, tx)
			<-evtChan // Wait for tx
			nodeUpdated, err = PCA.QueryNode(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.Equal(t, newBalance, nodeUpdated.StakedTokens)
			assert.Equal(t, signer, nodeUpdated.OutputAddress)

			cleanup()
			stopCli()
		})
	}
}

func TestUnstakeNode8(t *testing.T) {
	tt := []struct {
		name           string
		memoryNodeFn   func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		outputIsSigner bool
		*upgrades
	}{
		{name: "unstake after 8.0 upgrade; output is signer", memoryNodeFn: NewInMemoryTendermintNodeProto, outputIsSigner: true, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}, eight0Upgrade: upgrade{3}}},
		{name: "unstake after 8.0 upgrade; output is not signer", memoryNodeFn: NewInMemoryTendermintNodeProto, outputIsSigner: false, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}, eight0Upgrade: upgrade{3}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			// 8.0 release
			isAfter8 := false
			if tc.eight0Upgrade.height != 0 {
				codec.TestMode = -3
				codec.UpgradeFeatureMap[codec.NonCustodialUpdateKey] = tc.eight0Upgrade.height
				isAfter8 = true
			}
			var chains = []string{"0001"}
			gen, vals := twoValTwoNodeGenesisState8()
			_, kb, cleanup := tc.memoryNodeFn(t, gen)
			time.Sleep(1 * time.Second)
			kp, err := kb.GetCoinbase()
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var tx *sdk.TxResponse
			<-evtChan // Wait for block
			memCli, _, evtChan := subscribeTo(t, tmTypes.EventTx)
			_, err = PCA.QueryBalance(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			signer := kp.GetAddress()
			output := signer
			if tc.outputIsSigner {
				for _, val := range vals {
					if val.Address.String() == signer.String() {
						signer = val.OutputAddress
						output = signer
					}
				}
			} else {
				for _, val := range vals {
					if val.Address.String() == signer.String() {
						signer = val.Address
						output = val.OutputAddress
					}
				}
			}
			// set output address
			tx, err = nodes.StakeTx(memCodec(), memCli, kb, chains, "https://myPocketNode.com:8080", sdk.NewInt(1000000000000000), kp, output, "test", tc.upgrades.codecUpgrade.upgradeMod, isAfter8, signer)
			assert.Nil(t, err)
			assert.NotNil(t, tx)
			<-evtChan // Wait for tx
			tx, err = nodes.UnstakeTx(memCodec(), memCli, kb, kp.GetAddress(), signer, "test", tc.upgrades.codecUpgrade.upgradeMod, isAfter8)
			assert.Nil(t, err)
			assert.NotNil(t, tx)
			<-evtChan // Wait for tx
			_, _, evtChan = subscribeTo(t, tmTypes.EventNewBlockHeader)
			for {
				select {
				case res := <-evtChan:
					if len(res.Events["begin_unstake.module"]) == 1 {
						got, err := PCA.QueryNodes(PCA.LastBlockHeight(), nodeTypes.QueryValidatorsParams{StakingStatus: 1, JailedStatus: 0, Blockchain: "", Page: 1, Limit: 1}) // unstaking
						assert.Nil(t, err)
						res := got.Result.([]nodeTypes.Validator)
						assert.Equal(t, 1, len(res))
						got, err = PCA.QueryNodes(PCA.LastBlockHeight(), nodeTypes.QueryValidatorsParams{StakingStatus: 2, JailedStatus: 0, Blockchain: "", Page: 1, Limit: 1}) // staked
						assert.Nil(t, err)
						res = got.Result.([]nodeTypes.Validator)
						assert.Equal(t, 1, len(res))
						memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlockHeader)
						header := <-evtChan // Wait for header
						if len(header.Events["unstake.module"]) == 1 {
							got, err := PCA.QueryNodes(PCA.LastBlockHeight(), nodeTypes.QueryValidatorsParams{StakingStatus: 0, JailedStatus: 0, Blockchain: "", Page: 1, Limit: 1})
							assert.Nil(t, err)
							res := got.Result.([]nodeTypes.Validator)
							assert.Equal(t, 1, len(res))
							vals := got.Result.([]nodeTypes.Validator)
							addr := vals[0].Address
							balance, err := PCA.QueryBalance(addr.String(), PCA.LastBlockHeight())
							assert.Nil(t, err)
							assert.NotEqual(t, balance, sdk.ZeroInt())
							tx, err = nodes.StakeTx(memCodec(), memCli, kb, chains, "https://myPocketNode.com:8080", sdk.NewInt(10000000), kp, signer, "test", tc.upgrades.codecUpgrade.upgradeMod, isAfter8, signer)
							assert.Nil(t, err)
							assert.NotNil(t, tx)
							assert.Equal(t, tx.Code, uint32(0x0))
							cleanup()
							stopCli()

						}
						return
					}
				default:
					continue
				}
			}
		})
	}

}
func TestStakeNodeAfter8(t *testing.T) {
	tt := []struct {
		name           string
		memoryNodeFn   func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		outputIsSigner bool
		*upgrades
	}{
		{name: "stake after 8.0 upgrade; output is signer", memoryNodeFn: NewInMemoryTendermintNodeProto, outputIsSigner: true, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}, eight0Upgrade: upgrade{height: 2}}},
		{name: "stake after 8.0 upgrade; output is not signer", memoryNodeFn: NewInMemoryTendermintNodeProto, outputIsSigner: false, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}, eight0Upgrade: upgrade{height: 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform necessary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			// 8.0 release
			isAfter8 := false
			if tc.eight0Upgrade.height != 0 {
				codec.UpgradeFeatureMap[codec.NonCustodialUpdateKey] = tc.eight0Upgrade.height
				isAfter8 = true
			}
			gen, vals := twoValTwoNodeGenesisState8()
			_, kb, cleanup := tc.memoryNodeFn(t, gen)
			time.Sleep(1 * time.Second)
			kp, err := kb.GetCoinbase()
			signer := kp.GetAddress()
			if tc.outputIsSigner {
				for _, val := range vals {
					if val.Address.String() == signer.String() {
						signer = val.OutputAddress
					}
				}
			}
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var tx *sdk.TxResponse
			var chains = []string{"0001"}
			<-evtChan // Wait for block
			memCli, stopCli, _ := subscribeTo(t, tmTypes.EventTx)
			tx, err = nodes.StakeTx(memCodec(), memCli, kb, chains, "https://myPocketNode.com:8080", sdk.NewInt(10000000), kp, signer, "test", tc.upgrades.codecUpgrade.upgradeMod, isAfter8, signer)
			assert.Nil(t, err)
			assert.NotNil(t, tx)
			if isAfter8 == false && tc.outputIsSigner {
				assert.Equal(t, 4, int(tx.Code))
			} else {
				assert.Equal(t, 0, int(tx.Code))
			}
			cleanup()
			stopCli()

		})
	}
}
func TestEditStakeNode8(t *testing.T) {
	tt := []struct {
		name           string
		memoryNodeFn   func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		outputIsSigner bool
		*upgrades
	}{
		{name: "editStake after 8.0 upgrade", memoryNodeFn: NewInMemoryTendermintNodeProto, outputIsSigner: true, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}, eight0Upgrade: upgrade{height: 3}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			codec.TestMode = 0
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			// 8.0 upgrade
			isAfter8 := false
			if tc.eight0Upgrade.height != 0 {
				codec.TestMode = -3
				codec.UpgradeFeatureMap[codec.NonCustodialUpdateKey] = tc.eight0Upgrade.height
				isAfter8 = true
			} else {
				codec.TestMode = -2
			}
			var newChains = []string{"2121"}
			var newServiceURL = "https://newServiceUrl.com:8081"
			gen, vals := twoValTwoNodeGenesisState8()
			_, kb, cleanup := tc.memoryNodeFn(t, gen)
			time.Sleep(1 * time.Second)
			kp, err := kb.GetCoinbase()
			assert.Nil(t, err)
			signer := kp.GetAddress()
			assert.Nil(t, err)
			if tc.outputIsSigner {
				for _, val := range vals {
					if val.Address.String() == signer.String() {
						signer = val.OutputAddress
					}
				}
			}
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var tx *sdk.TxResponse
			<-evtChan // Wait for block
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			balance, err := PCA.QueryBalance(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			n, err := PCA.QueryNode(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			var newBalance = balance.Sub(sdk.NewInt(100000)).Add(n.StakedTokens)
			fmt.Println(signer.String())
			tx, err = nodes.StakeTx(memCodec(), memCli, kb, newChains, newServiceURL, newBalance, kp, signer, "test", tc.upgrades.codecUpgrade.upgradeMod, isAfter8, signer)
			assert.Nil(t, err)
			assert.NotNil(t, tx)
			<-evtChan // Wait for tx
			nodeUpdated, err := PCA.QueryNode(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			// assert not the same as the old node
			assert.NotEqual(t, nodeUpdated, n)
			// assert chains, serviceurl, and stake updated
			assert.Equal(t, newChains, nodeUpdated.Chains)
			// assert chains, serviceurl, and stake updated
			assert.Equal(t, newServiceURL, nodeUpdated.ServiceURL)
			// assert chains, serviceurl, and stake updated
			assert.Equal(t, newBalance, nodeUpdated.StakedTokens)
			cleanup()
			stopCli()
		})
	}
}

func TestSendTransaction(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "send tx from an amino account with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "send tx from a proto account with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			time.Sleep(1 * time.Second)
			cb, err := kb.GetCoinbase()
			assert.Nil(t, err)
			kp, err := kb.Create("test")
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var transferAmount = sdk.NewInt(1000)
			var tx *sdk.TxResponse

			<-evtChan // Wait for block
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			tx, err = nodes.Send(memCodec(), memCli, kb, cb.GetAddress(), kp.GetAddress(), "test", transferAmount, tc.upgrades.codecUpgrade.upgradeMod)
			assert.Nil(t, err)
			assert.NotNil(t, tx)
			assert.Equal(t, int(tx.Code), 0)

			<-evtChan // Wait for tx
			balance, err := PCA.QueryBalance(kp.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.True(t, balance.Equal(transferAmount))
			balance, err = PCA.QueryBalance(cb.GetAddress().String(), PCA.LastBlockHeight())
			assert.Nil(t, err)

			cleanup()
			stopCli()
		})
	}
}

func TestDuplicateTxWithRawTx(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "send duplicate tx from an amino account with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "send duplicate tx from a proto account with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}}, // TODO:  FULL PROTO SCENARIO
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			time.Sleep(1 * time.Second)
			cb, err := kb.GetCoinbase()
			assert.Nil(t, err)
			kp, err := kb.Create("test")
			assert.Nil(t, err)
			pk, err := kb.ExportPrivateKeyObject(cb.GetAddress(), "test")
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			// create the transaction
			txBz, err := types.DefaultTxEncoder(memCodec())(types.NewTestTx(sdk.Context{}.WithChainID("pocket-test"),
				&nodeTypes.MsgSend{
					FromAddress: cb.GetAddress(),
					ToAddress:   kp.GetAddress(),
					Amount:      sdk.NewInt(1),
				},
				pk,
				rand2.Int64(),
				sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(100000)))), -1)
			assert.Nil(t, err)
			// create the transaction
			_, err = types.DefaultTxEncoder(memCodec())(types.NewTestTx(sdk.Context{}.WithChainID("pocket-test"),
				&nodeTypes.MsgSend{
					FromAddress: cb.GetAddress(),
					ToAddress:   kp.GetAddress(),
					Amount:      sdk.NewInt(1),
				},
				pk,
				rand2.Int64(),
				sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(100000)))), -1)
			assert.Nil(t, err)

			<-evtChan // Wait for block
			memCli, _, evtChan := subscribeTo(t, tmTypes.EventTx)
			_, err = nodes.RawTx(memCodec(), memCli, cb.GetAddress(), txBz)
			assert.Nil(t, err)
			// next tx
			<-evtChan // Wait for tx
			_, _, evtChan = subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for  block
			memCli, stopCli, _ := subscribeTo(t, tmTypes.EventTx)
			txResp, err := nodes.RawTx(memCodec(), memCli, cb.GetAddress(), txBz)
			if err == nil && txResp.Code == 0 {
				t.Fatal("should fail on replay attack")
			}
			cleanup()
			stopCli()
		})
	}

}

func TestChangeParamsComplexTypeTx(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "change complex type params from an amino account with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "change complex type params from a proto account with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}}, // TODO: FIX !!
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			resetTestACL()
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			time.Sleep(1 * time.Second)
			cb, err := kb.GetCoinbase()
			assert.Nil(t, err)
			kps, err := kb.List()
			assert.Nil(t, err)
			kp2 := kps[1]
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			a := testACL
			a.SetOwner("gov/acl", kp2.GetAddress())
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			tx, err := gov.ChangeParamsTx(memCodec(), memCli, kb, cb.GetAddress(), "gov/acl", a, "test", 1000000, false)
			assert.Nil(t, err)
			assert.NotNil(t, tx)
			select {
			case _ = <-evtChan:
				//fmt.Println(res)
				acl, err := PCA.QueryACL(PCA.LastBlockHeight())
				assert.Nil(t, err)
				o := acl.GetOwner("gov/acl")
				assert.Equal(t, kp2.GetAddress().String(), o.String())
				cleanup()
				stopCli()
			}
		})
	}
}

func TestChangeParamsSimpleTx(t *testing.T) {

	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "change complex type params from an amino account with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "change complex type params from a proto account with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}}, // TODO: FULL PROTO SCENARIO
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			resetTestACL()
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			time.Sleep(1 * time.Second)
			cb, err := kb.GetCoinbase()
			assert.Nil(t, err)
			_, err = kb.List()
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			tx, err := gov.ChangeParamsTx(memCodec(), memCli, kb, cb.GetAddress(), "application/StabilityAdjustment", 100, "test", 1000000, false)
			assert.Nil(t, err)
			assert.NotNil(t, tx)
			select {
			case _ = <-evtChan:
				//fmt.Println(res)
				assert.Nil(t, err)
				o, _ := PCA.QueryParam(PCA.LastBlockHeight(), "application/StabilityAdjustment")
				assert.Equal(t, "100", o.Value)
				cleanup()
				stopCli()
			}
		})
	}
}

// NOTE: This is a single long test to show the lifecycle of the query paramerter feature
func TestBlockSize_ChangeParams(t *testing.T) {
	blockSizeKey := "pocketcore/BlockByteSize"

	// TODO: Understand why we need this?
	codec.TestMode = -2

	// Prepare governance parameters
	codec.UpgradeHeight = 2                               // Upgrade the network at height 2
	codec.UpgradeFeatureMap[codec.BlockSizeModifyKey] = 3 // Enable the feature at height 4
	_ = memCodecMod(true)
	resetTestACL()

	// Prepare the test network
	_, kb, cleanup := NewInMemoryTendermintNodeProto(t, oneAppTwoNodeGenesis())
	time.Sleep(1 * time.Second)

	// Get the address of the ACL owner (i.e. the DAO)
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	daoAddr := cb.GetAddress()

	// Subscribe to new events
	_, _, newBlockEventChan := subscribeTo(t, tmTypes.EventNewBlock)
	memCli, stopCli, txEventChan := subscribeTo(t, tmTypes.EventTx)

	// Wait for block 2
	<-newBlockEventChan

	// Before activation, the parameter does not exist and should be an empty string
	queryRes, err := PCA.QueryParam(PCA.LastBlockHeight(), blockSizeKey)
	assert.Nil(t, err)
	assert.Equal(t, "", queryRes.Value)

	// Wait for block 3
	<-newBlockEventChan

	// After activation, the parameter should be the default value
	queryRes, err = PCA.QueryParam(PCA.LastBlockHeight(), blockSizeKey)
	assert.Nil(t, err)
	assert.Equal(t, "4000000", queryRes.Value)

	// Changing the parameter after activation
	tx, err := gov.ChangeParamsTx(memCodec(), memCli, kb, daoAddr, blockSizeKey, "9000000", "test", 10000, false)
	assert.Nil(t, err)
	assert.NotNil(t, tx)

	// Parameter still does not exist and should be an empty string
	<-txEventChan
	queryRes, err = PCA.QueryParam(PCA.LastBlockHeight(), blockSizeKey)
	assert.Nil(t, err)
	assert.Equal(t, "9000000", queryRes.Value)

	// Wait for block 4
	<-newBlockEventChan

	// Verify the parameter maintains the new value
	queryRes, err = PCA.QueryParam(PCA.LastBlockHeight(), blockSizeKey)
	assert.Nil(t, err)
	assert.Equal(t, "9000000", queryRes.Value)

	// Cleanup
	cleanup()
	stopCli()
}

// NOTE: This is a hacky test prone to race conditions and global variable management to account for lack of
// fine-grained control of Tendermint's consensus engine (which runs in the background).
//
// How it works:
// 1. Configure the network & initialize block size updates
// 2. Start a background process that keeps sending transactions
// 3. Start a background process that doubles the block size every block
// 4. Repeat (2) & (3) N times
//
// VERIFICATION: The assertions in the test do minimal verification but evaluating the logs provides an additional guarantee.
// 1. Verify that smaller block sizes cannot be completely filled
// 2. Verify a sufficiently large block includes the backlog of transactions
// 3. Verify that at steady state (i.e. block is large enough to include all transactions), the number of txs in each block is the same
func TestBlockSize_MaximumSize(t *testing.T) {
	blockSizeKey := "pocketcore/BlockByteSize"

	// Prepare network configs
	codec.TestMode = -3                                   // Includes codec upgrade, validator split and non-custodial upgrade
	codec.UpgradeHeight = 2                               // Height at which codec was upgraded from amino to proto
	codec.UpgradeFeatureMap[codec.BlockSizeModifyKey] = 3 // Height at which to enable block size upgrades
	_ = memCodecMod(true)
	resetTestACL()

	// Pick a small initial genesis block size
	blockParamsMaxBytes := int64(2000)

	// Get the current global values used for testing
	globalTendermintTimeoutCommit := tendermintTimeoutCommit

	// On cleanup, revert global values so they are the same everywhere else
	t.Cleanup(func() {
		tendermintTimeoutCommit = globalTendermintTimeoutCommit
	})

	// Increase the Tendermint timeout commit: the min amount of time to wait between blocks.
	// Needs to be increased so there's enough time to send & accumulate transactions that fill the block.
	tendermintTimeoutCommit = time.Duration(5) * time.Second

	// Prepare the test network
	_, kb, cleanup := NewInMemoryTendermintNodeProto(t, oneAppTwoNodeGenesis())
	defer cleanup()

	// Wait for the first block
	_, _, newBlockEventChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-newBlockEventChan // block 1
	<-newBlockEventChan // block 2
	<-newBlockEventChan // block 3

	// Verify the block param is activated by comparing against the default genesis value
	queryRes, err := PCA.QueryParam(PCA.LastBlockHeight(), blockSizeKey)
	assert.Nil(t, err)
	assert.Equal(t, strconv.Itoa(int(pocketTypes.DefaultBlockByteSize)), queryRes.Value)

	// Get the address of the ACL owner (i.e. the DAO)
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)

	// Set the block size to a small initial value
	tx, err := gov.ChangeParamsTx(memCodec(), getInMemoryTMClient(), kb, cb.GetAddress(), blockSizeKey, fmt.Sprintf("%d", blockParamsMaxBytes), "test", 10000, false)
	assert.Nil(t, err)
	assert.NotNil(t, tx)

	<-newBlockEventChan // block 4

	// Verify the block size has been reduced
	queryRes, err = PCA.QueryParam(PCA.LastBlockHeight(), blockSizeKey)
	assert.Nil(t, err)
	assert.Equal(t, strconv.Itoa(int(blockParamsMaxBytes)), queryRes.Value)

	numBlockSizeDoubles := 15
	wg := sync.WaitGroup{}
	wg.Add(numBlockSizeDoubles) // 2kb, 4kb and 8kb, ...

	// Start a background process that sends as many transactions as possible
	txCtx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Create a new random keypair
				kp, err := kb.Create("test")
				assert.Nil(t, err)
				// Send funds from the DAO to the new address
				_, err = nodes.Send(memCodec(), getInMemoryTMClient(), kb, cb.GetAddress(), kp.GetAddress(), "test", sdk.NewInt(10000), false)
				assert.Nil(t, err)
			}
		}
	}(txCtx)
	defer cancel()

	// New blocks are created every `tendermintTimeoutCommit`.
	// Need to asynchronously capture new blocks, while txs are being sent above, and process appropriately
	newBlockCtx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			// _, _, newBlockEventChan = subscribeTo(t, tmTypes.EventNewBlock)
			select {
			case <-ctx.Done():
				return
			default:
				<-newBlockEventChan // Wait for block

				// Grab all the latest transactions from the block
				height := PCA.LastBlockHeight()
				res, err := PCA.QueryAllBlockTxs(height, 0, 1000) // we do not expect more than 1000 transactions to accumulate, hence one page is enough
				assert.Nil(t, err)

				// Count the size of all the transactions
				// NB: `CreateProposalBlock` in `tendermint/state/execution.go` can be seen that not the entire block
				// will be filled with transactions.
				total := 0
				for _, tx := range res.Txs {
					total += len(tx.Tx)
				}

				// Current max value
				queryRes, err := PCA.QueryParam(height, blockSizeKey)
				assert.Nil(t, err)
				t.Logf("Transactions accumulated at height=%d: totalTxBytes=%d, len(res.Txs)=%d, TotalTxCount=%d, blockSize=%s", height, total, len(res.Txs), res.TotalCount, queryRes.Value)

				// Double the block size
				blockParamsMaxBytes *= 2
				tx, err := gov.ChangeParamsTx(memCodec(), getInMemoryTMClient(), kb, cb.GetAddress(), blockSizeKey, fmt.Sprintf("%d", blockParamsMaxBytes), "test", 10000, false)
				assert.Nil(t, err)
				assert.NotNil(t, tx)

				wg.Done()
			}
		}
	}(newBlockCtx)
	defer cancel()

	// Wait for numBlockSizeDoubles to finish
	wg.Wait()
}

func TestChangepip22ParamsBeforeActivationHeight(t *testing.T) {

	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "change pip22 parameter before activation height", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}}, // TODO: FULL PROTO SCENARIO
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			codec.TestMode = -2
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			resetTestACL()
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			time.Sleep(1 * time.Second)
			cb, err := kb.GetCoinbase()
			assert.Nil(t, err)
			_, err = kb.List()
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			//Before Activation of the parameter ACL do not exist and the value and parameter should be 0 or nil
			firstquery, _ := PCA.QueryParam(PCA.LastBlockHeight(), "pos/ServicerStakeWeightMultiplier")
			assert.Equal(t, "", firstquery.Value)
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			//Tx wont modify anything as ACL is not configured (Txresult should be gov code 5)
			tx, err := gov.ChangeParamsTx(memCodec(), memCli, kb, cb.GetAddress(), "pos/ServicerStakeWeightMultiplier", 1, "test", 10000, false)
			assert.Nil(t, err)
			assert.NotNil(t, tx)
			select {
			case _ = <-evtChan:
				//fmt.Println(res)
				assert.Nil(t, err)
				o, _ := PCA.QueryParam(PCA.LastBlockHeight(), "pos/ServicerStakeWeightMultiplier")
				//value should be equal to the first query of the param
				assert.Equal(t, firstquery.Value, o.Value)
				cleanup()
				stopCli()
			}
		})
	}
}

func TestChangeParamspip22afterActivationHeight(t *testing.T) {

	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "change pip22 parameter past activation height", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}}, // TODO: FULL PROTO SCENARIO
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			codec.TestMode = -3
			codec.UpgradeFeatureMap[codec.RSCALKey] = tc.upgrades.codecUpgrade.height + 1
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			resetTestACL()
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			time.Sleep(1 * time.Second)
			cb, err := kb.GetCoinbase()
			assert.Nil(t, err)
			_, err = kb.List()
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			<-evtChan // Wait for block
			<-evtChan // Wait for another block
			//After Activation of the parameter ACL should be created(allowing modifying the value) and parameter should have default value of 4000000
			o, _ := PCA.QueryParam(PCA.LastBlockHeight(), "pos/ServicerStakeFloorMultiplier")
			assert.Equal(t, "15000000000", o.Value)
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			tx, err := gov.ChangeParamsTx(memCodec(), memCli, kb, cb.GetAddress(), "pos/ServicerStakeFloorMultiplier", 16000000000, "test", 10000, false)
			assert.Nil(t, err)
			assert.NotNil(t, tx)
			select {
			case _ = <-evtChan:
				//fmt.Println(res)
				assert.Nil(t, err)
				o, _ := PCA.QueryParam(PCA.LastBlockHeight(), "pos/ServicerStakeFloorMultiplier")
				assert.Equal(t, "16000000000", o.Value)
				cleanup()
				stopCli()
			}
		})
	}
}

func TestUpgrade(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "change complex type params from an amino account with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "change complex type params from a proto account with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			time.Sleep(1 * time.Second)
			cb, err := kb.GetCoinbase()
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var tx *sdk.TxResponse
			<-evtChan // Wait for block
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			tx, err = gov.UpgradeTx(memCodec(), memCli, kb, cb.GetAddress(), govTypes.Upgrade{
				Height:  1000,
				Version: "2.0.0",
			}, "test", 1000000, tc.codecUpgrade.upgradeMod)
			assert.Nil(t, err)
			assert.NotNil(t, tx)

			<-evtChan // Wait for tx
			u, err := PCA.QueryUpgrade(PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.True(t, u.UpgradeVersion() == "2.0.0")

			cleanup()
			stopCli()
		})
	}
}

func TestDAOTransfer(t *testing.T) {
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "change complex type params from an amino account with amino codec", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "change complex type params from a proto account with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			_, kb, cleanup := tc.memoryNodeFn(t, oneAppTwoNodeGenesis())
			time.Sleep(1 * time.Second)
			cb, err := kb.GetCoinbase()
			assert.Nil(t, err)
			_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
			var tx *sdk.TxResponse
			<-evtChan // Wait for block
			memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			tx, err = gov.DAOTransferTx(memCodec(), memCli, kb, cb.GetAddress(), nil, sdk.OneInt(), govTypes.DAOBurn.String(), "test", 1000000, tc.codecUpgrade.upgradeMod)
			assert.Nil(t, err)
			assert.NotNil(t, tx)

			<-evtChan // Wait for tx
			balance, err := PCA.QueryDaoBalance(PCA.LastBlockHeight())
			assert.Nil(t, err)
			assert.True(t, balance.Equal(sdk.NewInt(999)))

			cleanup()
			stopCli()
		})
	}
}

func TestClaimAminoTx(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "claim tx from amino with amino codec ", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		//{name: "claim tx from a proto with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 4}}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			genBz, _, validators, app := fiveValidatorsOneAppGenesis()
			kb := getInMemoryKeybase()
			_, _, cleanup := tc.memoryNodeFn(t, genBz)
			time.Sleep(1 * time.Second)
			for i := 0; i < 5; i++ {
				appPrivateKey, err := kb.ExportPrivateKeyObject(app.Address, "test")
				assert.Nil(t, err)
				// setup AAT
				aat := pocketTypes.AAT{
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
				proof := pocketTypes.RelayProof{
					Entropy:            int64(rand.Int()),
					RequestHash:        hex.EncodeToString(pocketTypes.Hash([]byte("fake"))),
					SessionBlockHeight: 1,
					ServicerPubKey:     validators[0].PublicKey.RawString(),
					Blockchain:         sdk.PlaceholderHash,
					Token:              aat,
					Signature:          "",
				}
				sig, err = appPrivateKey.Sign(proof.Hash())
				if err != nil {
					t.Fatal(err)
				}
				proof.Signature = hex.EncodeToString(sig)
				pocketTypes.SetProof(pocketTypes.SessionHeader{
					ApplicationPubKey:  appPrivateKey.PublicKey().RawString(),
					Chain:              sdk.PlaceholderHash,
					SessionBlockHeight: 1,
				}, pocketTypes.RelayEvidence, proof, sdk.NewInt(1000000), pocketTypes.GlobalEvidenceCache)
				assert.Nil(t, err)
			}
			_, _, evtChan := subscribeTo(t, tmTypes.EventTx)
			res := <-evtChan
			if res.Events["message.action"][0] != pocketTypes.EventTypeClaim {
				t.Fatal("claim message was not received first")
			}
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			res = <-evtChan
			if res.Events["message.action"][0] != pocketTypes.EventTypeProof {
				t.Fatal("proof message was not received afterward")
			}
			cleanup()
			stopCli()
		})
	}
}

func TestClaimProtoTx(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		//{name: "claim tx from amino with amino codec ", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "claim tx from a proto with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 5}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.upgrades != nil { // NOTE: Use to perform neccesary upgrades for test
				codec.UpgradeHeight = tc.upgrades.codecUpgrade.height
				_ = memCodecMod(tc.upgrades.codecUpgrade.upgradeMod)
			}
			genBz, _, validators, app := fiveValidatorsOneAppGenesis()
			kb := getInMemoryKeybase()
			_, _, cleanup := tc.memoryNodeFn(t, genBz)
			time.Sleep(1 * time.Second)
			for i := 0; i < 5; i++ {
				appPrivateKey, err := kb.ExportPrivateKeyObject(app.Address, "test")
				assert.Nil(t, err)
				// setup AAT
				aat := pocketTypes.AAT{
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
				proof := pocketTypes.RelayProof{
					Entropy:            int64(rand.Int()),
					RequestHash:        hex.EncodeToString(pocketTypes.Hash([]byte("fake"))),
					SessionBlockHeight: 1,
					ServicerPubKey:     validators[0].PublicKey.RawString(),
					Blockchain:         sdk.PlaceholderHash,
					Token:              aat,
					Signature:          "",
				}
				sig, err = appPrivateKey.Sign(proof.Hash())
				if err != nil {
					t.Fatal(err)
				}
				proof.Signature = hex.EncodeToString(sig)
				pocketTypes.SetProof(pocketTypes.SessionHeader{
					ApplicationPubKey:  appPrivateKey.PublicKey().RawString(),
					Chain:              sdk.PlaceholderHash,
					SessionBlockHeight: 1,
				}, pocketTypes.RelayEvidence, proof, sdk.NewInt(1000000), pocketTypes.GlobalEvidenceCache)
				assert.Nil(t, err)
			}
			_, _, evtChan := subscribeTo(t, tmTypes.EventTx)
			res := <-evtChan
			fmt.Println(res)
			if res.Events["message.action"][0] != pocketTypes.EventTypeClaim {
				t.Fatal("claim message was not received first")
			}
			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			res = <-evtChan
			if res.Events["message.action"][0] != pocketTypes.EventTypeProof {
				t.Fatal("proof message was not received afterward")
			}
			cleanup()
			stopCli()
		})
	}
}

func TestAminoClaimTxChallenge(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		{name: "challenge a claim tx from amino with amino codec ", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		//{name: "challenge a claim tx from a proto with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}}, // TODO: FULL PROT SCENARIO
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			genBz, keys, _, _ := fiveValidatorsOneAppGenesis()
			challenges := NewValidChallengeProof(t, keys, 5)
			_, _, cleanup := tc.memoryNodeFn(t, genBz)
			for _, c := range challenges {
				c.Store(sdk.NewInt(1000000), pocketTypes.GlobalEvidenceCache)
			}
			_, _, evtChan := subscribeTo(t, tmTypes.EventTx)
			res := <-evtChan // Wait for tx
			if res.Events["message.action"][0] != pocketTypes.EventTypeClaim {
				t.Fatal("claim message was not received first")
			}

			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			res = <-evtChan // Wait for tx
			if res.Events["message.action"][0] != pocketTypes.EventTypeProof {
				t.Fatal("proof message was not received afterward")
			}
			cleanup()
			stopCli()
		})
	}
}

func TestProtoClaimTxChallenge(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	tt := []struct {
		name         string
		memoryNodeFn func(t *testing.T, genesisState []byte) (tendermint *node.Node, keybase keys.Keybase, cleanup func())
		*upgrades
	}{
		//{name: "challenge a claim tx from amino with amino codec ", memoryNodeFn: NewInMemoryTendermintNodeAmino, upgrades: &upgrades{codecUpgrade: codecUpgrade{false, 7000}}},
		{name: "challenge a claim tx from a proto with proto codec", memoryNodeFn: NewInMemoryTendermintNodeProto, upgrades: &upgrades{codecUpgrade: codecUpgrade{true, 2}}}, // TODO: FULL PROT SCENARIO
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			genBz, keys, _, _ := fiveValidatorsOneAppGenesis()
			challenges := NewValidChallengeProof(t, keys, 5)
			_, _, cleanup := tc.memoryNodeFn(t, genBz)
			for _, c := range challenges {
				c.Store(sdk.NewInt(1000000), pocketTypes.GlobalEvidenceCache)
			}
			_, _, evtChan := subscribeTo(t, tmTypes.EventTx)
			res := <-evtChan // Wait for tx
			if res.Events["message.action"][0] != pocketTypes.EventTypeClaim {
				t.Fatal("claim message was not received first")
			}

			_, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
			res = <-evtChan // Wait for tx
			if res.Events["message.action"][0] != pocketTypes.EventTypeProof {
				t.Fatal("proof message was not received afterward")
			}
			cleanup()
			stopCli()
		})
	}
}

func NewValidChallengeProof(t *testing.T, privateKeys []crypto.PrivateKey, numOfChallenges int) (challenge []pocketTypes.ChallengeProofInvalidData) {

	appPrivateKey := privateKeys[1]
	servicerPrivKey1 := privateKeys[4]
	servicerPrivKey2 := privateKeys[2]
	servicerPrivKey3 := privateKeys[3]
	clientPrivateKey := servicerPrivKey3
	appPubKey := appPrivateKey.PublicKey().RawString()
	servicerPubKey := servicerPrivKey1.PublicKey().RawString()
	servicerPubKey2 := servicerPrivKey2.PublicKey().RawString()
	servicerPubKey3 := servicerPrivKey3.PublicKey().RawString()
	reporterPrivKey := privateKeys[0]
	reporterPubKey := reporterPrivKey.PublicKey()
	reporterAddr := reporterPubKey.Address()
	clientPubKey := clientPrivateKey.PublicKey().RawString()
	var proofs []pocketTypes.ChallengeProofInvalidData
	for i := 0; i < numOfChallenges; i++ {
		validProof := pocketTypes.RelayProof{
			Entropy:            int64(rand.Intn(500000)),
			SessionBlockHeight: 1,
			ServicerPubKey:     servicerPubKey,
			RequestHash:        clientPubKey, // fake
			Blockchain:         sdk.PlaceholderHash,
			Token: pocketTypes.AAT{
				Version:              "0.0.1",
				ApplicationPublicKey: appPubKey,
				ClientPublicKey:      clientPubKey,
				ApplicationSignature: "",
			},
			Signature: "",
		}
		appSignature, er := appPrivateKey.Sign(validProof.Token.Hash())
		if er != nil {
			t.Fatalf(er.Error())
		}
		validProof.Token.ApplicationSignature = hex.EncodeToString(appSignature)
		clientSignature, er := clientPrivateKey.Sign(validProof.Hash())
		if er != nil {
			t.Fatalf(er.Error())
		}
		validProof.Signature = hex.EncodeToString(clientSignature)
		// valid proof 2
		validProof2 := pocketTypes.RelayProof{
			Entropy:            0,
			SessionBlockHeight: 1,
			ServicerPubKey:     servicerPubKey2,
			RequestHash:        clientPubKey, // fake
			Blockchain:         sdk.PlaceholderHash,
			Token: pocketTypes.AAT{
				Version:              "0.0.1",
				ApplicationPublicKey: appPubKey,
				ClientPublicKey:      clientPubKey,
				ApplicationSignature: "",
			},
			Signature: "",
		}
		appSignature, er = appPrivateKey.Sign(validProof2.Token.Hash())
		if er != nil {
			t.Fatalf(er.Error())
		}
		validProof2.Token.ApplicationSignature = hex.EncodeToString(appSignature)
		clientSignature, er = clientPrivateKey.Sign(validProof2.Hash())
		if er != nil {
			t.Fatalf(er.Error())
		}
		validProof2.Signature = hex.EncodeToString(clientSignature)
		// valid proof 3
		validProof3 := pocketTypes.RelayProof{
			Entropy:            0,
			SessionBlockHeight: 1,
			ServicerPubKey:     servicerPubKey3,
			RequestHash:        clientPubKey, // fake
			Blockchain:         sdk.PlaceholderHash,
			Token: pocketTypes.AAT{
				Version:              "0.0.1",
				ApplicationPublicKey: appPubKey,
				ClientPublicKey:      clientPubKey,
				ApplicationSignature: "",
			},
			Signature: "",
		}
		appSignature, er = appPrivateKey.Sign(validProof3.Token.Hash())
		if er != nil {
			t.Fatalf(er.Error())
		}
		validProof3.Token.ApplicationSignature = hex.EncodeToString(appSignature)
		clientSignature, er = clientPrivateKey.Sign(validProof3.Hash())
		if er != nil {
			t.Fatalf(er.Error())
		}
		validProof3.Signature = hex.EncodeToString(clientSignature)
		// create responses
		majorityResponsePayload := `{"id":67,"jsonrpc":"2.0","result":"Mist/v0.9.3/darwin/go1.4.1"}`
		minorityResponsePayload := `{"id":67,"jsonrpc":"2.0","result":"Mist/v0.9.3/darwin/go1.4.2"}`
		// majority response 1
		majResp1 := pocketTypes.RelayResponse{
			Signature: "",
			Response:  majorityResponsePayload,
			Proof:     validProof,
		}
		sig, er := servicerPrivKey1.Sign(majResp1.Hash())
		if er != nil {
			t.Fatalf(er.Error())
		}
		majResp1.Signature = hex.EncodeToString(sig)
		// majority response 2
		majResp2 := pocketTypes.RelayResponse{
			Signature: "",
			Response:  majorityResponsePayload,
			Proof:     validProof2,
		}
		sig, er = servicerPrivKey2.Sign(majResp2.Hash())
		if er != nil {
			t.Fatalf(er.Error())
		}
		majResp2.Signature = hex.EncodeToString(sig)
		// minority response
		minResp := pocketTypes.RelayResponse{
			Signature: "",
			Response:  minorityResponsePayload,
			Proof:     validProof3,
		}
		sig, er = servicerPrivKey3.Sign(minResp.Hash())
		if er != nil {
			t.Fatalf(er.Error())
		}
		minResp.Signature = hex.EncodeToString(sig)
		// create valid challenge proof
		proofs = append(proofs, pocketTypes.ChallengeProofInvalidData{
			MajorityResponses: []pocketTypes.RelayResponse{
				majResp1,
				majResp2,
			},
			MinorityResponse: minResp,
			ReporterAddress:  sdk.Address(reporterAddr),
		})
	}
	return proofs
}
