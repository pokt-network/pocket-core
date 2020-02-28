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

func (p Proof) Validate(maxRelays int64, numberOfChains, sessionNodeCount int, sessionBlockHeight int64, hb HostedBlockchains, verifyPubKey string) sdk.Error {
	// validate the session block height
	if p.SessionBlockHeight != sessionBlockHeight {
		return NewInvalidBlockHeightError(ModuleName)
	}
	// validate blockchain
	if err := HashVerification(p.Blockchain); err != nil {
		return err
	}
	evidenceHeader := SessionHeader{
		ApplicationPubKey:  p.Token.ApplicationPublicKey,
		Chain:              p.Blockchain,
		SessionBlockHeight: p.SessionBlockHeight,
	}
	// validate not over service
	totalRelays := GetAllEvidences().GetTotalRelays(evidenceHeader)
	if !GetAllEvidences().IsUniqueProof(evidenceHeader, p) {
		return NewDuplicateProofError(ModuleName)
	}
	if totalRelays >= int64(math.Ceil(float64(maxRelays)/float64(numberOfChains))/(float64(sessionNodeCount))) {
		return NewOverServiceError(ModuleName)
	}
	// validate the public key correctness
	if p.ServicerPubKey != verifyPubKey {
		return NewInvalidNodePubKeyError(ModuleName) // the public key is not this nodes, so they would not get paid
	}
	// ensure the blockchain is supported
	if !hb.ContainsFromString(p.Blockchain) {
		return NewUnsupportedBlockchainNodeError(ModuleName)
	}
	// validate the Proof public key format
	if err := PubKeyVerification(p.ServicerPubKey); err != nil {
		return NewInvalidNodePubKeyError(ModuleName)
	}
	// validate the verify public key format
	if err := PubKeyVerification(verifyPubKey); err != nil {
		return NewInvalidNodePubKeyError(ModuleName)
	}
	// validate the service token
	if err := p.Token.Validate(); err != nil {
		return NewInvalidTokenError(ModuleName, err)
	}
	return SignatureVerification(p.Token.ClientPublicKey, p.HashString(), p.Signature)
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
func (p Proof) Bytes() []byte {
	res, err := json.Marshal(relayProof{
		Entropy:            p.Entropy,
		ServicerPubKey:     p.ServicerPubKey,
		Blockchain:         p.Blockchain,
		SessionBlockHeight: p.SessionBlockHeight,
		Signature:          "", // omit the signature
		Token:              p.Token.HashString(),
	})
	if err != nil {
		panic(fmt.Sprintf("an error occured converting the relay proof to bytes:\n%v", err))
	}
	return res
}

// convert the Proof to bytes
func (p Proof) BytesWithSignature() []byte {
	res, err := json.Marshal(relayProof{
		Entropy:            p.Entropy,
		ServicerPubKey:     p.ServicerPubKey,
		Blockchain:         p.Blockchain,
		SessionBlockHeight: p.SessionBlockHeight,
		Signature:          p.Signature,
		Token:              p.Token.HashString(),
	})
	if err != nil {
		panic(fmt.Sprintf("an error occured converting the relay proof to bytes with signature:\n%v", err))
	}
	return res
}

// addr the Proof bytes
func (p Proof) Hash() []byte {
	res := p.Bytes()
	return Hash(res)
}

// hex encode the Proof addr
func (p Proof) HashString() string {
	return hex.EncodeToString(p.Hash())
}

// addr the Proof bytes
func (p Proof) HashWithSignature() []byte {
	res := p.BytesWithSignature()
	return Hash(res)
}

// hex encode the Proof addr
func (p Proof) HashStringWithSignature() string {
	return hex.EncodeToString(p.HashWithSignature())
}
