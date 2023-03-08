package codec

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsBetweenNamedFeatureActivationHeight(t *testing.T) {
	tolerance := int64(5)
	codecUpgradeKey := ClearUnjailedValSessionKey

	// have to init global map / not mock friendly.
	UpgradeFeatureMap = map[string]int64{codecUpgradeKey: 20000}
	codec := NewCodec(nil)

	// Test zero values / out of bounds
	assert.False(t, codec.IsBetweenNamedFeatureActivationHeight(0, codecUpgradeKey, tolerance))
	assert.False(t, codec.IsBetweenNamedFeatureActivationHeight(19994, codecUpgradeKey, tolerance))
	assert.False(t, codec.IsBetweenNamedFeatureActivationHeight(20006, codecUpgradeKey, tolerance))

	// Test in bounds
	assert.True(t, codec.IsBetweenNamedFeatureActivationHeight(19995, codecUpgradeKey, tolerance))
	assert.True(t, codec.IsBetweenNamedFeatureActivationHeight(20005, codecUpgradeKey, tolerance))
	assert.True(t, codec.IsBetweenNamedFeatureActivationHeight(20000, codecUpgradeKey, tolerance))
}
