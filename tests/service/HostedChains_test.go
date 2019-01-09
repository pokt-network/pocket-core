package service

import (
	"github.com/pokt-network/pocket-core/node"
	"testing"
)

func TestHostedChains(t *testing.T) {
	json := []byte(
		"[" +
			"{\"blockchain\":" +
				"{\"name\": \"ethereum\",\"netid\": \"1\", \"version\": \"1.0\"}, " +
			"\"port\":\"8545\", \"medium\":\"rpc\"}," +
			"{\"blockchain\":" +
				"{\"name\": \"bitcoin\",\"netid\": \"1\", \"version\": \"1.0\"}, " +
			"\"port\":\"8333\", \"medium\":\"rpc\"}" +
		"]")
	node.UnmarshalChains(json)
	hc := node.GetHostedChains()
	if len(*hc)==0 {
		t.Fatalf("No hosted chains were found")
	}
	t.Log(hc)
}
