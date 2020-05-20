package keeper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"

	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/tendermint/tendermint/rpc/client"
)

// auto sends a proof transaction for the claim
func (k Keeper) SendProofTx(ctx sdk.Ctx, n client.Client, keybase keys.Keybase, proofTx func(cliCtx util.CLIContext, txBuilder auth.TxBuilder, branches [2]pc.MerkleProof, leafNode, cousin pc.Proof, evidenceType pc.EvidenceType) (*sdk.TxResponse, error)) {
	kp, err := k.GetPKFromFile(ctx)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("an error occured retrieving the pk from the file for the Proof Transaction:\n%v", err))
		return
	}
	// get the self address
	addr := sdk.Address(kp.PublicKey().Address())
	// get all mature (waiting period has passed) claims for your address
	claims, err := k.GetMatureClaims(ctx, addr)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("an error occured getting the mature claims in the Proof Transaction:\n%v", err))
		return
	}
	// for every claim of the mature set
	for _, claim := range claims {
		// if the claim is found to be verified in the world state, you can delete it from the cache and not send again
		if _, found := k.GetReceipt(ctx, addr, claim.SessionHeader, claim.EvidenceType); found {
			// remove from the local cache
			if err := pc.DeleteEvidence(claim.SessionHeader, claim.EvidenceType); err != nil {
				ctx.Logger().Debug(err.Error())
			}
			continue
		}
		// check to see if evidence is stored in cache
		evidence, err := pc.GetEvidence(claim.SessionHeader, claim.EvidenceType, sdk.ZeroInt())
		if err != nil || evidence.Proofs == nil || len(evidence.Proofs) == 0 {
			ctx.Logger().Info(fmt.Sprintf("the evidence object for evidence is not found, ignoring pending claim for app: %s, at sessionHeight: %d", claim.ApplicationPubKey, claim.SessionBlockHeight))
			continue
		}
		// get the session context
		sessionCtx, err := ctx.PrevCtx(claim.SessionBlockHeight)
		if err != nil {
			ctx.Logger().Info(fmt.Sprintf("could not get Session Context, ignoring pending claim for app: %s, at sessionHeight: %d", claim.ApplicationPubKey, claim.SessionBlockHeight))
			continue
		}
		// generate the needed pseudorandom index using the information found in the first transaction
		index, err := k.getPseudorandomIndex(ctx, claim.TotalProofs, claim.SessionHeader, sessionCtx)
		if err != nil {
			ctx.Logger().Error(err.Error())
			continue
		}
		// get the merkle proof object for the pseudorandom index
		branch, cousinIndex := evidence.GenerateMerkleProof(int(index))
		// get the leaf and cousin for the required pseudorandom index
		leaf := pc.GetProof(claim.SessionHeader, claim.EvidenceType, index)
		cousin := pc.GetProof(claim.SessionHeader, claim.EvidenceType, int64(cousinIndex))
		// generate the auto txbuilder and clictx
		txBuilder, cliCtx, err := newTxBuilderAndCliCtx(ctx, pc.MsgProofName, n, keybase, k)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("an error occured in the transaction process of the Proof Transaction:\n%v", err))
			return
		}
		// send the proof TX
		_, err = proofTx(cliCtx, txBuilder, branch, leaf, cousin, evidence.EvidenceType)
		if err != nil {
			ctx.Logger().Error(err.Error())
		}
	}
}

