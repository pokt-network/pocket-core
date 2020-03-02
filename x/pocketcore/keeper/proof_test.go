package keeper

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_ValidateProof(t *testing.T) { // happy path only todo
	ctx, _, _, _, keeper, keys := createTestInput(t, false)
	npk, memInvoices, header, _, _ := simulateRelays(t, keeper, &ctx, 5)

	// create a session header
	i, found := memInvoices.GetInvoice(header)
	if !found {
		t.Fatalf("Set invoice not found")
	}

	root := i.GenerateMerkleRoot()
	totalRelays := memInvoices.GetTotalRelays(header)
	assert.Equal(t, totalRelays, int64(5))
	// generate a claim message
	claimMsg := types.MsgClaim{
		SessionHeader: header,
		MerkleRoot:    root,
		TotalRelays:   5,
		FromAddress:   sdk.Address(sdk.Address(npk.Address())),
	}

	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("KVStore", keys["params"]).Return(ctx.KVStore(keys["params"]))
	mockCtx.On("MustGetPrevCtx", header.SessionBlockHeight +  keeper.ProofWaitingPeriod(ctx)*keeper.SessionFrequency(ctx)).Return(ctx)

	// generate the pseudorandom proof
	neededLeafIndex, er := keeper.GetPseudorandomIndex(mockCtx, totalRelays, header)
	assert.Nil(t, er)

	// create the proof message
	inv, found := memInvoices.GetInvoice(header)
	if !found {
		t.Fatalf("Set invoice not found 2")
	}
	merkleProofs, cousinIndex := inv.GenerateMerkleProof(int(neededLeafIndex))
	// get leaf and cousin node
	leafNode := memInvoices.GetProof(header, int(neededLeafIndex))
	// get leaf and cousin node
	cousinNode := memInvoices.GetProof(header, cousinIndex)
	// create proof message
	proofMsg := types.MsgProof{
		MerkleProofs: merkleProofs,
		Leaf:         leafNode,
		Cousin:       cousinNode,
	}
	// validate proof
	eror := keeper.ValidateProof(mockCtx, claimMsg, proofMsg)
	if eror != nil {
		t.Fatalf(eror.Error())
	}
}

func TestKeeper_GetPsuedorandomIndex(t *testing.T) {
	var totalRelays []int = []int{10, 1000, 10000000}
	for _, relays := range totalRelays{
		ctx, _, _, _, keeper, keys := createTestInput(t, false)
		_, _, header, _, _ := simulateRelays(t, keeper, &ctx, 999)

		mockCtx := new(Ctx)
		mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
		mockCtx.On("KVStore", keys["params"]).Return(ctx.KVStore(keys["params"]))
		mockCtx.On("MustGetPrevCtx", header.SessionBlockHeight +  keeper.ProofWaitingPeriod(ctx)*keeper.SessionFrequency(ctx)).Return(ctx)

		// generate the pseudorandom proof
		neededLeafIndex, err := keeper.GetPseudorandomIndex(mockCtx, int64(relays), header)
		assert.Nil(t, err)
		assert.LessOrEqual(t, neededLeafIndex, int64(relays))
	}
}

func TestKeeper_GetSetInvoice(t *testing.T) {
	ctx, _, _, _, keeper,_ := createTestInput(t, false)
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
		SessionBlockHeight: 1,
	}
	storedInvoice := types.StoredInvoice{
		SessionHeader:   validHeader,
		ServicerAddress: sdk.Address(npk.Address()).String(),
		TotalRelays:     2000,
	}
	addr :=  sdk.Address(sdk.Address(npk.Address()))
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("MustGetPrevCtx", validHeader.SessionBlockHeight).Return(ctx)
	keeper.SetInvoice(mockCtx, addr, storedInvoice)

	inv, found := keeper.GetInvoice(mockCtx, sdk.Address(sdk.Address(npk.Address())), validHeader)
	assert.True(t, found)
	assert.Equal(t, inv, storedInvoice)
}

