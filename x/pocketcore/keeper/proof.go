package keeper

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	nodeexported "github.com/pokt-network/pocket-core/x/nodes/exported"
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

type Proof pc.ProofBatch

// todo possible attacks around tricking the client and making the requests after the nexSessionBlockHash has been revealed?
func (k Keeper) GenerateProofs(relaysCompleted uint64, nextSessionBlockHash, sessionKey string) []uint64 { // todo created on the spot! need to audit. Possible to brute force?
	var result []uint64
	proofsHash := hex.EncodeToString(crypto.SHA3FromString(nextSessionBlockHash + sessionKey))
	proofsHash = proofsHash[:16] // take first 16 characters to fit int 64
	for i := 0; i < len(proofsHash); i++ {
		res, err := strconv.ParseUint(proofsHash[i:], 16, 64)
		if err != nil {
			panic(err)
		}
		if relaysCompleted > res {
			result = append(result, res)
		}
	}
	return result
}

func (p Proof) Validate(k Keeper, ctx sdk.Context, nodeVerify nodeexported.ValidatorI, allNodes []nodeexported.ValidatorI, maxNumRelays int, neededProofs map[int]struct{}) error {
	// verify that the node was a part of the session
	err := k.SessionVerification(ctx, nodeVerify, p.ApplicationPubKey, p.Chain, p.SessionBlockHash, allNodes)
	if err != nil {
		return err
	}
	// validate the the proper proofs were sent and their validity
	for _, proof := range p.Proofs {
		if _, contains := neededProofs[proof.Counter]; !contains {
			return pc.InvalidICError
		} else {
			delete(neededProofs, proof.Counter)
		}
		err := proof.Token.Validate()
		if err != nil {
			return err
		}
	}
	if len(neededProofs) != 0 {
		return pc.NotEveryICProvidedError
	}
	return nil
}
