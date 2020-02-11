package app

import (
	"encoding/hex"
	"fmt"
	apps "github.com/pokt-network/pocket-core/x/apps"
	"github.com/pokt-network/pocket-core/x/nodes"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/types"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/stretchr/testify/assert"
	tmTypes "github.com/tendermint/tendermint/types"
	"math/rand"
	"strings"
	"testing"
)

func TestUnstakeApp(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	kp, err := kb.GetCoinbase()
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	var tx *sdk.TxResponse
	var chains = []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}
	select {
	case <-evtChan:
		var err error
		memCli, stopCli, evtChan = subscribeTo(t, tmTypes.EventTx)
		tx, err = apps.StakeTx(memCodec(), memCli, kb, chains, sdk.NewInt(1000000), kp, "test")
		assert.Nil(t, err)
		assert.NotNil(t, tx)
	}
	select {
	case <-evtChan:
		got, err := apps.QueryApplications(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(got))
		memCli, stopCli, evtChan = subscribeTo(t, tmTypes.EventTx)
		tx, err = apps.UnstakeTx(memCodec(), memCli, kb, kp.GetAddress(), "test")
	}
	select {
	case <-evtChan:
		got, err := apps.QueryUnstakingApplications(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(got))
		got, err = apps.QueryStakedApplications(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(got))
	}
	cleanup()
	stopCli()
}

func TestUnstakeNode(t *testing.T) {
	var chains = []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}
	_, kb, cleanup := NewInMemoryTendermintNode(t, twoValTwoNodeGenesisState())
	kp, err := kb.GetCoinbase()
	assert.Nil(t, err)
	var balance1 sdk.Int
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	var tx *sdk.TxResponse
	select {
	case <-evtChan:
		var err error
		memCli, stopCli, evtChan = subscribeTo(t, tmTypes.EventTx)
		balance1, err = nodes.QueryAccountBalance(memCodec(), memCli, kp.GetAddress(), 0)
		tx, err = nodes.UnstakeTx(memCodec(), memCli, kb, kp.GetAddress(), "test")
		assert.Nil(t, err)
		assert.NotNil(t, tx)
	}
	select {
	case <-evtChan:
		memCli, stopCli, evtChan = subscribeTo(t, tmTypes.EventNewBlockHeader)
		for {
			select {
			case res := <-evtChan:
				if len(res.Events["begin_unstake.module"]) == 1 {
					got, err := nodes.QueryUnstakingValidators(memCodec(), memCli, 0)
					assert.Nil(t, err)
					assert.Equal(t, 1, len(got))
					got, err = nodes.QueryStakedValidators(memCodec(), memCli, 0)
					assert.Nil(t, err)
					assert.Equal(t, 1, len(got))
					memCli, stopCli, evtChan = subscribeTo(t, tmTypes.EventNewBlockHeader)
					select {
					case res := <-evtChan:
						if len(res.Events["unstake.module"]) == 1 {
							got, err := nodes.QueryUnstakedValidators(memCodec(), memCli, 0)
							assert.Nil(t, err)
							assert.Equal(t, 1, len(got))
							assert.Equal(t, got[0].StakedTokens.Int64(), int64(0))
							addr := got[0].Address
							balance, err := nodes.QueryAccountBalance(memCodec(), memCli, addr, 0)
							fmt.Println(balance1, balance)
							assert.NotZero(t, balance.Int64())
							tx, err = nodes.StakeTx(memCodec(), memCli, kb, chains, "https://myPocketNode:8080", sdk.NewInt(10000000), kp, "test")
							assert.Nil(t, err)
							assert.NotNil(t, tx)
							assert.True(t, strings.Contains(tx.Logs.String(), `"success":true`))
							cleanup()
							stopCli()
						}
					}
					return
				}
			default:
				continue
			}
		}
	}
}

func TestStakeNode(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, twoValTwoNodeGenesisState())
	kp, err := kb.GetCoinbase()
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	var tx *sdk.TxResponse
	var chains = []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}
	select {
	case <-evtChan:
		var err error
		memCli, stopCli, evtChan = subscribeTo(t, tmTypes.EventTx)
		tx, err = nodes.StakeTx(memCodec(), memCli, kb, chains, "https://myPocketNode:8080", sdk.NewInt(10000000), kp, "test")
		assert.Nil(t, err)
		assert.NotNil(t, tx)
		assert.True(t, strings.Contains(tx.Logs.String(), `"success":true`))
		cleanup()
		stopCli()
	}
}

