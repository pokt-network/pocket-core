package keeper

import (
	"encoding/hex"
	"encoding/json"
	merkle "github.com/pokt-network/merkle"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/tendermint/tendermint/node"
	"strconv"
)

// validate the zero knowledge range proof using the proof message and the claim message
func (k Keeper) ValidateProof(ctx sdk.Context, proof pc.MsgClaim, claim pc.MsgProof) error {
	// generate the needed pseudorandom proof index
	reqProof := k.GeneratePseudoRandomProof(ctx, proof.TotalRelays, proof.SessionHeader)
	// if the required proof index does not match the claim leafNode index
	if reqProof != claim.LeafNode.Index {
		return pc.NewInvalidProofsError(pc.ModuleName)
	}
	// do a merkle proof using the merkle proof, the previously submitted root, and the leafNode to ensure validity of the claim
	if !(merkle.NewTree(pc.Hasher.New()).VerifyProof(merkle.Proof(claim.MerkleProof), proof.Root, claim.LeafNode.Hash())) {
		return pc.NewInvalidMerkleVerifyError(pc.ModuleName)
	}
	// check the validity of the token
	if err := claim.LeafNode.Token.Validate(); err != nil {
		return err
	}
	// verify the client signature
	if err := pc.SignatureVerification(claim.LeafNode.Token.ClientPublicKey, claim.LeafNode.HashString(), claim.LeafNode.Signature); err != nil {
		return err
	}
	return nil
}

// generates the required pseudorandom index for the zero knowledge proof
func (k Keeper) GeneratePseudoRandomProof(ctx sdk.Context, totalRelays int64, header pc.SessionHeader) int64 {
	// get the context for the proof (the proof context is X sessions after the session began)
	proofContext := ctx.WithBlockHeight(header.SessionBlockHeight + int64(k.ProofWaitingPeriod(ctx))*k.SessionFrequency(ctx)) // next session block hash
	// get the pseudorandomGenerator json bytes
	r, err := json.Marshal(pseudorandomGenerator{
		blockHash: hex.EncodeToString(proofContext.BlockHeader().GetLastBlockId().Hash), // block hash
		header:    header.HashString(),                                                  // header hashstring
	})
	if err != nil {
		panic(err)
	}
	// hash the bytes and take the first 16 bytes of the string
	proofsHash := hex.EncodeToString(pc.Hash(r))[:16]
	// for each hex character of the hash
	for i := 0; i < 16; i++ {
		// parse the integer from this point of the hex string onward
		res, err := strconv.ParseInt(proofsHash[i:], 16, 64)
		if err != nil {
			panic(err)
		}
		// if the total relays is greater than the resulting integer, this is the pseudorandom chosen proof
		if totalRelays > res {
			// todo this leans towards the end
			return res
		}
	}
	return 0
}

// struct used for creating the psuedorandom index
type pseudorandomGenerator struct {
	blockHash string
	header    string
}

// auto sends a claim of work based on relays completed
func (k Keeper) SendClaimTx(ctx sdk.Context, n *node.Node, proofTx func(cdc *codec.Codec, cliCtx util.CLIContext, txBuilder auth.TxBuilder, header pc.SessionHeader, totalRelays int64, root []byte) error) { // todo should move tx to keeper?
	// get all the proofs held in memory
	proofs := pc.GetAllProofs()
	// for every proof in AllProofs
	for _, proof := range (*proofs).M {
		// if the blockchain in the proof is not supported then delete it because nodes don't get paid for unsupported blockchains
		if !k.IsPocketSupportedBlockchain(ctx.WithBlockHeight(proof.SessionBlockHeight), proof.Chain) && proof.TotalRelays > 0 {
			proofs.DeleteProofs(proof.SessionHeader)
			continue
		}
		// check the current state to see if the unverified proof has already been sent and processed (if so, then skip this proof)
		if _, found := k.GetUnverfiedProof(ctx, sdk.ValAddress(n.PrivValidator().GetPubKey().Address()), proof.SessionHeader); found {
			continue
		}
		// generate the auto txbuilder and clictx
		txBuilder, cliCtx := newTxBuilderAndCliCtx(ctx, n, k)
		// generate the merkle root for this proof
		root, sdkErr := proof.Tree.GetMerkleRoot()
		if sdkErr != nil {
			panic(sdkErr)
		}
		// send in the proof header, the total relays completed, and the merkle root (ensures data integrity)
		if err := proofTx(k.cdc, cliCtx, txBuilder, proof.SessionHeader, proof.TotalRelays, root); err != nil {
			panic(err)
		}
	}
}

