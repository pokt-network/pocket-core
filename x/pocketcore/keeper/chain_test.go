package keeper

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateChain(t *testing.T) {
	validNNChain := types.NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "geth",
		Inter:   "",
	}
	// no ticker
	invalidNoTick := validNNChain
	invalidNoTick.Ticker = ""
	// no net id
	invalidNoNetid := validNNChain
	invalidNoNetid.Netid = ""
	// no net id
	invalidNoVersion := validNNChain
	invalidNoVersion.Version = ""
	tests := []struct {
		name     string
		chain    types.NonNativeChain
		hasError bool
	}{
		{
			name:     "invalid generation: missing ticker",
			chain:    invalidNoTick,
			hasError: true,
		},
		{
			name:     "invalid generation: missing netid",
			chain:    invalidNoNetid,
			hasError: true,
		},
		{
			name:     "invalid generation: missing version",
			chain:    invalidNoVersion,
			hasError: true,
		},
		{
			name:     "valid generation",
			chain:    validNNChain,
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain, err := GenerateChain(tt.chain.Ticker, tt.chain.Netid, tt.chain.Version, tt.chain.Client, tt.chain.Inter)
			assert.Equal(t, err != nil, tt.hasError)
			if !tt.hasError {
				assert.Nil(t, types.HashVerification(chain))
			}
		})
	}
}
