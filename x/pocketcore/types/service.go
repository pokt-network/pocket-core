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
	"math"
	"net/http"
	"strings"
)

const DEFAULTHTTPMETHOD = "POST"

// a read / write API request from a hosted (non native) blockchain
type Relay struct {
	Payload Payload    `json:"payload"` // the data payload of the request
	Meta    RelayMeta  `json:"meta"`    // metadata for the relay request
	Proof   RelayProof `json:"proof"`   // the authentication scheme needed for work
}

func (r *Relay) Validate(ctx sdk.Ctx, node nodeexported.ValidatorI, hb HostedBlockchains, sessionBlockHeight int64,
	sessionNodeCount int, allNodes []nodeexported.ValidatorI, app appexported.ApplicationI) sdk.Error {
	// validate payload
	if err := r.Payload.Validate(); err != nil {
		return NewEmptyPayloadDataError(ModuleName)
	}
	// validate the metadata
	if err := r.Meta.Validate(ctx); err != nil {
		return err
	}
	// validate the relay hash = request hash
	if r.Proof.RequestHash != r.RequestHashString() {
		return NewRequestHashError(ModuleName)
	}
	// ensure the blockchain is supported locally
	if !hb.ContainsFromString(r.Proof.Blockchain) {
		return NewUnsupportedBlockchainNodeError(ModuleName)
	}
	evidenceHeader := SessionHeader{
		ApplicationPubKey:  r.Proof.Token.ApplicationPublicKey,
		Chain:              r.Proof.Blockchain,
		SessionBlockHeight: r.Proof.SessionBlockHeight,
	}
	// validate unique relay
	totalRelays := GetTotalProofs(evidenceHeader, RelayEvidence)
	// get evidence key by proof
	if !IsUniqueProof(evidenceHeader, r.Proof) {
		return NewDuplicateProofError(ModuleName)
	}
	// validate not over service
	if totalRelays >= int64(math.Ceil(float64(app.GetMaxRelays().Int64())/float64(len(app.GetChains())))/(float64(sessionNodeCount))) {
		return NewOverServiceError(ModuleName)
	}
	// validate the Proof
	if err := r.Proof.ValidateLocal(app.GetChains(), sessionNodeCount, sessionBlockHeight, node.GetPublicKey().RawString()); err != nil {
		return err
	}
	// get the sessionContext
	sessionContext, er := ctx.PrevCtx(sessionBlockHeight)
	if er != nil {
		return sdk.ErrInternal(er.Error())
	}
	// generate the header
	header := SessionHeader{
		ApplicationPubKey:  app.GetPublicKey().RawString(),
		Chain:              r.Proof.Blockchain,
		SessionBlockHeight: sessionBlockHeight,
	}
	// check cache
	session, found := GetSession(header)
	// if not found generate the session
	if !found {
		var err sdk.Error
		session, err = NewSession(header, BlockHash(sessionContext), allNodes, sessionNodeCount)
		if err != nil {
			return err
		}
		// add to cache
		SetSession(session)
	}
	// validate the session
	err := session.Validate(ctx, node, app, sessionNodeCount)
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
	url = strings.Trim(url, `/`) + "/" + strings.Trim(r.Payload.Path, `/`)
	// do basic http request on the relay
	res, er := executeHTTPRequest(r.Payload.Data, url, r.Payload.Method, r.Payload.Headers)
	if er != nil {
		return res, NewHTTPExecutionError(ModuleName, er)
	}
	return res, nil
}

func (r Relay) RequestHash() []byte {

	relay := struct {
		Payload Payload   `json:"payload"` // the data payload of the request
		Meta    RelayMeta `json:"meta"`    // metadata for the relay request
	}{r.Payload, r.Meta}

	res, err := json.Marshal(relay)
	if err != nil {
		panic(fmt.Sprintf("cannot marshal relay request hash: %s", err.Error()))
	}
	return res
}

func (r Relay) RequestHashString() string {
	return hex.EncodeToString(r.RequestHash())
}

// the payload of the relay
type Payload struct {
	Data    string            `json:"data"`              // the actual data string for the external chain
	Method  string            `json:"method"`            // the http CRUD method
	Path    string            `json:"path"`              // the REST Path
	Headers map[string]string `json:"headers,omitempty"` // http headers
}

func (p Payload) Bytes() []byte {
	bz, err := json.Marshal(Payload{Data: p.Data, Method: p.Method, Path: p.Path})
	if err != nil {
		panic(fmt.Sprintf("an error occured converting the payload to bytes:\n%v", err))
	}
	return bz
}

func (p Payload) Hash() []byte {
	return Hash(p.Bytes())
}

func (p Payload) HashString() string {
	return hex.EncodeToString(p.Hash())
}

func (p Payload) Validate() sdk.Error {
	if p.Data == "" && p.Path == "" {
		return NewEmptyPayloadDataError(ModuleName)
	}
	return nil
}

type payload struct {
	Data   string `json:"data"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

func (p Payload) MarshalJSON() ([]byte, error) {
	pay := payload{
		Data:   p.Data,
		Method: p.Method,
		Path:   p.Path,
	}
	return json.Marshal(pay)
}

type RelayMeta struct {
	BlockHeight int64 `json:"block_height"`
}

func (m RelayMeta) Validate(ctx sdk.Ctx) sdk.Error {
	if ctx.BlockHeight()+5 < m.BlockHeight || ctx.BlockHeight()-5 > m.BlockHeight {
		return NewOutOfSyncRequestError(ModuleName)
	}
	return nil
}

// response structure for the relay
type RelayResponse struct {
	Signature string     `json:"signature"` // signature from the node in hex
	Response  string     `json:"payload"`   // response to relay
	Proof     RelayProof `json:"Proof"`     // to be signed by the client
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
	Proof     string `json:"Proof"`
}

type ChallengeResponse struct {
	Response string `json:"response"`
}

type DispatchResponse struct {
	Session     Session `json:"session"`
	BlockHeight int64   `json:"block_height"`
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

func sortJSONResponse(response string) string {
	var rawJSON map[string]interface{}
	if err := json.Unmarshal([]byte(response), &rawJSON); err != nil {
		return response // couldn't unmarshal into json
	}
	bz, err := json.Marshal(rawJSON)
	if err != nil {
		return response
	}
	return string(bz)
}
