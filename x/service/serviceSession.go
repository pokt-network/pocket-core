package service

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/x/session"
)

// verifies the session for the node (self)
func SessionSelfVerification(appPubKey ServiceAppPubKey, blockchain ServiceBlockchain, blockHash string) error {
	sessionApp := session.SessionAppPubKey(appPubKey)
	sessionBlockchain := session.SessionBlockchain(blockchain)
	bh, err := hex.DecodeString(blockHash)
	if err != nil {
		return NewBlockHashHexDecodeError(err)
	}
	sessionBlockID := session.SessionBlockID{Hash: bh}
	sess, err := session.NewSession(sessionApp, sessionBlockchain, sessionBlockID)
	if err != nil {
		return NewServiceSessionGenerationError(err)
	}
	if !sess.Nodes.Contains(FAKESELFNODE) {
		return InvalidSessionError
	}
	return nil
}
