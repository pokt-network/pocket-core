package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	appexported "github.com/pokt-network/pocket-core/x/apps/exported"
	nodeexported "github.com/pokt-network/pocket-core/x/nodes/exported"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const DEFAULTHTTPMETHOD = "POST"

// "Relay" - A read / write API request from a hosted (non native) external blockchain
type Relay struct {
	Payload Payload    `json:"payload"` // the data payload of the request
	Meta    RelayMeta  `json:"meta"`    // metadata for the relay request
	Proof   RelayProof `json:"proof"`   // the authentication scheme needed for work
}

// "Validate" - Checks the validity of a relay request using store data
func (r *Relay) Validate(ctx sdk.Ctx, keeper PosKeeper, node nodeexported.ValidatorI, hb *HostedBlockchains, sessionBlockHeight int64,
	sessionNodeCount int, app appexported.ApplicationI) (maxPossibleRelays sdk.Int, err sdk.Error) {
	// validate payload
	if err := r.Payload.Validate(); err != nil {
		return sdk.ZeroInt(), NewEmptyPayloadDataError(ModuleName)
	}
	// validate the metadata
	if err := r.Meta.Validate(ctx); err != nil {
		return sdk.ZeroInt(), err
	}
	// validate the relay merkleHash = request merkleHash
	if r.Proof.RequestHash != r.RequestHashString() {
		return sdk.ZeroInt(), NewRequestHashError(ModuleName)
	}
	// ensure the blockchain is supported locally
	if !hb.Contains(r.Proof.Blockchain) {
		return sdk.ZeroInt(), NewUnsupportedBlockchainNodeError(ModuleName)
	}
	evidenceHeader := SessionHeader{
		ApplicationPubKey:  r.Proof.Token.ApplicationPublicKey,
		Chain:              r.Proof.Blockchain,
		SessionBlockHeight: r.Proof.SessionBlockHeight,
	}
	maxPossibleRelays = MaxPossibleRelays(app, int64(sessionNodeCount))
	// validate unique relay
	evidence, totalRelays := GetTotalProofs(evidenceHeader, RelayEvidence, maxPossibleRelays)
	// get evidence key by proof
	if !IsUniqueProof(r.Proof, evidence) {
		return sdk.ZeroInt(), NewDuplicateProofError(ModuleName)
	}
	// validate not over service
	if sdk.NewInt(totalRelays).GTE(maxPossibleRelays) {
		return sdk.ZeroInt(), NewOverServiceError(ModuleName)
	}
	// validate the Proof
	if err := r.Proof.ValidateLocal(app.GetChains(), sessionNodeCount, sessionBlockHeight, node.GetPublicKey().RawString()); err != nil {
		return sdk.ZeroInt(), err
	}
	// get the sessionContext
	sessionContext, er := ctx.PrevCtx(sessionBlockHeight)
	if er != nil {
		return sdk.ZeroInt(), sdk.ErrInternal(er.Error())
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
		session, err = NewSession(sessionContext, ctx, keeper, header, BlockHash(sessionContext), sessionNodeCount)
		if err != nil {
			return sdk.ZeroInt(), err
		}
		// add to cache
		SetSession(session)
	}
	// validate the session
	err = session.Validate(node, app, sessionNodeCount)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	// if the payload method is empty, set it to the default
	if r.Payload.Method == "" {
		r.Payload.Method = DEFAULTHTTPMETHOD
	}
	return maxPossibleRelays, nil
}

// "Execute" - Attempts to do a request on the non-native blockchain specified
func (r Relay) Execute(hostedBlockchains *HostedBlockchains) (string, sdk.Error) {
	// retrieve the hosted blockchain url requested
	chain, err := hostedBlockchains.GetChain(r.Proof.Blockchain)
	if err != nil {
		// metric track
		GlobalServiceMetric().AddErrorFor(r.Proof.Blockchain)
		return "", err
	}
	url := strings.Trim(chain.URL, `/`)
	if len(r.Payload.Path) > 0 {
		url = url + "/" + strings.Trim(r.Payload.Path, `/`)
	}
	// do basic http request on the relay
	res, er := executeHTTPRequest(r.Payload.Data, url, GlobalPocketConfig.UserAgent, chain.BasicAuth, r.Payload.Method, r.Payload.Headers)
	if er != nil {
		// metric track
		GlobalServiceMetric().AddErrorFor(r.Proof.Blockchain)
		return res, NewHTTPExecutionError(ModuleName, er)
	}
	return res, nil
}

// "Bytes" - Returns the bytes representation of the Relay
func (r Relay) Bytes() []byte {
	//Anonymous Struct used because of #742 empty proof object being marshalled
	relay := struct {
		Payload Payload   `json:"payload"` // the data payload of the request
		Meta    RelayMeta `json:"meta"`    // metadata for the relay request
	}{r.Payload, r.Meta}
	res, err := json.Marshal(relay)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot marshal relay request merkleHash: %s", err.Error()))
	}
	return res
}

// "Requesthash" - The cryptographic merkleHash representation of the request
func (r Relay) RequestHash() []byte {
	return Hash(r.Bytes())
}

