package mesh

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alitto/pond"
	"github.com/pokt-network/pocket-core/app"
	sdk "github.com/pokt-network/pocket-core/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/puzpuzpuz/xsync"
	"github.com/robfig/cron/v3"
	"io"
	"io/ioutil"
	log2 "log"
	"net/http"
	"time"
)

const (
	// Threshhold to determine when to expire an old session. (nodeBlockHeight - session
	sessionBlockHeightExpireThreshhold = 12
)
const (
	MarshallingError = iota
	NewRequestError
	ExecuteRequestError
	ReadAllBodyError
)

// DispatchSessionNode - app session node structure
type DispatchSessionNode struct {
	Address       string          `json:"address"`
	Chains        []string        `json:"chains"`
	Jailed        bool            `json:"jailed"`
	OutputAddress string          `json:"output_address"`
	PublicKey     string          `json:"public_key"`
	ServiceUrl    string          `json:"service_url"`
	Status        sdk.StakeStatus `json:"status"`
	Tokens        string          `json:"tokens"`
	UnstakingTime time.Time       `json:"unstaking_time"`
}

// DispatchSession - app session structure
type DispatchSession struct {
	Header pocketTypes.SessionHeader `json:"header"`
	Key    string                    `json:"key"`
	Nodes  []DispatchSessionNode     `json:"nodes"`
}

// DispatchResponse handle /v1/client/dispatch response due to was unable to inflate it using pocket core struct
// it was throwing an error about Nodes unmarshalling
type DispatchResponse struct {
	BlockHeight int64           `json:"block_height"`
	Session     DispatchSession `json:"session"`
}

// Contains - evaluate if the dispatch response contains passed address in their node list
func (sn DispatchResponse) Contains(addr string) bool {
	// if empty return
	if addr == "" {
		return false
	}
	// loop over the nodes
	for _, node := range sn.Session.Nodes {
		if node.Address == addr {
			return true
		}
	}

	return false
}

// ShouldKeep - evaluate if this dispatch response is one that we need to keep for the running mesh node.
func (sn DispatchResponse) ShouldKeep() bool {
	// loop over the nodes
	for _, node := range sn.Session.Nodes {
		if _, ok := servicerMap.Load(node.Address); ok {
			return true
		}
	}
	// if hit here, no one of in the map match the dispatch response nodes.
	return false
}

// GetSupportedNodes - return a list of the supported nodes of running mesh node from the DispatchResponse payload.
func (sn DispatchResponse) GetSupportedNodes() []string {
	nodes := make([]string, 0)
	// loop over the nodes
	for _, node := range sn.Session.Nodes {
		// There is reference to node address so that way we don't have to recreate address twice for pre-leanpokt
		if _, ok := servicerMap.Load(node.Address); ok {
			nodes = append(nodes, node.Address)
		}
	}
	// if hit here, no one of in the map match the dispatch response nodes.
	return nodes
}

// NodeSession - contains error/valid information for the node-session relation
// NOTE: NodeSession values are unsafely modified across multiple goroutines. It is a shared mutable data structure
// that can result inconsistent state across multiple routines since it lacks locks / atomicity.
// This should not be a problem as node sessions isValid is default true and only changed one way to be false
// and remaining relays is merely only a estimate, and query states can be finalized after a couple queries.
type NodeSession struct {
	Key string
	// session info
	Hash            string                 // session hash
	AppPubKey       string                 // application public key
	Chain           string                 // application-session-chain
	ServicerPubKey  string                 // servicer public key
	ServicerAddress string                 // servicer address
	BlockHeight     int64                  // session block height
	Dispatch        *DispatchResponse      // dispatch response from the fullNode
	Queue           bool                   // flag to know if the async validation of the session is already in queue
	Queried         bool                   // flag to know if the session was queried to know the validity of it
	RetryTimes      int                    // how many times the session validation was retried
	RelayMeta       *pocketTypes.RelayMeta // just a sample of any relay to get the dispatch
	// servicer related info
	ServicerNode    *servicer
	RemainingRelays int64             // how many relays the servicer can still service - todo: probably will remove and handle the overServiceError from the fullNode
	IsValid         bool              // if the session is or not valid
	Error           *SdkErrorResponse // in case session is not valid anymore, this will be the error to be returned.
}