// auto sends a proof transaction for the claim
func (k Keeper) SendProofTx(ctx sdk.Context, n *node.Node, claimTx func(cdc *codec.Codec, cliCtx util.CLIContext, txBuilder auth.TxBuilder, porBranch pc.MerkleProof, leafNode pc.Proof) error) {
	// get the self address
	addr := sdk.ValAddress(n.PrivValidator().GetPubKey().Address())
	// get all mature (waiting period has passed) proofs for your address
	proofs := k.GetMatureClaims(ctx, addr)
	// for every proof of the mature set
	for _, proof := range proofs {
		// if the proof is found to be verified in the world state, you can delete it from the cache and not send again
		if _, found := k.GetProof(ctx, addr, proof.SessionHeader); found {
			// remove from the local cache
			pc.GetAllProofs().DeleteProofs(proof.SessionHeader)
			// remove from the temporary world state
			k.DeleteClaim(ctx, addr, proof.SessionHeader)
			continue
		}
		// generate the auto txbuilder and clictx
		txBuilder, cliCtx := newTxBuilderAndCliCtx(ctx, n, k)
		// generate the proof of relay object using the found proof and local cache
		por := pc.ProofOfRelay{
			SessionHeader: proof.SessionHeader,
			TotalRelays:   proof.TotalRelays,
			Proofs:        pc.GetAllProofs().GetProofs(proof.SessionHeader),
			Tree:          pc.Tree(merkle.NewTree(pc.Hasher.New())),
		}
		// generate the needed pseudorandom proof using the information found in the first transaction
		reqProof := int(k.GeneratePseudoRandomProof(ctx, proof.TotalRelays, proof.SessionHeader))
		// generate the merkle tree from the por structure
		por.GenerateMerkleTree()
		// get the merkle proof object for the pseudorandom proof index
		branch, sdkErr := por.Tree.GetMerkleProof(reqProof)
		if sdkErr != nil {
			panic(sdkErr)
		}
		// get the leaf for the required pseudorandom proof index
		leaf := pc.GetAllProofs().GetProof(proof.SessionHeader, reqProof)
		// send the claim TX
		err := claimTx(k.cdc, cliCtx, txBuilder, branch, *leaf)
		if err != nil {
			panic(err)
		}
	}
}

// retrieve the verified proof
func (k Keeper) GetProof(ctx sdk.Context, address sdk.ValAddress, header pc.SessionHeader) (proof pc.StoredProof, found bool) {
	store := ctx.KVStore(k.storeKey)
	res := store.Get(pc.KeyForProof(ctx, address, header))
	if res == nil {
		return pc.StoredProof{}, false
	}
	k.cdc.MustUnmarshalBinaryBare(res, &proof)
	return proof, true
}

// set the verified proof
func (k Keeper) SetProof(ctx sdk.Context, address sdk.ValAddress, p pc.StoredProof) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(p)
	store.Set(pc.KeyForProof(ctx, address, p.SessionHeader), bz)
}

func (k Keeper) SetAllProofs(ctx sdk.Context, proofs []pc.StoredProof) {
	store := ctx.KVStore(k.storeKey)
	for _, proof := range proofs {
		bz := k.cdc.MustMarshalBinaryBare(proof)
		store.Set(pc.KeyForProof(ctx, sdk.ValAddress(proof.ServicerAddress), proof.SessionHeader), bz)
	}
}

// get all verified proofs for this address
func (k Keeper) GetAllProofs(ctx sdk.Context) (proofs []pc.StoredProof) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, pc.ProofKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var summary pc.StoredProof
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &summary)
		proofs = append(proofs, summary)
	}
	return
}

// get all verified proofs for this address
func (k Keeper) GetProofs(ctx sdk.Context, address sdk.ValAddress) (proofs []pc.StoredProof) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, pc.KeyForProofs(address))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var summary pc.StoredProof
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &summary)
		proofs = append(proofs, summary)
	}
	return
}

