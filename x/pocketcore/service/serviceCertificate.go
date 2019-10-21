package service

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
)

// Certificates is a slice of type `ServiceCertificate`
// which are individual proofs of work completed
type Certificates []ServiceCertificate

// the header of the relay batch that is used
// to identify the relay batch in the global map
type CertificatesHeader struct {
	SessionHash       string
	ApplicationPubKey string
}

// add proof of work completed (type service certificate) to the certificate structure
func (e Certificates) AddCertificate(sc ServiceCertificate) error {
	if e == nil || len(e) == 0 {
		return EmptyCertificatesError
	}
	// if the increment counter is less than the certificate slice
	if len(e) < sc.Counter {
		return InvalidCertificateSizeError
	}
	// if the certificate at index[increment counter] is not empty
	if e[sc.Counter].Signature != "" {
		return DuplicateCertificateError
	}
	// set certificate at index[service certificate] = proof of work completed (Service Certificate)
	e[sc.Counter] = sc
	return nil
}

// Type that signifies work completed
type ServiceCertificate struct {
	ServiceCertificatePayload
	Signature string `json:"signature"` // client's signature for the service in hex
}

type ServiceCertificatePayload struct {
	Counter       int          `json:"counter"`       // needed for pseudorandom certificate selection
	NodePublicKey string       `json:"nodePublicKey"` // needed to prove the node was the servicer
	ServiceToken  ServiceToken `json:"serviceToken"`  // needed to prove the client is authorized to use developers throughput
}

func (sap ServiceCertificatePayload) Hash() ([]byte, error) { // TODO true hash with amino encoding of payload
	return crypto.SHA3FromBytes(append([]byte(sap.NodePublicKey))), nil
}

func (sa ServiceCertificate) Validate() error {
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
	// validate the public key correctness // todo consider abstraction so that this function can also be used for external service certificate(s)
	if sa.NodePublicKey != "" { // TODO check for current public key
		return InvalidNodePubKeyError // the public key is not this nodes, so they would not get paid
	}
	hash, err := sa.ServiceCertificatePayload.Hash()
	if err != nil {
		return NewServiceCertificateHashError(err)
	}
	if !crypto.MockVerifySignature(cpkBytes, hash, []byte(sa.Signature)) { //todo change to real sig verification
		return InvalidICSignatureError
	}
	return nil
}
