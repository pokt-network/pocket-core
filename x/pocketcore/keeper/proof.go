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
		if _, found := k.GetReceipt(ctx, addr, claim.SessionHeader, claim.EvidenceType); found {
			// remove from the local cache
			pc.GetEvidenceMap().DeleteEvidence(claim.SessionHeader, claim.EvidenceType)
			continue
		}
		// check to see if evidence is stored in cache
		evidence, found := pc.GetEvidenceMap().GetEvidence(claim.SessionHeader, claim.EvidenceType)
		if !found || evidence.Proofs == nil || len(evidence.Proofs) == 0 {
			ctx.Logger().Info(fmt.Sprintf("the evidence object for evidence is not found, ignoring pending claim for app: %s, at sessionHeight: %d", claim.ApplicationPubKey, claim.SessionBlockHeight))
			continue
		}
		// generate the needed pseudorandom index using the information found in the first transaction
		index, err := k.getPseudorandomIndex(ctx, claim.TotalProofs, claim.SessionHeader)
		if err != nil {
			ctx.Logger().Error(err.Error())
			continue
		}
		// get the merkle proof object for the pseudorandom index
		branch, cousinIndex := evidence.GenerateMerkleProof(int(index))
		// get the leaf and cousin for the required pseudorandom index
		leaf := pc.GetEvidenceMap().GetProof(claim.SessionHeader, claim.EvidenceType, int(index))
		cousin := pc.GetEvidenceMap().GetProof(claim.SessionHeader, claim.EvidenceType, cousinIndex)
		// generate the auto txbuilder and clictx
		txBuilder, cliCtx, err := newTxBuilderAndCliCtx(ctx, pc.MsgProofName, n, keybase, k)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("an error occured in the transaction process of the ProofTX:\n%v", err))
			return
		}
		// send the proof TX
		_, err = proofTx(cliCtx, txBuilder, branch, leaf, cousin)
		if err != nil {
			ctx.Logger().Error(err.Error())
		}
	}
}

func (k Keeper) ValidateProof(ctx sdk.Ctx, proof pc.MsgProof) (servicerAddr sdk.Address, claim pc.MsgClaim, sdkError sdk.Error) {
	// get the public key from the claim
	addrs := proof.GetSigners()
	if len(addrs) < 1 {
		return nil, pc.MsgClaim{}, pc.NewEmptyAddressError(pc.ModuleName)
	}
	addr := addrs[0]
	// get the claim for the address
	claim, found := k.GetClaim(ctx, addr, proof.Leaf.SessionHeader(), proof.Leaf.EvidenceType())
	// if the claim is not found for this claim
	if !found {
		return nil, pc.MsgClaim{}, pc.NewClaimNotFoundError(pc.ModuleName)
	}
	// validate the proof
	ctx.Logger().Info(fmt.Sprintf("Generate psuedorandom proof with %d proofs, at session height of %d, for app: %s", claim.TotalProofs, claim.SessionBlockHeight, claim.ApplicationPubKey))
	reqProof, err := k.getPseudorandomIndex(ctx, claim.TotalProofs, claim.SessionHeader)
	if err != nil {
		return nil, pc.MsgClaim{}, sdk.ErrInternal(err.Error())
	}
	// if the required proof message index does not match the leaf node index
	if reqProof != int64(proof.MerkleProofs[0].Index) {
		return nil, pc.MsgClaim{}, pc.NewInvalidProofsError(pc.ModuleName)
	}
	// validate level count on claim by total relays
	levelCount := len(proof.MerkleProofs[0].HashSums)
	if levelCount != int(math.Ceil(math.Log2(float64(claim.TotalProofs)))) {
		return nil, pc.MsgClaim{}, pc.NewInvalidProofsError(pc.ModuleName)
	}
	// validate the merkle proof
	if !proof.MerkleProofs.Validate(claim.MerkleRoot, proof.Leaf, proof.Cousin, claim.TotalProofs) {
		return nil, pc.MsgClaim{}, pc.NewInvalidMerkleVerifyError(pc.ModuleName)
	}
	// get the session context
	sessionCtx, err := ctx.PrevCtx(claim.SessionBlockHeight)
	if err != nil {
		return nil, pc.MsgClaim{}, sdk.ErrInternal(err.Error())
	}
	// get the application
	application, found := k.GetAppFromPublicKey(ctx, claim.ApplicationPubKey)
	if !found {
		return nil, pc.MsgClaim{}, pc.NewAppNotFoundError(pc.ModuleName)
	}
	// validate the proof depending on the type of proof it is
	er := proof.Leaf.Validate(application.GetChains(), int(k.SessionNodeCount(sessionCtx)), claim.SessionBlockHeight)
	if er != nil {
		return nil, pc.MsgClaim{}, er
	}
	// return the needed info to the handler
	return addr, claim, nil
}

func (k Keeper) ExecuteProof(ctx sdk.Ctx, proof pc.MsgProof, claim pc.MsgClaim) sdk.Error {
	switch proof.Leaf.(type) {
	case pc.RelayProof:
		ctx.Logger().Info(fmt.Sprintf("reward coins to %s, for %d relays", claim.FromAddress.String(), claim.TotalProofs))
		k.AwardCoinsForRelays(ctx, claim.TotalProofs, claim.FromAddress)
		err := k.DeleteClaim(ctx, claim.FromAddress, claim.SessionHeader, pc.RelayEvidence)
		if err != nil {
			return sdk.ErrInternal(err.Error())
		}
	case pc.ChallengeProofInvalidData:
		ctx.Logger().Info(fmt.Sprintf("burning coins from %s, for %d valid challenges", claim.FromAddress.String(), claim.TotalProofs))
		k.BurnCoinsForChallenges(ctx, claim.TotalProofs, claim.FromAddress)
		err := k.DeleteClaim(ctx, claim.FromAddress, claim.SessionHeader, pc.ChallengeEvidence)
		if err != nil {
			return sdk.ErrInternal(err.Error())
		}
	}
	return nil
}

// struct used for creating the psuedorandom index
type pseudorandomGenerator struct {
	BlockHash string
	Header    string
}

// generates the required pseudorandom index for the zero knowledge proof
func (k Keeper) getPseudorandomIndex(ctx sdk.Ctx, totalRelays int64, header pc.SessionHeader) (int64, error) {
	// get the context for the proof (the proof context is X sessions after the session began)
	proofContext, err := ctx.PrevCtx(header.SessionBlockHeight + k.ClaimSubmissionWindow(ctx)*k.SessionFrequency(ctx)) // next session block hash
	if err != nil {
		return 0, err
	}
	// get the pseudorandomGenerator json bytes
	proofBlockHeader := proofContext.BlockHeader()
	blockHash := hex.EncodeToString(proofBlockHeader.GetLastBlockId().Hash)
	headerHash := header.HashString()
	pseudoGenerator := pseudorandomGenerator{blockHash, headerHash}
	r, err := json.Marshal(pseudoGenerator)
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

// todo move this once password management is fixed
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
