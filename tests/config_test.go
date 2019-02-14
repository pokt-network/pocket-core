package tests

import (
	"os"
	"testing"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/const"
)

func TestGlobalConfig(t *testing.T) {
	c := config.GlobalConfig()
	if c != nil {
		config.Print()
		return
	}
	t.Fatalf("The global configuration object returned nil")
}

func TestBuildConfig(t *testing.T) {
	config.Build()
	_, err := os.Stat(_const.DATADIR)
	if err != nil {
		t.Fatalf("Couldn't follow path")
	}
	if os.IsNotExist(err) {
		t.Fatalf("Datadir doesn't exist")
	}
}

func TestDataDir(t *testing.T) {
	datadir := config.GlobalConfig().DD
	if datadir == _const.DATADIR {
		t.Log(datadir)
		return
	}
	t.Fatalf("Data Directory: " + datadir + " is the incorrect value. \n Expected: " + _const.DATADIR)
}

func TestLogsDir(t *testing.T) {
	config.Build()
	_, err := os.Stat(_const.DATADIR + _const.FILESEPARATOR + "logs")
	if err != nil {
		t.Fatalf("Couldn't follow path")
	}
}
