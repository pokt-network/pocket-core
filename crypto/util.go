package crypto

import (
	"crypto/sha1"
	"fmt"
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

func NewSHA1Hash() (string, error) {
	randBytes, err := RandBytes(16)
	if err != nil {
		return "", err
	}
	hash := sha1.New()
	hash.Write(randBytes)
	bs := hash.Sum(nil)
	return fmt.Sprintf("%x", bs), nil
}