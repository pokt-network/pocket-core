package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"log"
)

// "Proof" - An interface representation of an economic proof of work/burn (relay or challenge)
type Proof interface {
	Hash() []byte                                                                                        // returns cryptographic hash bz
	HashString() string                                                                                  // returns the hex string representation of the hash
	ValidateBasic() sdk.Error                                                                            // storeless validation check for the object
	GetSigner() sdk.Address                                                                              // returns the main signer(s) for the proof (used in messages)
	SessionHeader() SessionHeader                                                                        // returns the session header
	Validate(appSupportedBlockchains []string, sessionNodeCount int, sessionBlockHeight int64) sdk.Error // validate the object
	Store(max sdk.Int)                                                                                   // handle the proof after validation
}

var _ Proof = RelayProof{} // ensure implements interface at compile time

// "RelayProof" - A proof object that represetns one relay finished
type RelayProof struct {
	RequestHash        string `json:"request_hash"`         // the hash of the request used for response comparison
	Entropy            int64  `json:"entropy"`              // a random int64 used for replay prevention on the node verification (merkle) side
	SessionBlockHeight int64  `json:"session_block_height"` // The height of the session block
	ServicerPubKey     string `json:"servicer_pub_key"`     // the public key of the servicer in hex
	Blockchain         string `json:"blockchain"`           // the non-native chain net id in hex
	Token              AAT    `json:"aat"`                  // the app auth token object
	Signature          string `json:"signature"`            // the signature in hex
}

// "ValidateLocal" - Validates the proof object, where the owner of the proof is the local node
func (rp RelayProof) ValidateLocal(appSupportedBlockchains []string, sessionNodeCount int, sessionBlockHeight int64, verifyPubKey string) sdk.Error {
	// validate the public key correctness
	if rp.ServicerPubKey != verifyPubKey {
		return NewInvalidNodePubKeyError(ModuleName) // the public key is not this nodes, so they would not get paid
	}
	// validate the verify public key format
	if err := PubKeyVerification(verifyPubKey); err != nil {
		return NewInvalidNodePubKeyError(ModuleName)
	}
	err := rp.Validate(appSupportedBlockchains, sessionNodeCount, sessionBlockHeight)
	if err != nil {
		return err
	}
	return nil
}

// "Validate" - Validates the relay proof object
func (rp RelayProof) Validate(appSupportedBlockchains []string, sessionNodeCount int, sessionBlockHeight int64) sdk.Error {
	//Basic Validations
	err := rp.ValidateBasic()
	if err != nil {
		return err
	}
	// validate the session block height
	if rp.SessionBlockHeight != sessionBlockHeight {
		return NewInvalidBlockHeightError(ModuleName)
	}
	// check for supported blockchain
	c1 := false
	for _, chain := range appSupportedBlockchains {
		if rp.Blockchain == chain {
			c1 = true
		}
	}
	if !c1 {
		return NewUnsupportedBlockchainAppError(ModuleName)
	}
	return nil
}

// "ValidateBasic" - Provides a lighter weight, storeless validation of the relay proof object
func (rp RelayProof) ValidateBasic() sdk.Error {
	// verify the session block height is positive
	if rp.SessionBlockHeight < 1 {
		return NewInvalidBlockHeightError(ModuleName)
	}
	// verify the public key format for the leaf
	if err := PubKeyVerification(rp.ServicerPubKey); err != nil {
		return err
	}
	// verify the blockchain addr format
	if err := NetworkIdentifierVerification(rp.Blockchain); err != nil {
		return err
	}
	// verify the request hash format
	if err := HashVerification(rp.RequestHash); err != nil {
		return err
	}
	// verify non negative index
	if rp.Entropy < 0 {
		return NewInvalidEntropyError(ModuleName)
	}
	// verify a valid token
	if err := rp.Token.Validate(); err != nil {
		return NewInvalidTokenError(ModuleName, err)
	}
	// verify the client signature on the Proof
	if err := SignatureVerification(rp.Token.ClientPublicKey, rp.HashString(), rp.Signature); err != nil {
		return err
	}
	return nil
}

