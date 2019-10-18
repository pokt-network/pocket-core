package service

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/session"
)

// verifies the session for the node (self)
func SessionSelfVerification(appPubKey ServiceAppPubKey, blockchain ServiceBlockchain, blockHash string, allActiveNodes types.Nodes) error {
	sessionApp := session.SessionAppPubKey(appPubKey)
	sessionBlockchain, err := hex.DecodeString(string(blockchain))
	if err != nil {
		return NewBlockHashHexDecodeError(err)
	}
	bh, err := hex.DecodeString(blockHash)
	if err != nil {
		return NewBlockHashHexDecodeError(err)
	}
	sessionBlockID := session.SessionBlockID{Hash: bh}
	sess, err := session.NewSession(sessionApp, sessionBlockchain, sessionBlockID, allActiveNodes)
	if err != nil {
		return NewServiceSessionGenerationError(err)
	}
	if !sess.Nodes.Contains(FAKESELFNODE) {
		return InvalidSessionError
	}
	return nil
}