func TestKeeper_GetSetInvoices(t *testing.T) {
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
	storedInvoice := types.StoredInvoice{
		SessionHeader:   validHeader,
		ServicerAddress: sdk.Address(npk.Address()).String(),
		TotalRelays:     2000,
	}
	storedInvoice2 := types.StoredInvoice{
		SessionHeader:   validHeader2,
		ServicerAddress: sdk.Address(npk.Address()).String(),
		TotalRelays:     2000,
	}
	invoices := []types.StoredInvoice{storedInvoice, storedInvoice2}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("MustGetPrevCtx", validHeader.SessionBlockHeight).Return(ctx)
	keeper.SetInvoices(mockCtx, invoices)
	inv, er := keeper.GetInvoices(mockCtx, sdk.Address(sdk.Address(npk.Address())))
	assert.Nil(t, er)
	assert.Contains(t, inv, storedInvoice)
	assert.Contains(t, inv, storedInvoice2)
}

func TestKeeper_GetAllInvoices(t *testing.T) {
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
	storedInvoice := types.StoredInvoice{
		SessionHeader:   validHeader,
		ServicerAddress: sdk.Address(npk.Address()).String(),
		TotalRelays:     2000,
	}
	storedInvoice2 := types.StoredInvoice{
		SessionHeader:   validHeader,
		ServicerAddress: sdk.Address(npk2.Address()).String(),
		TotalRelays:     2000,
	}
	invoices := []types.StoredInvoice{storedInvoice, storedInvoice2}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("MustGetPrevCtx", validHeader.SessionBlockHeight).Return(ctx)
	keeper.SetInvoices(mockCtx, invoices)
	inv := keeper.GetAllInvoices(mockCtx)
	assert.Contains(t, inv, storedInvoice)
	assert.Contains(t, inv, storedInvoice2)
}

func TestKeeper_GetSetClaim(t *testing.T) {
	ctx, _, _, _, keeper, _ := createTestInput(t, false)
	npk, inMemInvoices, header, _, _ := simulateRelays(t, keeper, &ctx, 5)
	i, found := inMemInvoices.GetInvoice(header)
	assert.True(t, found)
	claim := types.MsgClaim{
		SessionHeader: header,
		MerkleRoot:    i.GenerateMerkleRoot(),
		TotalRelays:   9,
		FromAddress:   sdk.Address(npk.Address()),
	}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("MustGetPrevCtx", header.SessionBlockHeight).Return(ctx)
	keeper.SetClaim(mockCtx, claim)
	c, found := keeper.GetClaim(mockCtx, sdk.Address(npk.Address()), header)
	assert.True(t, found)
	assert.Equal(t, claim, c)
}

func TestKeeper_GetSetDeleteClaims(t *testing.T) {
	ctx, _, _, _, keeper, _ := createTestInput(t, false)
	var claims []types.MsgClaim
	var pubKeys []crypto.PublicKey

	for i := 0; i<2; i++ {
		npk, inMemInvoices, header, _, _ := simulateRelays(t, keeper, &ctx, 5)
		i, found := inMemInvoices.GetInvoice(header)
		assert.True(t, found)
		claim := types.MsgClaim{
			SessionHeader: header,
			MerkleRoot:    i.GenerateMerkleRoot(),
			TotalRelays:   9,
			FromAddress:   sdk.Address(sdk.Address(npk.Address())),
		}
		claims = append(claims, claim)
		pubKeys = append(pubKeys, npk)
	}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("MustGetPrevCtx", claims[0].SessionBlockHeight).Return(ctx)
	keeper.SetClaims(mockCtx, claims)
	// todo store npk & headers
	for idx, pk := range pubKeys {
		c, err := keeper.GetClaims(mockCtx, sdk.Address(pk.Address()))
		assert.Nil(t, err)
		assert.Contains(t, c, claims[idx])
	}
	c := keeper.GetAllClaims(mockCtx)
	assert.Contains(t, c, claims[0])
	assert.Contains(t, c, claims[1])
	keeper.DeleteClaim(mockCtx, sdk.Address(pubKeys[0].Address()), claims[0].SessionHeader)
	c = keeper.GetAllClaims(ctx)
	assert.NotContains(t, c, claims[0])
	assert.Contains(t, c, claims[1])
}