// "SessionHeader" - Returns the session header corresponding with the proof
func (rp RelayProof) SessionHeader() SessionHeader {
	return SessionHeader{
		ApplicationPubKey:  rp.Token.ApplicationPublicKey,
		Chain:              rp.Blockchain,
		SessionBlockHeight: rp.SessionBlockHeight,
	}
}

// "relayProof" - A structure used to json marshal the RelayProof
type relayProof struct {
	Entropy            int64  `json:"entropy"`
	SessionBlockHeight int64  `json:"session_block_height"`
	ServicerPubKey     string `json:"servicer_pub_key"`
	Blockchain         string `json:"blockchain"`
	Signature          string `json:"signature"`
	Token              string `json:"token"`
	RequestHash        string `json:"request_hash"`
}

// "Bytes" - Converts the RelayProof to bytes
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
		log.Fatal(fmt.Errorf("an error occured converting the relay RelayProof to bytes:\n%v", err).Error())
	}
	return res
}

// "BytesWithSignature" - Convert the RelayProof to bytes
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
		log.Fatalf(fmt.Errorf("an error occured converting the relay RelayProof to bytesWithSignature:\n%v", err).Error())
	}
	return res
}

// "Hash" - Returns the cryptographic hash of the rp bytes
func (rp RelayProof) Hash() []byte {
	res := rp.Bytes()
	return Hash(res)
}

// "HashString" - Returns the hex encoded string of the rp hash
func (rp RelayProof) HashString() string {
	return hex.EncodeToString(rp.Hash())
}

// "HashWithSignature" - Returns the cryptographic hash of the rp bytes (with signature field)
func (rp RelayProof) HashWithSignature() []byte {
	res := rp.BytesWithSignature()
	return Hash(res)
}

// "HashStringWithSignature" - Returns the hex encoded string of the rp hash (with signature field)
func (rp RelayProof) HashStringWithSignature() string {
	return hex.EncodeToString(rp.HashWithSignature())
}

// "Store" - Handles the relay proof object by adding it to the cache
func (rp RelayProof) Store(maxRelays sdk.Int) {
	// add the Proof to the global (in memory) collection of proofs
	SetProof(rp.SessionHeader(), RelayEvidence, rp, maxRelays)
}

func (rp RelayProof) GetSigner() sdk.Address {
	pk, err := crypto.NewPublicKey(rp.ServicerPubKey)
	if err != nil {
		return nil
	}
	return sdk.Address(pk.Address())
}

// ---------------------------------------------------------------------------------------------------------------------

// "ChallengeProofInvalidData" - Is a challenge of response data using a majority consensus
type ChallengeProofInvalidData struct {
	MajorityResponses [2]RelayResponse `json:"majority_responses"` // the majority who agreed
	MinorityResponse  RelayResponse    `json:"minority_response"`  // the minority who disagreed
	ReporterAddress   sdk.Address      `json:"address"`            // the address of the reporter
}

var _ Proof = ChallengeProofInvalidData{} // compile time interface implementation

// "ValidateLocal" - Validate local is used to validate a challenge request directly from a client
func (c ChallengeProofInvalidData) ValidateLocal(h SessionHeader, maxRelays sdk.Int, supportedBlockchains []string, sessionNodeCount int, sessionNodes SessionNodes, selfAddr sdk.Address) sdk.Error {
	// check if verifyPubKey in session (must be in session to do challenges)
	if !sessionNodes.ContainsAddress(selfAddr) {
		return NewNodeNotInSessionError(ModuleName)
	}
	sessionblockHeight := h.SessionBlockHeight
	// calculate the maximum possible challenges
	maxPossibleChallenges := maxRelays.ToDec().Quo(sdk.NewDec(int64(len(supportedBlockchains)))).Quo(sdk.NewDec(int64(sessionNodeCount))).RoundInt()
	// check for overflow on # of proofs
	evidence, er := GetEvidence(h, ChallengeEvidence, maxPossibleChallenges)
	if er != nil {
		return sdk.ErrInternal(er.Error())
	}
	if evidence.NumOfProofs >= maxPossibleChallenges.Int64() {
		return NewOverServiceError(ModuleName)
	}
	err := c.Validate(supportedBlockchains, sessionNodeCount, sessionblockHeight)
	if err != nil {
		return err
	}
	return nil
}