// "RequestHashString" - The hex string representation of the request merkleHash
func (r Relay) RequestHashString() string {
	return hex.EncodeToString(r.RequestHash())
}

// "Payload" - A data being sent to the non-native chain
type Payload struct {
	Data    string            `json:"data"`              // the actual data string for the external chain
	Method  string            `json:"method"`            // the http CRUD method
	Path    string            `json:"path"`              // the REST Path
	Headers map[string]string `json:"headers,omitempty"` // http headers
}

// "Bytes" - The bytes reprentation of a payload object
func (p Payload) Bytes() []byte {
	bz, err := json.Marshal(p)
	if err != nil {
		log.Fatalf(fmt.Errorf("an error occured converting the payload to bytes:\n%v", err).Error())
	}
	return bz
}

// "Hash" - The cryptographic merkleHash representation of the payload object
func (p Payload) Hash() []byte {
	return Hash(p.Bytes())
}

// "HashString" - The hex encoded string representation of the payload object
func (p Payload) HashString() string {
	return hex.EncodeToString(p.Hash())
}

// "Validate" - Validity check for the payload object
func (p Payload) Validate() sdk.Error {
	if p.Data == "" && p.Path == "" {
		return NewEmptyPayloadDataError(ModuleName)
	}
	return nil
}

// "payload" - A structure used for custom json marshalling/unmarshalling
type payload struct {
	Data    string            `json:"data"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
}

// "MarshalJSON" - Overrides json marshalling
func (p Payload) MarshalJSON() ([]byte, error) {
	pay := payload{ //nolint:golint,gosimple
		Data:    p.Data,
		Method:  p.Method,
		Path:    p.Path,
		Headers: p.Headers,
	}
	return json.Marshal(pay)
}

// "RelayMeta" - Metadata that is included in the relay request
type RelayMeta struct {
	BlockHeight int64 `json:"block_height"` // the block height when the request is made
}

// "Validate" - Validates the relay meta object
func (m RelayMeta) Validate(ctx sdk.Ctx) sdk.Error {
	// ensures the block height is within the acceptable range
	if ctx.BlockHeight()+int64(GlobalPocketConfig.ClientBlockSyncAllowance) < m.BlockHeight || ctx.BlockHeight()-int64(GlobalPocketConfig.ClientBlockSyncAllowance) > m.BlockHeight {
		return NewOutOfSyncRequestError(ModuleName)
	}
	return nil
}

func InitClientBlockAllowance(allowance int) {
	GlobalPocketConfig.ClientBlockSyncAllowance = allowance
}

// response structure for the relay
type RelayResponse struct {
	Signature string     `json:"signature"` // signature from the node in hex
	Response  string     `json:"payload"`   // response to relay
	Proof     RelayProof `json:"proof"`     // to be signed by the client
}

// "Validate" - The node validates the response after signing
func (rr RelayResponse) Validate() sdk.Error {
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

// "Hash" - The cryptographic merkleHash representation of the relay response
func (rr RelayResponse) Hash() []byte {
	seed, err := json.Marshal(relayResponse{
		Signature: "",
		Response:  rr.Response,
		Proof:     rr.Proof.HashString(),
	})
	if err != nil {
		log.Fatalf(fmt.Errorf("an error occured hashing the relay response:\n%v", err).Error())
	}
	return Hash(seed)
}

// "HashString" - The hex string representation of the merkleHash
func (rr RelayResponse) HashString() string {
	return hex.EncodeToString(rr.Hash())
}

// "relayResponse" - a structure used for custom json
type relayResponse struct {
	Signature string `json:"signature"`
	Response  string `json:"payload"`
	Proof     string `json:"Proof"`
}

// "ChallengeReponse" - The response object used in challenges
type ChallengeResponse struct {
	Response string `json:"response"`
}

// "DispatchResponse" - The response object used in dispatching
type DispatchResponse struct {
	Session     Session `json:"session"`
	BlockHeight int64   `json:"block_height"`
}

// "executeHTTPRequest" takes in the raw json string and forwards it to the RPC endpoint
func executeHTTPRequest(payload, url, userAgent string, basicAuth BasicAuth, method string, headers map[string]string) (string, error) {
	// generate an http request
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}
	if basicAuth.Username != "" {
		req.SetBasicAuth(basicAuth.Username, basicAuth.Password)
	}
	if userAgent == "" {
		req.Header.Set("User-Agent", userAgent)
	}
	// add headers if needed
	if len(headers) == 0 {
		req.Header.Set("Content-Type", "application/json")
	} else {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	// execute the request
	resp, err := (&http.Client{Timeout: GetRPCTimeout() * time.Millisecond}).Do(req)
	if err != nil {
		return "", err
	}
	// read all bz
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if GlobalPocketConfig.JSONSortRelayResponses {
		body = []byte(sortJSONResponse(string(body)))
	}
	// return
	return string(body), nil
}

// "sortJSONResponse" - sorts json from a relay response
func sortJSONResponse(response string) string {
	var rawJSON map[string]interface{}
	// unmarshal into json
	if err := json.Unmarshal([]byte(response), &rawJSON); err != nil {
		return response
	}
	// marshal into json
	bz, err := json.Marshal(rawJSON)
	if err != nil {
		return response
	}
	return string(bz)
}
