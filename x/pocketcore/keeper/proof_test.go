package keeper

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"math/rand"
	"testing"
)

func TestKeeper_ValidateProof(t *testing.T) { // happy path only todo
	ctx, _, _, _, keeper := createTestInput(t, false)
	npk, validHeader := simulateRelays(t, 1)
	i, found := types.GetAllInvoices().GetInvoice(validHeader)
	if !found {
		t.Fatalf("Set invoice not found")
	}
	root := i.GenerateMerkleRoot()
	totalRelays := types.GetAllInvoices().GetTotalRelays(validHeader)
	assert.Equal(t, totalRelays, int64(9))
	// generate a claim message
	claimMsg := types.MsgClaim{
		SessionHeader: validHeader,
		MerkleRoot:    root,
		TotalRelays:   9,
		FromAddress:   npk.Address(),
	}
	// generate the pseudorandom proof
	neededLeafIndex := keeper.GeneratePseudoRandomProof(ctx, totalRelays, validHeader)
	// create the proof message
	inv, found := types.GetAllInvoices().GetInvoice(validHeader)
	if !found {
		t.Fatalf("Set invoice not found 2")
	}
	merkleProofs, cousinIndex := inv.GenerateMerkleProof(int(neededLeafIndex))
	// get leaf and cousin node
	leafNode := types.GetAllInvoices().GetProof(validHeader, int(neededLeafIndex))
	// get leaf and cousin node
	cousinNode := types.GetAllInvoices().GetProof(validHeader, cousinIndex)
	// create proof message
	proofMsg := types.MsgProof{
		MerkleProofs: merkleProofs,
		Leaf:         leafNode,
		Cousin:       cousinNode,
	}
	// validate proof
	eror := keeper.ValidateProof(ctx, claimMsg, proofMsg)
	if eror != nil {
		t.Fatalf(eror.Error())
	}
}