func (k Keeper) SetAllClaims(ctx sdk.Context, claims []pc.MsgClaim) {
	store := ctx.KVStore(k.storeKey)
	for _, msg := range claims {
		bz := k.cdc.MustMarshalBinaryBare(msg)
		store.Set(pc.KeyForClaim(ctx, msg.FromAddress, msg.SessionHeader), bz)
	}
}

// get all verified proofs
func (k Keeper) GetAllClaims(ctx sdk.Context) (proofs []pc.MsgClaim) {
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

// get all verified proofs for this address
func (k Keeper) GetClaims(ctx sdk.Context, address sdk.ValAddress) (proofs []pc.MsgClaim) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, pc.KeyForClaims(address))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var summary pc.MsgClaim
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &summary)
		proofs = append(proofs, summary)
	}
	return
}

// get the unverified proof for this address
func (k Keeper) GetUnverfiedProof(ctx sdk.Context, address sdk.ValAddress, header pc.SessionHeader) (msg pc.MsgClaim, found bool) {
	store := ctx.KVStore(k.storeKey)
	res := store.Get(pc.KeyForClaim(ctx, address, header))
	if res == nil {
		return pc.MsgClaim{}, false
	}
	k.cdc.MustUnmarshalBinaryBare(res, &msg)
	return msg, true
}

// set the unverified proof
func (k Keeper) SetClaim(ctx sdk.Context, msg pc.MsgClaim) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(msg)
	store.Set(pc.KeyForClaim(ctx, msg.FromAddress, msg.SessionHeader), bz)
}

func (k Keeper) DeleteClaim(ctx sdk.Context, address sdk.ValAddress, header pc.SessionHeader) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(pc.KeyForClaim(ctx, address, header))
}

// get the mature unverified proofs for this address
func (k Keeper) GetMatureClaims(ctx sdk.Context, address sdk.ValAddress) (matureProofs []pc.MsgClaim) {
	var msg = pc.MsgClaim{}
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, pc.KeyForClaims(address))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), msg)
		if k.ClaimIsMature(ctx, msg.SessionBlockHeight) {
			matureProofs = append(matureProofs, msg)
		}
	}
	return
}

// delete expired unverified proofs
func (k Keeper) DeleteExpiredClaims(ctx sdk.Context) {
	var msg = pc.MsgClaim{}
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, pc.ClaimKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), msg)
		sessionContext := ctx.WithBlockHeight(msg.SessionBlockHeight)
		// if more sessions has passed than the expiration of unverified pseudorandomGenerator, delete from set
		if (ctx.BlockHeight()-msg.SessionBlockHeight)/k.SessionFrequency(sessionContext) >= int64(k.ClaimExpiration(sessionContext)) { // todo confirm these contexts should be now and not when submitted
			store.Delete(iterator.Key())
		}
	}
}

// is the proof mature? able to be claimed because the `waiting period` has passed since the sessionBlock
func (k Keeper) ClaimIsMature(ctx sdk.Context, sessionBlockHeight int64) bool {
	waitingPeriodInBlocks := k.ProofWaitingPeriod(ctx) * k.SessionFrequency(ctx)
	if ctx.BlockHeight() >= waitingPeriodInBlocks+sessionBlockHeight {
		return true
	}
	return false
}

// todo this auto tx needs to be fixed
func newTxBuilderAndCliCtx(ctx sdk.Context, n *node.Node, k Keeper) (txBuilder auth.TxBuilder, cliCtx util.CLIContext) {
	fromAddr := sdk.AccAddress(n.PrivValidator().GetPubKey().Address())
	fee := auth.NewStdFee(9000, sdk.NewCoins(sdk.NewInt64Coin(k.StakeDenom(ctx), 0)))
	cliCtx = util.NewCLIContext(n, fromAddr, k.coinbasePassphrase).WithCodec(k.cdc)
	accGetter := auth.NewAccountRetriever(cliCtx)
	err := accGetter.EnsureExists(fromAddr)
	if err != nil {
		panic(err)
	}
	account, err := accGetter.GetAccount(fromAddr)
	if err != nil {
		panic(err)
	}
	txBuilder = auth.TxBuilder{
		auth.DefaultTxEncoder(k.cdc),
		k.keybase,
		account.GetAccountNumber(),
		account.GetSequence(),
		fee.Gas,
		1,
		false,
		n.GenesisDoc().ChainID,
		"",
		fee.Amount,
		fee.GasPrices(),
	}
	return
}
