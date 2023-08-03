package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateServiceURL(t *testing.T) {
	validURL := "https://foo.bar:8080"
	// missing prefix
	invalidURLNoPrefix := "foo.bar:8080"
	// wrong prefix
	invalidURLWrongPrefix := "ws://foo.bar:8080"
	// no port
	invalidURLNoPort := "ws://foo.bar"
	// bad port
	invalidURLBadPort := "ws://foo.bar:66666"
	// bad url
	invalidURLBad := "https://foobar:8080"
	assert.Nil(t, ValidateServiceURL(validURL))
	assert.NotNil(t, ValidateServiceURL(invalidURLNoPrefix), "invalid no prefix")
	assert.NotNil(t, ValidateServiceURL(invalidURLWrongPrefix), "invalid wrong prefix")
	assert.NotNil(t, ValidateServiceURL(invalidURLNoPort), "invalid no port")
	assert.NotNil(t, ValidateServiceURL(invalidURLBadPort), "invalid bad port")
	assert.NotNil(t, ValidateServiceURL(invalidURLBad), "invalid bad url")
}

func TestCompareStringMaps(t *testing.T) {
	m1 := map[string]int{}
	m2 := map[string]int{}
	assert.True(t, CompareStringMaps(m1, m2))

	m1["a"] = 10
	m1["b"] = 100
	assert.False(t, CompareStringMaps(m1, m2))

	m2["b"] = 100
	m2["a"] = 10
	assert.True(t, CompareStringMaps(m2, m1))

	m2 = nil
	assert.False(t, CompareStringMaps(m1, m2))
	assert.False(t, CompareStringMaps(nil, m1))

	m1 = nil
	assert.True(t, CompareStringMaps(m1, m2))

	assert.True(t, CompareStringMaps(nil, map[string]int{}))
}
