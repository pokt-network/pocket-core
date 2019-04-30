package core

import (
	"encoding/json"
)

// "Relay" is a JSON structure that specifies information to complete reads and writes to other blockchains
type Relay struct {
	Blockchain []byte `json:"blockchain"` // blockchain Hash (includes the version and the netid)
	Payload    []byte `json:"data"`       // the data payload
	DevID      []byte `json:"devid"`      // the id needed to confirm servicing
	Token      Token  `json:"token"`      // the token given from the developer
	Method     []byte `json:"method"`     // the HTTP method needed for the call (defaults to POST)
	Path       []byte `json:"URL"`        // optional param for REST
}

type Token struct {
	ExpDate []byte
}

type RelayMessage struct {
	Relay     Relay
	Signature []byte
}

// "Validate" validates the relay objects legitimacy
func (r *Relay) Validate() error {
	if !GetHostedChains().ContainsFromBytes(r.Blockchain) {
		return UnsupportedBlockchainError
	}
	// TODO invalid token
	return nil
}

// "ErrorCheck" checks the relay object for initialization errors
func (r *Relay) ErrorCheck() error {
	if r.Blockchain == nil || len(r.Blockchain) == 0 {
		return MissingBlockchainError
	}
	if r.Payload == nil || len(r.Payload) == 0 {
		return MissingPayloadError
	}
	if r.Token.ExpDate == nil || len(r.Token.ExpDate) == 0 {
		return InvalidTokenError
	}
	if r.DevID == nil || len(r.DevID) == 0 {
		return MissingDevidError
	}
	if r.Method == nil || len(r.Method) == 0 {
		r.Method = []byte(DefaultHTTPMethod)
	}
	if r.Path == nil || len(r.Path) == 0 {
		return MissingPathError
	}
	if len(r.DevID) != 33 {
		return InvalidDevIDError
	}
	return r.Validate()
}

func (rm *RelayMessage) ErrorCheck() error {
	if rm.Signature == nil || len(rm.Signature) == 0 {
		return MissingSignatureError
	}
	return rm.Relay.ErrorCheck()
}

// "JSONToToken" converts token to json
func JSONToToken(b []byte) (*Token, error) {
	var t Token
	err := json.Unmarshal(b, t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// "TokenToJSON" converts token to json object
func (t Token) TokenToJSON() ([]byte, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return b, nil
}
