package integration

import (
	"github.com/pokt-network/pocket-core/config"
	"testing"

	"github.com/pokt-network/pocket-core/db"
	"github.com/pokt-network/pocket-core/node"
)

func DummyNode() node.Node {
	chains := []node.Blockchain{{Name: "ETH", NetID: "4"}}
	n := node.Node{
		GID:         "test",
		IP:          "123",
		RelayPort:   "0",
		ClientID:    "0",
		CliVersion:  "0",
		Blockchains: chains,
	}
	return n
}

func TestDB(t *testing.T) {
	// if service node skip
	if !config.GlobalConfig().Dispatch{
		t.Skip()
	}
	n := DummyNode()
	_, err := db.DB().Add(n)
	if err != nil {
		t.Log(assumptions)
		t.Fatalf(err.Error())
	}
	_, err = db.DB().Remove(n)
	if err != nil {
		t.Log(assumptions)
		t.Fatalf(err.Error())
	}
}
