package keeper

import (
	"testing"

	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_GetSetClaim(t *testing.T) {
	ctx, _, _, _, keeper, _, _ := createTestInput(t, false)
	npk, header, _ := simulateRelays(t, keeper, &ctx, 5)
	evidence, err := types.GetEvidence(header, types.RelayEvidence, sdk.NewInt(100000))
	assert.Nil(t, err)
	claim := types.MsgClaim{
		SessionHeader: header,
		MerkleRoot:    evidence.GenerateMerkleRoot(0),
		TotalProofs:   9,
		FromAddress:   sdk.Address(npk.Address()),
		EvidenceType:  types.RelayEvidence,
	}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("BlockHeight").Return(int64(1))
	mockCtx.On("PrevCtx", header.SessionBlockHeight).Return(ctx, nil)
	mockCtx.On("BlockHash", header.SessionBlockHeight).Return(types.Hash([]byte("fake")), nil)
	err = keeper.SetClaim(mockCtx, claim)
	assert.Nil(t, err)
	c, found := keeper.GetClaim(mockCtx, sdk.Address(npk.Address()), header, types.RelayEvidence)
	assert.True(t, found)
	assert.Equal(t, claim.MerkleRoot, c.MerkleRoot)
}

func TestKeeper_GetSetDeleteClaims(t *testing.T) {
	ctx, _, _, _, keeper, _, _ := createTestInput(t, false)
	var claims []types.MsgClaim
	var pubKeys []crypto.PublicKey

	for i := 0; i < 2; i++ {
		npk, header, _ := simulateRelays(t, keeper, &ctx, 5)
		evidence, err := types.GetEvidence(header, types.RelayEvidence, sdk.NewInt(1000))
		assert.Nil(t, err)
		claim := types.MsgClaim{
			SessionHeader: header,
			MerkleRoot:    evidence.GenerateMerkleRoot(0),
			TotalProofs:   9,
			FromAddress:   sdk.Address(sdk.Address(npk.Address())),
			EvidenceType:  types.RelayEvidence,
		}
		claims = append(claims, claim)
		pubKeys = append(pubKeys, npk)
	}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("PrevCtx", claims[0].SessionHeader.SessionBlockHeight).Return(ctx, nil)
	mockCtx.On("BlockHeight").Return(int64(1))
	mockCtx.On("Logger").Return(ctx.Logger())
	keeper.SetClaims(mockCtx, claims)
	c := keeper.GetAllClaims(mockCtx)
	assert.Len(t, c, 2)
	_ = keeper.DeleteClaim(mockCtx, sdk.Address(pubKeys[0].Address()), claims[0].SessionHeader, types.RelayEvidence)
	_, err := keeper.GetClaim(ctx, claims[0].FromAddress, claims[0].SessionHeader, claims[0].EvidenceType)
	assert.NotNil(t, err)
}

func TestKeeper_GetMatureClaims(t *testing.T) {
	ctx, _, _, _, keeper, keys, _ := createTestInput(t, false)
	npk, header, _ := simulateRelays(t, keeper, &ctx, 5)
	npk2, header2, _ := simulateRelays(t, keeper, &ctx, 20)

	i, err := types.GetEvidence(header, types.RelayEvidence, sdk.NewInt(1000))
	assert.Nil(t, err)
	i2, err := types.GetEvidence(header2, types.RelayEvidence, sdk.NewInt(1000))
	assert.Nil(t, err)

	matureClaim := types.MsgClaim{
		SessionHeader: header,
		MerkleRoot:    i.GenerateMerkleRoot(0),
		TotalProofs:   9,
		FromAddress:   sdk.Address(npk.Address()),
		EvidenceType:  types.RelayEvidence,
	}
	immatureClaim := types.MsgClaim{
		SessionHeader: header2,
		MerkleRoot:    i2.GenerateMerkleRoot(0),
		TotalProofs:   9,
		FromAddress:   sdk.Address(npk2.Address()),
		EvidenceType:  types.RelayEvidence,
	}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("KVStore", keys["params"]).Return(ctx.KVStore(keys["params"]))
	mockCtx.On("PrevCtx", header.SessionBlockHeight).Return(ctx, nil)
	mockCtx.On("BlockHeight").Return(ctx.BlockHeight())

	claims := []types.MsgClaim{matureClaim, immatureClaim}
	keeper.SetClaims(mockCtx, claims)
	c1, err := keeper.GetMatureClaims(mockCtx, sdk.Address(npk.Address()))
	assert.Nil(t, err)

	mockCtx = new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("KVStore", keys["params"]).Return(ctx.KVStore(keys["params"]))
	mockCtx.On("PrevCtx", header.SessionBlockHeight).Return(ctx, nil)
	mockCtx.On("BlockHeight").Return(int64(1))

	c2, err := keeper.GetMatureClaims(mockCtx, sdk.Address(npk2.Address()))
	assert.Nil(t, err)
	assert.Len(t, c1, 1)
	assert.Nil(t, c2)
}

func TestKeeper_DeleteExpiredClaims(t *testing.T) {
	ctx, _, _, _, keeper, keys, _ := createTestInput(t, false)
	npk, header, _ := simulateRelays(t, keeper, &ctx, 5)
	npk2, header2, _ := simulateRelays(t, keeper, &ctx, 20)

	i, err := types.GetEvidence(header, types.RelayEvidence, sdk.NewInt(1000))
	assert.Nil(t, err)
	i2, err := types.GetEvidence(header2, types.RelayEvidence, sdk.NewInt(1000))
	assert.Nil(t, err)
	expiredClaim := types.MsgClaim{
		SessionHeader: header,
		MerkleRoot:    i.GenerateMerkleRoot(0),
		TotalProofs:   9,
		FromAddress:   sdk.Address(npk.Address()),
		EvidenceType:  types.RelayEvidence,
	}
	header2.SessionBlockHeight = int64(20) // NOTE start a later block than 1
	notExpired := types.MsgClaim{
		SessionHeader: header2,
		MerkleRoot:    i2.GenerateMerkleRoot(0),
		TotalProofs:   9,
		FromAddress:   sdk.Address(npk2.Address()),
		EvidenceType:  types.RelayEvidence,
	}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("KVStore", keys["params"]).Return(ctx.KVStore(keys["params"]))
	mockCtx.On("PrevCtx", header.SessionBlockHeight).Return(ctx, nil)
	mockCtx.On("PrevCtx", header2.SessionBlockHeight).Return(ctx, nil)
	mockCtx.On("BlockHeight").Return(int64(1))
	mockCtx.On("BlockHeight").Return(int64(2501)) // NOTE minimum height to start expiring from block 1

	claims := []types.MsgClaim{expiredClaim, notExpired}
	keeper.SetClaims(mockCtx, claims)
	keeper.DeleteExpiredClaims(mockCtx)
	c1 := keeper.GetAllClaims(mockCtx)
	notExpired.ExpirationHeight = 2501
	assert.Contains(t, c1, notExpired, "does not contain notExpired claim")
	assert.NotContains(t, c1, expiredClaim, "contains expired claim")
}
