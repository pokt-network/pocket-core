package keeper

import (
	"fmt"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/tendermint/tendermint/rpc/client"
)

// auto sends a claim of work based on relays completed
func (k Keeper) SendClaimTx(ctx sdk.Ctx, n client.Client, keybase keys.Keybase, claimTx func(keybase keys.Keybase, cliCtx util.CLIContext, txBuilder auth.TxBuilder, header pc.SessionHeader, totalProofs int64, root pc.HashSum, evidenceType pc.EvidenceType) (*sdk.TxResponse, error)) {
	kp, err := keybase.GetCoinbase()
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("an error occured retrieving the coinbase for the claimTX:\n%v", err))
		return
	}
	iter := pc.EvidenceIterator()
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		evidence := iter.Value()
		evidenceLength := len(evidence.Proofs)
		if evidenceLength == 0 {
			ctx.Logger().Error("evidence of length zero was found in evidence map")
			continue
		}
		evidenceType := evidence.Proofs[0].EvidenceType()
		if evidenceLength < 5 {
			pc.DeleteEvidence(evidence.SessionHeader, evidenceType)
			continue
		}
		// if the blockchain in the evidence is not supported then delete it because nodes don't get paid for unsupported blockchains
		if !k.IsPocketSupportedBlockchain(ctx.WithBlockHeight(evidence.SessionHeader.SessionBlockHeight), evidence.SessionHeader.Chain) && evidence.NumOfProofs > 0 {
			ctx.Logger().Info(fmt.Sprintf("claim for %s blockchain isn't pocket supported, so will not send. Deleting evidence\n", evidence.SessionHeader.Chain))
			pc.DeleteEvidence(evidence.SessionHeader, evidenceType)
			continue
		}
		// check the current state to see if the unverified evidence has already been sent and processed (if so, then skip this evidence)
		ctx.Logger().Info(fmt.Sprintf("get claim for address: %s", kp.GetAddress().String()))
		if _, found := k.GetClaim(ctx, sdk.Address(kp.GetAddress()), evidence.SessionHeader, evidenceType); found {
			continue
		}
		if k.ClaimIsMature(ctx, evidence.SessionBlockHeight) {
			pc.DeleteEvidence(evidence.SessionHeader, evidenceType)
			continue
		}
		// generate the merkle root for this evidence
		root := evidence.GenerateMerkleRoot()
		// generate the auto txbuilder and clictx
		txBuilder, cliCtx, err := newTxBuilderAndCliCtx(ctx, pc.MsgClaimName, n, keybase, k)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("an error occured retrieving the coinbase for the claimTX:\n%v", err))
			return
		}
		// send in the evidence header, the total relays completed, and the merkle root (ensures data integrity)
		if _, err := claimTx(keybase, cliCtx, txBuilder, evidence.SessionHeader, evidence.NumOfProofs, root, evidenceType); err != nil {
			ctx.Logger().Error(fmt.Sprintf("an error occured retrieving the coinbase for the claimTX:\n%v", err))
		}
	}
}

func (k Keeper) ValidateClaim(ctx sdk.Ctx, claim pc.MsgClaim) sdk.Error {
	if claim.EvidenceType == 0 {
		return pc.NewNoEvidenceTypeErr(pc.ModuleName)
	}
	// get the session context
	sessionContext, er := ctx.PrevCtx(claim.SessionBlockHeight)
	if er != nil {
		return sdk.ErrInternal(er.Error())
	}
	// if is not a pocket supported blockchain then return not supported error
	if !k.IsPocketSupportedBlockchain(sessionContext, claim.Chain) {
		return pc.NewChainNotSupportedErr(pc.ModuleName)
	}
	// get the node from the keeper at the time of the session
	node, found := k.GetNode(sessionContext, claim.FromAddress)
	// if not found return not found error
	if !found {
		return pc.NewNodeNotFoundErr(pc.ModuleName)
	}
	// get the application at the session context
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
		var err sdk.Error
		nodes := k.GetAllNodes(ctx)
		session, err = pc.NewSession(claim.SessionHeader, pc.BlockHash(sessionContext), nodes, sessionNodeCount)
		if err != nil {
			ctx.Logger().Error(fmt.Errorf("Could not generate session with public key: %s,  for chain: %s", app.GetPublicKey().RawString(), claim.Chain).Error())
			return err
		}
		// add to cache
		pc.SetSession(session)
	}
	// validate the session
	err := session.Validate(ctx, node, app, sessionNodeCount)
	if err != nil {
		return err
	}
	// check if the proof is ready to be claimed, if it's already ready to be claimed, then it's too late to submit cause the secret is revealed
	if k.ClaimIsMature(ctx, claim.SessionBlockHeight) {
		return pc.NewExpiredProofsSubmissionError(pc.ModuleName)
	}
	return nil
}