// "Validate" - validate is used to validate a challenge request
func (c ChallengeProofInvalidData) Validate(appSupportedBlockchains []string, sessionNodeCount int, sessionBlockHeight int64) sdk.Error {
	majResponse := c.MajorityResponses[0]
	majResponse2 := c.MajorityResponses[1]
	// check for duplicates
	if majResponse.Proof.ServicerPubKey == majResponse2.Proof.ServicerPubKey ||
		majResponse2.Proof.ServicerPubKey == c.MinorityResponse.Proof.ServicerPubKey ||
		c.MinorityResponse.Proof.ServicerPubKey == majResponse.Proof.ServicerPubKey {
		return NewDuplicatePublicKeyError(ModuleName)
	}
	// check for identical request hashes
	if majResponse.Proof.RequestHash != majResponse2.Proof.RequestHash ||
		majResponse2.Proof.RequestHash != c.MinorityResponse.Proof.RequestHash ||
		majResponse.Proof.RequestHash != c.MinorityResponse.Proof.RequestHash {
		return NewMismatchedRequestHashError(ModuleName)
	}
	// check for identical app public keys
	if majResponse.Proof.Token.ApplicationPublicKey != majResponse2.Proof.Token.ApplicationPublicKey ||
		majResponse2.Proof.Token.ApplicationPublicKey != c.MinorityResponse.Proof.Token.ApplicationPublicKey ||
		majResponse.Proof.Token.ApplicationPublicKey != c.MinorityResponse.Proof.Token.ApplicationPublicKey {
		return NewMismatchedAppPubKeyError(ModuleName)
	}
	// check for identical session heights
	if majResponse.Proof.SessionBlockHeight != majResponse2.Proof.SessionBlockHeight ||
		majResponse2.Proof.SessionBlockHeight != c.MinorityResponse.Proof.SessionBlockHeight ||
		majResponse.Proof.SessionBlockHeight != c.MinorityResponse.Proof.SessionBlockHeight {
		return NewMismatchedSessionHeightError(ModuleName)
	}
	// check for identical external blockchains
	if majResponse.Proof.Blockchain != majResponse2.Proof.Blockchain ||
		majResponse2.Proof.Blockchain != c.MinorityResponse.Proof.Blockchain ||
		majResponse.Proof.Blockchain != c.MinorityResponse.Proof.Blockchain {
		return NewMismatchedBlockchainsError(ModuleName)
	}
	// check for a true majority minority response
	majResp, majResp2, minResp := sortJSONResponse(majResponse.Response), sortJSONResponse(majResponse2.Response), sortJSONResponse(c.MinorityResponse.Response)
	if majResp != majResp2 || minResp == majResp {
		return NewNoMajorityResponseError(ModuleName)
	}
	// check for supported blockchain
	supported := false
	for _, chain := range appSupportedBlockchains {
		if majResponse.Proof.Blockchain == chain {
			supported = true
		}
	}
	if !supported {
		return NewUnsupportedBlockchainAppError(ModuleName)
	}
	// check signatures
	pubKey1, err := crypto.NewPublicKey(majResponse.Proof.ServicerPubKey)
	if err != nil {
		return NewPubKeyError(ModuleName, err)
	}
	pubKey2, err := crypto.NewPublicKey(majResponse2.Proof.ServicerPubKey)
	if err != nil {
		return NewPubKeyError(ModuleName, err)
	}
	pubKey3, err := crypto.NewPublicKey(c.MinorityResponse.Proof.ServicerPubKey)
	if err != nil {
		return NewPubKeyError(ModuleName, err)
	}
	sig1, err := hex.DecodeString(majResponse.Signature)
	if err != nil {
		return NewSignatureError(ModuleName, err)
	}
	sig2, err := hex.DecodeString(majResponse2.Signature)
	if err != nil {
		return NewSignatureError(ModuleName, err)
	}
	sig3, err := hex.DecodeString(c.MinorityResponse.Signature)
	if err != nil {
		return NewSignatureError(ModuleName, err)
	}
	if !pubKey1.VerifyBytes(majResponse.Hash(), sig1) {
		return NewInvalidSignatureError(ModuleName)
	}
	if !pubKey2.VerifyBytes(majResponse2.Hash(), sig2) {
		return NewInvalidSignatureError(ModuleName)
	}
	if !pubKey3.VerifyBytes(c.MinorityResponse.Hash(), sig3) {
		return NewInvalidSignatureError(ModuleName)
	}
	return nil
}

