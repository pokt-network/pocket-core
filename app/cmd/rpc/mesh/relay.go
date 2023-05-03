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
	logger.Debug(fmt.Sprintf("storing relay %s", relay.RequestHashString()))
	rb, err := json.Marshal(relay)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = relaysCacheDb.Put(hash, rb)
	if err != nil {
		logger.Error(fmt.Sprintf("error adding relay %s to cache. %s", relay.RequestHashString(), err.Error()))
	}

	return
}

// decodeCacheRelay - decode []byte relay from cache to pocketTypes.Relay
func decodeCacheRelay(body []byte) (relay *pocketTypes.Relay) {
	if err := json.Unmarshal(body, &relay); err != nil {
		logger.Error("error decoding cache relay")
		// todo: delete key from cache?
		return nil
	}
	return
}

// deleteCacheRelay - delete a key from relay cache
func deleteCacheRelay(relay *pocketTypes.Relay) {
	hash := relay.RequestHash()
	err := relaysCacheDb.Delete(hash)
	if err != nil {
		logger.Error("error deleting relay from cache %s", hex.EncodeToString(hash))
		return
	}

	return
}

// evaluateServicerError - this will change internalCache[hash].IsValid bool depending on the result of the evaluation
func evaluateServicerError(r *pocketTypes.Relay, err *SdkErrorResponse) (isSessionStillValid bool) {
	hash := getSessionHashFromRelay(r)

	isSessionStillValid = !ShouldInvalidateSession(err.Code) // we should not retry if is invalid

	if isSessionStillValid {
		return isSessionStillValid
	}

	servicerNode := getServicerFromPubKey(r.Proof.ServicerPubKey)

	if appSession, ok := servicerNode.LoadAppSession(hash); ok {
		appSession.IsValid = isSessionStillValid
		appSession.Error = err
		servicerNode.StoreAppSession(hash, appSession)
	} else {
		logger.Error(
			fmt.Sprintf(
				"missing session hash=%s from cache; it should be there but if u see this after a restart it's ok.",
				hex.EncodeToString(hash),
			),
		)
	}

	return
}

// notifyServicer - call servicer to ack about the processed relay.
func notifyServicer(r *pocketTypes.Relay) {
	relayTimeStart := time.Now()
	// discard this relay at the end of this function, to end this function the servicer will be retried N times
	defer deleteCacheRelay(r)

	result := RPCRelayResponse{}
	ctx := context.WithValue(context.Background(), "result", &result)
	jsonData, err := json.Marshal(r)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"error encoding relay %s for servicer: %s",
				r.RequestHashString(),
				err.Error(),
			),
		)
		return
	}

	servicerAddress, err := GetAddressFromPubKeyAsString(r.Proof.ServicerPubKey)

	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"unable to decode service public key from relay %s to address",
				r.RequestHashString(),
			),
		)
		return
	}

	logger.Debug(
		fmt.Sprintf(
			"delivery relay %s notification to servicer %s",
			r.RequestHashString(),
			servicerAddress,
		),
	)

	servicerNode, ok := servicerMap.Load(servicerAddress)

	if !ok {
		logger.Error(
			fmt.Sprintf(
				"unable to find servicer with address=%s to notify relay %s",
				servicerAddress,
				r.RequestHashString(),
			),
		)
		return
	}

	requestURL := fmt.Sprintf(
		"%s%s?chain=%s&app=%s",
		servicerNode.Node.URL,
		ServicerRelayEndpoint,
		r.Proof.Blockchain,
		r.Proof.Token.ApplicationPublicKey,
	)
	req, err := retryablehttp.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(jsonData))
	req.Header.Set(AuthorizationHeader, servicerAuthToken.Value)
	if err != nil {
		logger.Error(fmt.Sprintf("error formatting Servicer URL: %s", err.Error()))
		servicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, &servicerNode.Address, true, NotifyStatusType, "500")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(ServicerHeader, servicerAddress)
	if app.GlobalMeshConfig.UserAgent != "" {
		req.Header.Set("User-Agent", app.GlobalMeshConfig.UserAgent)
	}
	resp, err := relaysClient.Do(req)

	if err != nil {
		logger.Error(fmt.Sprintf("error dispatching relay to Servicer: %s", err.Error()))
		servicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, &servicerNode.Address, true, NotifyStatusType, "500")
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}(resp.Body)

	// read the body just to allow http 1.x be able to reuse the connection
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Couldn't parse response body.", "err", err)
		servicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, &servicerNode.Address, true, NotifyStatusType, fmt.Sprintf("%d", resp.StatusCode))
		return
	}

	isSuccess := resp.StatusCode == 200

	if result.Dispatch != nil && result.Dispatch.BlockHeight > servicerNode.Node.Status.Height {
		servicerNode.Node.Status.Height = result.Dispatch.BlockHeight
	}

	if !isSuccess {
		logger.Debug(
			fmt.Sprintf(
				"servicer %s reject relay %s\n: CODE=%d\nCODESPACE=%s\nMESSAGE=%s",
				servicerAddress, r.RequestHashString(),
				result.Error.Code, result.Error.Codespace, result.Error.Error,
			),
		)
		servicerNode.Node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, &servicerNode.Address, true, NotifyStatusType, fmt.Sprintf("%d", resp.StatusCode))
		evaluateServicerError(r, result.Error)
	} else {
		logger.Debug(fmt.Sprintf("servicer processed relay %s successfully", r.RequestHashString()))

		header := pocketTypes.SessionHeader{
			ApplicationPubKey:  r.Proof.Token.ApplicationPublicKey,
			Chain:              r.Proof.Blockchain,
			SessionBlockHeight: r.Proof.SessionBlockHeight,
		}

		hash := header.Hash()
		if appSession, ok := servicerNode.LoadAppSession(hash); ok {
			appSession.RemainingRelays -= 1
			logger.Debug(
				fmt.Sprintf(
					"servicer %s has %d remaining relays to process for app %s at blockchain %s",
					servicerNode.Address.String(),
					appSession.RemainingRelays,
					r.Proof.Token.ApplicationPublicKey,
					r.Proof.Blockchain,
				),
			)
			if appSession.RemainingRelays <= 0 {
				logger.Debug(
					fmt.Sprintf(
						"servicer %s exhaust relays for app %s at blockchain %s",
						servicerNode.Address.String(),
						r.Proof.Token.ApplicationPublicKey,
						r.Proof.Blockchain,
					),
				)
				appSession.IsValid = false
				appSession.Error = NewSdkErrorFromPocketSdkError(pocketTypes.NewOverServiceError(ModuleName))
			}
			servicerNode.StoreAppSession(hash, appSession)
		}

		// track the notify relay time
		relayDuration := time.Since(relayTimeStart)
		servicerNode.Node.MetricsWorker.AddServiceMetricRelayFor(r, &servicerNode.Address, relayDuration, true)
	}

	return
}

