package app

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
	Our versioning system has two possible prefixes: BETA or RC.
	Besides that, it has at least 3, possibly 4 numbers: 0.X.Y or 0.X.Y.Z.
        As used in comparison for transaction upgrade, it does not support double digits in any of the numbers.

*/
func TestAppVersionIsSensible(t *testing.T) {
	r, e := regexp.MatchString(`(BETA|RC)-0\.\d\.\d[\.\d]?`, AppVersion)
	require.True(t, r)
	require.Nil(t, e)
}
