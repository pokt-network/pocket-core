package keeper

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	"strconv"
	"sync"
)

var (
	globalProofBatches *pc.ProofBatches // map[ProofsHeader]ProofBatch
	pbOnce             sync.Once
)

func (k Keeper) GetProofBatches() *pc.ProofBatches {
	pbOnce.Do(func() {
		globalProofBatches = (*pc.ProofBatches)(types.NewList())
	})
	return globalProofBatches
}

// todo possible attacks around tricking the client and making the requests after the nexSessionBlockHash has been revealed?
func (k Keeper) GenerateProofs(ctx sdk.Context, sessionBlockHeight, relaysCompleted int64, sessionKey string) []int64 { // todo created on the spot! need to audit. Possible to brute force?
	var result []int64
	nextSessionContext := ctx.WithBlockHeight(sessionBlockHeight + k.posKeeper.SessionBlockFrequency(ctx)) // next session block hash
	blockHash := crypto.HexEncodeToString(nextSessionContext.BlockHeader().GetLastBlockId().Hash)
	proofsHash := hex.EncodeToString(crypto.SHA3FromString(blockHash + sessionKey)) // makes it unique for each session!
	proofsHash = proofsHash[:16]                                                    // take first 16 characters to fit int 64
	for i := 0; i < len(proofsHash); i++ {
		res, err := strconv.ParseInt(proofsHash[i:], 16, 64)
		if err != nil {
			panic(err)
		}
		if relaysCompleted > res {
			result = append(result, res)
		}
	}
	return result
}

func (k Keeper) AreProofsValid(sessionContext sdk.Context, proofsIndex []int64, proofs pc.Proofs) bool {
	if len(proofsIndex) != len(proofs) {
		return false
	}
	clientPubKeyBz, err := hex.DecodeString(proofs[0].Token.ClientPublicKey)
	if err != nil {
		return false
	}
	verifyNodePubKey := proofs[0].NodePublicKey
	nextSessionStartTime := sessionContext.WithBlockHeight(sessionContext.BlockHeight() + k.posKeeper.SessionBlockFrequency(sessionContext)).BlockTime()
	for i, proof := range proofs {
		if proof.Counter != int(proofsIndex[i]) {
			return false
		}
		if proof.NodePublicKey != verifyNodePubKey {
			return false
		}
		if proof.Timestamp.After(nextSessionStartTime) {
			return false
		}
		if err := proof.Token.Validate(); err != nil {
			return false
		}
		proofSigBz, err := hex.DecodeString(proof.Signature)
		if err != nil {
			return false
		}
		messageHash := crypto.SHA3FromString(strconv.Itoa(proof.Counter) + proof.Timestamp.String() + proof.Token.HashString() + proof.NodePublicKey)
		if !crypto.MockVerifySignature(clientPubKeyBz, messageHash, proofSigBz) {
			return false
		}
	}
	return true
}

func (k Keeper) GetProofsSummary(ctx sdk.Context, address sdk.ValAddress, header pc.ProofsHeader) (summary pc.ProofSummary) {
	store := ctx.KVStore(k.storeKey)
	res := store.Get(pc.KeyForNodeProofSummary(address, header))
	k.cdc.MustUnmarshalBinaryBare(res, &summary)
	return
}

func (k Keeper) GetAllProofSummaries(ctx sdk.Context, address sdk.ValAddress) (summaries []pc.ProofSummary) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, pc.KeyForNodeProofSummaries(address))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var summary pc.ProofSummary
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &summary)
		summaries = append(summaries, summary)
	}
	return
}

func (k Keeper) GetAllProofSummariesForApp(ctx sdk.Context, address sdk.ValAddress, appPubKeyHex string) (summaries []pc.ProofSummary) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, pc.KeyForNodeProofSummariesForApp(address, appPubKeyHex))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var summary pc.ProofSummary
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &summary)
		summaries = append(summaries, summary)
	}
	return
}

func (k Keeper) SetProofSummary(ctx sdk.Context, address sdk.ValAddress, summary pc.ProofSummary) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(summary)
	store.Set(pc.KeyForNodeProofSummary(address, summary.ProofsHeader), bz)
}
