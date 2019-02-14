package tests

import (
	"testing"

	"github.com/pokt-network/pocket-core/service"
)

func TestReport(t *testing.T) {
	if _, err := service.HandleReport(&service.Report{GID: "TestReport", Message: "This is a test report"}); err != nil {
		t.Fatalf(err.Error())
	}
}
