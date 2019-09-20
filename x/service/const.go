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
	RPC ServiceType = iota + 1
	REST

	// TODO remove and get from world state
	FAKEAPPPUBKEY = "043bcf236d393dc09e9526d01d295e898ab771550dce52d4bd806d03e7b7fc3b8684a5f5ac7532a8fb2e6b4b58b9cc78d2dd528eebbd43dd26035ba3364c7f47d4"
)

var(
	FAKENODEPRIVKEY, _ = crypto.NewPrivateKey() // todo replace with global key (self)
	FAKESELFNODE = fixtures.GenerateAliveNode() // todo replace with global node (self)
)
