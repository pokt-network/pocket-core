package keeper

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_GetSetClaim(t *testing.T) {
	ctx, _, _, _, keeper, _ := createTestInput(t, false)
	npk, evidences, header, _, _ := simulateRelays(t, keeper, &ctx, 5)
	evidence, found := evidences.GetEvidence(header, types.RelayEvidence)
	assert.True(t, found)
	claim := types.MsgClaim{
		SessionHeader: header,
		MerkleRoot:    evidence.GenerateMerkleRoot(),
		TotalProofs:   9,
		FromAddress:   sdk.Address(npk.Address()),
		EvidenceType:  types.RelayEvidence,
	}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("PrevCtx", header.SessionBlockHeight).Return(ctx, nil)
	err := keeper.SetClaim(mockCtx, claim)
	assert.Nil(t, err)
	c, found := keeper.GetClaim(mockCtx, sdk.Address(npk.Address()), header, types.RelayEvidence)
	assert.True(t, found)
	assert.Equal(t, claim, c)
}

func TestKeeper_GetSetDeleteClaims(t *testing.T) {
	ctx, _, _, _, keeper, _ := createTestInput(t, false)
	var claims []types.MsgClaim
	var pubKeys []crypto.PublicKey

	for i := 0; i < 2; i++ {
		npk, evidences, header, _, _ := simulateRelays(t, keeper, &ctx, 5)
		evidence, found := evidences.GetEvidence(header, types.RelayEvidence)
		assert.True(t, found)
		claim := types.MsgClaim{
			SessionHeader: header,
			MerkleRoot:    evidence.GenerateMerkleRoot(),
			TotalProofs:   9,
			FromAddress:   sdk.Address(sdk.Address(npk.Address())),
			EvidenceType:  types.RelayEvidence,
		}
		claims = append(claims, claim)
		pubKeys = append(pubKeys, npk)
	}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("PrevCtx", claims[0].SessionBlockHeight).Return(ctx, nil)
	mockCtx.On("Logger").Return(ctx.Logger())
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
	keeper.DeleteClaim(mockCtx, sdk.Address(pubKeys[0].Address()), claims[0].SessionHeader, types.RelayEvidence)
	c = keeper.GetAllClaims(ctx)
	assert.NotContains(t, c, claims[0])
	assert.Contains(t, c, claims[1])
}

func TestKeeper_GetMatureClaims(t *testing.T) {
	ctx, _, _, _, keeper, keys := createTestInput(t, false)
	npk, evidences, header, _, _ := simulateRelays(t, keeper, &ctx, 5)
	npk2, evidences2, header2, _, _ := simulateRelays(t, keeper, &ctx, 999)

	i, found := evidences.GetEvidence(header, types.RelayEvidence)
	assert.True(t, found)
	i2, found := evidences2.GetEvidence(header2, types.RelayEvidence)
	assert.True(t, found)

	matureClaim := types.MsgClaim{
		SessionHeader: header,
		MerkleRoot:    i.GenerateMerkleRoot(),
		TotalProofs:   9,
		FromAddress:   sdk.Address(npk.Address()),
		EvidenceType:  types.RelayEvidence,
	}
	immatureClaim := types.MsgClaim{
		SessionHeader: header2,
		MerkleRoot:    i2.GenerateMerkleRoot(),
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
	assert.Contains(t, c1, matureClaim)
	assert.Nil(t, c2)
}

func TestKeeper_DeleteExpiredClaims(t *testing.T) {
	ctx, _, _, _, keeper, keys := createTestInput(t, false)
	npk, inevidenceMap, header, _, _ := simulateRelays(t, keeper, &ctx, 5)
	npk2, inevidenceMap2, header2, _, _ := simulateRelays(t, keeper, &ctx, 999)

	i, found := inevidenceMap.GetEvidence(header, types.RelayEvidence)
	assert.True(t, found)
	i2, found := inevidenceMap2.GetEvidence(header2, types.RelayEvidence)
	assert.True(t, found)
	expiredClaim := types.MsgClaim{
		SessionHeader: header,
		MerkleRoot:    i.GenerateMerkleRoot(),
		TotalProofs:   9,
		FromAddress:   sdk.Address(npk.Address()),
		EvidenceType:  types.RelayEvidence,
	}
	header2.SessionBlockHeight = int64(20) // NOTE start a later block than 1
	notExpired := types.MsgClaim{
		SessionHeader: header2,
		MerkleRoot:    i2.GenerateMerkleRoot(),
		TotalProofs:   9,
		FromAddress:   sdk.Address(npk2.Address()),
		EvidenceType:  types.RelayEvidence,
	}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("KVStore", keys["params"]).Return(ctx.KVStore(keys["params"]))
	mockCtx.On("PrevCtx", header.SessionBlockHeight).Return(ctx, nil)
	mockCtx.On("PrevCtx", header2.SessionBlockHeight).Return(ctx, nil)
	mockCtx.On("BlockHeight").Return(int64(2501)) // NOTE minimum height to start expiring from block 1

	claims := []types.MsgClaim{expiredClaim, notExpired}
	keeper.SetClaims(mockCtx, claims)
	keeper.DeleteExpiredClaims(mockCtx)
	c1 := keeper.GetAllClaims(mockCtx)

	assert.Contains(t, c1, notExpired, "does not contain notExpired claim")
	assert.NotContains(t, c1, expiredClaim, "contains expired claim")
}
