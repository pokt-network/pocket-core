package keeper

import (
	"bytes"
	"encoding/hex"
	"github.com/pokt-network/pocket-core/types"
	appexported "github.com/pokt-network/pocket-core/x/apps/exported"
	nodeexported "github.com/pokt-network/pocket-core/x/nodes/exported"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	"io/ioutil"
	"net/http"
)

func (k Keeper) ExecuteRelay(ctx sdk.Context, selfNode nodeexported.ValidatorI, relay pc.Relay, hostedBlockchains types.HostedBlockchains, sessionBlockIDHex string, allActiveNodes []nodeexported.ValidatorI, app appexported.ApplicationI) (string, error) {
	return Relay(relay).Execute(k, ctx, selfNode, hostedBlockchains, sessionBlockIDHex, allActiveNodes, app)
}

// a read / write API request from a hosted (non native) blockchain
type Relay pc.Relay

// executes the relay on the non-native blockchain specified
func (r Relay) Execute(keeper Keeper, ctx sdk.Context, selfNode nodeexported.ValidatorI, hostedBlockchains types.HostedBlockchains, sessionBlockIDHex string, allActiveNodes []nodeexported.ValidatorI, app appexported.ApplicationI) (string, error) {
	if err := r.Validate(keeper, ctx, selfNode, hostedBlockchains, sessionBlockIDHex, allActiveNodes, app); err != nil {
		return "", err
	}
	switch r.Payload.Type() {
	case pc.HTTP:
		url, err := hostedBlockchains.GetChainURL(r.Blockchain)
		if err != nil {
			return "", err
		}
		return executeHTTPRequest(r.Payload.Data, url, r.Payload.Method)
	}
	return "", pc.UnsupportedPayloadTypeError
}

func (r Relay) Validate(keeper Keeper, ctx sdk.Context, nodeVerify nodeexported.ValidatorI, hostedBlockchains types.HostedBlockchains, sessionBlockIDHex string, allActiveNodes []nodeexported.ValidatorI, app appexported.ApplicationI) error {
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
	if err := keeper.SessionVerification(ctx, nodeVerify,
		hex.EncodeToString(app.GetConsPubKey().Bytes()), // todo not absolutely sure that this conversion to hex string is accurate
		r.Blockchain,
		sessionBlockIDHex,
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
func (r Relay) StoreProofs(sessionBlockIDHex, chain string, maxNumberOfRelays int) error {
	// grab the relay batch container
	rbs := GetProofBatches()
	// add the proof to the proper batch
	return rbs.AddProof(r.Proof, sessionBlockIDHex, chain, maxNumberOfRelays)
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
