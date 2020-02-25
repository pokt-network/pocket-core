package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNonNativeChain_HashString(t *testing.T) {
	chainMissingTicker := NonNativeChain{
		Ticker:  "",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "",
		Inter:   "",
	}
	chainMissingNetID := NonNativeChain{
		Ticker:  "eth",
		Netid:   "",
		Version: "v1.9.9",
		Client:  "",
		Inter:   "",
	}
	chainMissingVersion := NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "",
		Client:  "",
		Inter:   "",
	}
	chainValid := NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "",
		Inter:   "",
	}
	tests := []struct {
		name     string
		chain    NonNativeChain
		hasError bool
	}{
		{
			name:     "Chain doesn't have a ticker",
			chain:    chainMissingTicker,
			hasError: true,
		},
		{
			name:     "Chain doesn't have a netid",
			chain:    chainMissingNetID,
			hasError: true,
		},
		{
			name:     "Chain doesn't have a version",
			chain:    chainMissingVersion,
			hasError: true,
		},
		{
			name:     "Chain doesn't have a version",
			chain:    chainValid,
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.chain.Bytes()
			assert.Equal(t, err != nil, tt.hasError)
			_, err = tt.chain.Hash()
			assert.Equal(t, err != nil, tt.hasError)
			_, err = tt.chain.HashString()
			assert.Equal(t, err != nil, tt.hasError)
		})
	}
}
