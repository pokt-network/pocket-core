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
type NodeSession struct {
	PubKey string
	// todo: if we can figure out a way to check this, otherwise receive/notify relays until fullNode return evidence sealed.
	// ^^ that could work with problem, just few "free relays"
	RemainingRelays int64
	RelayMeta       *pocketTypes.RelayMeta
	Queried         bool
	RetryTimes      int
	IsValid         bool
	Error           *SdkErrorResponse
	Session         *Session
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

	address, _ := GetAddressFromPubKeyAsString(ns.PubKey)

	logger.Debug(
		fmt.Sprintf(
			"servicer=%s exhaust relays for app=%s chain=%s sessionHeight=%d",
			address,
			ns.Session.AppPublicKey,
			ns.Session.Chain,
			ns.Session.BlockHeight,
		),
	)

	ns.IsValid = false
	// this is what pocket core return on their code, should 1 time they maybe return over-service but then all the next
	// request will be sealed evidence.
	ns.Error = NewSdkErrorFromPocketSdkError(pocketTypes.NewSealedEvidenceError(ModuleName))

	return false
}

// Session - Contains general app session information
type Session struct {
	Hash         string
	AppPublicKey string
	Chain        string
	BlockHeight  int64
	Dispatch     *DispatchResponse
	Nodes        *xsync.MapOf[string, *NodeSession]
}

func (ns *NodeSession) ReScheduleValidationTask(session *Session) {
	if ns.RetryTimes > app.GlobalMeshConfig.SessionStorageValidateRetryMaxTimes {
		// todo: what other thing we can do here?
		ns.IsValid = false
		ns.Error = NewSdkErrorFromPocketSdkError(
			sdk.ErrInternal(
				fmt.Sprintf(
					"unable to verify session=%s app=%s chain=%s blockHeight=%d servicer=%s",
					session.Hash,
					session.AppPublicKey,
					session.Chain,
					session.BlockHeight,
					ns.PubKey,
				),
			),
		)
		return
	}

	address, _ := GetAddressFromPubKeyAsString(ns.PubKey)

	sessionStorage.Metrics.AddSessionStorageMetricQueueFor(session, address, true)
	sessionStorage.ValidationWorker.Submit(session.ValidateSessionTask(ns.PubKey))
}

func (s *Session) GetNodeSessionByPubKey(servicerPubKey string) (*NodeSession, error) {
	var nodeSession *NodeSession

	if v, ok := s.Nodes.Load(servicerPubKey); !ok {
		// in theory this should never be hit
		return nil, errors.New(fmt.Sprintf(
			"unable to locate servicer %s on session hash=%s app=%s chain=%s. Please report it to Geo-Mesh developers.",
			servicerPubKey,
			s.Hash,
			s.AppPublicKey,
			s.Chain,
		))
	} else {
		nodeSession = v
	}

	return nodeSession, nil
}

