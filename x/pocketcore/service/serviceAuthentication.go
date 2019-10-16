package service

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
)

type ServiceCertificate struct {
	ServiceCertificatePayload
	Signature string `json:"signature"` // client's signature for the service in hex
}

type ServiceCertificatePayload struct {
	Counter       int          `json:"counter"`       // needed for pseudorandom evidence selection
	NodePublicKey string       `json:"nodePublicKey"` // needed to prove the node was the servicer
	ServiceToken  ServiceToken `json:"serviceToken"`  // needed to prove the client is authorized to use developers throughput
}

func (sap ServiceCertificatePayload) Hash() ([]byte, error) { // TODO true hash with amino encoding of payload
	return crypto.SHA3FromBytes(append([]byte(sap.NodePublicKey))), nil
}

func (sa ServiceCertificate) Validate() error {
	if sa.Counter < 0 {
		return NegativeICCounterError
	}
	cpkBytes, err := hex.DecodeString(sa.ServiceToken.AATMessage.ClientPublicKey)
	if err != nil {
		return NewClientPubKeyDecodeError(err)
	}
	// validate the service token
	if err := sa.ServiceToken.Validate(); err != nil {
		return NewInvalidTokenError(err)
	}
	// validate the public key correctness
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
