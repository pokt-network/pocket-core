package keeper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/rpc/client"
	"math"
	"strconv"
)

func BeginBlocker(ctx sdk.Ctx, _ abci.RequestBeginBlock, k Keeper) {
	// delete the proofs held within the world state for too long
	k.DeleteExpiredClaims(ctx)
}

// validate the zero knowledge range proof using the proof message and the claim message
func (k Keeper) ValidateProof(ctx sdk.Ctx, claimMsg pc.MsgClaim, proofMsg pc.MsgProof) error {
	// generate the needed pseudorandom claimMsg index
	ctx.Logger().Info(fmt.Sprintf("Generate psuedorandom proof with %d, at session height of %d", claimMsg.TotalRelays, claimMsg.SessionBlockHeight))
	reqProof, err := k.GetPseudorandomIndex(ctx, claimMsg.TotalRelays, claimMsg.SessionHeader)
	if err != nil {
		return err
	}
	// if the required proof message index does not match the leaf node index
	if reqProof != int64(proofMsg.MerkleProofs[0].Index) {
		return pc.NewInvalidProofsError(pc.ModuleName)
	}
	// validate level count on claimMsg by total relays
	levelCount := len(proofMsg.MerkleProofs[0].HashSums)
	if levelCount != int(math.Ceil(math.Log2(float64(claimMsg.TotalRelays)))) {
		return pc.NewInvalidProofsError(pc.ModuleName)
	}
	if !proofMsg.MerkleProofs.Validate(claimMsg.MerkleRoot, proofMsg.Leaf, proofMsg.Cousin, claimMsg.TotalRelays) {
		return pc.NewInvalidMerkleVerifyError(pc.ModuleName)
	}
	// check the validity of the token
	if err := proofMsg.Leaf.Token.Validate(); err != nil {
		return err
	}
	// verify the client signature
	if err := pc.SignatureVerification(proofMsg.Leaf.Token.ClientPublicKey, proofMsg.Leaf.HashString(), proofMsg.Leaf.Signature); err != nil {
		return err
	}
	return nil
}

// struct used for creating the psuedorandom index
type pseudorandomGenerator struct {
	BlockHash string
	Header    string
}

// generates the required pseudorandom index for the zero knowledge proof
func (k Keeper) GetPseudorandomIndex(ctx sdk.Ctx, totalRelays int64, header pc.SessionHeader) (int64, error) {
	// get the context for the proof (the proof context is X sessions after the session began)
	proofContext := ctx.MustGetPrevCtx(header.SessionBlockHeight + k.ClaimSubmissionWindow(ctx)*k.SessionFrequency(ctx)) // next session block hash
	// get the pseudorandomGenerator json bytes
	proofBlockHeader := proofContext.BlockHeader()
	r, err := json.Marshal(pseudorandomGenerator{
		blockHash: hex.EncodeToString(proofBlockHeader.GetLastBlockId().Hash), // block hash
		header:    header.HashString(),                                        // header hashstring
	})
	if err != nil {
		return 0, err
	}
	// hash the bytes and take the first 15 characters of the string
	proofsHash := hex.EncodeToString(pc.Hash(r))[:15]
	var maxValue int64
	// for each hex character of the hash
	for i := 15; i > 0; i-- {
		// parse the integer from this point of the hex string onward
		maxValue, err = strconv.ParseInt(string(proofsHash[:i]), 16, 64)
		if err != nil {
			return 0, err

		}
		// if the total relays is greater than the resulting integer, this is the pseudorandom chosen proof
		if totalRelays > maxValue {
			firstCharacter, err := strconv.ParseInt(string(proofsHash[0]), 16, 64)
			if err != nil {
				return 0, err
			}
			selection := firstCharacter%int64(i) + 1
			// parse the integer from this point of the hex string onward
			index, err := strconv.ParseInt(proofsHash[:selection], 16, 64)
			if err != nil {
				return 0, err
			}
			return index, err
		}
	}
	return 0, nil
}

