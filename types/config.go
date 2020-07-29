package types

import (
	"sync"
)

// TmConfig is the structure that holds the SDK configuration parameters.
// This could be used to initialize certain configuration parameters for the SDK.
type Config struct {
	mtx             sync.RWMutex
	sealed          bool
	txEncoder       TxEncoder
	addressVerifier func([]byte) error
}

var (
	// Initializing an instance of TmConfig
	sdkConfig = &Config{
		sealed:    false,
		txEncoder: nil,
	}
)

// GetConfig returns the config instance for the SDK.
func GetConfig() *Config {
	return sdkConfig
}

func (config *Config) assertNotSealed() {
	config.mtx.Lock()
	defer config.mtx.Unlock()

	if config.sealed {
		panic("TmConfig is sealed")
	}
}

// SetTxEncoder builds the TmConfig with TxEncoder used to marshal StdTx to bytes
func (config *Config) SetTxEncoder(encoder TxEncoder) {
	config.assertNotSealed()
	config.txEncoder = encoder
}

// SetAddressVerifier builds the TmConfig with the provided function for verifying that addresses
// have the correct format
func (config *Config) SetAddressVerifier(addressVerifier func([]byte) error) {
	config.assertNotSealed()
	config.addressVerifier = addressVerifier
}

// Set the BIP-0044 CoinType code on the config
func (config *Config) SetCoinType(coinType uint32) {
	config.assertNotSealed()
}

// Seal seals the config such that the config state could not be modified further
func (config *Config) Seal() *Config {
	config.mtx.Lock()
	defer config.mtx.Unlock()

	config.sealed = true
	return config
}

// GetTxEncoder return function to encode transactions
func (config *Config) GetTxEncoder() TxEncoder {
	return config.txEncoder
}

// GetAddressVerifier returns the function to verify that addresses have the correct format
func (config *Config) GetAddressVerifier() func([]byte) error {
	return config.addressVerifier
}
