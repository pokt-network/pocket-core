package mesh

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pokt-network/pocket-core/app"
	sdk "github.com/pokt-network/pocket-core/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

// storeRelay - persist relay to disk
func storeRelay(relay *pocketTypes.Relay) {
	hash := relay.RequestHash()
	LogRelay(relay, "storing relay", LogLvlDebug)
	rb, err := json.Marshal(relay)
	if err != nil {
		LogRelay(relay, fmt.Sprintf("error=%s marshaling relay", CleanError(err.Error())), LogLvlError)
		return
	}

	err = relaysCacheDb.Put(hash, rb)
	if err != nil {
		LogRelay(relay, fmt.Sprintf("error=%s adding relay to cache", CleanError(err.Error())), LogLvlError)
	}

	return
}

// decodeCacheRelay - decode []byte relay from cache to pocketTypes.Relay
func decodeCacheRelay(body []byte) (relay *pocketTypes.Relay) {
	if err := json.Unmarshal(body, &relay); err != nil {
		LogRelay(relay, fmt.Sprintf("error=%s decoding cache relay", CleanError(err.Error())), LogLvlError)
		// todo: delete key from cache?
		deleteCacheRelay(relay) // because is malformed, probably.
		return nil
	}
	return
}

// deleteCacheRelay - delete a key from relay cache
func deleteCacheRelay(relay *pocketTypes.Relay) {
	hash := relay.RequestHash()
	err := relaysCacheDb.Delete(hash)
	if err != nil {
		LogRelay(relay, fmt.Sprintf("error=%s deleting relay from cache", CleanError(err.Error())), LogLvlError)
		return
	}

	return
}

// evaluateServicerError - this will change internalCache[hash].IsValid bool depending on the result of the evaluation
func evaluateServicerError(r *pocketTypes.Relay, err *SdkErrorResponse) (isSessionStillValid bool) {
	isSessionStillValid = !ShouldInvalidateSession(err.Code) // we should not retry if is invalid

	if isSessionStillValid {
		return isSessionStillValid
	}

	sessionStorage.InvalidateNodeSession(r, err)

	return
}

// addRelayToQueue - add relay to worker queue to be notified
func addRelayToQueue(r *pocketTypes.Relay) func() {
	return func() {
		notifyServicer(r)
	}
}