// "ValidateBasic" - Provides a lightweight, storeless validity check
func (c ChallengeProofInvalidData) ValidateBasic() sdk.Error {
	// ensure address is not empty
	if c.ReporterAddress == nil {
		return NewEmptyAddressError(ModuleName)
	}
	// ensure can decode from hex
	majResp, majResp2 := c.MajorityResponses[0], c.MajorityResponses[1]
	if _, err := hex.DecodeString(majResp.Signature); err != nil {
		return NewSigDecodeError(ModuleName)
	}
	if _, err := hex.DecodeString(majResp2.Signature); err != nil {
		return NewSigDecodeError(ModuleName)
	}
	if _, err := hex.DecodeString(c.MinorityResponse.Signature); err != nil {
		return NewSigDecodeError(ModuleName)
	}
	// validate each response individuall
	if err := majResp.Validate(); err != nil {
		return err
	}
	if err := majResp2.Validate(); err != nil {
		return err
	}
	if err := c.MinorityResponse.Validate(); err != nil {
		return err
	}
	// validate the proofs individually
	if err := majResp.Proof.ValidateBasic(); err != nil {
		return err
	}
	if err := majResp2.Proof.ValidateBasic(); err != nil {
		return err
	}
	if err := c.MinorityResponse.Proof.ValidateBasic(); err != nil {
		return err
	}
	// compare the responses and ensure minority is in disagreement w/ the majority responses
	if c.MinorityResponse.Proof.RequestHash != majResp.Proof.RequestHash || majResp.Proof.RequestHash != majResp2.Proof.RequestHash {
		return NewMismatchedRequestHashError(ModuleName)
	}
	return nil
}

// "SessionHeader" - Returns the session header for the challenge proof
func (c ChallengeProofInvalidData) SessionHeader() SessionHeader {
	return SessionHeader{
		ApplicationPubKey:  c.MinorityResponse.Proof.Token.ApplicationPublicKey,
		Chain:              c.MinorityResponse.Proof.Blockchain,
		SessionBlockHeight: c.MinorityResponse.Proof.SessionBlockHeight,
	}
}

// "challengeProofInvalidData" - is used to marshal / unmarshal json
type challengeProofInvalidData struct {
	MajorityResponses [2]relayResponse
	MinorityResponse  relayResponse
}

// "Bytes" - Bytes representaiton fo the challenge proof object
func (c ChallengeProofInvalidData) Bytes() []byte {
	majResp, majResp2 := c.MajorityResponses[0], c.MajorityResponses[1]
	bz, err := json.Marshal(challengeProofInvalidData{
		MajorityResponses: [2]relayResponse{
			{
				Signature: majResp.Signature,
				Response:  majResp.Response,
				Proof:     majResp.Proof.HashStringWithSignature(),
			},
			{
				Signature: majResp2.Signature,
				Response:  majResp2.Response,
				Proof:     majResp2.Proof.HashStringWithSignature(),
			},
		},
		MinorityResponse: relayResponse{
			Signature: c.MinorityResponse.Signature,
			Response:  c.MinorityResponse.Response,
			Proof:     c.MinorityResponse.Proof.HashStringWithSignature(),
		},
	})
	if err != nil {
		log.Fatalf(fmt.Errorf("an error occured converting the challengeproof to bytes\n%v", err).Error())
	}
	return bz
}

// "Hash" - The cryptographic hash representation of the challenge bytes
func (c ChallengeProofInvalidData) Hash() []byte {
	return Hash(c.Bytes())
}

// "HashString" - The hex encoded string representation fo the challenge hash
func (c ChallengeProofInvalidData) HashString() string {
	return hex.EncodeToString(c.Hash())
}

// "GetSigners" - Returns the signer(s) for the message
func (c ChallengeProofInvalidData) GetSigner() sdk.Address {
	return c.ReporterAddress
}

// "Store" - Stores the challenge proof (stores in cache)
func (c ChallengeProofInvalidData) Store(maxChallenges sdk.Int) {
	// add the Proof to the global (in memory) collection of proofs
	SetProof(c.SessionHeader(), ChallengeEvidence, c, maxChallenges)
}
