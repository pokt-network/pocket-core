package unit

import (
	"testing"

	"github.com/pokt-network/pocket-core/logs"
)

func TestLogs(t *testing.T) {
	if err := logs.Log("Unit test for the log functionality", logs.InfoLevel, logs.JSONLogFormat); err != nil {
		t.Fatalf(err.Error())
	}
}