// auto sends a proof transaction for the claim
func (k Keeper) SendProofTx(ctx sdk.Ctx, n client.Client, keybase keys.Keybase, proofTx func(cliCtx util.CLIContext, txBuilder auth.TxBuilder, branches [2]pc.MerkleProof, leafNode, cousin pc.Proof) (*sdk.TxResponse, error)) {
	kp, err := keybase.GetCoinbase()
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("an error occured retrieving the coinbase for the ProofTX:\n%v", err))
		return
	}
	// get the self address
	addr := kp.GetAddress()
	// get all mature (waiting period has passed) claims for your address
	claims, err := k.GetMatureClaims(ctx, addr)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("an error occured getting the mature claims in the ProofTX:\n%v", err))
		return
	}
	// for every claim of the mature set
	for _, claim := range claims {
		// if the claim is found to be verified in the world state, you can delete it from the cache and not send again
		if _, found := k.GetReceipt(ctx, addr, claim.SessionHeader); found {
			// remove from the local cache
			pc.GetEvidenceMap().DeleteEvidence(claim.SessionHeader)
			continue
		}
		// check to see if evidence is stored in cache
		evidence, found := pc.GetEvidenceMap().GetEvidence(claim.SessionHeader)
		if !found || evidence.Proofs == nil || len(evidence.Proofs) == 0 {
			ctx.Logger().Info(fmt.Sprintf("the evidence object for evidence is not found, ignoring pending claim"))
			continue
		}
		// generate the needed pseudorandom index using the information found in the first transaction
		index, err := k.GetPseudorandomIndex(ctx, claim.TotalRelays, claim.SessionHeader)
		if err != nil {
			ctx.Logger().Error(err.Error())
		}
		// get the merkle proof object for the pseudorandom index
		branch, cousinIndex := evidence.GenerateMerkleProof(int(index))
		// get the leaf and cousin for the required pseudorandom index
		leaf := pc.GetEvidenceMap().GetProof(claim.SessionHeader, int(index))
		cousin := pc.GetEvidenceMap().GetProof(claim.SessionHeader, cousinIndex)
		// generate the auto txbuilder and clictx
		txBuilder, cliCtx, err := newTxBuilderAndCliCtx(ctx, pc.MsgProofName, n, keybase, k)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("an error occured in the transaction process of the ProofTX:\n%v", err))
			return
		}
		// send the claim TX
		_, err = proofTx(cliCtx, txBuilder, branch, leaf, cousin)
		if err != nil {
			ctx.Logger().Error(err.Error())
		}
	}
}

func newTxBuilderAndCliCtx(ctx sdk.Ctx, msgType string, n client.Client, keybase keys.Keybase, k Keeper) (txBuilder auth.TxBuilder, cliCtx util.CLIContext, err error) {
	// get the coinbase, as it is the sender of the automatic message
	kp, err := keybase.GetCoinbase()
	if err != nil {
		return txBuilder, cliCtx, err
	}
	// get the from address from the coinbase
	fromAddr := kp.GetAddress()
	// get the genesis doc from the node for the chainID
	genDoc, err := n.Genesis()
	if err != nil {
		return txBuilder, cliCtx, err
	}
	// create a client context for sending
	cliCtx = util.NewCLIContext(n, fromAddr, k.coinbasePassphrase).WithCodec(k.cdc)
	// broadcast synchronously
	cliCtx.BroadcastMode = util.BroadcastSync
	// get the account to ensure balance
	accGetter := auth.NewAccountRetriever(cliCtx)
	// retrieve the account for a balance check (and ensure it exists)
	account, err := accGetter.GetAccount(fromAddr)
	if err != nil {
		return txBuilder, cliCtx, err
	}
	// check the fee amount
	fee := sdk.NewInt(pc.PocketFeeMap[msgType])
	if account.GetCoins().AmountOf(k.posKeeper.StakeDenom(ctx)).LTE(fee) {
		ctx.Logger().Error(fmt.Sprintf("insufficient funds for the auto %s transaction: the fee needed is %v ", msgType, fee))
	}
	// ensure that the tx builder has the correct tx encoder, chainID, fee, and keybase
	txBuilder = auth.NewTxBuilder(
		auth.DefaultTxEncoder(k.cdc),
		genDoc.Genesis.ChainID,
		"",
		sdk.NewCoins(sdk.NewCoin(k.posKeeper.StakeDenom(ctx), fee)),
	).WithKeybase(keybase)
	return
}
