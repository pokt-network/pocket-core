package keeper

import (
	"fmt"

	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/tendermint/tendermint/rpc/client"
)

// "SendClaimTx" - Automatically sends a claim of work/challenge based on relays or challenges stored.
func (k Keeper) SendClaimTx(ctx sdk.Ctx, n client.Client, keybase keys.Keybase, claimTx func(pk crypto.PrivateKey, cliCtx util.CLIContext, txBuilder auth.TxBuilder, header pc.SessionHeader, totalProofs int64, root pc.HashSum, evidenceType pc.EvidenceType) (*sdk.TxResponse, error)) {
	// get the private val key (main) account from the keybase
	kp, err := k.GetPKFromFile(ctx)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("an error occured retrieving the private key from file for the claim transaction:\n%s", err.Error()))
		return
	}
	// retrieve the iterator to go through each piece of evidence in storage
	iter := pc.EvidenceIterator()
	defer iter.Close()
	// loop through each evidence
	for ; iter.Valid(); iter.Next() {
		evidence := iter.Value()
		evidenceLength := len(evidence.Proofs)
		// if the number of proofs in the evidence object is zero
		if evidenceLength == 0 {
			ctx.Logger().Error("evidence of length zero was found in evidence storage")
			continue
		}
		// get the type of the first piece of evidence to know if we are dealing with challenge or relays
		evidenceType := evidence.Proofs[0].EvidenceType()
		// if the evidence length is less than 5, it would not satisfy our merkle tree needs
		if evidenceLength < 5 {
			if err := pc.DeleteEvidence(evidence.SessionHeader, evidenceType); err != nil {
				ctx.Logger().Debug(err.Error())
			}
			continue
		}
		// get the session context
		sessionCtx, er := ctx.PrevCtx(evidence.SessionHeader.SessionBlockHeight)
		if er != nil {
			ctx.Logger().Error("could not get sessionCtx")
			continue
		}
		// if the blockchain in the evidence is not supported then delete it because nodes don't get paid/challenged for unsupported blockchains
		if !k.IsPocketSupportedBlockchain(sessionCtx.WithBlockHeight(evidence.SessionHeader.SessionBlockHeight), evidence.SessionHeader.Chain) && evidence.NumOfProofs > 0 {
			ctx.Logger().Info(fmt.Sprintf("claim for %s blockchain isn't pocket supported, so will not send. Deleting evidence\n", evidence.SessionHeader.Chain))
			if err := pc.DeleteEvidence(evidence.SessionHeader, evidenceType); err != nil {
				ctx.Logger().Debug(err.Error())
			}
			continue
		}
		// check the current state to see if the unverified evidence has already been sent and processed (if so, then skip this evidence)
		if _, found := k.GetClaim(ctx, sdk.Address(kp.PublicKey().Address()), evidence.SessionHeader, evidenceType); found {
			continue
		}
		// if the claim is mature, delete it because we cannot submit a mature claim
		if k.ClaimIsMature(ctx, evidence.SessionBlockHeight) {
			if err := pc.DeleteEvidence(evidence.SessionHeader, evidenceType); err != nil {
				ctx.Logger().Debug(err.Error())
			}
			continue
		}
		// generate the merkle root for this evidence
		root := evidence.GenerateMerkleRoot()
		// generate the auto txbuilder and clictx
		txBuilder, cliCtx, err := newTxBuilderAndCliCtx(ctx, pc.MsgClaimName, n, keybase, k)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("an error occured creating the tx builder for the claim tx:\n%s", err.Error()))
			return
		}
		// send in the evidence header, the total relays completed, and the merkle root (ensures data integrity)
		if _, err := claimTx(kp, cliCtx, txBuilder, evidence.SessionHeader, evidence.NumOfProofs, root, evidenceType); err != nil {
			ctx.Logger().Error(fmt.Sprintf("an error occured executing the claim transaciton: \n%s", err.Error()))
		}
	}
}