func TestStakeApp(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	kp, err := kb.GetCoinbase()
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	var tx *sdk.TxResponse
	var chains = []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}

	select {
	case <-evtChan:
		var err error
		memCli, stopCli, evtChan = subscribeTo(t, tmTypes.EventTx)
		tx, err = apps.StakeTx(memCodec(), memCli, kb, chains, sdk.NewInt(1000000), kp, "test")
		assert.Nil(t, err)
		assert.NotNil(t, tx)
	}
	select {
	case <-evtChan:
		got, err := apps.QueryApplications(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(got))
	}
	stopCli()
	cleanup()
}

func TestSendTransaction(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	kp, err := kb.Create("test")
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	var transferAmount = sdk.NewInt(1000)
	var tx *sdk.TxResponse
	select {
	case <-evtChan:
		var err error
		memCli, stopCli, evtChan = subscribeTo(t, tmTypes.EventTx)
		tx, err = nodes.Send(memCodec(), memCli, kb, cb.GetAddress(), kp.GetAddress(), "test", transferAmount)
		assert.Nil(t, err)
		assert.NotNil(t, tx)
	}
	select {
	case <-evtChan:
		balance, err := nodes.QueryAccountBalance(memCodec(), memCli, kp.GetAddress(), 0)
		assert.Nil(t, err)
		assert.True(t, balance.Equal(transferAmount))
		balance, err = nodes.QueryAccountBalance(memCodec(), memCli, cb.GetAddress(), 0)
		assert.Nil(t, err)
	}
	cleanup()
	stopCli()
}

func TestSendRawTx(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	kp, err := kb.Create("test")
	assert.Nil(t, err)
	pk, err := kb.ExportPrivateKeyObject(cb.GetAddress(), "test")
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	// create the transaction
	txBz, err := types.DefaultTxEncoder(memCodec())(types.NewTestTx(sdk.Context{}.WithChainID("pocket-test"),
		[]sdk.Msg{bank.MsgSend{
			FromAddress: cb.GetAddress(),
			ToAddress:   kp.GetAddress(),
			Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1))),
		}},
		[]crypto.PrivateKey{pk},
		[]uint64{0},
		[]uint64{0},
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(100000)))))
	assert.Nil(t, err)
	select {
	case <-evtChan:
		memCli, stopCli, evtChan = subscribeTo(t, tmTypes.EventTx)
		var err error
		txResp, err := nodes.RawTx(memCodec(), memCli, cb.GetAddress(), txBz)
		fmt.Println(txResp.Logs)
		assert.Nil(t, err)
		assert.NotNil(t, txResp)
	}
	// next block
	select {
	case <-evtChan:
		res, err := nodes.QueryAccountBalance(memCodec(), memCli, cb.GetAddress(), 0)
		assert.Nil(t, err)
		assert.Equal(t, int64(999899999), res.Int64())
	}
	cleanup()
	stopCli()
}

func TestClaimTx(t *testing.T) {
	genBz, validators, app := fiveValidatorsOneAppGenesis()
	kb := getInMemoryKeybase()
	for i := 0; i < 8; i++ {
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
			SessionBlockHeight: 1,
			ServicerPubKey:     validators[0].PublicKey.RawString(),
			Blockchain:         dummyChainsHash,
			Token:              aat,
			Signature:          "",
		}
		sig, err = appPrivateKey.Sign(proof.Hash())
		if err != nil {
			t.Fatal(err)
		}
		proof.Signature = hex.EncodeToString(sig)
		err = pocketTypes.GetAllInvoices().AddToInvoice(pocketTypes.SessionHeader{
			ApplicationPubKey:  appPrivateKey.PublicKey().RawString(),
			Chain:              dummyChainsHash,
			SessionBlockHeight: 1,
		}, proof)
		assert.Nil(t, err)
	}
	_, _, cleanup := NewInMemoryTendermintNode(t, genBz)
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
	select {
	case res := <-evtChan:
		if res.Events["message.action"][0] != pocketTypes.EventTypeClaim {
			t.Fatal("claim message was not received first")
		}
		_, stopCli, evtChan = subscribeTo(t, tmTypes.EventTx)
		select {
		case res := <-evtChan:
			if res.Events["message.action"][0] != pocketTypes.EventTypeProof {
				t.Fatal("proof message was not received afterward")
			}
			cleanup()
			stopCli()
		}
	}
}
