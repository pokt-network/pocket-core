package types

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetHostedChains(t *testing.T) {
	hc := GetHostedChains()
	assert.NotNil(t, hc)
	assert.NotNil(t, hc.M)
}

func TestHostedBlockchains_AddAndGetChain(t *testing.T) {
	hc := GetHostedChains()
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
		Hash: ethereum,
		URL:  "https://www.google.com",
	}
	testHostedBlockchain2 := HostedBlockchain{
		Hash: bitcoin,
		URL:  "https://www.yahoo.com",
	}
	hc.Add(testHostedBlockchain)
	hc.Add(testHostedBlockchain2)
	// add duplicate
	hc.Add(testHostedBlockchain)
	hc2 := GetHostedChains()
	res, err := hc2.GetChain(testHostedBlockchain.Hash)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal(t, res, testHostedBlockchain)
	res2, err := hc2.GetChain(testHostedBlockchain2.Hash)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal(t, res2, testHostedBlockchain2)
}

func TestHostedBlockchains_Delete(t *testing.T) {
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
		Hash: ethereum,
		URL:  "https://www.google.com",
	}
	GetHostedChains().Add(testHostedBlockchain)
	GetHostedChains().Delete(testHostedBlockchain)
	_, err = GetHostedChains().GetChain(testHostedBlockchain.Hash)
	assert.Equal(t, err, NewErrorChainNotHostedError(ModuleName))
}

func TestHostedBlockchains_LenAndClear(t *testing.T) {
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
		Hash: ethereum,
		URL:  "https://www.google.com",
	}
	testHostedBlockchain2 := HostedBlockchain{
		Hash: bitcoin,
		URL:  "https://www.yahoo.com",
	}
	GetHostedChains().Add(testHostedBlockchain)
	GetHostedChains().Add(testHostedBlockchain2)
	assert.Equal(t, GetHostedChains().Len(), 2)
	GetHostedChains().Clear()
	assert.Len(t, GetHostedChains().M, 0)
	assert.Equal(t, GetHostedChains().Len(), 0)
}

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
		Hash: ethereum,
		URL:  url,
	}
	GetHostedChains().Add(testHostedBlockchain)
	u, err := GetHostedChains().GetChainURL(ethereum)
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
		Hash: ethereum,
		URL:  url,
	}
	GetHostedChains().Add(testHostedBlockchain)
	assert.True(t, GetHostedChains().ContainsFromString(ethereum))
	assert.False(t, GetHostedChains().ContainsFromString(bitcoin))
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
		Hash: ethereum,
		URL:  url,
	}
	HCNoURL := HostedBlockchain{
		Hash: ethereum,
		URL:  "",
	}
	HCNoHash := HostedBlockchain{
		Hash: "",
		URL:  url,
	}
	HCInvalidHash := HostedBlockchain{
		Hash: hex.EncodeToString([]byte("bad")),
		URL:  url,
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
			name:     "Invalid HostedBlockchain, invalid Hash",
			hc:       HostedBlockchains{M: map[string]HostedBlockchain{HCInvalidHash.URL: HCInvalidHash}},
			hasError: true,
		},
		{
			name:     "Valid HostedBlockchain",
			hc:       HostedBlockchains{M: map[string]HostedBlockchain{testHostedBlockchain.Hash: testHostedBlockchain}},
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.hc.Validate() != nil, tt.hasError)
		})
	}
}