func (s *Session) ValidateSessionTask(servicerPubKey string) func() {
	return func() {
		logger.Debug(fmt.Sprintf(
			"running an optimistic session validation for app=%s chain=%s session_height=%d servicer=%s",
			s.AppPublicKey,
			s.Chain,
			s.BlockHeight,
			servicerPubKey,
		))

		nodeSession, e := s.GetNodeSessionByPubKey(servicerPubKey)
		if e != nil {
			// in theory this should never be hit
			logger.Error(fmt.Sprintf(
				"error=%s getting servicer=%s on session for app=%s chain=%s session_height=%d",
				e.Error(),
				servicerPubKey,
				s.AppPublicKey,
				s.Chain,
				s.BlockHeight,
			))
			return
		}

		if nodeSession == nil {
			logger.Error(
				fmt.Sprintf(
					"servicer=%s not found with ValidateSessionTask.GetNodeSessionByPubKey for session hash=%s app=%s chain=%s session_height=%d",
					servicerPubKey,
					s.Hash,
					s.AppPublicKey,
					s.Chain,
					s.BlockHeight,
				))
			return
		}

		servicerNode := getServicerFromPubKey(nodeSession.PubKey)

		if servicerNode == nil {
			logger.Error(
				fmt.Sprintf(
					"servicer=%s not found with ValidateSessionTask.getServicerFromPubKey for session hash=%s app=%s chain=%s session_height=%d",
					servicerPubKey,
					s.Hash,
					s.AppPublicKey,
					s.Chain,
					s.BlockHeight,
				))
			return
		}

		if s.BlockHeight > servicerNode.Node.Status.Height {
			// reschedule this session check because the node is not still on the expected block
			logger.Debug(
				fmt.Sprintf(
					"servicer=%s session_height=%d greater than servicer node last know height=%d so this validation will be requeue.",
					servicerNode.Address.String(),
					s.BlockHeight,
					servicerNode.Node.Status.Height,
				))
			nodeSession.ReScheduleValidationTask(s)
			return
		}

		result, statusCode, e := s.GetDispatch(nodeSession)

		if e != nil {
			logger.Error(fmt.Sprintf("error getting session disatch err=%s", e.Error()))
			// -5 = read body issue
			// StatusOK = 200
			// StatusUnauthorized = 401 - maybe after few retries node runner fix the issue and this will move? should we retry this?
			if statusCode == ReadAllBodyError || statusCode == http.StatusOK || statusCode == http.StatusUnauthorized {
				// this will re queue this.
				nodeSession.ReScheduleValidationTask(s)
			}
			return
		}

		isSuccess := statusCode == 200
		nodeSession.Queried = true // no mater result, this was checked among the fullNode
		logger.Debug(fmt.Sprintf("session=%s servicer=%s queried done", s.Hash, servicerNode.Address.String()))

		if isSuccess {
			// dispatch response about session - across nodes
			s.Dispatch = result.Dispatch
			// node-session specific
			remainingRelays, _ := result.RemainingRelays.Int64()
			nodeSession.RemainingRelays = remainingRelays
			if result.Error != nil {
				nodeSession.IsValid = !ShouldInvalidateSession(result.Error.Code)
				if !nodeSession.IsValid {
					nodeSession.Error = result.Error
				}
			} else {
				nodeSession.IsValid = result.Success && remainingRelays > 0
			}
		} else if result.Error != nil {
			nodeSession.IsValid = !ShouldInvalidateSession(result.Error.Code)
		} else {
			nodeSession.ReScheduleValidationTask(s)
		}
	}
}

func (s *Session) GetDispatch(nodeSession *NodeSession) (result *RPCSessionResult, statusCode int, e error) {
	servicerAddress, _ := GetAddressFromPubKeyAsString(nodeSession.PubKey)

	servicerNode := getServicerFromPubKey(nodeSession.PubKey)

	payload := pocketTypes.MeshSession{
		SessionHeader: pocketTypes.SessionHeader{
			ApplicationPubKey:  s.AppPublicKey,
			Chain:              s.Chain,
			SessionBlockHeight: s.BlockHeight,
		},
		Meta:               *nodeSession.RelayMeta,
		ServicerPubKey:     nodeSession.PubKey,
		Blockchain:         s.Chain,
		SessionBlockHeight: s.BlockHeight,
	}

	logger.Debug(
		fmt.Sprintf(
			"session store - reading session for app=%s chain=%s blockHeight=%d servicer=%s",
			s.AppPublicKey,
			s.Chain,
			s.BlockHeight,
			servicerAddress,
		),
	)
	jsonData, e1 := json.Marshal(payload)
	if e1 != nil {
		statusCode = MarshallingError
		e = e1
		return
	}

	requestURL := fmt.Sprintf(
		"%s%s",
		servicerNode.Node.URL,
		ServicerSessionEndpoint,
	)
	req, e2 := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if e2 != nil {
		// should we retry? because here exists for sure a "config" setup issue with the requestURL
		statusCode = NewRequestError
		e = errors.New(fmt.Sprintf(
			"error creating check session request app=%s chain=%s blockHeight=%d servicer=%s err=%s",
			s.AppPublicKey, s.Chain, s.BlockHeight, servicerAddress, e1.Error(),
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
			s.AppPublicKey, s.Chain, s.BlockHeight, servicerAddress, e2.Error(),
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
			s.AppPublicKey, s.Chain, s.BlockHeight, servicerAddress, e3.Error(),
		))

		return
	}

	result = &RPCSessionResult{}
	e5 := json.Unmarshal(body, result)
	if e5 != nil {
		e = errors.New(fmt.Sprintf(
			"error unmarshalling check session response to RPCSessionResult app=%s chain=%s blockHeight=%d servicer=%s err=%s",
			s.AppPublicKey, s.Chain, s.BlockHeight, servicerAddress, e4.Error(),
		))
		return
	}

	return
}

