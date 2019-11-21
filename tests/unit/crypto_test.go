package unit

import (
	"crypto/elliptic"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pokt-network/pocket-core/crypto"
	"reflect"
	"testing"
)

func TestS256(t *testing.T) {
	tests := []struct {
		name string
		want elliptic.Curve
	}{
		// TODO: Add test cases.
		{"Generates secp256k1 elliptic.Curve", secp256k1.S256()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := crypto.S256(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("S256() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA3FromBytes(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want []byte
	}{
		{"SHA3-256 From Bytes", []byte{}, []byte{167, 255, 198, 248, 191, 30, 215, 102, 81, 193, 71, 86, 160, 97, 214, 98, 245, 128, 255, 77, 228, 59, 73, 250, 130, 216, 10, 75, 128, 248, 67, 74}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := crypto.SHA3FromBytes(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SHA3FromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA3FromString(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want []byte
	}{
		{"SHA3-256 From HashHex", "foo", []byte{118, 211, 188, 65, 201, 245, 136, 247, 252, 208, 213, 191, 71, 24, 248, 248, 75, 28, 65, 178, 8, 130, 112, 49, 0, 185, 235, 148, 19, 128, 124, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := crypto.SHA3FromString(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SHA3FromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPrivateKey(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Test Private Key Generation Succesful"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := crypto.NewPrivateKey()
			if err != nil {
				t.Errorf("NewPrivateKey() error = %v", err)
				return
			}
		})
	}
}

func TestPrivateKey_GetPublicKey(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Test Public Key retrieval from Private Key"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privKey, err := crypto.NewPrivateKey()
			if err != nil {
				t.Errorf("Error calling NewPrivateKey() error = %v", err)
				return
			}
			_ = privKey.GetPublicKey()
		})
	}
}

func TestPrivateKey_Bytes(t *testing.T) {
	tests := []struct {
		name         string
		wantedLength int
	}{
		{"Test Private Key Bytes Length of 32 bytes", 32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privKey, err := crypto.NewPrivateKey()
			if err != nil {
				t.Errorf("Error calling NewPrivateKey() error = %v", err)
				return
			}
			got := privKey.Bytes()
			if len(got) != tt.wantedLength {
				t.Errorf("Bytes() error, wantedLength = %v, got = %v", tt.wantedLength, len(got))
				return
			}
		})
	}
}

func TestPublicKey_Bytes(t *testing.T) {
	tests := []struct {
		name         string
		wantedLength int
	}{
		{"Test Private Key Bytes of 65 bytes", 65},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, pubKey, err := crypto.NewKeypair()
			if err != nil {
				t.Errorf("Error calling NewPrivateKey() error = %v", err)
				return
			}
			got := pubKey.Bytes()
			if len(got) != tt.wantedLength {
				t.Errorf("Bytes() error, wantedLength = %v, got = %v", tt.wantedLength, len(got))
				return
			}
		})
	}
}

func TestPrivateKey_Sign(t *testing.T) {
	tests := []struct {
		name         string
		wantedLength int
	}{
		{"Test Private Key ApplicationSignature Length of 65 bytes", 65},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privKey, _, err := crypto.NewKeypair()
			if err != nil {
				t.Errorf("NewKeypair() error = %v", err)
				return
			}
			signature, errSign := privKey.Sign(crypto.SHA3FromString("foo"))
			if errSign != nil {
				t.Errorf("Sign() error = %v", errSign)
				return
			}
			if len(signature) != tt.wantedLength {
				t.Errorf("Sign() produced invalid length signature, got = %v, expected = %v", len(signature), tt.wantedLength)
				return
			}
		})
	}
}

func TestNewKeypair(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Test Key Pair Generation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := crypto.NewKeypair()
			if err != nil {
				t.Errorf("NewKeypair() error = %v", err)
				return
			}
		})
	}
}

