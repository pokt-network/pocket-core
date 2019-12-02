package types

import (
	"bytes"
	"encoding/hex"
	appexported "github.com/pokt-network/pocket-core/x/apps/exported"
	nodeexported "github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"io/ioutil"
	"net/http"
)

const DEFAULTHTTPMETHOD = "POST"

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

func (r Relay) Validate(ctx sdk.Context, nodeVerify nodeexported.ValidatorI, hostedBlockchains HostedBlockchains,
	sessionBlockHeight int64, sessionNodeCount int, allActiveNodes []nodeexported.ValidatorI, app appexported.ApplicationI) sdk.Error {
	// check to see if the blockchain is empty
	if len(r.Blockchain) == 0 {
		return NewEmptyBlockchainError(ModuleName)
	}
	// check to see if the payload is empty
	if r.Payload.Data == "" || len(r.Payload.Data) == 0 {
		return NewEmptyPayloadDataError(ModuleName)
	}
	// ensure the blockchain is supported
	if !hostedBlockchains.ContainsFromString(r.Blockchain) {
		return NewUnsupportedBlockchainNodeError(ModuleName)
	}
	// check to see if non-native blockchain is staked for by the developer
	if _, contains := app.GetChains()[r.Blockchain]; !contains {
		return NewNotStakedBlockchainError(ModuleName)
	}
	// verify that node verify is of this session
	if _, err := SessionVerification(ctx, nodeVerify, app, r.Blockchain, sessionBlockHeight, sessionNodeCount, allActiveNodes); err != nil {
		return err
	}
	// check to see if the service proof is valid
	if err := r.Proof.Validate(app.GetMaxRelays().Int64(), hex.EncodeToString(nodeVerify.GetConsPubKey().Bytes())); err != nil {
		return err
	}
	// payload type to handle correctly
	if len((r.Payload).Method) == 0 {
		r.Payload.Method = DEFAULTHTTPMETHOD
	}
	return nil
}

// executes the relay on the non-native blockchain specified
func (r Relay) Execute(hostedBlockchains HostedBlockchains) (string, sdk.Error) {
	// retrieve the hosted blockchain url requested
	url, err := hostedBlockchains.GetChainURL(r.Blockchain)
	if err != nil {
		return "", err
	}
	// do basic http request on the relay
	res, er := executeHTTPRequest(r.Payload.Data, url, r.Payload.Method)
	return res, NewHTTPExecutionError(ModuleName, er)
}

// "executeHTTPRequest" takes in the raw json string and forwards it to the RPC endpoint
// todo improved http responses
func executeHTTPRequest(payload string, url string, method string) (string, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", NewHTTPStatusCodeError(ModuleName, resp.StatusCode)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

// store the proofs of work done for the relay batch
func (r Relay) HandleProof(ctx sdk.Context, sessionBlockHeight int64, maxNumberOfRelays int) sdk.Error {
	// grab the relay batch container
	rbs := GetAllProofs()
	header := PORHeader{
		ApplicationPubKey:  r.Proof.Token.ApplicationPublicKey,
		Chain:              r.Blockchain,
		SessionBlockHeight: sessionBlockHeight,
	}
	// add the proof to the proper batch
	return rbs.AddProof(header, r.Proof, maxNumberOfRelays)
}

type RelayResponse struct {
	Signature string // signature from the node in hex
	Response  string // response to relay
	Proof     Proof  // to be signed by the client
}

// node validates the response after signing
func (rr RelayResponse) Validate() sdk.Error {
	// the counter for the authorization must be >=0
	if rr.Proof.Index < 0 {
		return NewInvalidIncrementCounterError(ModuleName)
	}
	// cannot contain empty response
	if rr.Response == "" {
		return NewEmptyResponseError(ModuleName)
	}
	// cannot contain empty signature (nodes must be accountable)
	if rr.Signature == "" || len(rr.Signature) == crypto.SignatureSize {
		return NewResponseSignatureError(ModuleName)
	}
	return nil
}

// node signs the response before validating back
func (rr RelayResponse) Hash() []byte {
	return SHA3FromString(rr.Response + rr.Proof.HashString()) // todo standardize
}

// node signs the response before validating back
func (rr RelayResponse) HashString() string {
	return hex.EncodeToString(SHA3FromString(rr.Response + rr.Proof.HashString())) // todo standardize
}
