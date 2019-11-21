package types

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
)

// a read / write API request from a hosted (non native) blockchain
type Relay struct {
	Blockchain string  `json:"blockchain"`       // the non-native blockchain needed to service
	Payload    Payload `json:"payload"`          // the data payload of the request
	Proof      Proof   `json:"incrementCounter"` // the authentication scheme needed for work
}

type Payload struct {
	Data   string `json:"data"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

type PayloadType int

const HTTP PayloadType = iota + 1

func (p Payload) Type() PayloadType {
	return HTTP
}

type RelayResponse struct {
	Signature   string       // signature from the node in hex
	Response    string       // response to relay
	ServiceAuth Proof // to be signed by the client
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
	privateKey := FAKENODEPRIVKEY // todo get node private key
	// sign the hash of the response body
	// todo should the node sign the service auth as well to preserve the counter? Is this a possible attack? YES! To prevent unauthorized answers and proof for bad acting
	sig, err := privateKey.Sign(crypto.SHA3FromString(rr.Response))
	if err != nil {
		return NewSignatureError(err)
	}
	rr.Signature = hex.EncodeToString(sig)
	return nil
}
