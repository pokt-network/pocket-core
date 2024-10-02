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

func TestCompareSlices(t *testing.T) {
	assert.True(t, CompareSlices([]string{"1"}, []string{"1"}))
	assert.True(t, CompareSlices([]int{3, 1}, []int{3, 1}))
	assert.False(t, CompareSlices([]int{3, 1}, []int{3, 2}))
	assert.False(t, CompareSlices([]int{3, 1}, []int{3}))

	// Empty and nil slices are identical
	assert.True(t, CompareSlices([]int{}, nil))
	assert.True(t, CompareSlices(nil, []int{}))
	assert.True(t, CompareSlices([]int{}, []int{}))
}
