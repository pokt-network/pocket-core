package db

import (
	"testing"
	
	"github.com/pokt-network/pocket-core/db"
	"github.com/pokt-network/pocket-core/node"
)

func DummyNode() node.Node {
	chains := []node.Blockchain{{Name: "ethereum", NetID: "1", Version: "1"}}
	n := node.Node{
		GID:         "test",
		IP:          "123",
		RelayPort:   "0",
		ClientPort:  "0",
		ClientID:    "0",
		CliVersion:  "0",
		Blockchains: chains,
	}
	return n
}

func TestPut(t *testing.T) {
	d := db.NewDB()
	_, err := d.Add(DummyNode())
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestRemove(t *testing.T) {
	d := db.NewDB()
	_, err := d.Remove(DummyNode())
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestGetAll(t *testing.T) {
	d := db.NewDB()
	_, err := d.GetAll()
	if err != nil {
		t.Fatalf(err.Error())
	}
}
