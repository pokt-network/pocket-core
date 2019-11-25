package keeper

import (
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/pokt-network/pocket-core/types"
	appexported "github.com/pokt-network/pocket-core/x/apps/exported"
	nodeexported "github.com/pokt-network/pocket-core/x/nodes/exported"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	"io/ioutil"
	"net/http"
)

func (k Keeper) HandleRelay(ctx sdk.Context, relay pc.Relay) (*pc.RelayResponse, error) {
	// get the latest session block
	sessionBlock := k.GetLatestSessionBlock(ctx)
	sessionBlockHash := hex.EncodeToString(sessionBlock.BlockHeader().GetLastBlockId().Hash)
	sessionBlockHeight := sessionBlock.BlockHeight()
	// retrieve all nodes available
	allNodes := k.GetAllNodes(ctx)
	// get self node
	selfNode := k.GetSelfNode(ctx)
	// get hosted blockchains
	hb := k.GetHostedBlockchains(ctx)
	// get application that staked for client
	app, found := k.GetAppFromPublicKey(ctx, relay.Proof.Token.ApplicationPublicKey)
	if !found {
		return nil, errors.New("no app found for addr")
	}
	// validate the next relay
	if err := Relay(relay).Validate(k, ctx, selfNode, hb, sessionBlockHash, sessionBlockHeight, allNodes, app); err != nil {
		return nil, err
	}
	// get proofs for this session

	// store the previous proof
	if err := Relay(relay).StoreProofs(k, ctx, sessionBlockHash, sessionBlockHeight, int(app.GetMaxRelays().Int64())); err != nil {
		return nil, err
	}
	// attempt to execute
	respPayload, err := Relay(relay).Execute(hb)
	if err != nil {
		return nil, err
	}
	// generate response object
	resp := &pc.RelayResponse{
		Response: respPayload,
		ServiceAuth: pc.Proof{
			Counter: 0, // todo
		},
	}
	// sign the response
	err = resp.Sign()
	if err != nil {
		return nil, err
	}
	// complete
	return resp, nil
}

// a read / write API request from a hosted (non native) blockchain
type Relay pc.Relay

// executes the relay on the non-native blockchain specified
func (r Relay) Execute(hostedBlockchains types.HostedBlockchains) (string, error) {
	// next check to see what type of payload it has
	switch r.Payload.Type() {
	case pc.HTTP:
		// retrieve the hosted blockchain url requested
		url, err := hostedBlockchains.GetChainURL(r.Blockchain)
		if err != nil {
			return "", err
		}
		// do basic http request on the relay
		return executeHTTPRequest(r.Payload.Data, url, r.Payload.Method)
	}
	return "", pc.UnsupportedPayloadTypeError
}

func (r Relay) Validate(keeper Keeper, ctx sdk.Context, nodeVerify nodeexported.ValidatorI, hostedBlockchains types.HostedBlockchains, sessionBlockIDHex string, sessionBlockHeight int64, allActiveNodes []nodeexported.ValidatorI, app appexported.ApplicationI) error {
	// check to see if the blockchain is empty
	if len(r.Blockchain) == 0 {
		return pc.EmptyBlockchainError
	}
	// check to see if the payload is empty
	if r.Payload.Data == "" || len(r.Payload.Data) == 0 {
		return pc.EmptyPayloadDataError
	}
	// ensure the blockchain is supported
	if !hostedBlockchains.ContainsFromString(r.Blockchain) {
		return pc.UnsupportedBlockchainError
	}
	// check to see if non native blockchain is staked for by the developer
	if _, contains := app.GetChains()[r.Blockchain]; !contains {
		return pc.NotStakedBlockchainError
	}
	// verify that node (self) is of this session
	if _, err := keeper.SessionVerification(ctx, nodeVerify,
		app,
		r.Blockchain,
		sessionBlockIDHex,
		sessionBlockHeight,
		allActiveNodes); err != nil {
		return err
	}
	// check to see if the service proof is valid
	if err := r.Proof.Validate(app.GetMaxRelays().Int64()); err != nil {
		return pc.NewServiceProofError(err)
	}
	if r.Payload.Type() == pc.HTTP {
		if len((r.Payload).Method) == 0 {
			r.Payload.Method = pc.DEFAULTHTTPMETHOD
		}
	}
	return nil
}

// store the proofs of work done for the relay batch
func (r Relay) StoreProofs(k Keeper, ctx sdk.Context, sessionBlockIDHex string, sessionBlockHeight int64, maxNumberOfRelays int) error {
	// grab the relay batch container
	rbs := k.GetProofBatches()
	// add the proof to the proper batch
	return rbs.AddProof(r.Proof, sessionBlockIDHex, r.Blockchain, sessionBlockHeight, maxNumberOfRelays)
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
		return "", pc.NewHTTPStatusCodeError(resp.StatusCode)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