func (ns *NodeSession) CountRelay() bool {
	if !ns.Queried {
		// if this session is not validated yet, will keep been optimistic
		return true
	}

	if ns.RemainingRelays > 0 {
		ns.RemainingRelays -= 1
		return true // still can send relays
	}

	ns.Log("exhaust relays for ", LogLvlDebug)

	ns.IsValid = false
	// this is what pocket core return on their code, should 1 time they maybe return over-service but then all the next
	// request will be sealed evidence.
	ns.Error = NewSdkErrorFromPocketSdkError(pocketTypes.NewSealedEvidenceError(ModuleName))

	return false
}

func (ns *NodeSession) ValidateSessionTask() {
	// add a little time between validation based on the amount of retries it already made.
	// todo: probably we need a better handling for this "requeue" because this sleep will hang the worker for the amount of time.
	// to avoid this a bigger workers for this pool is a good idea.
	sleepDuration := time.Duration(ns.RetryTimes) * time.Second
	ns.Log(fmt.Sprintf("sleeping %d seconds before validate", sleepDuration), LogLvlDebug)
	time.Sleep(sleepDuration)

	ns.Log("running session validate task for ", LogLvlDebug)
	if ns.Queried {
		ns.Log("already validated", LogLvlDebug)
		ns.Queue = false
		return
	}

	// if the node is not still in the session height does not make sense to proceed, so just requeue it.
	if ns.BlockHeight > ns.ServicerNode.Node.Status.Height {
		// reschedule this session check because the node is not still on the expected block
		ns.Log(fmt.Sprintf("servicer latest node session_height=%d lower than relay", ns.ServicerNode.Node.GetLatestSessionBlockHeight()), LogLvlDebug)
		sessionStorage.SubmitSessionToValidate(ns, true)
		return
	}

	result, statusCode, e := ns.GetDispatch()
	if e != nil {
		ns.Log(fmt.Sprintf("error getting session disatch error=%s", e.Error()), LogLvlError)
		// StatusOK = 200
		// StatusUnauthorized = 401 - maybe after few retries node runner fix the issue and this will move? should we retry this?
		if statusCode == ReadAllBodyError || statusCode == http.StatusOK || statusCode == http.StatusUnauthorized {
			// this will re queue this.
			ns.Log(fmt.Sprintf("session requeue after get error with status_code=%d", statusCode), LogLvlDebug)
			sessionStorage.SubmitSessionToValidate(ns, true)
		}
		return
	}

	isSuccess := statusCode == 200

	if !isSuccess && result.Error == nil {
		// not success but could be due to network issue, so it will "retry" sending it again to queue
		sessionStorage.SubmitSessionToValidate(ns, true)
		return
	}

	ns.Log("session queried done", LogLvlInfo)
	if isSuccess {
		// dispatch response about session - across nodes
		ns.Dispatch = result.Dispatch
		// node-session specific
		remainingRelays, _ := result.RemainingRelays.Int64()
		ns.RemainingRelays = remainingRelays
	} else if result.Error != nil {
		ns.IsValid = !ShouldInvalidateSession(result.Error.Code)
	}

	ns.Queue = false
	ns.Queried = true
}

