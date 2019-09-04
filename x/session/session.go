package session

type Session struct {
	SessionKey     SessionKey         `json:"sessionkey"`
	Application    SessionApplication `json:"application"`
	NonNativeChain SessionBlockchain  `json:"nonnativechain"`
	BlockID        SessionBlockID     `json:"latestBlock"`
	Nodes          SessionNodes       `json:"sessionNodes"`
}

// Create a new session from seed data
func NewSession(app SessionApplication, nonNativeChain SessionBlockchain, blockID SessionBlockID) (*Session, error) {
	// first generate session key
	sessionKey, err := NewSessionKey(app, nonNativeChain, blockID)
	if err != nil {
		return nil, err
	}
	// then generate the service nodes for that session
	sessionNodes, err := NewSessionNodes(nonNativeChain, sessionKey)
	if err != nil {
		return nil, err
	}
	// then populate the structure and return
	return &Session{SessionKey: sessionKey, Application: app, NonNativeChain: nonNativeChain, BlockID: blockID, Nodes: sessionNodes}, nil
}
