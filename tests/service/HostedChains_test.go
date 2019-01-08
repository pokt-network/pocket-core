package service

import (
	"github.com/pokt-network/pocket-core/service"
	"testing"
)

func TestHostedChains(t *testing.T) {
	json := []byte(
		"[{\"name\": \"ethereum\",\"netid\": \"1\", \"version\": \"1.0\", \"port\":\"8545\", \"medium\":\"rpc\"}," +
		"{\"name\": \"bitcoin\",\"netid\": \"1\", \"version\": \"1.0\", \"port\":\"8333\", \"medium\":\"rpc\"}]")
	service.UnmarshalChains(json)
	hc := service.GetHostedChains()
	if len(*hc)==0 {
		t.Fatalf("No hosted chains were found")
	}
	t.Log(hc)
}