// execute - Attempts to do a request on the non-native blockchain specified
func execute(r *pocketTypes.Relay, hostedBlockchains *pocketTypes.HostedBlockchains, address *sdk.Address) (string, sdk.Error) {
	node := GetNodeFromAddress(address.String())

	if node == nil {
		node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, address, false, InternalStatusType, "500")
		return "", sdk.ErrInternal("Failed to find correct servicer PK")
	}

	// retrieve the hosted blockchain url requested
	chain, err := hostedBlockchains.GetChain(r.Proof.Blockchain)
	if err != nil {
		// metric track
		node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, address, false, InternalStatusType, "500")
		return "", err
	}
	_url := strings.Trim(chain.URL, `/`)
	if len(r.Payload.Path) > 0 {
		_url = _url + "/" + strings.Trim(r.Payload.Path, `/`)
	}

	// do basic http request on the relay
	res, er, code := ExecuteBlockchainHTTPRequest(
		r.Payload.Data, _url,
		app.GlobalMeshConfig.UserAgent, chain.BasicAuth,
		r.Payload.Method, r.Payload.Headers,
	)
	if er != nil {
		// metric track
		node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, address, false, InternalStatusType, fmt.Sprintf("%d", code))
		return res, pocketTypes.NewHTTPExecutionError(ModuleName, er)
	}

	if code >= 400 {
		node.MetricsWorker.AddServiceMetricErrorFor(r.Proof.Blockchain, address, false, ChainStatusType, fmt.Sprintf("%d", code))
	}

	return res, nil
}

