package service

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
)

type RelayResponse struct {
	Signature   string             // signature from the node in hex
	Response    string             // response to relay
	ServiceAuth ServiceCertificate // to be signed by the client
}

// node validates the response after signing
func (rr RelayResponse) Validate() error {
	// the counter for the authorization must be >=0
	if rr.ServiceAuth.Counter < 0 {
		return InvalidIncrementCounterError
	}
	// cannot contain empty response
	if rr.Response == "" {
		return EmptyResponseError
	}
	// cannot contain empty signature (nodes must be accountable)
	if rr.Signature == "" {
		return ResponseSignatureError
	}
	return nil
}

// node signs the response before validating back
func (rr RelayResponse) Sign() error {
	privateKey := FAKENODEPRIVKEY
	// sign the hash of the response body
	// todo should the node sign the service auth as well to preserve the counter? Is this a possible attack?
	sig, err := privateKey.Sign(crypto.SHA3FromString(rr.Response))
	if err != nil {
		return NewSignatureError(err)
	}
	rr.Signature = hex.EncodeToString(sig)
	return nil
}