// "ValidateClaim" - Validates a claim message and returns an sdk error if invalid
func (k Keeper) ValidateClaim(ctx sdk.Ctx, claim pc.MsgClaim) (err sdk.Error) {
	// check to see if evidence type is included in the message
	if claim.EvidenceType == 0 {
		return pc.NewNoEvidenceTypeErr(pc.ModuleName)
	}
	// get the session context (state info at the beginning of the session)
	sessionContext, er := ctx.PrevCtx(claim.SessionBlockHeight)
	if er != nil {
		return sdk.ErrInternal(er.Error())
	}
	// ensure that session ended
	sessionEndHeight := sessionContext.BlockHeight() + k.BlocksPerSession(sessionContext) - 1
	if ctx.BlockHeight() <= sessionEndHeight {
		return pc.NewInvalidBlockHeightError(pc.ModuleName)
	}
	// if is not a pocket supported blockchain then return not supported error
	if !k.IsPocketSupportedBlockchain(sessionContext, claim.Chain) {
		return pc.NewChainNotSupportedErr(pc.ModuleName)
	}
	// get the node from the keeper (at the state of the start of the session)
	node, found := k.GetNode(sessionContext, claim.FromAddress)
	// if not found return not found error
	if !found {
		return pc.NewNodeNotFoundErr(pc.ModuleName)
	}
	// get the application (at the state of the start of the session)
	app, found := k.GetAppFromPublicKey(sessionContext, claim.ApplicationPubKey)
	// if not found return not found error
	if !found {
		return pc.NewAppNotFoundError(pc.ModuleName)
	}
	// get the session node count for the time of the session
	sessionNodeCount := int(k.SessionNodeCount(sessionContext))
	// check cache
	session, found := pc.GetSession(claim.SessionHeader)
	// if not found generate the session
	if !found {
		// use the session end context to ensure that people who were jailed mid session do not get to submit claims
		sessionEndCtx, er := ctx.PrevCtx(sessionEndHeight)
		if er != nil {
			return sdk.ErrInternal("could not get prev context: " + er.Error())
		}
		// create a new session to validate
		session, err = pc.NewSession(sessionContext, sessionEndCtx, k.posKeeper, claim.SessionHeader, pc.BlockHash(sessionContext), sessionNodeCount)
		if err != nil {
			ctx.Logger().Error(fmt.Errorf("could not generate session with public key: %s, for chain: %s", app.GetPublicKey().RawString(), claim.Chain).Error())
			return err
		}
	}
	// validate the session
	err = session.Validate(node, app, sessionNodeCount)
	if err != nil {
		return err
	}
	// check if the proof is ready to be claimed, if it's already ready to be claimed, then it's too late to submit cause the secret is revealed
	if k.ClaimIsMature(ctx, claim.SessionBlockHeight) {
		return pc.NewExpiredProofsSubmissionError(pc.ModuleName)
	}
	return nil
}

// "SetClaim" - Sets the claim message in the state storage
func (k Keeper) SetClaim(ctx sdk.Ctx, msg pc.MsgClaim) error {
	// retrieve the store
	store := ctx.KVStore(k.storeKey)
	// generate the store key
	key, err := pc.KeyForClaim(ctx, msg.FromAddress, msg.SessionHeader, msg.EvidenceType)
	if err != nil {
		return err
	}
	// generate the expiration height upon setting
	if msg.ExpirationHeight == 0 {
		sessionCtx, err := ctx.PrevCtx(msg.SessionBlockHeight)
		if err != nil {
			return err
		}
		msg.ExpirationHeight = ctx.BlockHeight() + k.ClaimExpiration(sessionCtx)*k.BlocksPerSession(sessionCtx)
	}
	// marshal the message into amino
	bz := k.cdc.MustMarshalBinaryBare(msg)
	// set in the store
	store.Set(key, bz)
	return nil
}

// "GetClaim" - Retrieves the claim message from the store, requires the evidence type and header to return the proper claim message
func (k Keeper) GetClaim(ctx sdk.Ctx, address sdk.Address, header pc.SessionHeader, evidenceType pc.EvidenceType) (msg pc.MsgClaim, found bool) {
	// retrieve the store
	store := ctx.KVStore(k.storeKey)
	// generate the store key
	key, err := pc.KeyForClaim(ctx, address, header, evidenceType)
	if err != nil {
		ctx.Logger().Error("an error occured getting the claim:\n", msg)
		return pc.MsgClaim{}, false
	}
	// get the claim msg from the store
	res := store.Get(key)
	if res == nil {
		return pc.MsgClaim{}, false
	}
	// unmarshal into message object
	k.cdc.MustUnmarshalBinaryBare(res, &msg)
	// return the object
	return msg, true
}

