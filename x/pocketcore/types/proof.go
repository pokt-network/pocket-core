package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"math"
)

type Proof interface {
	Validate(maxRelays int64, numberOfChains, sessionNodeCount int, sessionBlockHeight int64, hb HostedBlockchains, verifyPubKey string) sdk.Error
	Hash() []byte
	HashString() string
	HashWithSignature() []byte
	HashStringWithSignature() string
}

// RelayProof per relay
type RelayProof struct {
	Entropy            int64  `json:"entropy"`
	SessionBlockHeight int64  `json:"session_block_height"`
	ServicerPubKey     string `json:"servicer_pub_key"`
	Blockchain         string `json:"blockchain"`
	Token              AAT    `json:"aat"`
	Signature          string `json:"signature"`
}

func (p RelayProof) Validate(maxRelays int64, numberOfChains, sessionNodeCount int, sessionBlockHeight int64, hb HostedBlockchains, verifyPubKey string) sdk.Error {
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
	totalRelays := GetEvidenceMap().GetTotalRelays(evidenceHeader)
	if !GetEvidenceMap().IsUniqueProof(evidenceHeader, p) {
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
	// validate the RelayProof public key format
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

// structure used to json marshal the RelayProof
type relayRelayProof struct {
	Entropy            int64  `json:"entropy"`
	SessionBlockHeight int64  `json:"session_block_height"`
	ServicerPubKey     string `json:"servicer_pub_key"`
	Blockchain         string `json:"blockchain"`
	Signature          string `json:"signature"`
	Token              string `json:"token"`
}

// convert the RelayProof to bytes
func (p RelayProof) Bytes() []byte {
	res, err := json.Marshal(relayRelayProof{
		Entropy:            p.Entropy,
		ServicerPubKey:     p.ServicerPubKey,
		Blockchain:         p.Blockchain,
		SessionBlockHeight: p.SessionBlockHeight,
		Signature:          "", // omit the signature
		Token:              p.Token.HashString(),
	})
	if err != nil {
		panic(fmt.Sprintf("an error occured converting the relay RelayProof to bytes:\n%v", err))
	}
	return res
}

// convert the RelayProof to bytes
func (p RelayProof) BytesWithSignature() []byte {
	res, err := json.Marshal(relayRelayProof{
		Entropy:            p.Entropy,
		ServicerPubKey:     p.ServicerPubKey,
		Blockchain:         p.Blockchain,
		SessionBlockHeight: p.SessionBlockHeight,
		Signature:          p.Signature,
		Token:              p.Token.HashString(),
	})
	if err != nil {
		panic(fmt.Sprintf("an error occured converting the relay RelayProof to bytes with signature:\n%v", err))
	}
	return res
}

// addr the RelayProof bytes
func (p RelayProof) Hash() []byte {
	res := p.Bytes()
	return Hash(res)
}

// hex encode the RelayProof addr
func (p RelayProof) HashString() string {
	return hex.EncodeToString(p.Hash())
}

// addr the RelayProof bytes
func (p RelayProof) HashWithSignature() []byte {
	res := p.BytesWithSignature()
	return Hash(res)
}

// hex encode the RelayProof addr
func (p RelayProof) HashStringWithSignature() string {
	return hex.EncodeToString(p.HashWithSignature())
}
