package session

import "github.com/pokt-network/pocket-core/types"

type Session struct {
	SessionKey     SessionKey        `json:"sessionkey"`
	Application    SessionAppPubKey  `json:"appPubKey"`
	NonNativeChain SessionBlockchain `json:"nonnativechain"`
	BlockID        SessionBlockID    `json:"latestBlock"`
	Nodes          SessionNodes      `json:"sessionNodes"`
}

// Create a new session from seed data
func NewSession(app SessionAppPubKey, nonNativeChain SessionBlockchain, blockID SessionBlockID, allActiveNodes types.Nodes) (*Session, error) { // todo possibly convert block id to block hash
	// first generate session key
	sessionKey, err := NewSessionKey(app, nonNativeChain, blockID)
	if err != nil {
		return nil, err
	}
	// then generate the service nodes for that session
	sessionNodes, err := NewSessionNodes(nonNativeChain, sessionKey, allActiveNodes)
	if err != nil {
		return nil, err
	}
	// then populate the structure and return
	return &Session{SessionKey: sessionKey, Application: app, NonNativeChain: nonNativeChain, BlockID: blockID, Nodes: sessionNodes}, nil
}