func (k Keeper) ValidateProof(ctx sdk.Ctx, proof pc.MsgProof) (servicerAddr sdk.Address, claim pc.MsgClaim, sdkError sdk.Error) {
	// get the public key from the claim
	addr := proof.GetSigner()
	// get the claim for the address
	claim, found := k.GetClaim(ctx, addr, proof.Leaf.SessionHeader(), proof.EvidenceType)
	// if the claim is not found for this claim
	if !found {
		return nil, pc.MsgClaim{}, pc.NewClaimNotFoundError(pc.ModuleName)
	}
	// get the session context
	sessionCtx, err := ctx.PrevCtx(claim.SessionBlockHeight)
	if err != nil {
		return nil, pc.MsgClaim{}, sdk.ErrInternal(err.Error())
	}
	// validate the proof
	ctx.Logger().Info(fmt.Sprintf("Generate psuedorandom proof with %d proofs, at session height of %d, for app: %s", claim.TotalProofs, claim.SessionBlockHeight, claim.ApplicationPubKey))
	reqProof, err := k.getPseudorandomIndex(ctx, claim.TotalProofs, claim.SessionHeader, sessionCtx)
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
	// validate number of proofs
	if params := k.GetParams(ctx); claim.TotalProofs < params.MinimumNumberOfProofs {
		return nil, pc.MsgClaim{}, pc.NewInvalidMerkleVerifyError(pc.ModuleName)
	}
	// validate the merkle proofs
	isValid, isReplayAttack := proof.MerkleProofs.Validate(claim.MerkleRoot, proof.Leaf, proof.Cousin, claim.TotalProofs)
	// if the merkle proof is a replay attack
	if isReplayAttack {
		return addr, claim, pc.NewReplayAttackError(pc.ModuleName)
	}
	// if is not valid for other reasons
	if !isValid {
		return nil, pc.MsgClaim{}, pc.NewInvalidMerkleVerifyError(pc.ModuleName)
	}
	// get the application
	application, found := k.GetAppFromPublicKey(sessionCtx, claim.ApplicationPubKey)
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
		proof, ok := proof.Leaf.(pc.ChallengeProofInvalidData)
		if !ok {
			return pc.NewInvalidProofsError(pc.ModuleName)
		}
		pk := proof.MinorityResponse.Proof.ServicerPubKey
		pubKey, err := crypto.NewPublicKey(pk)
		if err != nil {
			return sdk.ErrInvalidPubKey(err.Error())
		}
		k.BurnCoinsForChallenges(ctx, claim.TotalProofs, sdk.Address(pubKey.Address()))
		err = k.DeleteClaim(ctx, claim.FromAddress, claim.SessionHeader, pc.ChallengeEvidence)
		if err != nil {
			return sdk.ErrInternal(err.Error())
		}
		// small reward for the challenge proof invalid data
		k.AwardCoinsForRelays(ctx, claim.TotalProofs/100, claim.FromAddress)
	}
	return nil
}

// struct used for creating the psuedorandom index
type pseudorandomGenerator struct {
	BlockHash string
	Header    string
}

// generates the required pseudorandom index for the zero knowledge proof
func (k Keeper) getPseudorandomIndex(ctx sdk.Ctx, totalRelays int64, header pc.SessionHeader, sessionCtx sdk.Ctx) (int64, error) {
	// get the context for the proof (the proof context is X sessions after the session began)
	proofContext, err := ctx.PrevCtx(header.SessionBlockHeight + k.ClaimSubmissionWindow(sessionCtx)*k.BlocksPerSession(sessionCtx)) // next session block hash
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
	return pc.PseudoRandomGeneration(totalRelays, r)
	//return 0, nil
}

func (k Keeper) HandleReplayAttack(ctx sdk.Ctx, address sdk.Address, numberOfChallenges sdk.Int) {
	ctx.Logger().Error(fmt.Sprintf("Replay Attack Detected: By %s, for %v proofs", address.String(), numberOfChallenges))
	k.posKeeper.BurnForChallenge(ctx, numberOfChallenges.Mul(sdk.NewInt(k.ReplayAttackBurnMultiplier(ctx))), address)
}

func newTxBuilderAndCliCtx(ctx sdk.Ctx, msgType string, n client.Client, keybase keys.Keybase, k Keeper) (txBuilder auth.TxBuilder, cliCtx util.CLIContext, err error) {
	// get the pk, as it is the sender of the automatic message
	kp, err := k.GetPKFromFile(ctx)
	if err != nil {
		return txBuilder, cliCtx, err
	}
	// get the from address from the pkf
	fromAddr := sdk.Address(kp.PublicKey().Address())
	// get the genesis doc from the node for the chainID
	genDoc, err := n.Genesis()
	if err != nil {
		return txBuilder, cliCtx, err
	}
	// create a client context for sending
	cliCtx = util.NewCLIContext(n, fromAddr, "").WithCodec(k.cdc)
	pk, err := k.GetPKFromFile(ctx)
	if err != nil {
		return txBuilder, cliCtx, err
	}
	cliCtx.PrivateKey = pk
	// broadcast synchronously
	cliCtx.BroadcastMode = util.BroadcastSync
	// get the account to ensure balance
	// retrieve the account for a balance check (and ensure it exists)
	account, err := cliCtx.GetAccount(fromAddr)
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
		auth.DefaultTxDecoder(k.cdc),
		genDoc.Genesis.ChainID,
		"",
		sdk.NewCoins(sdk.NewCoin(k.posKeeper.StakeDenom(ctx), fee)),
	).WithKeybase(keybase)
	return
}