// notifyServicer - call servicer to ack about the processed relay.
func notifyServicer(r *pocketTypes.Relay) {
	relayTimeStart := time.Now()
	requeue := false

	// discard this relay at the end of this function, to end this function the servicer will be retried N times
	defer func(_r *pocketTypes.Relay, requeue *bool) {
		if *requeue {
			LogRelay(r, "avoid deleting relay from cache because it is re-queued", LogLvlInfo)
			return
		}
		deleteCacheRelay(r)
	}(r, &requeue)

	result := RPCRelayResponse{}
	ctx := context.WithValue(context.Background(), "result", &result)
	jsonData, e1 := json.Marshal(r)
	if e1 != nil {
		LogRelay(r, fmt.Sprintf("notify - error=%s encoding relay", CleanError(e1.Error())), LogLvlError)
		return
	}

	ns, e2 := sessionStorage.GetSession(r)

	if e2 != nil {
		// just to have a different text that provide better understanding of what is going on the logs.
		LogRelay(r, fmt.Sprintf("notify - unable to notify due to error=%s getting session", e2.Error), LogLvlError)
		ns.ServicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(
			r.Proof.Blockchain, &ns.ServicerNode.Address,
			true, GetSessionErrorType, "500",
		)
		return
	}

	// Add relay to bloom filter if it exists
	if ns.bloomFilter != nil {
		ns.bloomFilter.Add(r.Proof.Bytes())
	}

	// Check if we need already queried the session and it was invalid
	if ns.Queried && !ns.IsValid {
		LogRelay(r, fmt.Sprintf("notify - unable to notify because session was invalidated with error=%s", CleanError(e1.Error())), LogLvlError)
		ns.ServicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(
			r.Proof.Blockchain, &ns.ServicerNode.Address,
			true, InvalidSessionType, fmt.Sprintf("%d", ns.Error.Code),
		)
		return
	}

	// check if we haven't queried it before and it's not inflight to be queued
	if !ns.Queried && !ns.Queue {
		// we should re-schedule it to later on the session is validated, or even validated but invalid to be discarded.
		LogRelay(r, "notify - relay re-queued because session is not queried and not in queue neither", LogLvlInfo)
		requeue = true
		ns.ServicerNode.Node.Worker.Submit(addRelayToQueue(r))
		return
	}

	// Safety measure to not ask for an app session within range
	if !ns.ServicerNode.Node.CanHandleRelayWithinTolerance(r.Proof.SessionBlockHeight) {
		LogRelay(r, fmt.Sprintf(
			"notify - unable to delivery because relay session height is not within tolerance of fullNode session_height=%d",
			ns.ServicerNode.Node.GetLatestSessionBlockHeight(),
		), LogLvlError)
		ns.ServicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(
			r.Proof.Blockchain, &ns.ServicerNode.Address,
			true, SessionHeightOutOfRangeType, fmt.Sprintf("%d", ns.Error.Code),
		)
		return
	}

	requestURL := fmt.Sprintf(
		"%s%s?chain=%s&app=%s",
		ns.ServicerNode.Node.URL,
		ServicerRelayEndpoint,
		r.Proof.Blockchain,
		r.Proof.Token.ApplicationPublicKey,
	)
	req, e3 := retryablehttp.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(jsonData))
	req.Header.Set(AuthorizationHeader, servicerAuthToken.Value)
	if e3 != nil {
		LogRelay(r, fmt.Sprintf("notify - error=%s formatting url to call fullNode of servicer", e3.Error()), LogLvlError)
		ns.ServicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(
			r.Proof.Blockchain, &ns.ServicerNode.Address,
			true, NotifyRequestErrorType, "500",
		)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(ServicerHeader, ns.ServicerAddress)
	if app.GlobalMeshConfig.UserAgent != "" {
		req.Header.Set("User-Agent", app.GlobalMeshConfig.UserAgent)
	}
	resp, e4 := relaysClient.Do(req)

	if e4 != nil {
		LogRelay(r, fmt.Sprintf("notify - error=%s dispatching relay to fullNode", CleanError(e4.Error())), LogLvlError)
		ns.ServicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(
			r.Proof.Blockchain, &ns.ServicerNode.Address,
			true, NotifyRequestErrorType, "500",
		)
		return
	}

	defer func(Body io.ReadCloser) {
		e5 := Body.Close()
		if e5 != nil {
			LogRelay(r, fmt.Sprintf(
				"notify - error=%s closing dispatch notification response body",
				CleanError(e5.Error()),
			), LogLvlError)
			return
		}
	}(resp.Body)

	// read the body just to allow http 1.x be able to reuse the connection
	_, e6 := ioutil.ReadAll(resp.Body)
	if e6 != nil {
		LogRelay(r, fmt.Sprintf(
			"notify - error=%s parsing response from endpoint=%s at fullNode",
			CleanError(e1.Error()), ServicerRelayEndpoint,
		), LogLvlError)
		ns.ServicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(
			r.Proof.Blockchain, &ns.ServicerNode.Address,
			true, NotifyResponseErrorType, fmt.Sprintf("%d", resp.StatusCode),
		)
		return
	}

	isSuccess := resp.StatusCode == 200

	if result.Dispatch != nil && result.Dispatch.BlockHeight > ns.ServicerNode.Node.Status.Height {
		ns.ServicerNode.Node.Status.Height = result.Dispatch.BlockHeight
	}

	if !isSuccess {
		LogRelay(r, fmt.Sprintf(
			"notify - relay rejected by fullNode with message=%s code=%d codespace=%s",
			result.Error.Error, result.Error.Code, result.Error.Codespace,
		), LogLvlError)
		ns.ServicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(
			r.Proof.Blockchain, &ns.ServicerNode.Address,
			true, NotifyResponseErrorType, fmt.Sprintf("%d", resp.StatusCode),
		)
		evaluateServicerError(r, result.Error)
	} else {
		LogRelay(r, "notify - servicer processed relay successfully", LogLvlDebug)

		exhausted := ns.CountRelay()

		if exhausted {
			LogRelay(r, "notify - servicer exhaust relays", LogLvlDebug)
		} else {
			LogRelay(r, fmt.Sprintf("notify - servicer has %d remaining relays", ns.RemainingRelays), LogLvlDebug)
		}

		// track the notify relay time
		relayDuration := time.Since(relayTimeStart)
		ns.ServicerNode.Node.MetricsWorker.AddServiceMetricRelayFor(
			r, &ns.ServicerNode.Address,
			relayDuration, true,
		)
	}

	return
}

