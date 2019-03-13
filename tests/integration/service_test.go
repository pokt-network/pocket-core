package integration

import (
	"github.com/pokt-network/pocket-core/config"
	"testing"
)

const (
	register   = "register"
	unregister = "unregister"
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
