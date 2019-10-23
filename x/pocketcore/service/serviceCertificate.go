package service

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
)

// Proofs is a slice of type `ServiceProof`
// which are individual proofs of work completed
type Proofs []ServiceProof

// the header of the relay batch that is used
// to identify the relay batch in the global map
type ProofsHeader struct {
	SessionHash       string
	ApplicationPubKey string
}

// add proof of work completed (type service proof) to the proof structure
func (e Proofs) AddProof(sc ServiceProof) error {
	if e == nil || len(e) == 0 {
		return EmptyProofsError
	}
	// if the increment counter is less than the proof slice
	if len(e) < sc.Counter {
		return InvalidProofSizeError
	}
	// if the proof at index[increment counter] is not empty
	if e[sc.Counter].Signature != "" {
		return DuplicateProofError
	}
	// set proof at index[service proof] = proof of work completed (Service Proof)
	e[sc.Counter] = sc
	return nil
}

// Type that signifies work completed
type ServiceProof struct {
	ServiceProofPayload
	Signature string `json:"signature"` // client's signature for the service in hex
}

type ServiceProofPayload struct {
	Counter       int          `json:"counter"`       // needed for pseudorandom proof selection
	NodePublicKey string       `json:"nodePublicKey"` // needed to prove the node was the servicer
	ServiceToken  ServiceToken `json:"serviceToken"`  // needed to prove the client is authorized to use developers throughput
}

func (sap ServiceProofPayload) Hash() ([]byte, error) { // TODO true hash with amino encoding of payload
	return crypto.SHA3FromBytes(append([]byte(sap.NodePublicKey))), nil
}

func (sa ServiceProof) Validate() error {
	// check for negative counter
	if sa.Counter < 0 {
		return NegativeICCounterError
	}
	// todo validate max counter
	// if sa.Counter > app.maxNumberOfRelays(){}
	// decode the client public key todo remove and have sig verification convert to bytes
	cpkBytes, err := hex.DecodeString(sa.ServiceToken.AATMessage.ClientPublicKey)
	if err != nil {
		return NewClientPubKeyDecodeError(err)
	}
	// validate the service token
	if err := sa.ServiceToken.Validate(); err != nil {
		return NewInvalidTokenError(err)
	}
	// validate the public key correctness // todo consider abstraction so that this function can also be used for external service proof(s)
	if sa.NodePublicKey != "" { // TODO check for current public key
		return InvalidNodePubKeyError // the public key is not this nodes, so they would not get paid
	}
	hash, err := sa.ServiceProofPayload.Hash()
	if err != nil {
		return NewServiceProofHashError(err)
	}
	if !crypto.MockVerifySignature(cpkBytes, hash, []byte(sa.Signature)) { //todo change to real sig verification
		return InvalidICSignatureError
	}
	return nil
}
