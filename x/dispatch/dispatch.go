package dispatch

import (
	"github.com/pokt-network/pocket-core/x/blockchain"
	"github.com/pokt-network/pocket-core/x/session"
)

func DispatchPeers() {
	// TODO returns tendermint peers cross checked with world state
}

// essentially a wrapper for session.NewSession()
func DispatchSession(application session.SessionAppPubKey, nonNativeChain session.SessionBlockchain) (*session.Session, error) {
	// runs session generation
	sess, err := session.NewSession(application, nonNativeChain, session.SessionBlockID(blockchain.GetLatestSessionBlock()))
	if err != nil {
		return sess, NewSessionGenerationError(err)
	}
	return sess, err
}
