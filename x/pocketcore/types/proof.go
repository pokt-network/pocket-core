package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"math"
)

// Proof per relay
type Proof struct {
	Entropy            int64  `json:"entropy"`
	SessionBlockHeight int64  `json:"session_block_height"`
	ServicerPubKey     string `json:"servicer_pub_key"`
	Blockchain         string `json:"blockchain"`
	Token              AAT    `json:"aat"`
	Signature          string `json:"signature"`
}

func (rp Proof) Validate(maxRelays int64, numberOfChains, sessionNodeCount int, sessionBlockHeight int64, hb HostedBlockchains, verifyPubKey string) sdk.Error {
	// validate the session block height
	if rp.SessionBlockHeight != sessionBlockHeight {
		return NewInvalidBlockHeightError(ModuleName)
	}
	// validate blockchain
	if err := HashVerification(rp.Blockchain); err != nil {
		return err
	}
	invoiceHeader := SessionHeader{
		ApplicationPubKey:  rp.Token.ApplicationPublicKey,
		Chain:              rp.Blockchain,
		SessionBlockHeight: rp.SessionBlockHeight,
	}
	// validate not over service
	totalRelays := GetAllInvoices().GetTotalRelays(invoiceHeader)
	if !GetAllInvoices().IsUniqueProof(invoiceHeader, rp) {
		return NewDuplicateProofError(ModuleName)
	}
	if totalRelays >= int64(math.Ceil(float64(maxRelays)/float64(numberOfChains))/(float64(sessionNodeCount))) {
		return NewOverServiceError(ModuleName)
	}
	// validate the public key correctness
	if rp.ServicerPubKey != verifyPubKey {
		return NewInvalidNodePubKeyError(ModuleName) // the public key is not this nodes, so they would not get paid
	}
	// ensure the blockchain is supported
	if !hb.ContainsFromString(rp.Blockchain) {
		return NewUnsupportedBlockchainNodeError(ModuleName)
	}
	// validate the Proof public key format
	if err := PubKeyVerification(rp.ServicerPubKey); err != nil {
		return NewInvalidNodePubKeyError(ModuleName)
	}
	// validate the verify public key format
	if err := PubKeyVerification(verifyPubKey); err != nil {
		return NewInvalidNodePubKeyError(ModuleName)
	}
	// validate the service token
	if err := rp.Token.Validate(); err != nil {
		return NewInvalidTokenError(ModuleName, err)
	}
	return SignatureVerification(rp.Token.ClientPublicKey, rp.HashString(), rp.Signature)
}

// structure used to json marshal the Proof
type relayProof struct {
	Entropy            int64  `json:"entropy"`
	SessionBlockHeight int64  `json:"session_block_height"`
	ServicerPubKey     string `json:"servicer_pub_key"`
	Blockchain         string `json:"blockchain"`
	Signature          string `json:"signature"`
	Token              string `json:"token"`
}

// convert the Proof to bytes
func (rp Proof) Bytes() []byte {
	res, err := json.Marshal(relayProof{
		Entropy:            rp.Entropy,
		ServicerPubKey:     rp.ServicerPubKey,
		Blockchain:         rp.Blockchain,
		SessionBlockHeight: rp.SessionBlockHeight,
		Signature:          "", // omit the signature
		Token:              rp.Token.HashString(),
	})
	if err != nil {
		panic(fmt.Sprintf("an error occured converting the relay proof to bytes:\n%v", err))
	}
	return res
}

// convert the Proof to bytes
func (rp Proof) BytesWithSignature() []byte {
	res, err := json.Marshal(relayProof{
		Entropy:            rp.Entropy,
		ServicerPubKey:     rp.ServicerPubKey,
		Blockchain:         rp.Blockchain,
		SessionBlockHeight: rp.SessionBlockHeight,
		Signature:          rp.Signature,
		Token:              rp.Token.HashString(),
	})
	if err != nil {
		panic(fmt.Sprintf("an error occured converting the relay proof to bytes with signature:\n%v", err))
	}
	return res
}

// addr the Proof bytes
func (rp Proof) Hash() []byte {
	res := rp.Bytes()
	return Hash(res)
}

// hex encode the Proof addr
func (rp Proof) HashString() string {
	return hex.EncodeToString(rp.Hash())
}

// addr the Proof bytes
func (rp Proof) HashWithSignature() []byte {
	res := rp.BytesWithSignature()
	return Hash(res)
}

// hex encode the Proof addr
func (rp Proof) HashStringWithSignature() string {
	return hex.EncodeToString(rp.HashWithSignature())
}
