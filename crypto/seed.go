// This package is for cryptography that is used in Pocket Core.
package crypto

import (
	"math/rand"
	"time"
)

// "seed.go" holds cryptographic seed generation.

/*
"GenerateSeed" creates the random seed from nanosecond.
*/
func GenerateSeed() {
	rand.Seed(time.Now().UTC().UnixNano()) // used for random generators.
}
