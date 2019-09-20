package dispatch

import (
	"github.com/pokt-network/pocket-core/tests/fixtures"
	"github.com/pokt-network/pocket-core/x/session"
)

var (
	validApplication      = session.SessionAppPubKey(fixtures.GenerateApplication().PubKey)
	validNonNativeChain   = fixtures.GenerateNonNativeBlockchain()
	invalidApplication    = ""
	invalidNonNativeChain = session.SessionBlockchain{}
)
