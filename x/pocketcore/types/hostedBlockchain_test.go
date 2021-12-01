package types

import (
	"encoding/hex"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostedBlockchains_GetChainURL(t *testing.T) {
	url := "https://www.google.com:443"
	ethereum := hex.EncodeToString([]byte{01})
	testHostedBlockchain := HostedBlockchain{
		ID:  ethereum,
		URL: url,
	}
	hb := HostedBlockchains{
		M: map[string]HostedBlockchain{testHostedBlockchain.ID: testHostedBlockchain},
		L: sync.Mutex{},
	}
	u, err := hb.GetChainURL(ethereum)
	assert.Nil(t, err)
	assert.Equal(t, u, url)
}

func TestHostedBlockchains_ContainsFromString(t *testing.T) {
	url := "https://www.google.com:443"
	ethereum := hex.EncodeToString([]byte{01})
	bitcoin := hex.EncodeToString([]byte{02})
	testHostedBlockchain := HostedBlockchain{
		ID:  ethereum,
		URL: url,
	}
	hb := HostedBlockchains{
		M: map[string]HostedBlockchain{testHostedBlockchain.ID: testHostedBlockchain},
		L: sync.Mutex{},
	}
	assert.True(t, hb.Contains(ethereum))
	assert.False(t, hb.Contains(bitcoin))
}

func TestHostedBlockchains_Validate(t *testing.T) {
	url := "https://www.google.com:443"
	ethereum := hex.EncodeToString([]byte{01})
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
		ID:  hex.EncodeToString([]byte("badlksajfljasdfklj")),
		URL: url,
	}
	tests := []struct {
		name     string
		hc       *HostedBlockchains
		hasError bool
	}{
		{
			name:     "Invalid HostedBlockchain, no URL",
			hc:       &HostedBlockchains{M: map[string]HostedBlockchain{HCNoURL.URL: HCNoURL}, L: sync.Mutex{}},
			hasError: true,
		},
		{
			name:     "Invalid HostedBlockchain, no URL",
			hc:       &HostedBlockchains{M: map[string]HostedBlockchain{HCNoHash.URL: HCNoHash}, L: sync.Mutex{}},
			hasError: true,
		},
		{
			name:     "Invalid HostedBlockchain, invalid ID",
			hc:       &HostedBlockchains{M: map[string]HostedBlockchain{HCInvalidHash.URL: HCInvalidHash}, L: sync.Mutex{}},
			hasError: true,
		},
		{
			name:     "Valid HostedBlockchain",
			hc:       &HostedBlockchains{M: map[string]HostedBlockchain{testHostedBlockchain.ID: testHostedBlockchain}, L: sync.Mutex{}},
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.hc.Validate() != nil, tt.hasError)
		})
	}
}