func (ns *NodeSession) GetDispatch() (result *RPCSessionResult, statusCode int, e error) {
	payload := pocketTypes.MeshSession{
		SessionHeader: pocketTypes.SessionHeader{
			ApplicationPubKey:  ns.AppPubKey,
			Chain:              ns.Chain,
			SessionBlockHeight: ns.BlockHeight,
		},
		Meta:               *ns.RelayMeta,
		ServicerPubKey:     ns.ServicerPubKey,
		Blockchain:         ns.Chain,
		SessionBlockHeight: ns.BlockHeight,
	}

	ns.Log("reading session from store", LogLvlDebug)
	jsonData, e1 := json.Marshal(payload)
	if e1 != nil {
		statusCode = MarshallingError
		e = e1
		return
	}

	requestURL := fmt.Sprintf(
		"%s%s",
		ns.ServicerNode.Node.URL,
		ServicerSessionEndpoint,
	)
	req, e2 := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if e2 != nil {
		// should we retry? because here exists for sure a "config" setup issue with the requestURL
		statusCode = NewRequestError
		e = errors.New(fmt.Sprintf(
			"error creating check session request app=%s chain=%s blockHeight=%d servicer=%s err=%s",
			ns.AppPubKey, ns.Chain, ns.BlockHeight, ns.ServicerAddress, e1.Error(),
		))
		return
	}

	req.Header.Set(AuthorizationHeader, servicerAuthToken.Value)
	req.Header.Set("Content-Type", "application/json")
	if app.GlobalMeshConfig.UserAgent != "" {
		req.Header.Set("User-Agent", app.GlobalMeshConfig.UserAgent)
	}

	resp, e3 := servicerClient.Do(req)
	if e3 != nil {
		statusCode = ExecuteRequestError
		e = errors.New(fmt.Sprintf(
			"error calling check session request app=%s chain=%s blockHeight=%d servicer=%s err=%s",
			ns.AppPubKey, ns.Chain, ns.BlockHeight, ns.ServicerAddress, e2.Error(),
		))
		return
	}

	statusCode = resp.StatusCode

	defer func(Body io.ReadCloser) {
		err1 := Body.Close()
		if err1 != nil {
			return // add log here
		}
	}(resp.Body)

	// read the body just to allow http 1.x be able to reuse the connection
	body, e4 := ioutil.ReadAll(resp.Body)

	if e4 != nil {
		statusCode = ReadAllBodyError // override this to allow caller know when the error was
		e = errors.New(fmt.Sprintf(
			"error reading check session response body app=%s chain=%s blockHeight=%d servicer=%s err=%s",
			ns.AppPubKey, ns.Chain, ns.BlockHeight, ns.ServicerPubKey, e3.Error(),
		))

		return
	}

	result = &RPCSessionResult{}
	e5 := json.Unmarshal(body, result)
	if e5 != nil {
		e = errors.New(fmt.Sprintf(
			"error unmarshalling check session response to RPCSessionResult app=%s chain=%s blockHeight=%d servicer=%s err=%s",
			ns.AppPubKey, ns.Chain, ns.BlockHeight, ns.ServicerAddress, e4.Error(),
		))
		return
	}

	return
}

func (ns *NodeSession) ToString() string {
	str := fmt.Sprintf(
		"session_hash=%s session_height=%d app=%s chain=%s servicer=%s remaining_relays=%d queried=%v queue=%v is_valid=%v",
		ns.Hash, ns.BlockHeight,
		ns.AppPubKey, ns.Chain,
		ns.ServicerAddress,
		ns.RemainingRelays,
		ns.Queried, ns.Queue, ns.IsValid,
	)

	if !ns.IsValid {
		str = fmt.Sprintf(
			"%s error=%s code=%d codespace=%s",
			str, ns.Error.Error, ns.Error.Code, ns.Error.Codespace,
		)
	}

	return str
}

// Log - log the entity valuable properties with a message. Internally call Log method from logger.go file.
func (ns *NodeSession) Log(msg string, level string) {
	Log(fmt.Sprintf("msg=%s %s", msg, ns.ToString()), level)
}

type SessionStorage struct {
	Sessions         *xsync.MapOf[string, *NodeSession]
	ValidationWorker *pond.WorkerPool
	Metrics          *Metrics
}

var (
	sessionStorage SessionStorage
)

