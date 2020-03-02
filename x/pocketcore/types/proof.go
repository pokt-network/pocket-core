package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"math"
)

type Proof interface {
	Validate(maxRelays int64, numberOfChains, sessionNodeCount int, sessionBlockHeight int64, hb HostedBlockchains, payloadHash, verifyPubKey string) sdk.Error
	Hash() []byte
	HashString() string
	HashWithSignature() []byte
	HashStringWithSignature() string
}

type ChallengeProofInvalidData struct {
}

func (c ChallengeProofInvalidData) Validate(maxRelays int64, numberOfChains, sessionNodeCount int, sessionBlockHeight int64, hb HostedBlockchains, payloadHash, verifyPubKey string) sdk.Error {
	panic("implement me")
}

func (c ChallengeProofInvalidData) Hash() []byte {
	panic("implement me")
}

func (c ChallengeProofInvalidData) HashString() string {
	panic("implement me")
}

func (c ChallengeProofInvalidData) HashWithSignature() []byte {
	panic("implement me")
}

func (c ChallengeProofInvalidData) HashStringWithSignature() string {
	panic("implement me")
}

type ChallengeProofCorruptedRequest struct {
}

func (c ChallengeProofCorruptedRequest) Validate(maxRelays int64, numberOfChains, sessionNodeCount int, sessionBlockHeight int64, hb HostedBlockchains, payloadHash, verifyPubKey string) sdk.Error {
	panic("implement me")
}

func (c ChallengeProofCorruptedRequest) Hash() []byte {
	panic("implement me")
}

func (c ChallengeProofCorruptedRequest) HashString() string {
	panic("implement me")
}

func (c ChallengeProofCorruptedRequest) HashWithSignature() []byte {
	panic("implement me")
}

func (c ChallengeProofCorruptedRequest) HashStringWithSignature() string {
	panic("implement me")
}

// RelayProof per relay
type RelayProof struct {
	RequestHash        string `json:"request_hash"`
	Entropy            int64  `json:"entropy"`
	SessionBlockHeight int64  `json:"session_block_height"`
	ServicerPubKey     string `json:"servicer_pub_key"`
	Blockchain         string `json:"blockchain"`
	Token              AAT    `json:"aat"`
	Signature          string `json:"signature"`
}

func (rp RelayProof) Validate(maxRelays int64, numberOfChains, sessionNodeCount int, sessionBlockHeight int64, hb HostedBlockchains, payloadHash, verifyPubKey string) sdk.Error {
	// validate the session block height
	if rp.SessionBlockHeight != sessionBlockHeight {
		return NewInvalidBlockHeightError(ModuleName)
	}
	// validate blockchain
	if err := HashVerification(rp.Blockchain); err != nil {
		return err
	}
	// validate request hash
	if payloadHash != rp.RequestHash {
		return NewRequestHashError(ModuleName)
	}
	evidenceHeader := SessionHeader{
		ApplicationPubKey:  rp.Token.ApplicationPublicKey,
		Chain:              rp.Blockchain,
		SessionBlockHeight: rp.SessionBlockHeight,
	}
	// validate not over service
	totalRelays := GetEvidenceMap().GetTotalRelays(evidenceHeader)
	if !GetEvidenceMap().IsUniqueProof(evidenceHeader, rp) {
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
	// validate the RelayProof public key format
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

// structure used to json marshal the RelayProof
type relayProof struct {
	Entropy            int64  `json:"entropy"`
	SessionBlockHeight int64  `json:"session_block_height"`
	ServicerPubKey     string `json:"servicer_pub_key"`
	Blockchain         string `json:"blockchain"`
	Signature          string `json:"signature"`
	Token              string `json:"token"`
	RequestHash        string `json:"request_hash"`
}

// convert the RelayProof to bytes
func (rp RelayProof) Bytes() []byte {
	res, err := json.Marshal(relayProof{
		Entropy:            rp.Entropy,
		RequestHash:        rp.RequestHash,
		ServicerPubKey:     rp.ServicerPubKey,
		Blockchain:         rp.Blockchain,
		SessionBlockHeight: rp.SessionBlockHeight,
		Signature:          "", // omit the signature
		Token:              rp.Token.HashString(),
	})
	if err != nil {
		panic(fmt.Sprintf("an error occured converting the relay RelayProof to bytes:\n%v", err))
	}
	return res
}

// convert the RelayProof to bytes
func (rp RelayProof) BytesWithSignature() []byte {
	res, err := json.Marshal(relayProof{
		Entropy:            rp.Entropy,
		ServicerPubKey:     rp.ServicerPubKey,
		RequestHash:        rp.RequestHash,
		Blockchain:         rp.Blockchain,
		SessionBlockHeight: rp.SessionBlockHeight,
		Signature:          rp.Signature,
		Token:              rp.Token.HashString(),
	})
	if err != nil {
		panic(fmt.Sprintf("an error occured converting the relay RelayProof to bytes with signature:\n%v", err))
	}
	return res
}

// addr the RelayProof bytes
func (rp RelayProof) Hash() []byte {
	res := rp.Bytes()
	return Hash(res)
}

// hex encode the RelayProof addr
func (rp RelayProof) HashString() string {
	return hex.EncodeToString(rp.Hash())
}

// addr the RelayProof bytes
func (rp RelayProof) HashWithSignature() []byte {
	res := rp.BytesWithSignature()
	return Hash(res)
}

// hex encode the RelayProof addr
func (rp RelayProof) HashStringWithSignature() string {
	return hex.EncodeToString(rp.HashWithSignature())
}
