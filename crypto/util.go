package crypto

import (
	"math/rand"
)

// "RandBytes" returns a random string of bytes.
func RandBytes(n int) ([]byte, error) {
	output := make([]byte, n)
	_, err := rand.Read(output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
