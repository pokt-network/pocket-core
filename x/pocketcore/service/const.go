package service

import (
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/tests/fixtures"
)

const (
	// The http method defaults to POST for Relays
	// because JSON RPC uses POST for all calls
	DEFAULTHTTPMETHOD = "POST"
	// HTTP Service Types
	HTTP ServiceType = iota + 1
)

var (
	FAKENODEPRIVKEY, _ = crypto.NewPrivateKey()       // todo replace with global key (self)
	FAKESELFNODE       = fixtures.GenerateAliveNode() // todo replace with global node (self)
)
