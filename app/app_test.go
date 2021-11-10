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
	exp := regexp.MustCompile(`^(BETA|RC)-0\.\d\.\d(\.\d)?$`)

	// first, we verify that the regexp filters unwanted formats
	nonmatchers := []string{"1.0", "BETA1.0", "BETA-1.0", "RC-1.0", "RC-0.11.0", "RC-0.0.11", "RC-0.6.0.11", "RC-0.6.0."}
	for _, matcher := range nonmatchers {
		require.False(t, exp.MatchString(matcher))
	}

	// then we check for some of the desired formats and current version
	matchers := []string{"BETA-0.6.0", "RC-0.6.0", "BETA-0.6.0.0", AppVersion}
	for _, matcher := range matchers {
		require.True(t, exp.MatchString(matcher))
	}
}
