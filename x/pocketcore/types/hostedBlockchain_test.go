package types

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestHostedBlockchains_GetChainURL(t *testing.T) {
	url := "https://www.google.com"
	ethereum, err := NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	testHostedBlockchain := HostedBlockchain{
		ID:  ethereum,
		URL: url,
	}
	hb := HostedBlockchains{
		M: map[string]HostedBlockchain{testHostedBlockchain.ID: testHostedBlockchain},
		l: sync.Mutex{},
		o: sync.Once{},
	}
	u, err := hb.GetChainURL(ethereum)
	assert.Nil(t, err)
	assert.Equal(t, u, url)
}

func TestHostedBlockchains_ContainsFromString(t *testing.T) {
	url := "https://www.google.com"
	ethereum, err := NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	bitcoin, err := NonNativeChain{
		Ticker:  "btc",
		Netid:   "1",
		Version: "0.19.0.1",
		Client:  "",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	testHostedBlockchain := HostedBlockchain{
		ID:  ethereum,
		URL: url,
	}
	hb := HostedBlockchains{
		M: map[string]HostedBlockchain{testHostedBlockchain.ID: testHostedBlockchain},
		l: sync.Mutex{},
		o: sync.Once{},
	}
	assert.True(t, hb.Contains(ethereum))
	assert.False(t, hb.Contains(bitcoin))
}

func TestHostedBlockchains_Validate(t *testing.T) {
	url := "https://www.google.com"
	ethereum, err := NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	testHostedBlockchain := HostedBlockchain{
		ID:  ethereum,
		URL: url,
	}
	HCNoURL := HostedBlockchain{
		ID:  ethereum,
		URL: "",
	}
	HCNoHash := HostedBlockchain{
		ID:  "",
		URL: url,
	}
	HCInvalidHash := HostedBlockchain{
		ID:  hex.EncodeToString([]byte("bad")),
		URL: url,
	}
	tests := []struct {
		name     string
		hc       HostedBlockchains
		hasError bool
	}{
		{
			name:     "Invalid HostedBlockchain, no URL",
			hc:       HostedBlockchains{M: map[string]HostedBlockchain{HCNoURL.URL: HCNoURL}},
			hasError: true,
		},
		{
			name:     "Invalid HostedBlockchain, no URL",
			hc:       HostedBlockchains{M: map[string]HostedBlockchain{HCNoHash.URL: HCNoHash}},
			hasError: true,
		},
		{
			name:     "Invalid HostedBlockchain, invalid ID",
			hc:       HostedBlockchains{M: map[string]HostedBlockchain{HCInvalidHash.URL: HCInvalidHash}},
			hasError: true,
		},
		{
			name:     "Valid HostedBlockchain",
			hc:       HostedBlockchains{M: map[string]HostedBlockchain{testHostedBlockchain.ID: testHostedBlockchain}},
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.hc.Validate() != nil, tt.hasError)
		})
	}
}