func InitializeSessionStorage() {
	name := "session-storage"
	sessionStorage = SessionStorage{
		Sessions: xsync.NewMapOf[*NodeSession](),
		ValidationWorker: NewWorkerPool(
			name,
			"lazy", // app.GlobalMeshConfig.MetricsWorkerStrategy,
			5,      // app.GlobalMeshConfig.MetricsMaxWorkers,
			10000,  // app.GlobalMeshConfig.MetricsMaxWorkersCapacity,
			30000,  // app.GlobalMeshConfig.MetricsWorkersIdleTimeout,
		),
	}

	// add metrics worker and session storage pool metrics to monitor it and understand how is working.
	sessionStorage.Metrics = NewWorkerPoolMetrics(name, sessionStorage.ValidationWorker)
	sessionStorage.Metrics.Start()
}

func (ss *SessionStorage) Stop() {
	sessionStorage.ValidationWorker.Stop()
	sessionStorage.Metrics.Stop()
}

func (ss *SessionStorage) NewNodeSessionFromRelay(relay *pocketTypes.Relay) (*NodeSession, *SdkErrorResponse) {
	sessionHash := ss.GetSessionHashFromRelay(relay)
	servicerNode := getServicerFromPubKey(relay.Proof.ServicerPubKey)
	servicerAddress, _ := GetAddressFromPubKeyAsString(relay.Proof.ServicerPubKey)

	if servicerNode == nil {
		// servicer was not found on the keys loaded.
		return nil, NewSdkErrorFromPocketSdkError(pocketTypes.NewSelfNotFoundError(ModuleName))
	}

	return &NodeSession{
		Key:             ss.GetSessionKeyByRelay(relay),
		Hash:            sessionHash,
		AppPubKey:       relay.Proof.Token.ApplicationPublicKey,
		Chain:           relay.Proof.Blockchain,
		ServicerPubKey:  relay.Proof.ServicerPubKey,
		ServicerAddress: servicerAddress,
		BlockHeight:     relay.Proof.SessionBlockHeight,
		Dispatch:        nil,
		Queue:           false,
		Queried:         false,
		RetryTimes:      0,
		RelayMeta:       &relay.Meta,
		ServicerNode:    servicerNode,
		RemainingRelays: -1,   // means that is unlimited until check it
		IsValid:         true, // true until fullNode negate this
		Error:           nil,
	}, nil
}

func (ss *SessionStorage) InvalidateNodeSession(relay *pocketTypes.Relay, e *SdkErrorResponse) *SdkErrorResponse {
	if e == nil {
		e1 := errors.New("SessionStorage.InvalidateNodeSession called without sdk error")
		logger.Error(e1.Error())
		return NewSdkErrorFromPocketSdkError(sdk.ErrInternal(e1.Error()))
	}

	if ns, ok := ss.Sessions.Load(ss.GetSessionKeyByRelay(relay)); !ok {
		e1 := errors.New(fmt.Sprintf(
			"session not found on mesh for app=%s chain=%s session_height=%d servicer=%s",
			ns.AppPubKey,
			ns.Chain,
			ns.BlockHeight,
			ns.ServicerAddress,
		))
		return NewSdkErrorFromPocketSdkError(sdk.ErrInternal(e1.Error()))
	} else {
		ns.Log("invalidating session", LogLvlInfo)
		// queue and queried set to avoid requeue and understand we are invalidating the session even with query it
		// like if the session is too old or too far in the future.
		ns.Queue = false
		ns.Queried = true
		ns.IsValid = false
		ns.Error = e
	}

	return nil
}

func (ss *SessionStorage) GetSessionHashFromRelay(relay *pocketTypes.Relay) string {
	sessionHeader := pocketTypes.SessionHeader{
		ApplicationPubKey:  relay.Proof.Token.ApplicationPublicKey,
		Chain:              relay.Proof.Blockchain,
		SessionBlockHeight: relay.Proof.SessionBlockHeight,
	}
	return hex.EncodeToString(sessionHeader.Hash())
}

