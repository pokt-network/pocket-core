package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should never be run in parallel because this replaces
// the global variable `UpgradeFeatureMap`.
func TestCodec_IsOnNamedFeatureActivationHeightWithTolerance(t *testing.T) {
	upgradeKey := ClearUnjailedValSessionKey
	upgradeHeight := int64(20000)
	tolerance := int64(5)

	// have to init global map / not mock friendly.
	originalUpgradeFeatureMap := UpgradeFeatureMap
	UpgradeFeatureMap = map[string]int64{upgradeKey: upgradeHeight}
	t.Cleanup(func() {
		UpgradeFeatureMap = originalUpgradeFeatureMap
	})

	codec := NewCodec(nil)

	// Test zero values / out of bounds
	assert.False(t, codec.IsOnNamedFeatureActivationHeightWithTolerance(
		0,
		upgradeKey,
		tolerance,
	))
	assert.False(t, codec.IsOnNamedFeatureActivationHeightWithTolerance(
		upgradeHeight-tolerance-1,
		upgradeKey,
		tolerance,
	))
	assert.False(t, codec.IsOnNamedFeatureActivationHeightWithTolerance(
		upgradeHeight+tolerance+1,
		upgradeKey,
		tolerance,
	))

	// Test in bounds
	assert.True(t, codec.IsOnNamedFeatureActivationHeightWithTolerance(
		upgradeHeight-tolerance,
		upgradeKey,
		tolerance,
	))
	assert.True(t, codec.IsOnNamedFeatureActivationHeightWithTolerance(
		upgradeHeight,
		upgradeKey,
		tolerance,
	))
	assert.True(t, codec.IsOnNamedFeatureActivationHeightWithTolerance(
		upgradeHeight+tolerance,
		upgradeKey,
		tolerance,
	))
}