// "SetClaims" - Sets all the claim messages in the state storage.
// (Needed for genesis initializing)
func (k Keeper) SetClaims(ctx sdk.Ctx, claims []pc.MsgClaim) {
	// loop through all of the claim messages one by one and set them
	for _, msg := range claims {
		err := k.SetClaim(ctx, msg)
		if err != nil {
			ctx.Logger().Error("an error occurred setting the claim:\n", msg)
		}
	}
}

// "GetClaims" - Gets all of the claim messages in the state storage for an address
func (k Keeper) GetClaims(ctx sdk.Ctx, address sdk.Address) (claims []pc.MsgClaim, err error) {
	// retrieve the store
	store := ctx.KVStore(k.storeKey)
	// generate the key for the claims
	key, err := pc.KeyForClaims(address)
	if err != nil {
		return nil, err
	}
	// iterate through all of the kv pairs and unmarshal into claim objects
	iterator := sdk.KVStorePrefixIterator(store, key)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var claim pc.MsgClaim
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &claim)
		claims = append(claims, claim)
	}
	return
}

// "GetAllClaims" - Gets all of the claim messages held in the state storage.
func (k Keeper) GetAllClaims(ctx sdk.Ctx) (claims []pc.MsgClaim) {
	// retrieve the store
	store := ctx.KVStore(k.storeKey)
	// iterate through the kv in the state and unmarshal into claim objects
	iterator := sdk.KVStorePrefixIterator(store, pc.ClaimKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var claim pc.MsgClaim
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &claim)
		claims = append(claims, claim)
	}
	return
}

// "DeleteClaim" - Removes a claim object for a certain key
func (k Keeper) DeleteClaim(ctx sdk.Ctx, address sdk.Address, header pc.SessionHeader, evidenceType pc.EvidenceType) error {
	// retrieve the store
	store := ctx.KVStore(k.storeKey)
	// generate the key for the claim
	key, err := pc.KeyForClaim(ctx, address, header, evidenceType)
	if err != nil {
		return err
	}
	// delete it from the state storage
	store.Delete(key)
	return nil
}

// "GetMatureClaims" - Returns the mature (ready to be proved, past its security waiting period)
func (k Keeper) GetMatureClaims(ctx sdk.Ctx, address sdk.Address) (matureProofs []pc.MsgClaim, err error) {
	// retrieve the store
	store := ctx.KVStore(k.storeKey)
	// generate the key for the claim
	key, err := pc.KeyForClaims(address)
	if err != nil {
		return nil, err
	}
	// iterate through all kv and see if the claim is mature for each
	iterator := sdk.KVStorePrefixIterator(store, key)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var msg pc.MsgClaim
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		// if the claim is mature, add it to the list
		if k.ClaimIsMature(ctx, msg.SessionBlockHeight) {
			matureProofs = append(matureProofs, msg)
		}
	}
	return
}

// "ClaimIsMature" - Returns if the claim is past its security waiting period
func (k Keeper) ClaimIsMature(ctx sdk.Ctx, sessionBlockHeight int64) bool {
	waitingPeriodInBlocks := k.ClaimSubmissionWindow(ctx) * k.BlocksPerSession(ctx)
	if ctx.BlockHeight() > waitingPeriodInBlocks+sessionBlockHeight {
		return true
	}
	return false
}

// "DeleteExpiredClaims" - Deletes the expired (claim expiration > # of session passed since claim genesis) claims
func (k Keeper) DeleteExpiredClaims(ctx sdk.Ctx) {
	var msg = pc.MsgClaim{}
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, pc.ClaimKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		// if more sessions has passed than the expiration of the claim's genesis, delete it from the set
		if msg.ExpirationHeight <= ctx.BlockHeight() {
			store.Delete(iterator.Key())
		}
	}
}