// execute - Attempts to do a request on the non-native blockchain specified
func execute(r *pocketTypes.Relay, hostedBlockchains *pocketTypes.HostedBlockchains, servicerNode *servicer) (string, sdk.Error) {
	address := &servicerNode.Address
	start := time.Now()
	code := 0
	defer func(c string, t time.Time, code *int, s *servicer) {
		if *code == 0 {
			// omit metric on non-called chains.
			return
		}
		s.Node.MetricsWorker.AddChainMetricFor(c, time.Since(t), *code)
	}(r.Proof.Blockchain, start, &code, servicerNode)
	node := GetNodeFromAddress(address.String())

	if node == nil {
		node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, address, false, ServicerNotFoundStatusType, "500")
		return "", sdk.ErrInternal("failed to find correct servicer PK")
	}

	// retrieve the hosted blockchain url requested
	chain, err := hostedBlockchains.GetChain(r.Proof.Blockchain)
	if err != nil {
		// metric track
		node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, address, false, ChainNotFoundStatusType, "500")
		return "", err
	}

	// do basic http request on the relay
	res, er, c := ExecuteBlockchainHTTPRequest(r.Payload, chain)
	if er != nil {
		// metric track
		node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, address, false, ChainStatusType, fmt.Sprintf("%d", code))
		return res, pocketTypes.NewHTTPExecutionError(ModuleName, er)
	}

	code = c

	if code >= 400 {
		node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, address, false, ChainStatusType, fmt.Sprintf("%d", code))
	}

	return res, nil
}

// processRelay - call execute and create RelayResponse or Error in case. Also trigger relay metrics.
func processRelay(relay *pocketTypes.Relay) (*pocketTypes.RelayResponse, sdk.Error) {
	LogRelay(relay, "handler - processing relay", LogLvlDebug)

	servicerAddress, e := GetAddressFromPubKeyAsString(relay.Proof.ServicerPubKey)
	if e != nil {
		return nil, sdk.ErrInternal("could not convert servicer hex to public key")
	}

	servicerNode, ok := servicerMap.Load(servicerAddress)
	if !ok {
		return nil, sdk.ErrInternal("failed to find correct servicer PK")
	}

	// attempt to execute
	respPayload, err := execute(relay, chains, servicerNode)
	if err != nil {
		LogRelay(relay, fmt.Sprintf("handler - call blockchain return error=%s", CleanError(err.Error())), LogLvlError)
		return nil, err
	}

	// generate response object
	resp := &pocketTypes.RelayResponse{
		Response: respPayload,
		Proof:    relay.Proof,
	}

	// sign the response
	sig, er := servicerNode.PrivateKey.Sign(resp.Hash())
	if er != nil {
		LogRelay(relay, fmt.Sprintf("handler - unable to sign relay response due to error=%s", CleanError(err.Error())), LogLvlError)
		return nil, pocketTypes.NewKeybaseError(pocketTypes.ModuleName, er)
	}
	// attach the signature in hex to the response
	resp.Signature = hex.EncodeToString(sig)
	return resp, nil
}

// validate - evaluate relay to understand if should or not processed.
func validate(r *pocketTypes.Relay) sdk.Error {
	logger.Debug(fmt.Sprintf("validating relay %s", r.RequestHashString()))
	// validate payload
	if err := r.Payload.Validate(); err != nil {
		return pocketTypes.NewEmptyPayloadDataError(ModuleName)
	}
	// validate appPubKey
	if err := pocketTypes.PubKeyVerification(r.Proof.Token.ApplicationPublicKey); err != nil {
		return err
	}
	// validate chain
	if err := pocketTypes.NetworkIdentifierVerification(r.Proof.Blockchain); err != nil {
		return pocketTypes.NewEmptyChainError(ModuleName)
	}
	// validate the relay merkleHash = request merkleHash
	if r.Proof.RequestHash != r.RequestHashString() {
		return pocketTypes.NewRequestHashError(ModuleName)
	}
	// ensure the blockchain is supported locally
	if !chains.Contains(r.Proof.Blockchain) {
		return pocketTypes.NewUnsupportedBlockchainNodeError(ModuleName)
	}
	// validate servicer public key
	servicerAddress, e := GetAddressFromPubKeyAsString(r.Proof.ServicerPubKey)
	if e != nil {
		return sdk.ErrInternal("could not convert servicer hex to public key")
	}
	// load servicer from servicer map, if not there maybe the servicer is pk is not loaded
	if servicerNode, ok := servicerMap.Load(servicerAddress); !ok {
		return sdk.ErrInternal("failed to find correct servicer PK")
	} else if r.Proof.SessionBlockHeight <= servicerNode.Node.GetLatestSessionBlockHeight() && !servicerNode.Node.CanHandleRelayWithinTolerance(r.Proof.SessionBlockHeight) {
		return pocketTypes.NewInvalidBlockHeightError(ModuleName)
	}

	ns, e1 := sessionStorage.GetSession(r)

	if e1 != nil {
		LogRelay(r, "handler - unable to get session from session storage", LogLvlError)
		return pocketTypes.NewInvalidSessionKeyError(ModuleName, errors.New(e1.Error))
	}

	if ns == nil {
		err := errors.New(fmt.Sprintf(
			"session not found for app=%s chain=%s servicer=%s ",
			r.Proof.Token.ApplicationPublicKey,
			r.Proof.Blockchain,
			servicerAddress,
		))
		LogRelay(r, "handler - session not found", LogLvlError)
		return pocketTypes.NewInvalidSessionKeyError(ModuleName, err)
	}

	if ns.bloomFilter != nil && ns.bloomFilter.Test(r.Proof.Bytes()) {
		return pocketTypes.NewDuplicateProofError(ModuleName)
	}

	if ns.IsValid {
		return nil
	}

	if ns.Error != nil {
		return NewPocketSdkErrorFromSdkError(ns.Error)
	}

	// Fallback invalid session
	e2 := errors.New(fmt.Sprintf(
		"invalid session=%s session_height=%d for app=%s chain=%s servicer=%s",
		ns.Key,
		ns.BlockHeight,
		r.Proof.Token.ApplicationPublicKey,
		r.Proof.Blockchain,
		ns.ServicerAddress,
	))
	return pocketTypes.NewInvalidSessionKeyError(ModuleName, e2)

}