func (ss *SessionStorage) SubmitSessionToValidate(ns *NodeSession, isRequeue bool) {
	if ns.RetryTimes > app.GlobalMeshConfig.SessionStorageValidateRetryMaxTimes {
		// todo: what other thing we can do here?
		ns.IsValid = false
		ns.Error = NewSdkErrorFromPocketSdkError(
			sdk.ErrInternal(
				fmt.Sprintf(
					"unable to verify after %d tries; session=%s app=%s chain=%s blockHeight=%d servicer=%s",
					ns.RetryTimes,
					ns.Hash,
					ns.AppPubKey,
					ns.Chain,
					ns.BlockHeight,
					ns.ServicerPubKey,
				),
			),
		)
		return
	}

	ns.Queue = true
	if isRequeue {
		ns.RetryTimes++
	}
	ss.ValidationWorker.Submit(func() {
		ns.ValidateSessionTask()
	})
	ss.Metrics.AddSessionStorageMetricQueueFor(ns, isRequeue)
}

func (ss *SessionStorage) GetSession(relay *pocketTypes.Relay) (*NodeSession, *SdkErrorResponse) {
	nNs, e := ss.NewNodeSessionFromRelay(relay) // new node session
	if e != nil {
		return nil, e
	}

	var ns *NodeSession

	// load or store - avoid race conditions loading and then adding
	if value, loaded := ss.Sessions.LoadOrStore(nNs.Key, nNs); loaded {
		ns = value
	} else {
		// this was stored, so we need to sent it to validate
		ns = nNs
	}

	// if the session was already queried, just return it, no mater if is valid or not.
	if ns.Queried {
		return ns, nil
	}

	// 1. Optimistic:
	// check if the session block height is +1 than our node, if yes that mean our node is not still on the height,
	// so we will be optimistic about this session and trust on the incoming relay, this session anyway will be moved to a validation
	// worker well it will keep checking the node status height so once the node is on the same or greater height it will check
	// the validity on the received session.
	// 2. Current Session:
	// allow it trusting that is a good session, so we lower the latency on the behind of a session.
	if ns.ServicerNode.Node.ShouldAssumeOptimisticSession(relay.Proof.SessionBlockHeight) || ns.ServicerNode.Node.CanHandleRelayWithinTolerance(relay.Proof.SessionBlockHeight) {
		// be optimistic about this session
		if !ns.Queue && !ns.Queried {
			// avoid re queue a session that is already set in queue
			ss.SubmitSessionToValidate(ns, false)
		}
		return ns, nil
	}

	// this session will be invalidated - and this is the same instance invalidate will use.
	sessionStorage.InvalidateNodeSession(relay, NewSdkErrorFromPocketSdkError(pocketTypes.NewSealedEvidenceError(ModuleName)))

	return ns, nil
}

func (ss *SessionStorage) GetSessionKeyByRelay(relay *pocketTypes.Relay) string {
	return fmt.Sprintf("%s-%s", ss.GetSessionHashFromRelay(relay), relay.Proof.ServicerPubKey)
}

// cleanOldSessions - clean up sessions that are longer than sessionBlockHeightExpireThreshhold blocks (just to be sure they are not needed)
func cleanOldSessions(c *cron.Cron) {
	_, err := c.AddFunc(fmt.Sprintf("@every %ds", app.GlobalMeshConfig.SessionCacheCleanUpInterval), func() {
		sessionsToDelete := make([]string, 0)
		servicerMap.Range(func(_ string, servicerNode *servicer) bool {
			sessionStorage.Sessions.Range(func(key string, ns *NodeSession) bool {
				if (servicerNode.Node.Status.Height - ns.BlockHeight) >= sessionBlockHeightExpireThreshhold {
					sessionsToDelete = append(sessionsToDelete, key)
				}
				return true
			})
			return true
		})
	})

	if err != nil {
		log2.Fatal(err)
	}
}