// processRelay - call execute and create RelayResponse or Error in case. Also trigger relay metrics.
func processRelay(relay *pocketTypes.Relay) (*pocketTypes.RelayResponse, sdk.Error) {
	relayTimeStart := time.Now()
	logger.Debug(fmt.Sprintf("processing relay %s", relay.RequestHashString()))

	servicerAddress, e := GetAddressFromPubKeyAsString(relay.Proof.ServicerPubKey)

	if e != nil {
		return nil, sdk.ErrInternal("could not convert servicer hex to public key")
	}

	servicerNode, ok := servicerMap.Load(servicerAddress)
	if !ok {
		return nil, sdk.ErrInternal("failed to find correct servicer PK")
	}

	// attempt to execute
	respPayload, err := execute(relay, chains, &servicerNode.Address)
	if err != nil {
		logger.Error(fmt.Sprintf("could not send relay %s with error: %s", relay.RequestHashString(), err.Error()))
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
		logger.Error(
			fmt.Sprintf("could not sign response for address: %s with hash: %v, with error: %s",
				servicerAddress, resp.HashString(), er.Error()),
		)
		return nil, pocketTypes.NewKeybaseError(pocketTypes.ModuleName, er)
	}
	// attach the signature in hex to the response
	resp.Signature = hex.EncodeToString(sig)
	// track the relay time
	relayTime := time.Since(relayTimeStart)
	// queue metric of relay
	servicerNode.Node.MetricsWorker.AddServiceMetricRelayFor(relay, &servicerNode.Address, relayTime, false)
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
	// validate servicer public key
	servicerAddress, e := GetAddressFromPubKeyAsString(r.Proof.ServicerPubKey)
	if e != nil {
		return sdk.ErrInternal("could not convert servicer hex to public key")
	}
	// load servicer from servicer map, if not there maybe the servicer is pk is not loaded
	if _, ok := servicerMap.Load(servicerAddress); !ok {
		return pocketTypes.NewInvalidSessionError(ModuleName)
	}

	header := pocketTypes.SessionHeader{
		ApplicationPubKey:  r.Proof.Token.ApplicationPublicKey,
		Chain:              r.Proof.Blockchain,
		SessionBlockHeight: r.Proof.SessionBlockHeight,
	}

	hash := header.Hash()

	servicerNode := getServicerFromPubKey(r.Proof.ServicerPubKey)

	if appSession, ok := servicerNode.LoadAppSession(hash); !ok {
		result := &RPCSessionResult{}
		e2 := getAppSession(r, result)

		if e2 != nil {
			logger.Error("Failure getting session information", "err", e2)
			return NewPocketSdkErrorFromSdkError(e2)
		}

		if !result.Success {
			logger.Error("Failure getting session information", "result", result)
			if result.Error != nil {
				return NewPocketSdkErrorFromSdkError(result.Error)
			}
			return pocketTypes.NewInvalidBlockHeightError(ModuleName)
		}

		remainingRelays, _ := result.RemainingRelays.Int64()

		isValid := result.Success && remainingRelays > 0 && result.Error == nil

		servicerNode.StoreAppSession(header.Hash(), &AppSessionCache{
			PublicKey:       header.ApplicationPubKey,
			Chain:           header.Chain,
			Dispatch:        result.Dispatch,
			RemainingRelays: remainingRelays,
			IsValid:         isValid,
			Error:           result.Error,
		})
	} else {
		if !appSession.IsValid {
			if appSession.Error != nil {
				return NewPocketSdkErrorFromSdkError(appSession.Error)
			} else {
				return sdk.ErrInternal("invalid session")
			}
		}
	}

	// is needed we call the node and validate if there is not a validation already in place get done by the cron?
	return nil
}

// HandleRelay - evaluate node status, validate relay payload and call processRelay
func HandleRelay(r *pocketTypes.Relay) (res *pocketTypes.RelayResponse, dispatch *DispatchResponse, err error) {
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
		logger.Error(
			fmt.Sprintf(
				"could not validate relay %s for app: %s, for chainID %v on node %s, at session height: %v, with error: %s",
				r.RequestHashString(),
				r.Proof.Token.ApplicationPublicKey,
				r.Proof.Blockchain,
				servicerAddress,
				r.Meta.BlockHeight,
				err.Error(),
			),
		)

		return
	}

	// store relay on cache; once we hit this point this relay will be processed so should be notified to servicer even
	// if process is shutdown
	storeRelay(r)

	res, err = processRelay(r)

	if err != nil && pocketTypes.ErrorWarrantsDispatch(err) {
		header := pocketTypes.SessionHeader{
			ApplicationPubKey:  r.Proof.Token.ApplicationPublicKey,
			Chain:              r.Proof.Blockchain,
			SessionBlockHeight: r.Proof.SessionBlockHeight,
		}

		hash := header.Hash()

		if appSession, ok := servicerNode.LoadAppSession(hash); !ok {
			response := RPCSessionResult{}
			err1 := getAppSession(r, &response)
			if err1 != nil {
				logger.Error(
					fmt.Sprintf(
						"error getting app %s session; hash %s",
						r.Proof.Token.ApplicationPublicKey,
						hash,
					),
				)
			} else {
				dispatch = response.Dispatch
			}
		} else {
			dispatch = appSession.Dispatch
		}
	}

	// add to task group pool
	if servicerNode.Node.Worker.Stopped() {
		// this should not happen, but just in case avoid a panic here.
		logger.Error(fmt.Sprintf("Worker of node %s was already stopped", servicerNode.Node.URL))
		return
	}

	servicerNode.Node.Worker.Submit(func() {
		notifyServicer(r)
	})

	return
}