func (k Keeper) SetClaim(ctx sdk.Ctx, msg pc.MsgClaim) error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(msg)
	key, err := pc.KeyForClaim(ctx, msg.FromAddress, msg.SessionHeader, msg.EvidenceType)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

func (k Keeper) GetClaim(ctx sdk.Ctx, address sdk.Address, header pc.SessionHeader, evidenceType pc.EvidenceType) (msg pc.MsgClaim, found bool) {
	store := ctx.KVStore(k.storeKey)
	key, err := pc.KeyForClaim(ctx, address, header, evidenceType)
	if err != nil {
		ctx.Logger().Error("an error occured getting the claim:\n", msg)
		return pc.MsgClaim{}, false
	}
	res := store.Get(key)
	if res == nil {
		return pc.MsgClaim{}, false
	}
	k.cdc.MustUnmarshalBinaryBare(res, &msg)
	return msg, true
}

func (k Keeper) SetClaims(ctx sdk.Ctx, claims []pc.MsgClaim) {
	store := ctx.KVStore(k.storeKey)
	for _, msg := range claims {
		bz := k.cdc.MustMarshalBinaryBare(msg)
		key, err := pc.KeyForClaim(ctx, msg.FromAddress, msg.SessionHeader, msg.EvidenceType)
		if err != nil {
			panic(fmt.Sprintf("an error occured setting the claims:\n%v", err))
		}
		store.Set(key, bz)
	}
}

func (k Keeper) GetClaims(ctx sdk.Ctx, address sdk.Address) (proofs []pc.MsgClaim, err error) {
	store := ctx.KVStore(k.storeKey)
	key, err := pc.KeyForClaims(address)
	if err != nil {
		return nil, err
	}
	iterator := sdk.KVStorePrefixIterator(store, key)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var summary pc.MsgClaim
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &summary)
		proofs = append(proofs, summary)
	}
	return
}

func (k Keeper) GetAllClaims(ctx sdk.Ctx) (proofs []pc.MsgClaim) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, pc.ClaimKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var summary pc.MsgClaim
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &summary)
		proofs = append(proofs, summary)
	}
	return
}

func (k Keeper) DeleteClaim(ctx sdk.Ctx, address sdk.Address, header pc.SessionHeader, evidenceType pc.EvidenceType) error {
	store := ctx.KVStore(k.storeKey)
	key, err := pc.KeyForClaim(ctx, address, header, evidenceType)
	if err != nil {
		return err
	}
	store.Delete(key)
	return nil
}

// get the mature unverified proofs for this address
func (k Keeper) GetMatureClaims(ctx sdk.Ctx, address sdk.Address) (matureProofs []pc.MsgClaim, err error) {
	store := ctx.KVStore(k.storeKey)
	key, err := pc.KeyForClaims(address)
	if err != nil {
		return nil, err
	}
	iterator := sdk.KVStorePrefixIterator(store, key)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var msg pc.MsgClaim
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		if k.ClaimIsMature(ctx, msg.SessionBlockHeight) {
			matureProofs = append(matureProofs, msg)
		}
	}
	return
}

// is the claim mature? able to be proved because the `waiting period` has passed since the sessionBlock
func (k Keeper) ClaimIsMature(ctx sdk.Ctx, sessionBlockHeight int64) bool {
	waitingPeriodInBlocks := k.ClaimSubmissionWindow(ctx) * k.SessionFrequency(ctx)
	if ctx.BlockHeight() > waitingPeriodInBlocks+sessionBlockHeight {
		return true
	}
	return false
}

// delete expired claims
func (k Keeper) DeleteExpiredClaims(ctx sdk.Ctx) {
	var msg = pc.MsgClaim{}
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, pc.ClaimKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		sessionContext, err := ctx.PrevCtx(msg.SessionBlockHeight)
		if err != nil {
			ctx.Logger().Error("a claim in the world state had a context that was unable to be retrieved, using current params as expiration")
			sessionContext = ctx.(sdk.Context)
		}
		// if more sessions has passed than the expiration of unverified pseudorandomGenerator, delete from set
		if (ctx.BlockHeight()-msg.SessionBlockHeight)/k.SessionFrequency(sessionContext) >= k.ClaimExpiration(sessionContext) {
			store.Delete(iterator.Key())
		}
	}
}
