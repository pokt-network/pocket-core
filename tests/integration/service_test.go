package integration

import (
	"github.com/pokt-network/pocket-core/config"
	"testing"
)

const (
	register   = "register"
	unregister = "unregister"
	whitelist  = "whitelist"
)

func TestRegister(t *testing.T) {
	// if dispatch node skip
	if config.GlobalConfig().Dispatch {
		t.Skip()
	}
	resp, err := requestFromFile(register)
	if err != nil {
		t.Log(assumptions)
		t.Fatalf(err.Error())
	}
	t.Log(resp)
}

func TestUnRegister(t *testing.T) {
	// if dispatch node skip
	if config.GlobalConfig().Dispatch {
		t.Skip()
	}
	resp, err := requestFromFile(unregister)
	if err != nil {
		t.Log(assumptions)
		t.Fatalf(err.Error())
	}
	t.Log(resp)
}

func TestWhiteList(t *testing.T) {
	if config.GlobalConfig().Dispatch {
		t.Skip()
	}
	resp, err := requestFromFile(whitelist)
	if err != nil {
		t.Log(assumptions)
		t.Fatalf(err.Error())
	}
	t.Log(resp)
}
