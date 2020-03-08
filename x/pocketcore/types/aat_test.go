package types

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAAT_VersionIsIncluded(t *testing.T) {
	appPrivKey := GetRandomPrivateKey()
	clientPrivKey := GetRandomPrivateKey()
	var AATNoVersion = AAT{
		Version:              "",
		ApplicationPublicKey: appPrivKey.PublicKey().RawString(),
		ClientPublicKey:      clientPrivKey.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	var AATWithVersion = AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPrivKey.PublicKey().RawString(),
		ClientPublicKey:      clientPrivKey.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	tests := []struct {
		name     string
		aat      AAT
		expected bool
	}{
		{
			name:     "AAT is missing the version",
			aat:      AATNoVersion,
			expected: false,
		},
		{
			name:     "AAT has the version",
			aat:      AATWithVersion,
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.aat.VersionIsIncluded(), tt.expected)
		})
	}
}

func TestAAT_VersionIsSupported(t *testing.T) {
	appPrivKey := GetRandomPrivateKey()
	clientPrivKey := GetRandomPrivateKey()
	var AATNotSupportedVersion = AAT{
		Version:              "0.0.11",
		ApplicationPublicKey: appPrivKey.PublicKey().RawString(),
		ClientPublicKey:      clientPrivKey.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	var AATSupported = AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPrivKey.PublicKey().RawString(),
		ClientPublicKey:      clientPrivKey.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	tests := []struct {
		name     string
		aat      AAT
		expected bool
	}{
		{
			name:     "AAT doesn't not have a supported version",
			aat:      AATNotSupportedVersion,
			expected: false,
		},
		{
			name:     "AAT has a supported version",
			aat:      AATSupported,
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.aat.VersionIsSupported(), tt.expected)
		})
	}
}

func TestAAT_ValidateVersion(t *testing.T) {
	appPrivKey := GetRandomPrivateKey()
	clientPrivKey := GetRandomPrivateKey()
	var AATVersionMissing = AAT{
		Version:              "",
		ApplicationPublicKey: appPrivKey.PublicKey().RawString(),
		ClientPublicKey:      clientPrivKey.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	var AATNotSupportedVersion = AAT{
		Version:              "0.0.11",
		ApplicationPublicKey: appPrivKey.PublicKey().RawString(),
		ClientPublicKey:      clientPrivKey.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	var AATSupported = AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPrivKey.PublicKey().RawString(),
		ClientPublicKey:      clientPrivKey.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	tests := []struct {
		name     string
		aat      AAT
		hasError bool
	}{
		{
			name:     "AAT is missing the version",
			aat:      AATVersionMissing,
			hasError: true,
		},
		{
			name:     "AAT doesn't not have a supported version",
			aat:      AATNotSupportedVersion,
			hasError: true,
		},
		{
			name:     "AAT has a supported version",
			aat:      AATSupported,
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.aat.ValidateVersion() != nil, tt.hasError)
		})
	}
}

func TestAAT_ValidateMessage(t *testing.T) {
	appPrivKey := GetRandomPrivateKey()
	clientPubKey := getRandomPubKey()
	var AATInvalidAppPubKey = AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPrivKey.PubKey().Address().String(),
		ClientPublicKey:      clientPubKey.RawString(),
		ApplicationSignature: "",
	}
	var AATInvalidClientPubKey = AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPrivKey.PublicKey().RawString(),
		ClientPublicKey:      clientPubKey.Address().String(),
		ApplicationSignature: "",
	}
	var AATValidMessage = AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPrivKey.PublicKey().RawString(),
		ClientPublicKey:      clientPubKey.RawString(),
		ApplicationSignature: "",
	}
	tests := []struct {
		name     string
		aat      AAT
		hasError bool
	}{
		{
			name:     "AAT doesn't have a valid app pub key",
			aat:      AATInvalidAppPubKey,
			hasError: true,
		},
		{
			name:     "AAT doesn't have a valid client pub key",
			aat:      AATInvalidClientPubKey,
			hasError: true,
		},
		{
			name:     "AAT has a valid message",
			aat:      AATValidMessage,
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.aat.ValidateMessage() != nil, tt.hasError)
		})
	}
}

func TestAAT_ValidateSignature(t *testing.T) {
	appPrivKey := GetRandomPrivateKey()
	clientPrivKey := GetRandomPrivateKey()
	var AATMissingSignature = AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPrivKey.PublicKey().RawString(),
		ClientPublicKey:      clientPrivKey.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	var AATInvalidSignature = AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPrivKey.PublicKey().RawString(),
		ClientPublicKey:      clientPrivKey.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	// sign with the client (invalid)
	clientSignature, err := clientPrivKey.Sign(AATInvalidSignature.Hash())
	if err != nil {
		t.Fatalf(err.Error())
	}
	AATInvalidSignature.ApplicationSignature = hex.EncodeToString(clientSignature)
	// sign with the application
	var AATValidSignature = AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPrivKey.PublicKey().RawString(),
		ClientPublicKey:      clientPrivKey.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	appSignature, err := appPrivKey.Sign(AATValidSignature.Hash())
	if err != nil {
		t.Fatalf(err.Error())
	}
	AATValidSignature.ApplicationSignature = hex.EncodeToString(appSignature)
	tests := []struct {
		name     string
		aat      AAT
		hasError bool
	}{
		{
			name:     "AAT doesn't have a signature",
			aat:      AATMissingSignature,
			hasError: true,
		},
		{
			name:     "AAT doesn't have a valid signature",
			aat:      AATInvalidSignature,
			hasError: true,
		},
		{
			name:     "AAT has a valid signature",
			aat:      AATValidSignature,
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.aat.ValidateSignature() != nil, tt.hasError)
		})
	}
}

func TestAAT_HashString(t *testing.T) {
	appPrivKey := GetRandomPrivateKey()
	clientPrivKey := GetRandomPrivateKey()
	var AAT = AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPrivKey.PublicKey().RawString(),
		ClientPublicKey:      clientPrivKey.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	assert.True(t, len(AAT.Hash()) == HashLength)
	assert.True(t, HashVerification(AAT.HashString()) == nil)
}

func TestAAT_Validate(t *testing.T) {
	appPrivKey := GetRandomPrivateKey()
	clientPrivKey := GetRandomPrivateKey()
	var AAT = AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPrivKey.PublicKey().RawString(),
		ClientPublicKey:      clientPrivKey.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	// sign with the client (invalid)
	applicationSignature, err := appPrivKey.Sign(AAT.Hash())
	if err != nil {
		t.Fatalf(err.Error())
	}
	AAT.ApplicationSignature = hex.EncodeToString(applicationSignature)
	assert.Nil(t, AAT.Validate())
}