func simulateRelays(t *testing.T, blockHeight int64) (nodePublicKey crypto.PublicKey, validHeader types.SessionHeader) {
	clientPrivateKey := getRandomPrivateKey()
	clientPubKey := crypto.PublicKey(clientPrivateKey.PubKey().(ed25519.PubKeyEd25519)).String()
	appPrivateKey := getRandomPrivateKey()
	appPubKey := crypto.PublicKey(appPrivateKey.PubKey().(ed25519.PubKeyEd25519)).String()
	npk := getRandomPubKey()
	nodePubKey := hex.EncodeToString(crypto.PublicKey(npk).Bytes())
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
	validRelay1 := types.Relay{
		Payload: types.Payload{
			Data:    "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[],\"id\":67}",
			Method:  "",
			Path:    "",
			Headers: nil,
		},
		Proof: types.RelayProof{
			Entropy:            1,
			SessionBlockHeight: blockHeight,
			ServicerPubKey:     nodePubKey,
			Blockchain:         ethereum,
			Token: types.AAT{
				Version:              "0.0.1",
				ApplicationPublicKey: appPubKey,
				ClientPublicKey:      clientPubKey,
				ApplicationSignature: "",
			},
			Signature: "",
		},
	}
	appSig, er := appPrivateKey.Sign(validRelay1.Proof.Token.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay1.Proof.Token.ApplicationSignature = hex.EncodeToString(appSig)
	clientSig, er := clientPrivateKey.Sign(validRelay1.Proof.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay1.Proof.Signature = hex.EncodeToString(clientSig)
	// valid relay 2
	validRelay2 := validRelay1
	validRelay2.Proof.Entropy = validRelay2.Proof.Entropy + int64(rand.Int())
	clientSig2, er := clientPrivateKey.Sign(validRelay2.Proof.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay2.Proof.Signature = hex.EncodeToString(clientSig2)
	// valid relay 3
	validRelay3 := validRelay1
	validRelay3.Proof.Entropy = validRelay3.Proof.Entropy + int64(rand.Int())
	clientSig3, er := clientPrivateKey.Sign(validRelay3.Proof.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay3.Proof.Signature = hex.EncodeToString(clientSig3)
	// valid relay 4
	validRelay4 := validRelay1
	validRelay4.Proof.Entropy = validRelay4.Proof.Entropy + int64(rand.Int())
	clientSig4, er := clientPrivateKey.Sign(validRelay4.Proof.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay4.Proof.Signature = hex.EncodeToString(clientSig4)
	// valid relay 5
	validRelay5 := validRelay1
	validRelay5.Proof.Entropy = validRelay5.Proof.Entropy + int64(rand.Int())
	clientSig5, er := clientPrivateKey.Sign(validRelay5.Proof.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay5.Proof.Signature = hex.EncodeToString(clientSig5)
	// valid relay 6
	validRelay6 := validRelay1
	validRelay6.Proof.Entropy = validRelay6.Proof.Entropy + int64(rand.Int())
	clientSig6, er := clientPrivateKey.Sign(validRelay6.Proof.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay6.Proof.Signature = hex.EncodeToString(clientSig6)
	// valid relay 7
	validRelay7 := validRelay1
	validRelay7.Proof.Entropy = validRelay7.Proof.Entropy + int64(rand.Int())
	clientSig7, er := clientPrivateKey.Sign(validRelay7.Proof.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay7.Proof.Signature = hex.EncodeToString(clientSig7)
	// valid relay 8
	validRelay8 := validRelay1
	validRelay8.Proof.Entropy = validRelay8.Proof.Entropy + int64(rand.Int())
	clientSig8, er := clientPrivateKey.Sign(validRelay8.Proof.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay8.Proof.Signature = hex.EncodeToString(clientSig8)
	// valid relay 9
	validRelay9 := validRelay1
	validRelay9.Proof.Entropy = validRelay9.Proof.Entropy + int64(rand.Int())
	clientSig9, er := clientPrivateKey.Sign(validRelay9.Proof.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay9.Proof.Signature = hex.EncodeToString(clientSig9)
	// create a session header
	validHeader = types.SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              ethereum,
		SessionBlockHeight: blockHeight,
	}
	err = types.GetAllInvoices().AddToInvoice(validHeader, validRelay1.Proof)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = types.GetAllInvoices().AddToInvoice(validHeader, validRelay2.Proof)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = types.GetAllInvoices().AddToInvoice(validHeader, validRelay3.Proof)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = types.GetAllInvoices().AddToInvoice(validHeader, validRelay4.Proof)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = types.GetAllInvoices().AddToInvoice(validHeader, validRelay5.Proof)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = types.GetAllInvoices().AddToInvoice(validHeader, validRelay6.Proof)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = types.GetAllInvoices().AddToInvoice(validHeader, validRelay7.Proof)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = types.GetAllInvoices().AddToInvoice(validHeader, validRelay8.Proof)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = types.GetAllInvoices().AddToInvoice(validHeader, validRelay9.Proof)
	if err != nil {
		t.Fatalf(err.Error())
	}
	return crypto.PublicKey(npk), validHeader
}

func TestKeeper_GetSetInvoice(t *testing.T) {
	ctx, _, _, _, keeper := createTestInput(t, false)
	appPrivateKey := getRandomPrivateKey()
	appPubKey := crypto.PublicKey(appPrivateKey.PubKey().(ed25519.PubKeyEd25519)).String()
	npk := getRandomPubKey()
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
	// create a session header
	validHeader := types.SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              ethereum,
		SessionBlockHeight: 1,
	}
	storedInvoice := types.StoredInvoice{
		SessionHeader:   validHeader,
		ServicerAddress: npk.Address().String(),
		TotalRelays:     2000,
	}
	keeper.SetInvoice(ctx, sdk.ValAddress(npk.Address()), storedInvoice)
	inv, found := keeper.GetInvoice(ctx, sdk.ValAddress(npk.Address()), validHeader)
	assert.True(t, found)
	assert.Equal(t, inv, storedInvoice)
}

func TestKeeper_GetSetInvoices(t *testing.T) {
	ctx, _, _, _, keeper := createTestInput(t, false)
	appPrivateKey := getRandomPrivateKey()
	appPubKey := crypto.PublicKey(appPrivateKey.PubKey().(ed25519.PubKeyEd25519)).String()
	appPrivateKey2 := getRandomPrivateKey()
	appPubKey2 := crypto.PublicKey(appPrivateKey2.PubKey().(ed25519.PubKeyEd25519)).String()
	npk := crypto.PublicKey(getRandomPubKey())
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
	// create a session header
	validHeader := types.SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              ethereum,
		SessionBlockHeight: 1,
	}
	// create a session header
	validHeader2 := types.SessionHeader{
		ApplicationPubKey:  appPubKey2,
		Chain:              ethereum,
		SessionBlockHeight: 1,
	}
	storedInvoice := types.StoredInvoice{
		SessionHeader:   validHeader,
		ServicerAddress: npk.Address().String(),
		TotalRelays:     2000,
	}
	storedInvoice2 := types.StoredInvoice{
		SessionHeader:   validHeader2,
		ServicerAddress: npk.Address().String(),
		TotalRelays:     2000,
	}
	invoices := []types.StoredInvoice{storedInvoice, storedInvoice2}
	keeper.SetInvoices(ctx, invoices)
	inv := keeper.GetInvoices(ctx, npk.Address())
	assert.Contains(t, inv, storedInvoice)
	assert.Contains(t, inv, storedInvoice2)
}

func TestKeeper_GetAllInvoices(t *testing.T) {
	ctx, _, _, _, keeper := createTestInput(t, false)
	appPrivateKey := getRandomPrivateKey()
	appPubKey := crypto.PublicKey(appPrivateKey.PubKey().(ed25519.PubKeyEd25519)).String()
	npk := crypto.PublicKey(getRandomPubKey())
	npk2 := crypto.PublicKey(getRandomPubKey())
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
	// create a session header
	validHeader := types.SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              ethereum,
		SessionBlockHeight: 1,
	}
	storedInvoice := types.StoredInvoice{
		SessionHeader:   validHeader,
		ServicerAddress: npk.Address().String(),
		TotalRelays:     2000,
	}
	storedInvoice2 := types.StoredInvoice{
		SessionHeader:   validHeader,
		ServicerAddress: npk2.Address().String(),
		TotalRelays:     2000,
	}
	invoices := []types.StoredInvoice{storedInvoice, storedInvoice2}
	keeper.SetInvoices(ctx, invoices)
	inv := keeper.GetAllInvoices(ctx)
	assert.Contains(t, inv, storedInvoice)
	assert.Contains(t, inv, storedInvoice2)
}

func TestKeeper_GetSetClaim(t *testing.T) {
	ctx, _, _, _, keeper := createTestInput(t, false)
	npk, validHeader := simulateRelays(t, 1)
	i, found := types.GetAllInvoices().GetInvoice(validHeader)
	assert.True(t, found)
	claim := types.MsgClaim{
		SessionHeader: validHeader,
		MerkleRoot:    i.GenerateMerkleRoot(),
		TotalRelays:   9,
		FromAddress:   npk.Address(),
	}
	keeper.SetClaim(ctx, claim)
	c, found := keeper.GetClaim(ctx, npk.Address(), validHeader)
	assert.True(t, found)
	assert.Equal(t, claim, c)
}

func TestKeeper_GetSetDeleteClaims(t *testing.T) {
	ctx, _, _, _, keeper := createTestInput(t, false)
	npk, validHeader := simulateRelays(t, 1)
	npk2, validHeader2 := simulateRelays(t, 1)
	i, found := types.GetAllInvoices().GetInvoice(validHeader)
	assert.True(t, found)
	i2, found := types.GetAllInvoices().GetInvoice(validHeader2)
	assert.True(t, found)
	claim1 := types.MsgClaim{
		SessionHeader: validHeader,
		MerkleRoot:    i.GenerateMerkleRoot(),
		TotalRelays:   9,
		FromAddress:   npk.Address(),
	}
	claim2 := types.MsgClaim{
		SessionHeader: validHeader2,
		MerkleRoot:    i2.GenerateMerkleRoot(),
		TotalRelays:   9,
		FromAddress:   npk2.Address(),
	}
	claims := []types.MsgClaim{claim1, claim2}
	keeper.SetClaims(ctx, claims)
	c := keeper.GetClaims(ctx, npk.Address())
	assert.Contains(t, c, claim1)
	c2 := keeper.GetClaims(ctx, npk2.Address())
	assert.Contains(t, c2, claim2)
	c3 := keeper.GetAllClaims(ctx)
	assert.Contains(t, c3, claim1)
	assert.Contains(t, c3, claim2)
	keeper.DeleteClaim(ctx, npk.Address(), validHeader)
	c4 := keeper.GetAllClaims(ctx)
	assert.NotContains(t, c4, claim1)
	assert.Contains(t, c4, claim2)
}

func TestKeeper_GetMatureClaims(t *testing.T) {
	ctx, _, _, _, keeper := createTestInput(t, false)
	npk, validHeader := simulateRelays(t, 1)
	npk2, validHeader2 := simulateRelays(t, 999)
	i, found := types.GetAllInvoices().GetInvoice(validHeader)
	assert.True(t, found)
	i2, found := types.GetAllInvoices().GetInvoice(validHeader2)
	assert.True(t, found)
	matureClaim := types.MsgClaim{
		SessionHeader: validHeader,
		MerkleRoot:    i.GenerateMerkleRoot(),
		TotalRelays:   9,
		FromAddress:   npk.Address(),
	}
	immatureClaim := types.MsgClaim{
		SessionHeader: validHeader2,
		MerkleRoot:    i2.GenerateMerkleRoot(),
		TotalRelays:   9,
		FromAddress:   npk2.Address(),
	}
	claims := []types.MsgClaim{matureClaim, immatureClaim}
	keeper.SetClaims(ctx, claims)
	c1 := keeper.GetMatureClaims(ctx, npk.Address())
	c2 := keeper.GetMatureClaims(ctx, npk2.Address())
	assert.Contains(t, c1, matureClaim)
	assert.Nil(t, c2)
}

func TestKeeper_DeleteExpiredClaims(t *testing.T) {
	ctx, _, _, _, keeper := createTestInput(t, false)
	npk, validHeader := simulateRelays(t, 1)
	npk2, validHeader2 := simulateRelays(t, 999)
	i, found := types.GetAllInvoices().GetInvoice(validHeader)
	assert.True(t, found)
	i2, found := types.GetAllInvoices().GetInvoice(validHeader2)
	assert.True(t, found)
	expiredClaim := types.MsgClaim{
		SessionHeader: validHeader,
		MerkleRoot:    i.GenerateMerkleRoot(),
		TotalRelays:   9,
		FromAddress:   npk.Address(),
	}
	notExpired := types.MsgClaim{
		SessionHeader: validHeader2,
		MerkleRoot:    i2.GenerateMerkleRoot(),
		TotalRelays:   9,
		FromAddress:   npk2.Address(),
	}
	claims := []types.MsgClaim{expiredClaim, notExpired}
	keeper.SetClaims(ctx, claims)
	keeper.DeleteExpiredClaims(ctx)
	c1 := keeper.GetAllClaims(ctx)
	assert.Contains(t, c1, notExpired)
	assert.NotContains(t, c1, expiredClaim)
}
