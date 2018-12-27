package config

import (
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/const"
	"testing"
)

func TestDataDir(t *testing.T) {
	config.InitializeConfiguration()
	config.PrintConfiguration()
	datadir := config.GetConfigInstance().Datadir
	if datadir == _const.DATADIR {
		t.Log(datadir)
	} else {
		t.Errorf("Data Directory: " + datadir + " is the incorrect value. \n Expected: " + _const.DATADIR)
	}
}
