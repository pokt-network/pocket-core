package config

import (
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/const"
	"os"
	"testing"
)

func TestBuildConfig(t *testing.T) {
	config.BuildConfiguration()
	_, err := os.Stat(_const.DATADIR)
	if err != nil {
		t.Fatalf("Couldn't follow path")
	}
	if os.IsNotExist(err) {
		t.Fatalf("Datadir doesn't exist")
	}
}

func TestLogsDir(t *testing.T) {
	config.BuildConfiguration()
	_, err := os.Stat(_const.DATADIR + _const.FILESEPARATOR + "logs")
	if err != nil {
		t.Fatalf("Couldn't follow path")
	}
	if os.IsNotExist(err) {
		t.Fatalf("Datadir doesn't exist")
	}
}
