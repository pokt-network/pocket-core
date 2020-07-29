package mintkey

import (
	"fmt"
	"github.com/tendermint/tendermint/crypto"
	"golang.org/x/crypto/scrypt"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkBcryptGenerateFromPassword(b *testing.B) {
	passphrase := []byte("passphrase")
	for securityParam := 9; securityParam < 16; securityParam++ {
		param := securityParam
		b.Run(fmt.Sprintf("benchmark-security-param-%d", param), func(b *testing.B) {
			saltBytes := crypto.CRandBytes(16)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := scrypt.Key(passphrase, saltBytes, n, r, p, klen)
				require.Nil(b, err)
			}
		})
	}
}
