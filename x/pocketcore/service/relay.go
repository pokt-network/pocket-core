package service

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/types"
)

// a read / write API request from a hosted (non native) blockchain
type Relay struct {
	Blockchain         ServiceBlockchain  `json:"blockchain"`       // the non-native blockchain needed to service
	Payload            ServicePayload     `json:"payload"`          // the data payload of the request
	ServiceCertificate ServiceCertificate `json:"incrementCounter"` // the authentication scheme needed for work
}

func (r Relay) Validate(hostedBlockchains ServiceBlockchains, sessionBlockIDHex string, allActiveNodes types.Nodes) error {
	// check to see if the blockchain is empty
	if r.Blockchain == nil || len(r.Blockchain) == 0 {
		return EmptyBlockchainError
	}
	// check to see if the payload is empty
	if r.Payload.Data.Bytes() == nil || len(r.Payload.Data.Bytes()) == 0 {
		return EmptyPayloadDataError
	}
	// ensure the blockchain is supported
	if !hostedBlockchains.Contains(hex.EncodeToString(r.Blockchain)) {
		return UnsupportedBlockchainError
	}
	// todo check to see if non native blockchain is staked for by the developer
	// getApplication().GetStakedBlockchains()
	// verify that node (self) is of this session
	if err := SessionSelfVerification(FAKEAPPPUBKEY,
		r.Blockchain,
		sessionBlockIDHex,
		allActiveNodes); err != nil {
		return err
	}
	// check to see if the service certificate is valid
	if err := r.ServiceCertificate.Validate(); err != nil {
		return NewServiceCertificateError(err)
	}
	if r.Payload.Type() == HTTP {
		if len((r.Payload).Method) == 0 {
			r.Payload.Method = DEFAULTHTTPMETHOD
		}
	}
	return nil
}

// store the evidence of work done for the relay batch
func (r Relay) StoreEvidence(sessionBlockIDHex string) error {
	// grab the relay batch container
	rbs := GetGlobalRelayBatches()
	// add the evidence to the proper batch
	return rbs.AddEvidence(r.ServiceCertificate, sessionBlockIDHex)
}

// executes the relay on the non-native blockchain specified
func (r Relay) Execute(hostedBlockchains ServiceBlockchains, sessionBlockIDHex string, allActiveNodes types.Nodes) (string, error) {
	if err := r.Validate(hostedBlockchains, sessionBlockIDHex, allActiveNodes); err != nil {
		return "", err
	}
	// handle the relay payload based on the type
	switch r.Payload.Type() {
	case HTTP:
		url, err := r.Blockchain.GetHostedChainURL(hostedBlockchains)
		if err != nil {
			return "", err
		}
		return executeHTTPRequest(r.Payload.Data.Bytes(), url, r.Payload.Method)
	}
	return "", UnsupportedPayloadTypeError
}
