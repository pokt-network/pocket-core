package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"math"
)

type Proof interface {
	Hash() []byte
	HashString() string
	ValidateBasic() sdk.Error
	GetSigners() []sdk.Address
	SessionHeader() SessionHeader
	Validate(appSupportedBlockchains []string, sessionNodeCount int, sessionBlockHeight int64) sdk.Error
	Handle() sdk.Error
	EvidenceType() EvidenceType
}

var _ Proof = RelayProof{}

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

func (rp RelayProof) Validate(appSupportedBlockchains []string, sessionNodeCount int, sessionBlockHeight int64) sdk.Error {
	// validate the session block height
	if rp.SessionBlockHeight != sessionBlockHeight {
		return NewInvalidBlockHeightError(ModuleName)
	}
	// validate blockchain
	if err := HashVerification(rp.Blockchain); err != nil {
		return err
	}
	// validate the RelayProof public key format
	if err := PubKeyVerification(rp.ServicerPubKey); err != nil {
		return NewInvalidNodePubKeyError(ModuleName)
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
	// validate the service token
	if err := rp.Token.Validate(); err != nil {
		return NewInvalidTokenError(ModuleName, err)
	}
	return SignatureVerification(rp.Token.ClientPublicKey, rp.HashString(), rp.Signature)
}

func (rp RelayProof) ValidateBasic() sdk.Error {
	// verify the session block height is positive
	if rp.SessionBlockHeight < 0 {
		return NewInvalidBlockHeightError(ModuleName)
	}
	// verify the public key format for the leaf
	if err := PubKeyVerification(rp.ServicerPubKey); err != nil {
		return err
	}
	// verify the blockchain addr format
	if err := HashVerification(rp.Blockchain); err != nil {
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

func (rp RelayProof) SessionHeader() SessionHeader {
	return SessionHeader{
		ApplicationPubKey:  rp.Token.ApplicationPublicKey,
		Chain:              rp.Blockchain,
		SessionBlockHeight: rp.SessionBlockHeight,
	}
}

func (rp RelayProof) EvidenceType() EvidenceType {
	return RelayEvidence
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

func (rp RelayProof) Handle() sdk.Error {
	// add the Proof to the global (in memory) collection of proofs
	return GetEvidenceMap().AddToEvidence(rp.SessionHeader(), rp)
}

func (rp RelayProof) GetSigners() []sdk.Address {
	pk, err := crypto.NewPublicKey(rp.ServicerPubKey)
	if err != nil {
		panic(fmt.Sprintf("an error occured getting the signer for the RelayProof: \n%v", err))
	}
	return []sdk.Address{sdk.Address(pk.Address())}
}

type ChallengeProofInvalidData struct {
	MajorityResponses [2]RelayResponse `json:"majority_responses"`
	MinorityResponse  RelayResponse    `json:"minority_response"`
	ReporterAddress   sdk.Address      `json:"address"`
}

var _ Proof = ChallengeProofInvalidData{}

// validate local is used to validate a challenge request directly from a client
func (c ChallengeProofInvalidData) ValidateLocal(maxRelays, sessionblockHeight int64, supportedBlockchains []string, sessionNodeCount int, sessionNodes SessionNodes, selfAddr sdk.Address) sdk.Error {
	// get the header to retrieve the evidence object
	h := SessionHeader{
		ApplicationPubKey:  c.MinorityResponse.Proof.Token.ApplicationPublicKey,
		Chain:              c.MinorityResponse.Proof.Blockchain,
		SessionBlockHeight: c.MinorityResponse.Proof.SessionBlockHeight,
	}
	// check for overflow on # of proofs
	evidence, _ := GetEvidenceMap().GetEvidence(h, ChallengeEvidence)
	if evidence.NumOfProofs >= int64(math.Ceil(float64(maxRelays)/float64(len(supportedBlockchains)))/(float64(sessionNodeCount))) {
		return NewOverServiceError(ModuleName)
	}
	// check if verifyPubKey in session (must be in session to do challenges)
	if !sessionNodes.ContainsAddress(selfAddr) {
		return NewNodeNotInSessionError(ModuleName)
	}
	err := c.Validate(supportedBlockchains, sessionNodeCount, sessionblockHeight)
	if err != nil {
		return err
	}
	return nil
}

// validate is used to validate a challenge request
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
	c1, c2, c3 := false, false, false
	for _, chain := range appSupportedBlockchains {
		if majResponse.Proof.Blockchain == chain {
			c1 = true
		}
		if majResponse2.Proof.Blockchain == chain {
			c2 = true
		}
		if c.MinorityResponse.Proof.Blockchain == chain {
			c3 = true
		}
	}
	if !c1 || !c2 || !c3 {
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

func (c ChallengeProofInvalidData) ValidateBasic() sdk.Error {
	if c.ReporterAddress == nil {
		return NewEmptyAddressError(ModuleName)
	}
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
	if err := majResp.Validate(); err != nil {
		return err
	}
	if err := majResp2.Validate(); err != nil {
		return err
	}
	if err := c.MinorityResponse.Validate(); err != nil {
		return err
	}
	if err := majResp.Proof.ValidateBasic(); err != nil {
		return err
	}
	if err := majResp2.Proof.ValidateBasic(); err != nil {
		return err
	}
	if err := c.MinorityResponse.Proof.ValidateBasic(); err != nil {
		return err
	}
	if c.MinorityResponse.Proof.RequestHash != majResp.Proof.RequestHash || majResp.Proof.RequestHash != majResp2.Proof.RequestHash {
		return NewMismatchedRequestHashError(ModuleName)
	}
	return nil
}

func (c ChallengeProofInvalidData) SessionHeader() SessionHeader {
	return SessionHeader{
		ApplicationPubKey:  c.MinorityResponse.Proof.Token.ApplicationPublicKey,
		Chain:              c.MinorityResponse.Proof.Blockchain,
		SessionBlockHeight: c.MinorityResponse.Proof.SessionBlockHeight,
	}
}

type challengeProofInvalidData struct {
	MajorityResponses [2]relayResponse
	MinorityResponse  relayResponse
}

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
		panic(fmt.Sprintf("an error occured converting the challengeproof to bytes\n%v", err))
	}
	return bz
}

func (c ChallengeProofInvalidData) Hash() []byte {
	return Hash(c.Bytes())
}

func (c ChallengeProofInvalidData) HashString() string {
	return hex.EncodeToString(c.Hash())
}

func (c ChallengeProofInvalidData) GetSigners() []sdk.Address {
	return []sdk.Address{c.ReporterAddress}
}

func (c ChallengeProofInvalidData) Handle() sdk.Error {
	// add the Proof to the global (in memory) collection of proofs
	return GetEvidenceMap().AddToEvidence(c.SessionHeader(), c)
}

func (c ChallengeProofInvalidData) EvidenceType() EvidenceType {
	return ChallengeEvidence
}
