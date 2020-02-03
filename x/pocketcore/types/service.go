package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	Payload Payload    `json:"payload"`    // the data payload of the request
	Proof   RelayProof `json:"relayProof"` // the authentication scheme needed for work
}

func (r *Relay) Validate(ctx sdk.Context, node nodeexported.ValidatorI, hb HostedBlockchains, sessionBlockHeight int64,
	sessionNodeCount int, allNodes []nodeexported.ValidatorI, app appexported.ApplicationI) sdk.Error {
	// validate payload
	if err := r.Payload.Validate(); err != nil {
		return NewEmptyPayloadDataError(ModuleName)
	}
	// validate the RelayProof
	if err := r.Proof.Validate(app.GetMaxRelays().Int64(), len(app.GetChains()), sessionNodeCount, hb, node.GetPublicKey().RawString()); err != nil {
		return err
	}
	// generate the session
	session, err := NewSession(app.GetPublicKey().RawString(), r.Proof.Blockchain, BlockHashFromBlockHeight(ctx, sessionBlockHeight), sessionBlockHeight, allNodes, sessionNodeCount)
	if err != nil {
		return err
	}
	// validate the session
	err = session.Validate(ctx, node, app, sessionNodeCount)
	if err != nil {
		return err
	}
	// if the payload method is empty, set it to the default
	if r.Payload.Method == "" {
		r.Payload.Method = DEFAULTHTTPMETHOD
	}
	return nil
}

// executes the relay on the non-native blockchain specified
func (r Relay) Execute(hostedBlockchains HostedBlockchains) (string, sdk.Error) {
	// retrieve the hosted blockchain url requested
	url, err := hostedBlockchains.GetChainURL(r.Proof.Blockchain)
	if err != nil {
		return "", err
	}
	// do basic http request on the relay
	res, er := executeHTTPRequest(r.Payload.Data, url, r.Payload.Method, r.Payload.Headers)
	if er != nil {
		return res, NewHTTPExecutionError(ModuleName, er)
	}
	return res, nil
}

// store the proofs of work done for the relay batch
func (r Relay) HandleProof(ctx sdk.Context, sessionBlockHeight int64) sdk.Error {
	// add the RelayProof to the global (in memory) collection of proofs
	return GetAllInvoices().AddToInvoice(SessionHeader{
		ApplicationPubKey:  r.Proof.Token.ApplicationPublicKey,
		Chain:              r.Proof.Blockchain,
		SessionBlockHeight: sessionBlockHeight,
	}, r.Proof)
}

// the payload of the relay
type Payload struct {
	Data    string            `json:"data"`    // the actual data string for the external chain
	Method  string            `json:"method"`  // the http CRUD method
	Path    string            `json:"path"`    // the REST Pathx
	Headers map[string]string `json:"headers"` // http headers
}

func (p Payload) Validate() sdk.Error {
	if p.Data == "" && p.Path == "" {
		return NewEmptyPayloadDataError(ModuleName)
	}
	return nil
}

// response structure for the relay
type RelayResponse struct {
	Signature string     `json:"signature"`  // signature from the node in hex
	Response  string     `json:"payload"`    // response to relay
	Proof     RelayProof `json:"RelayProof"` // to be signed by the client
}

// node validates the response after signing
func (rr RelayResponse) Validate() sdk.Error { // todo more validaton
	// cannot contain empty response
	if rr.Response == "" {
		return NewEmptyResponseError(ModuleName)
	}
	// cannot contain empty signature (nodes must be accountable)
	if rr.Signature == "" || len(rr.Signature) == crypto.Ed25519SignatureSize {
		return NewResponseSignatureError(ModuleName)
	}
	return nil
}

// node signs the response before validating back
func (rr RelayResponse) Hash() []byte {
	seed, err := json.Marshal(relayResponse{
		Signature: "",
		Response:  rr.Response,
		Proof:     rr.Proof.HashString(),
	})
	if err != nil {
		panic(fmt.Sprintf("an error occured hashing the relay response:\n%v", err))
	}
	return Hash(seed)
}

// node signs the response before validating back
func (rr RelayResponse) HashString() string {
	return hex.EncodeToString(rr.Hash())
}

type relayResponse struct {
	Signature string `json:"signature"`
	Response  string `json:"payload"`
	Proof     string `json:"RelayProof"`
}

// "executeHTTPRequest" takes in the raw json string and forwards it to the RPC endpoint
func executeHTTPRequest(payload string, url string, method string, headers map[string]string) (string, error) { // todo improved http responses
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}
	if len(headers) == 0 { // def to json
		req.Header.Set("Content-Type", "application/json")
	} else {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
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