// HandleRelay - evaluate node status, validate relay payload and call processRelay
func HandleRelay(r *pocketTypes.Relay) (res *pocketTypes.RelayResponse, dispatch *DispatchResponse, err error) {
	relayTimeStart := time.Now()
	servicerAddress, e := GetAddressFromPubKeyAsString(r.Proof.ServicerPubKey)

	if e != nil {
		return nil, nil, errors.New("could not convert servicer hex to public key")
	}

	servicerNode, ok := servicerMap.Load(servicerAddress)
	if !ok {
		return nil, nil, errors.New("failed to find correct servicer PK")
	}

	if servicerNode.Node.Status == nil {
		return nil, nil, fmt.Errorf("pocket node is currently unavailable")
	}

	if servicerNode.Node.Status.IsStarting {
		return nil, nil, fmt.Errorf("pocket node is unable to retrieve synced status from tendermint node, cannot service in this state")
	}

	if servicerNode.Node.Status.IsCatchingUp {
		return nil, nil, fmt.Errorf("pocket node is currently syncing to the blockchain, cannot service in this state")
	}

	err = validate(r)

	if err != nil {
		errStr := fmt.Sprintf("handler - could not validate relay/session due to error=%s", strings.Replace(CleanError(err.Error()), "\n", " ", -1))

		if app.GlobalMeshConfig.LogRelayRequest {
			// just if the setting is set to true, which by default is false, it will attach to the error the request_body of the relay, so this could help us provide context
			// on incoming error from portals.
			rb, ex := json.Marshal(r)
			if ex == nil {
				errStr = fmt.Sprintf("%s relay_request=%s", errStr, rb)
			}
		}

		code := "400"
		if castedError, kk := err.(sdk.Error); kk {
			code = fmt.Sprintf("%d", castedError.Code())
		}
		// help to track how many bad request mesh filter.
		servicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, &servicerNode.Address, false, BadRequest, code)

		LogRelay(r, errStr, LogLvlError)
		return
	}

	// store relay on cache; once we hit this point this relay will be processed so should be notified to servicer even
	// if process is shutdown
	storeRelay(r)

	blockChainCallStart := time.Now()
	res, err = processRelay(r)
	blockChainCallEnd := time.Since(blockChainCallStart)

	if err != nil && pocketTypes.ErrorWarrantsDispatch(err) {
		session, e1 := sessionStorage.GetSession(r)
		if e1 != nil {
			LogRelay(r, fmt.Sprintf("handler - error=%s getting session from session storage", e1.Error), LogLvlError)
		} else {
			dispatch = session.Dispatch
		}
	}

	// add to task group pool
	if servicerNode.Node.Worker.Stopped() {
		// this should not happen, but just in case avoid a panic here.
		LogRelay(r, fmt.Sprintf("handler - queue worker of node=%s is stoppend", servicerNode.Node.URL), LogLvlError)
		return
	}

	// 50k - 50
	servicerNode.Node.Worker.Submit(addRelayToQueue(r))

	// track the relay time (with chain)
	relayTimeDuration := time.Since(relayTimeStart)
	// queue metric of relay
	servicerNode.Node.MetricsWorker.AddServiceMetricRelayFor(r, &servicerNode.Address, relayTimeDuration, false)
	// handler time without call blockchain duration
	servicerNode.Node.MetricsWorker.AddServiceHandlerMetricRelayFor(
		r,
		&servicerNode.Address,
		relayTimeDuration-blockChainCallEnd,
		false,
	)
	return
}
