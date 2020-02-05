package types

import (
	"encoding/hex"
	"fmt"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"reflect"
)

// RouterKey is the module name router key
const (
	RouterKey    = ModuleName
	MsgClaimName = "claim"
	MsgProofName = "proof"
)

// MsgClaim claims that you completed `TotalRelays` and provides the merkle root for data integrity
type MsgClaim struct {
	SessionHeader `json:"header"` // header information for identification
	MerkleRoot    HashSum         `json:"merkle_root"`  // merkle root for data integrity
	TotalRelays   int64           `json:"total_relays"` // total number of relays
	FromAddress   sdk.Address     `json:"from_address"` // claimant
}

func (msg MsgClaim) Route() string { return RouterKey }
func (msg MsgClaim) Type() string  { return MsgClaimName }
func (msg MsgClaim) ValidateBasic() sdk.Error {
	// validate a non empty chain
	if msg.Chain == "" {
		return NewEmptyChainError(ModuleName)
	}
	// basic validation on the session block height
	if msg.SessionBlockHeight < 1 {
		return NewEmptyBlockIDError(ModuleName)
	}
	// validate non negative total relays
	if msg.TotalRelays <= 0 {
		return NewEmptyProofsError(ModuleName)
	}
	// validate the public key format
	if err := PubKeyVerification(msg.ApplicationPubKey); err != nil {
		return NewPubKeyError(ModuleName, err)
	}
	// validate the address format
	if err := AddressVerification(msg.FromAddress.String()); err != nil {
		return NewInvalidHashError(ModuleName, err)
	}
	// validate the root format
	if err := HashVerification(hex.EncodeToString(msg.MerkleRoot.Hash)); err != nil {
		return err
	}
	// ensure non zero root sum
	if msg.MerkleRoot.Sum == 0 {
		return NewInvalidRootError(ModuleName)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgClaim) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.FromAddress)}
}

// ---------------------------------------------------------------------------------------------------------------------

// MsgProof proves the previous claim by providing the merkle RelayProof and the leaf node
type MsgProof struct {
	MerkleProofs MerkleProofs `json:"merkle_proofs"` // the merkleProof needed to verify the proofs
	Leaf         RelayProof   `json:"leaf"`          // the needed to verify the RelayProof
	Cousin       RelayProof   `json:"cousin"`        // the cousin needed to verify the RelayProof
}

func (msg MsgProof) Route() string { return RouterKey }
func (msg MsgProof) Type() string  { return MsgProofName }
func (msg MsgProof) ValidateBasic() sdk.Error {
	// verify valid number of levels for merkle proofs
	if len(msg.MerkleProofs[0].HashSums) < 3 || len(msg.MerkleProofs[0].HashSums) != len(msg.MerkleProofs[1].HashSums) {
		return NewInvalidLeafCousinProofsComboError(ModuleName)
	}
	// ensure the two indices are not equal
	if msg.MerkleProofs[0].Index == msg.MerkleProofs[1].Index {
		return NewInvalidLeafCousinProofsComboError(ModuleName)
	}
	// verify the session block height is positive
	if msg.Leaf.SessionBlockHeight < 0 {
		return NewInvalidBlockHeightError(ModuleName)
	}
	// verify the session block height is positive
	if msg.Cousin.SessionBlockHeight < 0 {
		return NewInvalidBlockHeightError(ModuleName)
	}
	// verify the public key format for the leaf
	if err := PubKeyVerification(msg.Leaf.ServicerPubKey); err != nil {
		return err
	}
	// verify the public key format for the cousin
	if err := PubKeyVerification(msg.Cousin.ServicerPubKey); err != nil {
		return err
	}
	// verify the blockchain addr format
	if err := HashVerification(msg.Leaf.Blockchain); err != nil {
		return err
	}
	// verify the blockchain addr format
	if err := HashVerification(msg.Cousin.Blockchain); err != nil {
		return err
	}
	// verify non negative index
	if msg.Leaf.Entropy <= 0 {
		return NewInvalidIncrementCounterError(ModuleName)
	}
	// verify non negative index
	if msg.Cousin.Entropy <= 0 {
		return NewInvalidIncrementCounterError(ModuleName)
	}
	// ensure leaf does not equal cousin
	if reflect.DeepEqual(msg.Leaf, msg.Cousin) {
		return NewCousinLeafEquivalentError(ModuleName)
	}
	// ensure leaf relayProof does not equal cousin relayProof
	if reflect.DeepEqual(msg.MerkleProofs[0].HashSums, msg.MerkleProofs[1].HashSums) {
		return NewCousinLeafEquivalentError(ModuleName)
	}
	// verify a valid token
	if err := msg.Leaf.Token.Validate(); err != nil {
		return NewInvalidTokenError(ModuleName, err)
	}
	// verify the client signature on the RelayProof
	if err := SignatureVerification(msg.Leaf.Token.ClientPublicKey, msg.Leaf.HashString(), msg.Leaf.Signature); err != nil {
		return err
	}
	// verify the client signature on the RelayProof
	if err := SignatureVerification(msg.Cousin.Token.ClientPublicKey, msg.Cousin.HashString(), msg.Cousin.Signature); err != nil {
		return err
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgProof) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgProof) GetSigners() []sdk.Address {
	pk, err := crypto.NewPublicKey(msg.Leaf.ServicerPubKey)
	if err != nil {
		panic(fmt.Sprintf("an error occured getting the signer for the proof message, %v", err))
	}
	return []sdk.Address{sdk.Address(pk.Address())}
}
