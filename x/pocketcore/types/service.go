package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	appexported "github.com/pokt-network/pocket-core/x/apps/exported"
	nodeexported "github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"io/ioutil"
	"net/http"
)

// a read / write API request from a hosted (non native) blockchain
type Relay struct {
	Blockchain string  `json:"blockchain"`       // the non-native blockchain needed to service
	Payload    Payload `json:"payload"`          // the data payload of the request
	Proof      Proof   `json:"incrementCounter"` // the authentication scheme needed for work
}

func (r Relay) Validate(ctx sdk.Context, node nodeexported.ValidatorI, hb HostedBlockchains, sessionBlockHeight int64,
	sessionNodeCount int, allNodes []nodeexported.ValidatorI, app appexported.ApplicationI) sdk.Error {
	// validate blockchain
	if err := HashVerification(r.Blockchain); err != nil {
		return err
	}
	// validate payload
	if err := r.Payload.Validate(); err != nil {
		return NewEmptyPayloadDataError(ModuleName)
	}
	// validate the proof
	if err := r.Proof.Validate(app.GetMaxRelays().Int64(), hex.EncodeToString(node.GetConsPubKey().Bytes())); err != nil {
		return err
	}
	// ensure the blockchain is supported
	if !hb.ContainsFromString(r.Blockchain) {
		return NewUnsupportedBlockchainNodeError(ModuleName)
	}
	// generate the session
	session, err := NewSession(hex.EncodeToString(app.GetConsPubKey().Bytes()), r.Blockchain, BlockHashFromBlockHeight(ctx, sessionBlockHeight), sessionBlockHeight, allNodes, sessionNodeCount)
	if err != nil {
		return err
	}
	// validate the session
	err = session.Validate(ctx, node, app, sessionNodeCount)
	if err != nil {
		return err
	}
	// payload type to handle correctly
	if len((r.Payload).Method) == 0 {
		r.Payload.Method = DEFAULTHTTPMETHOD
	}
	return nil
}

const DEFAULTHTTPMETHOD = "POST"

// executes the relay on the non-native blockchain specified
func (r Relay) Execute(hostedBlockchains HostedBlockchains) (string, sdk.Error) {
	// retrieve the hosted blockchain url requested
	url, err := hostedBlockchains.GetChainURL(r.Blockchain)
	if err != nil {
		return "", err
	}
	// do basic http request on the relay
	res, er := executeHTTPRequest(r.Payload.Data, url, r.Payload.Method)
	if er != nil {
		return res, NewHTTPExecutionError(ModuleName, er)
	}
	return res, nil
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

// the payload of the relay
type Payload struct {
	Data   string `json:"data"`   // the actual data string for the external chain
	Method string `json:"method"` // the http CRUD method
	Path   string `json:"path"`   // the REST Pathx
}

func (p Payload) Validate() sdk.Error {
	if p.Data == "" && p.Path == "" {
		return NewEmptyPayloadDataError(ModuleName)
	}
	return nil
}

// response structure for the relay
type RelayResponse struct {
	Signature string // signature from the node in hex
	Response  string // response to relay
	Proof     Proof  // to be signed by the client
}

type relayResponse struct {
	sig   string
	resp  string
	proof string
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
	seed, err := json.Marshal(relayResponse{
		sig:   "",
		resp:  rr.Response,
		proof: rr.Proof.HashString(),
	})
	if err != nil {
		panic(err)
	}
	return SHA3FromBytes(seed)
}

// node signs the response before validating back
func (rr RelayResponse) HashString() string {
	return hex.EncodeToString(rr.Hash())
}

// "executeHTTPRequest" takes in the raw json string and forwards it to the RPC endpoint
func executeHTTPRequest(payload string, url string, method string) (string, error) { // todo improved http responses
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
