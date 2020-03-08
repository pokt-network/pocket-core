package keeper

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_ValidateProof(t *testing.T) { // happy path only todo
	types.GetEvidenceMap().Clear()
	ctx, _, _, _, keeper, keys := createTestInput(t, false)
	npk, evidenceMap, header, _, _ := simulateRelays(t, keeper, &ctx, 5)
	evidence, found := evidenceMap.GetEvidence(header, types.RelayEvidence)
	if !found {
		t.Fatalf("Set invoice not found")
	}

	root := evidence.GenerateMerkleRoot()
	totalRelays := evidenceMap.GetTotalRelays(header)
	assert.Equal(t, totalRelays, int64(5))
	// generate a claim message
	claimMsg := types.MsgClaim{
		SessionHeader: header,
		MerkleRoot:    root,
		TotalProofs:   5,
		FromAddress:   sdk.Address(npk.Address()),
		EvidenceType:  types.RelayEvidence,
	}

	mockCtx := &Ctx{}
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("KVStore", keys["params"]).Return(ctx.KVStore(keys["params"]))
	mockCtx.On("KVStore", keys["application"]).Return(ctx.KVStore(keys["application"]))
	mockCtx.On("Logger").Return(ctx.Logger())
	mockCtx.On("PrevCtx", header.SessionBlockHeight).Return(ctx, nil)
	mockCtx.On("PrevCtx", header.SessionBlockHeight+keeper.ClaimSubmissionWindow(ctx)*keeper.SessionFrequency(ctx)).Return(ctx, nil)

	// generate the pseudorandom proof
	neededLeafIndex, er := keeper.getPseudorandomIndex(mockCtx, totalRelays, header)
	assert.Nil(t, er)

	// create the proof message
	ev, found := evidenceMap.GetEvidence(header, types.RelayEvidence)
	if !found {
		t.Fatalf("Set evidence not found 2")
	}
	merkleProofs, cousinIndex := ev.GenerateMerkleProof(int(neededLeafIndex))
	// get leaf and cousin node
	leafNode := evidenceMap.GetProof(header, types.RelayEvidence, int(neededLeafIndex))
	// get leaf and cousin node
	cousinNode := evidenceMap.GetProof(header, types.RelayEvidence, cousinIndex)
	// create proof message
	proofMsg := types.MsgProof{
		MerkleProofs: merkleProofs,
		Leaf:         leafNode.(types.RelayProof),
		Cousin:       cousinNode.(types.RelayProof),
	}
	err := keeper.SetClaim(mockCtx, claimMsg)
	if err != nil {
		t.Fatal(err)
	}
	// validate proof
	_, _, err = keeper.ValidateProof(mockCtx, proofMsg)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestKeeper_GetPsuedorandomIndex(t *testing.T) {
	var totalRelays []int = []int{10, 1000, 10000000}
	for _, relays := range totalRelays {
		ctx, _, _, _, keeper, keys := createTestInput(t, false)
		_, _, header, _, _ := simulateRelays(t, keeper, &ctx, 999)

		mockCtx := new(Ctx)
		mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
		mockCtx.On("KVStore", keys["params"]).Return(ctx.KVStore(keys["params"]))
		mockCtx.On("PrevCtx", header.SessionBlockHeight+keeper.ClaimSubmissionWindow(ctx)*keeper.SessionFrequency(ctx)).Return(ctx, nil)

		// generate the pseudorandom proof
		neededLeafIndex, err := keeper.getPseudorandomIndex(mockCtx, int64(relays), header)
		assert.Nil(t, err)
		assert.LessOrEqual(t, neededLeafIndex, int64(relays))
	}
}

func TestKeeper_GetSetReceipt(t *testing.T) {
	ctx, _, _, _, keeper, _ := createTestInput(t, false)
	appPrivateKey := getRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
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
		SessionBlockHeight: 976,
	}
	receipt := types.Receipt{
		SessionHeader:   validHeader,
		ServicerAddress: sdk.Address(npk.Address()).String(),
		Total:           2000,
		EvidenceType:    types.RelayEvidence,
	}
	addr := sdk.Address(sdk.Address(npk.Address()))
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("PrevCtx", validHeader.SessionBlockHeight).Return(ctx, nil)
	mockCtx.On("Logger").Return(ctx.Logger())
	keeper.SetReceipt(mockCtx, addr, receipt)

	inv, found := keeper.GetReceipt(mockCtx, sdk.Address(npk.Address()), validHeader, receipt.EvidenceType)
	assert.True(t, found)
	assert.Equal(t, inv, receipt)
}

func TestKeeper_GetSetReceipts(t *testing.T) {
	ctx, _, _, _, keeper, _ := createTestInput(t, false)
	appPrivateKey := getRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	appPrivateKey2 := getRandomPrivateKey()
	appPubKey2 := appPrivateKey2.PublicKey().RawString()
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
	// create a session header
	validHeader2 := types.SessionHeader{
		ApplicationPubKey:  appPubKey2,
		Chain:              ethereum,
		SessionBlockHeight: 1,
	}
	receipt := types.Receipt{
		SessionHeader:   validHeader,
		ServicerAddress: sdk.Address(npk.Address()).String(),
		Total:           2000,
		EvidenceType:    types.RelayEvidence,
	}
	receipt2 := types.Receipt{
		SessionHeader:   validHeader2,
		ServicerAddress: sdk.Address(npk.Address()).String(),
		Total:           2000,
		EvidenceType:    types.RelayEvidence,
	}
	receipts := []types.Receipt{receipt, receipt2}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("PrevCtx", validHeader.SessionBlockHeight).Return(ctx, nil)
	mockCtx.On("Logger").Return(ctx.Logger())
	keeper.SetReceipts(mockCtx, receipts)
	inv, er := keeper.GetReceipts(mockCtx, sdk.Address(sdk.Address(npk.Address())))
	assert.Nil(t, er)
	assert.Contains(t, inv, receipt)
	assert.Contains(t, inv, receipt2)
}

func TestKeeper_GetAllReceipts(t *testing.T) {
	ctx, _, _, _, keeper, _ := createTestInput(t, false)
	appPrivateKey := getRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	npk := getRandomPubKey()
	npk2 := getRandomPubKey()
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
	receipt := types.Receipt{
		SessionHeader:   validHeader,
		ServicerAddress: sdk.Address(npk.Address()).String(),
		Total:           2000,
		EvidenceType:    types.RelayEvidence,
	}
	receipt2 := types.Receipt{
		SessionHeader:   validHeader,
		ServicerAddress: sdk.Address(npk2.Address()).String(),
		Total:           2000,
		EvidenceType:    types.RelayEvidence,
	}
	receipts := []types.Receipt{receipt, receipt2}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("PrevCtx", validHeader.SessionBlockHeight).Return(ctx, nil)
	mockCtx.On("Logger").Return(ctx.Logger())
	keeper.SetReceipts(mockCtx, receipts)
	inv := keeper.GetAllReceipts(mockCtx)
	assert.Contains(t, inv, receipt)
	assert.Contains(t, inv, receipt2)
}
