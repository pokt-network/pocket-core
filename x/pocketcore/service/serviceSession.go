package service

import (
	"encoding/hex"
	session2 "github.com/pokt-network/pocket-core/x/pocketcore/session"
)

// verifies the session for the node (self)
func SessionSelfVerification(appPubKey ServiceAppPubKey, blockchain ServiceBlockchain, blockHash string) error {
	sessionApp := session2.SessionAppPubKey(appPubKey)
	sessionBlockchain := session2.SessionBlockchain(blockchain)
	bh, err := hex.DecodeString(blockHash)
	if err != nil {
		return NewBlockHashHexDecodeError(err)
	}
	sessionBlockID := session2.SessionBlockID{Hash: bh}
	sess, err := session2.NewSession(sessionApp, sessionBlockchain, sessionBlockID)
	if err != nil {
		return NewServiceSessionGenerationError(err)
	}
	if !sess.Nodes.Contains(FAKESELFNODE) {
		return InvalidSessionError
	}
	return nil
}
