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
	Validated       bool
	RetryTimes      int
	IsValid         bool
	Error           *SdkErrorResponse
	Session         *Session
}

func (ns *NodeSession) CountRelay() (*NodeSession, bool) {
	if !ns.Validated {
		// if this session is not validated yet, will keep been optimistic
		return ns, true
	}

	ns.RemainingRelays -= 1

	if ns.RemainingRelays > 0 {
		return ns, true // still can send relays
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
	ns.Error = NewSdkErrorFromPocketSdkError(pocketTypes.NewOverServiceError(ModuleName))

	return ns, false
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

func (ns *NodeSession) ReScheduleValidationTask(session *Session, servicerPubKey string) {
	if ns.RetryTimes > 20 {
		// todo: what other thing we can do here?
		ns.IsValid = false
		ns.Error = NewSdkErrorFromPocketSdkError(
			sdk.ErrInternal(
				fmt.Sprintf(
					"unable to verify session=%s app=%s chain=%s blockHeight=%s servicer=%s",
					session.Hash,
					session.AppPublicKey,
					session.Chain,
					session.BlockHeight,
					servicerPubKey,
				),
			),
		)
		return
	}

	sessionStorage.ValidationWorker.Submit(session.ValidateSessionTask(servicerPubKey))
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
		nodeSession, e := s.GetNodeSessionByPubKey(servicerPubKey)
		if e != nil {
			// in theory this should never be hit
			logger.Error(e.Error())
			return
		}

		servicerNode := getServicerFromPubKey(nodeSession.PubKey)

		if s.BlockHeight > servicerNode.Node.Status.Height {
			// reschedule this session check because the node is not still on the expected block
			nodeSession.ReScheduleValidationTask(s, servicerPubKey)
			return
		}

		result, statusCode, e := s.GetDispatch(nodeSession)

		if e != nil {
			logger.Error(e.Error())
			// -5 = read body issue
			// StatusOK = 200
			// StatusUnauthorized = 401 - maybe after few retries node runner fix the issue and this will move? should we retry this?
			if statusCode == -5 || statusCode == http.StatusOK || statusCode == http.StatusUnauthorized {
				// this will re queue this.
				nodeSession.ReScheduleValidationTask(s, servicerPubKey)
			}
			return
		}

		isSuccess := statusCode == 200
		nodeSession.Validated = true // no mater result, this was checked among the fullNode

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
			nodeSession.ReScheduleValidationTask(s, ServicerSessionEndpoint)
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
		statusCode = -2
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
		statusCode = -3
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
		statusCode = -4
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
		statusCode = -5 // override this to allow caller know when the error was
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
		Validated:       false,
		Error:           nil,
		Session:         s,
	}
}

type SessionStorage struct {
	Sessions         *xsync.MapOf[string, *Session]
	ValidationWorker *pond.WorkerPool
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
	NewWorkerPoolMetrics(name, sessionStorage.ValidationWorker)
}

func (ss *SessionStorage) GetSession(relay *pocketTypes.Relay) (*Session, *SdkErrorResponse) {
	servicerAddress, _ := GetAddressFromPubKeyAsString(relay.Proof.ServicerPubKey)

	sessionHeader := pocketTypes.SessionHeader{
		ApplicationPubKey:  relay.Proof.Token.ApplicationPublicKey,
		Chain:              relay.Proof.Blockchain,
		SessionBlockHeight: relay.Proof.SessionBlockHeight,
	}

	sessionHash := hex.EncodeToString(sessionHeader.Hash())

	// if the session is already here
	if s, sOk := ss.Sessions.Load(sessionHash); sOk {
		hasDispatch := s.Dispatch != nil
		servicerInSession := hasDispatch && s.Dispatch.Contains(servicerAddress)
		_, servicerInLocalSession := s.Nodes.Load(relay.Proof.ServicerPubKey)

		if hasDispatch && servicerInSession || !hasDispatch && servicerInSession {
			if !servicerInLocalSession {
				nodeSession := s.NewNodeFromRelay(relay)
				s.Nodes.Store(relay.Proof.ServicerPubKey, nodeSession)
				s.ValidateSessionTask(relay.Proof.ServicerPubKey)
			}
			return s, nil
		} else if hasDispatch && !servicerInSession {
			// return session because
			return s, NewSdkErrorFromPocketSdkError(pocketTypes.NewSelfNotFoundError(ModuleName))
		}

		return s, nil
	}

	servicerNode := getServicerFromPubKey(relay.Proof.ServicerPubKey)

	// check if the relay is at the behind of a session
	// check if the session block height is +1 than our node, if yes that mean our node is not still on the height,
	// so we will be optimistic about this session and trust on the incoming relay, this session anyway will be moved to a validation
	// worker well it will keep checking the node status height so once the node is on the same or greater height it will check
	// the validity on the received session.
	if ss.ShouldAssumeOptimisticSession(relay, servicerNode.Node) {
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
	nodeSession.Validated = true // no mater result, this was checked among the fullNode

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
		}
	} else if result.Error != nil {
		nodeSession.IsValid = !ShouldInvalidateSession(result.Error.Code)
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
		Validated:       false,
		Error:           nil,
		Session:         &session,
	})

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
				Validated:       false, // mean this was not checked with fullnode yet
				IsValid:         true,  // true until node say the opposite
				Error:           nil,
				Session:         session,
			})
		}
	} else {
		session = ss.NewSessionFromRelay(relay)
		sessionStorage.Sessions.Store(hash, session)
	}

	// add this node/app/session relation to validate
	ss.ValidationWorker.Submit(session.ValidateSessionTask(relay.Proof.ServicerPubKey))

	return session, nil
}

func (ss *SessionStorage) ShouldAssumeOptimisticSession(relay *pocketTypes.Relay, servicerNode *fullNode) bool {
	sessionBlockHeight := relay.Proof.SessionBlockHeight
	fullNodeHeight := servicerNode.Status.Height
	blocksPerSession := servicerNode.BlocksPerSession
	return (sessionBlockHeight >= fullNodeHeight &&
		(fullNodeHeight%blocksPerSession == 0 || fullNodeHeight%blocksPerSession == 1)) &&
		relay.Proof.SessionBlockHeight-servicerNode.GetLatestSessionBlockHeight() <= 1
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