func TestKeeper_GetMatureClaims(t *testing.T) {
	ctx, _, _, _, keeper, keys := createTestInput(t, false)
	npk, inMemInvoices, header, _, _ := simulateRelays(t, keeper, &ctx, 5)
	npk2, inMemInvoices2, header2, _, _ := simulateRelays(t, keeper, &ctx, 999)

	i, found := inMemInvoices.GetInvoice(header)
	assert.True(t, found)
	i2, found := inMemInvoices2.GetInvoice(header2)
	assert.True(t, found)

	matureClaim := types.MsgClaim{
		SessionHeader: header,
		MerkleRoot:    i.GenerateMerkleRoot(),
		TotalRelays:   9,
		FromAddress:   sdk.Address(npk.Address()),
	}
	immatureClaim := types.MsgClaim{
		SessionHeader: header2,
		MerkleRoot:    i2.GenerateMerkleRoot(),
		TotalRelays:   9,
		FromAddress:   sdk.Address(npk2.Address()),
	}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("KVStore", keys["params"]).Return(ctx.KVStore(keys["params"]))
	mockCtx.On("MustGetPrevCtx", header.SessionBlockHeight).Return(ctx)
	mockCtx.On("BlockHeight").Return(ctx.BlockHeight())

	claims := []types.MsgClaim{matureClaim, immatureClaim}
	keeper.SetClaims(mockCtx, claims)
	c1, err := keeper.GetMatureClaims(mockCtx, sdk.Address(npk.Address()))
	assert.Nil(t, err)

	mockCtx = new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("KVStore", keys["params"]).Return(ctx.KVStore(keys["params"]))
	mockCtx.On("MustGetPrevCtx", header.SessionBlockHeight).Return(ctx)
	mockCtx.On("BlockHeight").Return(int64(1))

	c2, err := keeper.GetMatureClaims(mockCtx, sdk.Address(npk2.Address()))
	assert.Nil(t, err)
	assert.Contains(t, c1, matureClaim)
	assert.Nil(t, c2)
}

func TestKeeper_DeleteExpiredClaims(t *testing.T) {
	ctx, _, _, _, keeper, keys := createTestInput(t, false)
	npk, inMemInvoices, header, _, _ := simulateRelays(t, keeper, &ctx, 5)
	npk2, inMemInvoices2, header2, _, _ := simulateRelays(t, keeper, &ctx, 999)

	i, found := inMemInvoices.GetInvoice(header)
	assert.True(t, found)
	i2, found := inMemInvoices2.GetInvoice(header2)
	assert.True(t, found)
	expiredClaim := types.MsgClaim{
		SessionHeader: header,
		MerkleRoot:    i.GenerateMerkleRoot(),
		TotalRelays:   9,
		FromAddress:   sdk.Address(npk.Address()),
	}
	header2.SessionBlockHeight = int64(20) // NOTE start a later block than 1
	notExpired := types.MsgClaim{
		SessionHeader: header2,
		MerkleRoot:    i2.GenerateMerkleRoot(),
		TotalRelays:   9,
		FromAddress:   sdk.Address(npk2.Address()),
	}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("KVStore", keys["params"]).Return(ctx.KVStore(keys["params"]))
	mockCtx.On("MustGetPrevCtx", header.SessionBlockHeight).Return(ctx)
	mockCtx.On("MustGetPrevCtx", header2.SessionBlockHeight).Return(ctx)
	mockCtx.On("BlockHeight").Return(int64(2501)) // NOTE minimum height to start expiring from block 1

	claims := []types.MsgClaim{expiredClaim, notExpired}
	keeper.SetClaims(mockCtx, claims)
	keeper.DeleteExpiredClaims(mockCtx)
	c1 := keeper.GetAllClaims(mockCtx)

	assert.Contains(t, c1, notExpired, "does not contain notExpired claim")
	assert.NotContains(t, c1, expiredClaim, "contains expired claim")
}