func (s *Session) InvalidateNodeSession(servicerPubKey string, e *SdkErrorResponse) *SdkErrorResponse {
	nodeSession, e1 := s.GetNodeSessionByPubKey(servicerPubKey)

	if e1 != nil {
		return NewSdkErrorFromPocketSdkError(pocketTypes.NewInvalidSessionKeyError(ModuleName, e1))
	}

	nodeSession.IsValid = false
	nodeSession.Error = e

	return nil
}

func (s *Session) NewNodeFromRelay(relay *pocketTypes.Relay) *NodeSession {
	return &NodeSession{
		PubKey:          relay.Proof.ServicerPubKey,
		RemainingRelays: -1, // means that is unlimited until check it
		RelayMeta:       &relay.Meta,
		IsValid:         true, // true until node say the opposite
		Queried:         false,
		Error:           nil,
		Session:         s,
	}
}

type SessionStorage struct {
	Sessions         *xsync.MapOf[string, *Session]
	ValidationWorker *pond.WorkerPool
	Metrics          *Metrics
}

var (
	sessionStorage SessionStorage
)

func InitializeSessionStorage() {
	name := "session-storage"
	sessionStorage = SessionStorage{
		Sessions: xsync.NewMapOf[*Session](),
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

func (ss *SessionStorage) GetSessionHashFromRelay(relay *pocketTypes.Relay) string {
	sessionHeader := pocketTypes.SessionHeader{
		ApplicationPubKey:  relay.Proof.Token.ApplicationPublicKey,
		Chain:              relay.Proof.Blockchain,
		SessionBlockHeight: relay.Proof.SessionBlockHeight,
	}
	return hex.EncodeToString(sessionHeader.Hash())
}

func (ss *SessionStorage) GetSession(relay *pocketTypes.Relay) (*Session, *SdkErrorResponse) {
	servicerAddress, _ := GetAddressFromPubKeyAsString(relay.Proof.ServicerPubKey)

	sessionHash := ss.GetSessionHashFromRelay(relay)

	// if the session is already here
	if s, sOk := ss.Sessions.Load(sessionHash); sOk {
		hasDispatch := s.Dispatch != nil
		servicerInSession := hasDispatch && s.Dispatch.Contains(servicerAddress)
		_, servicerInLocalSession := s.Nodes.Load(relay.Proof.ServicerPubKey)

		if (hasDispatch && servicerInSession) || servicerInLocalSession {
			if !servicerInLocalSession {
				nodeSession := s.NewNodeFromRelay(relay)
				s.Nodes.Store(relay.Proof.ServicerPubKey, nodeSession)
				ss.ValidationWorker.Submit(s.ValidateSessionTask(relay.Proof.ServicerPubKey))
			}
			return s, nil
		} else if hasDispatch && !servicerInSession {
			// return session because
			return s, NewSdkErrorFromPocketSdkError(pocketTypes.NewSelfNotFoundError(ModuleName))
		}

		return s, nil
	}

	servicerNode := getServicerFromPubKey(relay.Proof.ServicerPubKey)

	// 1. Optimistic:
	// check if the relay is at the behind of a session
	// check if the session block height is +1 than our node, if yes that mean our node is not still on the height,
	// so we will be optimistic about this session and trust on the incoming relay, this session anyway will be moved to a validation
	// worker well it will keep checking the node status height so once the node is on the same or greater height it will check
	// the validity on the received session.
	// 2. Current Session:
	// allow it trusting that is a good session, so we lower the latency on the behind of a session.
	if ss.ShouldAssumeOptimisticSession(relay, servicerNode.Node) || servicerNode.Node.GetLatestSessionBlockHeight() == relay.Proof.SessionBlockHeight {
		// be optimistic about this session
		s, e := ss.AddSessionToValidate(relay)
		if e != nil {
			return nil, NewSdkErrorFromPocketSdkError(sdk.ErrInternal(e.Error()))
		}
		return s, nil
	}

	session := ss.NewSessionFromRelay(relay)

	nodeSession, e := session.GetNodeSessionByPubKey(relay.Proof.ServicerPubKey)
	if e != nil {
		// in theory this should never be hit
		return nil, NewSdkErrorFromPocketSdkError(sdk.ErrInternal(e.Error()))
	}

	result, statusCode, e := session.GetDispatch(nodeSession)
	if e != nil {
		return nil, NewSdkErrorFromPocketSdkError(sdk.ErrInternal(e.Error()))
	}

	isSuccess := statusCode == 200
	nodeSession.Queried = true // no mater result, this was checked among the fullNode

	if isSuccess {
		// dispatch response about session - across nodes
		session.Dispatch = result.Dispatch
		// node-session specific
		remainingRelays, _ := result.RemainingRelays.Int64()
		nodeSession.RemainingRelays = remainingRelays
		if result.Error == nil {
			nodeSession.IsValid = result.Success && remainingRelays > 0
		} else {
			nodeSession.IsValid = !ShouldInvalidateSession(result.Error.Code)
			nodeSession.Error = result.Error
		}
	} else if result.Error != nil {
		nodeSession.IsValid = !ShouldInvalidateSession(result.Error.Code)
		nodeSession.Error = result.Error
	}

	// return session as it is read, could be or not a valid one.
	return session, nil
}

func (ss *SessionStorage) GetNodeSession(relay *pocketTypes.Relay) (*NodeSession, *SdkErrorResponse) {
	session, e1 := ss.GetSession(relay)

	if e1 != nil {
		return nil, e1
	}

	nodeSession, e2 := session.GetNodeSessionByPubKey(relay.Proof.ServicerPubKey)

	if e2 != nil {
		return nil, NewSdkErrorFromPocketSdkError(pocketTypes.NewInvalidSessionKeyError(ModuleName, e2))
	}

	return nodeSession, nil
}

func (ss *SessionStorage) NewSessionFromRelay(relay *pocketTypes.Relay) *Session {
	sessionHeader := pocketTypes.SessionHeader{
		ApplicationPubKey:  relay.Proof.Token.ApplicationPublicKey,
		Chain:              relay.Proof.Blockchain,
		SessionBlockHeight: relay.Proof.SessionBlockHeight,
	}
	hash := hex.EncodeToString(sessionHeader.Hash())

	session := Session{
		Hash:         hash,
		AppPublicKey: sessionHeader.ApplicationPubKey,
		Chain:        sessionHeader.Chain,
		BlockHeight:  relay.Proof.SessionBlockHeight,
		Nodes:        xsync.NewMapOf[*NodeSession](),
		Dispatch:     nil,
	}
	session.Nodes.Store(relay.Proof.ServicerPubKey, &NodeSession{
		PubKey:          relay.Proof.ServicerPubKey,
		RemainingRelays: -1, // means that is unlimited until check it
		RelayMeta:       &relay.Meta,
		IsValid:         true, // true until node say the opposite
		Queried:         false,
		Error:           nil,
		Session:         &session,
	})
	ss.Sessions.Store(hash, &session)

	return &session
}

func (ss *SessionStorage) AddSessionToValidate(relay *pocketTypes.Relay) (*Session, error) {
	sessionHeader := pocketTypes.SessionHeader{
		ApplicationPubKey:  relay.Proof.Token.ApplicationPublicKey,
		Chain:              relay.Proof.Blockchain,
		SessionBlockHeight: relay.Proof.SessionBlockHeight,
	}

	var session *Session
	hash := hex.EncodeToString(sessionHeader.Hash())

	logger.Debug(fmt.Sprintf("adding session=%s to validate", hash))
	if v, ok := ss.Sessions.Load(hash); ok {
		session = v
		// add node to session if not there, but we already have the session
		// this could happen because multiple nodes on the mesh are working for the same session,
		// but the sessions are initialized by a relays, trusting on the incoming session retrieved.
		// so each time a servicer require this session, is not in the session nodes list, it is added for then
		// be sent to validate
		if _, nodeOk := session.Nodes.Load(relay.Proof.ServicerPubKey); !nodeOk {
			session.Nodes.Store(relay.Proof.ServicerPubKey, &NodeSession{
				PubKey:          relay.Proof.ServicerPubKey,
				RemainingRelays: -1, // means that is unlimited until check it
				RelayMeta:       &relay.Meta,
				Queried:         false, // mean this was not checked with fullnode yet
				IsValid:         true,  // true until node say the opposite
				Error:           nil,
				Session:         session,
			})
		}
	} else {
		session = ss.NewSessionFromRelay(relay)
		sessionStorage.Sessions.Store(hash, session)
	}

	servicerAddress, _ := GetAddressFromPubKeyAsString(relay.Proof.ServicerPubKey)

	// add this node/app/session relation to validate
	ss.ValidationWorker.Submit(session.ValidateSessionTask(relay.Proof.ServicerPubKey))
	sessionStorage.Metrics.AddSessionStorageMetricQueueFor(session, servicerAddress, false)
	return session, nil
}

// ShouldAssumeOptimisticSession - This will evaluate if the node is 1 block behind (at the end of the latest session he knows)
// and the received relay is on the immediate next session.
func (ss *SessionStorage) ShouldAssumeOptimisticSession(relay *pocketTypes.Relay, servicerNode *fullNode) bool {
	dispatcherSessionBlockHeight := relay.Proof.SessionBlockHeight
	fullNodeHeight := servicerNode.Status.Height
	blocksPerSession := servicerNode.BlocksPerSession
	servicerNodeSessionBlockHeight := servicerNode.GetLatestSessionBlockHeight()

	// check if the relay is on a "future" session in relation to the latest known block of the node
	isDispatcherAhead := dispatcherSessionBlockHeight >= fullNodeHeight
	// check if not is at the end of his session
	isFullNodeAtEndOfSession := (fullNodeHeight % blocksPerSession) == 0
	// check if the difference between fullNode and the relay session height is really close to avoid someone could abuse
	// of this optimistic approach.
	isDispatcherWithinTolerance := (dispatcherSessionBlockHeight - servicerNodeSessionBlockHeight) <= blocksPerSession

	return isDispatcherAhead && isFullNodeAtEndOfSession && isDispatcherWithinTolerance
}

// cleanOldSessions - clean up sessions that are longer than 50 blocks (just to be sure they are not needed)
func cleanOldSessions(c *cron.Cron) {
	_, err := c.AddFunc(fmt.Sprintf("@every %ds", app.GlobalMeshConfig.SessionCacheCleanUpInterval), func() {
		sessionsToDelete := make([]string, 0)
		servicerMap.Range(func(_ string, servicerNode *servicer) bool {
			sessionStorage.Sessions.Range(func(hash string, session *Session) bool {
				if session.BlockHeight < (servicerNode.Node.Status.Height - 6) {
					sessionsToDelete = append(sessionsToDelete, hash)
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
