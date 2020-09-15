package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pokt-network/pocket-core/crypto"
	tmCrypto "github.com/tendermint/tendermint/crypto"
	"gopkg.in/yaml.v2"
)

const (
	AddrLen = 20
)

// Address is a common interface for different types of addresses used by the SDK
type AddressI interface {
	Equals(Address) bool
	Empty() bool
	Marshal() ([]byte, error)
	MarshalJSON() ([]byte, error)
	Bytes() []byte
	String() string
	Format(s fmt.State, verb rune)
}

// Ensure that different address types implement the interface
var (
	_ AddressI       = Address{}
	_ yaml.Marshaler = Address{}
)

func VerifyAddressFormat(bz []byte) error {
	verifier := GetConfig().GetAddressVerifier()
	if verifier != nil {
		return verifier(bz)
	}
	if len(bz) != AddrLen {
		return errors.New("Incorrect address length")
	}
	return nil
}

// Address a wrapper around bytes meant to represent an address.
// When marshaled to a string or JSON.
type Address tmCrypto.Address

// Address a wrapper around bytes meant to represent an address.
// When marshaled to a string or JSON.
type addresses []Address

// AddressFromHex creates an Address from a hex string.
func AddressFromHex(address string) (addr Address, err error) {
	if len(address) == 0 {
		return Address{}, nil
	}

	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}
	err = VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

// Returns boolean for whether two Addresses are Equal
func (aa Address) Equals(aa2 Address) bool {
	if aa.Empty() && aa2.Empty() {
		return true
	}

	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Returns boolean for whether an Address is empty
func (aa Address) Empty() bool {
	if aa == nil {
		return true
	}

	aa2 := Address{}
	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (aa Address) Marshal() ([]byte, error) {
	return aa, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (aa *Address) Unmarshal(data []byte) error {
	*aa = data
	return nil
}

// MarshalJSON marshals to JSON.
func (aa Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(aa.String())
}

// MarshalYAML marshals to YAML.
func (aa Address) MarshalYAML() (interface{}, error) {
	return aa.String(), nil
}

// UnmarshalJSON unmarshals from JSON.
func (aa *Address) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	aa2, err := AddressFromHex(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// UnmarshalYAML unmarshals from JSON
func (aa *Address) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	aa2, err := AddressFromHex(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// RawBytes returns the raw address bytes.
func (aa Address) Bytes() []byte {
	return aa
}

// String implements the Stringer interface.
func (aa Address) String() string {
	if aa.Empty() {
		return ""
	}

	str := hex.EncodeToString(aa.Bytes())

	return str
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (aa Address) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = s.Write([]byte(aa.String()))
	case 'p':
		_, _ = s.Write([]byte(fmt.Sprintf("%p", aa)))
	default:
		_, _ = s.Write([]byte(fmt.Sprintf("%X", []byte(aa))))
	}
}

// get Address from pubkey
func GetAddress(pubkey crypto.PublicKey) Address {
	return Address(pubkey.Address())
}