func TestGetPublicKeyBytesFromSignature(t *testing.T) {
	tests := []struct {
		name         string
		wantedLength int
	}{
		{"Test Public Key Retrieval From ApplicationSignature", 65},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privKey, privKeyErr := crypto.NewPrivateKey()
			if privKeyErr != nil {
				t.Errorf("NewPrivateKey() error = %v", privKeyErr)
				return
			}

			messageHash := crypto.SHA3FromString("foo")
			signature, signError := privKey.Sign(messageHash)
			if signError != nil {
				t.Errorf("Sign() error = %v", signError)
				return
			}

			pubKeyBytes, pubKeyErr := crypto.GetPublicKeyBytesFromSignature(crypto.SHA3FromBytes(messageHash), signature)
			if pubKeyErr != nil {
				t.Errorf("GetPublicKeyFromSignature() error = %v", pubKeyErr)
				return
			}
			if len(pubKeyBytes) != tt.wantedLength {
				t.Errorf("GetPublicKeyFromSignature() Bytes Length = %v, want = %v", len(pubKeyBytes), tt.wantedLength)
				return
			}
		})
	}
}

func TestVerifySignature(t *testing.T) {
	tests := []struct {
		name                     string
		wantedLength             int
		wantedVerificationResult bool
	}{
		{"Test Public Key Retrieval From ApplicationSignature", 65, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privKey, privKeyErr := crypto.NewPrivateKey()
			if privKeyErr != nil {
				t.Errorf("NewPrivateKey() error = %v", privKeyErr)
				return
			}

			messageHash := crypto.SHA3FromString("foo")
			signature, signError := privKey.Sign(messageHash)
			if signError != nil {
				t.Errorf("Sign() error = %v", signError)
				return
			}

			pubKeyBytes, pubKeyErr := crypto.GetPublicKeyBytesFromSignature(messageHash, signature)
			if pubKeyErr != nil {
				t.Errorf("GetPublicKeyFromSignature() error = %v", pubKeyErr)
				return
			}
			if len(pubKeyBytes) != tt.wantedLength {
				t.Errorf("GetPublicKeyFromSignature() Bytes Length = %v, want = %v", len(pubKeyBytes), tt.wantedLength)
				return
			}

			verificationResult := crypto.VerifySignature(privKey.GetPublicKey(), messageHash, signature)
			if verificationResult != tt.wantedVerificationResult {
				t.Errorf("VerifySignature() error produced, wanted = %v, got = %v", tt.wantedVerificationResult, verificationResult)
			}
		})
	}
}

func TestVerifySignatureWithPubKeyBytes(t *testing.T) {
	tests := []struct {
		name                     string
		wantedLength             int
		wantedVerificationResult bool
	}{
		{"Test Public Key Retrieval From ApplicationSignature", 65, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privKey, privKeyErr := crypto.NewPrivateKey()
			if privKeyErr != nil {
				t.Errorf("NewPrivateKey() error = %v", privKeyErr)
				return
			}

			messageHash := crypto.SHA3FromString("foo")
			signature, signError := privKey.Sign(messageHash)
			if signError != nil {
				t.Errorf("Sign() error = %v", signError)
				return
			}

			pubKeyBytes, pubKeyErr := crypto.GetPublicKeyBytesFromSignature(messageHash, signature)
			if pubKeyErr != nil {
				t.Errorf("GetPublicKeyFromSignature() error = %v", pubKeyErr)
				return
			}
			if len(pubKeyBytes) != tt.wantedLength {
				t.Errorf("GetPublicKeyFromSignature() Bytes Length = %v, want = %v", len(pubKeyBytes), tt.wantedLength)
				return
			}

			verificationResult := crypto.VerifySignatureWithPubKeyBytes(pubKeyBytes, messageHash, signature)
			if verificationResult != tt.wantedVerificationResult {
				t.Errorf("VerifySignature() error produced, wanted = %v, got = %v", tt.wantedVerificationResult, verificationResult)
			}
		})
	}
}
