package keeper

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	"strconv"
	"sync"
)

var (
	globalAllProofs *pc.AllProofs
	apOnce          sync.Once
)

func (Keeper) GetAllProofs() *pc.AllProofs {
	apOnce.Do(func() {
		if globalAllProofs == nil {
			*globalAllProofs = make(map[string]pc.ProofOfRelay)
		}
	})
	return globalAllProofs
}

func (k Keeper) GenerateProofs(ctx sdk.Context, totalRelays int64, header pc.PORHeader) int64 {
	proofSessBlockContext := ctx.WithBlockHeight(header.SessionBlockHeight + int64(k.ProofWaitingPeriod(ctx))*k.SessionFrequency(ctx)) // next session block hash
	blockHash := crypto.HexEncodeToString(proofSessBlockContext.BlockHeader().GetLastBlockId().Hash)
	proofsHash := hex.EncodeToString(crypto.SHA3FromString(blockHash + header.String()))[:16] // makes it unique for each session!
	length := len(proofsHash)
	for i := 0; i < length; i++ {
		res, err := strconv.ParseInt(proofsHash[i:], 16, 64)
		if err != nil {
			panic(err)
		}
		if totalRelays > res {
			return res // todo created on the spot! need to audit. Possible to brute force?
		}
	}
	return 0
}

func (k Keeper) TrucateUnnecessaryProofs(porReq int, por pc.ProofOfRelay) pc.ProofOfRelay {
	kept := por.Proofs[porReq]
	por.Proofs = make([]pc.Proof, 1)
	por.Proofs = append(por.Proofs, kept)
	return por
}

func (k Keeper) AreProofsValid(ctx sdk.Context, verifyServicerPubKey string, por pc.ProofOfRelay) bool {
	if 1 != len(por.Proofs) {
		return false
	}
	reqProof := k.GenerateProofs(ctx, por.TotalRelays, por.PORHeader)
	proof := por.Proofs[0]
	if proof.Index != int(reqProof) {
		return false
	}
	if proof.ServicerPubKey != verifyServicerPubKey {
		return false
	}
	if err := proof.Token.Validate(); err != nil {
		return false
	}
	messageHash := hex.EncodeToString(crypto.SHA3FromString(strconv.Itoa(proof.Index) + proof.ServicerPubKey + proof.Token.HashString() + strconv.Itoa(int(proof.SessionBlockHeight)))) // todo standardize
	if !crypto.MockVerifySignature(por.Proofs[0].Token.ClientPublicKey, messageHash, proof.Signature) {
		return false

	}
	return true
}

func (k Keeper) GetProofsSummary(ctx sdk.Context, address sdk.ValAddress, header pc.PORHeader) (summary pc.ProofOfRelay) {
	store := ctx.KVStore(k.storeKey)
	res := store.Get(pc.KeyForProofOfRelay(ctx, address, header))
	k.cdc.MustUnmarshalBinaryBare(res, &summary)
	return
}

func (k Keeper) GetAllProofSummaries(ctx sdk.Context, address sdk.ValAddress) (summaries []pc.ProofOfRelay) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, pc.KeyForProofOfRelays(address))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var summary pc.ProofOfRelay
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &summary)
		summaries = append(summaries, summary)
	}
	return
}

func (k Keeper) GetAllProofSummariesForApp(ctx sdk.Context, address sdk.ValAddress, appPubKeyHex string) (summaries []pc.ProofOfRelay) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, pc.KeyForProofOfRelaysApp(address, appPubKeyHex))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var summary pc.ProofOfRelay
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &summary)
		summaries = append(summaries, summary)
	}
	return
}

func (k Keeper) SetProofOfRelay(ctx sdk.Context, address sdk.ValAddress, summary pc.ProofOfRelay) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(summary)
	store.Set(pc.KeyForProofOfRelay(ctx, address, summary.PORHeader), bz)
}
