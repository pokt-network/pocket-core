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
	"github.com/tendermint/tendermint/libs/common"
	tmTypes "github.com/tendermint/tendermint/types"
	db "github.com/tendermint/tm-db"
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

func TestDuplicateTxWithRawTx(t *testing.T) {
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
		common.RandInt64(),
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(100000)))))
	assert.Nil(t, err)
	// create the transaction
	txBz2, err := types.DefaultTxEncoder(memCodec())(types.NewTestTx(sdk.Context{}.WithChainID("pocket-test"),
		[]sdk.Msg{bank.MsgSend{
			FromAddress: cb.GetAddress(),
			ToAddress:   kp.GetAddress(),
			Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1))),
		}},
		[]crypto.PrivateKey{pk},
		common.RandInt64(),
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(100000)))))
	txBz2 = txBz2
	assert.Nil(t, err)
	select {
	case <-evtChan:
		memCli, stopCli, evtChan = subscribeTo(t, tmTypes.EventTx)
		var err error
		_, err = nodes.RawTx(memCodec(), memCli, cb.GetAddress(), txBz)
		assert.Nil(t, err)
	}
	// next tx
	select {
	case <-evtChan:
		memCli, stopCli, evtChan = subscribeTo(t, tmTypes.EventNewBlock)
		select {
		case <-evtChan:
			memCli, stopCli, evtChan = subscribeTo(t, tmTypes.EventTx)
			var err error
			txResp, err := nodes.RawTx(memCodec(), memCli, cb.GetAddress(), txBz)
			if err == nil && txResp.Code == 0 {
				t.Fatal("should fail on replay attack")
			}
			cleanup()
			stopCli()
		}
	}
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
		common.RandInt64(),
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
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	// init cache in memory
	pocketTypes.InitCache("data", "data", db.MemDBBackend, db.MemDBBackend, 100, 100)
	genBz, _, validators, app := fiveValidatorsOneAppGenesis()
	kb := getInMemoryKeybase()
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
			Blockchain:         dummyChainsHash,
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
			Chain:              dummyChainsHash,
			SessionBlockHeight: 1,
		}, pocketTypes.RelayEvidence, proof)
		assert.Nil(t, err)
	}
	_, _, cleanup := NewInMemoryTendermintNode(t, genBz)
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
	select {
	case res := <-evtChan:
		fmt.Println(res)
		if res.Events["message.action"][0] != pocketTypes.EventTypeClaim {
			t.Fatal("claim message was not received first")
		}
		_, stopCli, evtChan = subscribeTo(t, tmTypes.EventTx)
		select {
		case res := <-evtChan:
			fmt.Println(res)
			if res.Events["message.action"][0] != pocketTypes.EventTypeProof {
				t.Fatal("proof message was not received afterward")
			}
			cleanup()
			stopCli()
		}
	}
}

func TestClaimTxChallenge(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	pocketTypes.InitCache("data", "data", db.MemDBBackend, db.MemDBBackend, 100, 100)
	genBz, keys, _, _ := fiveValidatorsOneAppGenesis()
	challenges := NewValidChallengeProof(t, keys, 5)
	for _, c := range challenges {
		c.Handle()
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
			Blockchain:         dummyChainsHash,
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
			Blockchain:         dummyChainsHash,
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
			Blockchain:         dummyChainsHash,
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
			MajorityResponses: [2]pocketTypes.RelayResponse{
				majResp1,
				majResp2,
			},
			MinorityResponse: minResp,
			ReporterAddress:  sdk.Address(reporterAddr),
		})
	}
	return proofs
}
