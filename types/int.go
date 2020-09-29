package types

import (
	"encoding"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	"math/big"
)

const maxBitLen = 255

func newIntegerFromString(s string) (*big.Int, bool) {
	return new(big.Int).SetString(s, 0)
}

func equal(i *big.Int, i2 *big.Int) bool { return i.Cmp(i2) == 0 }

func gt(i *big.Int, i2 *big.Int) bool { return i.Cmp(i2) == 1 }

func gte(i *big.Int, i2 *big.Int) bool { return i.Cmp(i2) >= 0 }

func lt(i *big.Int, i2 *big.Int) bool { return i.Cmp(i2) == -1 }

func lte(i *big.Int, i2 *big.Int) bool { return i.Cmp(i2) <= 0 }

func add(i *big.Int, i2 *big.Int) *big.Int { return new(big.Int).Add(i, i2) }

func sub(i *big.Int, i2 *big.Int) *big.Int { return new(big.Int).Sub(i, i2) }

func mul(i *big.Int, i2 *big.Int) *big.Int { return new(big.Int).Mul(i, i2) }

func div(i *big.Int, i2 *big.Int) *big.Int { return new(big.Int).Quo(i, i2) }

func mod(i *big.Int, i2 *big.Int) *big.Int { return new(big.Int).Mod(i, i2) }

func neg(i *big.Int) *big.Int { return new(big.Int).Neg(i) }

func min(i *big.Int, i2 *big.Int) *big.Int {
	if i.Cmp(i2) == 1 {
		return new(big.Int).Set(i2)
	}

	return new(big.Int).Set(i)
}

func max(i *big.Int, i2 *big.Int) *big.Int {
	if i.Cmp(i2) == -1 {
		return new(big.Int).Set(i2)
	}

	return new(big.Int).Set(i)
}

func unmarshalText(i *big.Int, text string) error {
	if err := i.UnmarshalText([]byte(text)); err != nil {
		return err
	}

	if i.BitLen() > maxBitLen {
		return fmt.Errorf("integer out of range: %s", text)
	}

	return nil
}

var _ CustomProtobufType = (*BigInt)(nil)
var _ codec.ProtoMarshaler = &BigInt{}

// BigInt wraps integer with 256 bit range bound
// Checks overflow, underflow and division by zero
// Exists in range from -(2^maxBitLen-1) to 2^maxBitLen-1
type BigInt struct {
	i *big.Int
}

// NewInt constructs BigInt from int64
func NewInt(n int64) BigInt {
	return BigInt{big.NewInt(n)}
}

// NewIntFromUint64 constructs an BigInt from a uint64.
func NewIntFromUint64(n uint64) BigInt {
	b := big.NewInt(0)
	b.SetUint64(n)
	return BigInt{b}
}

// NewIntFromBigInt constructs BigInt from big.BigInt
func NewIntFromBigInt(i *big.Int) BigInt {
	if i.BitLen() > maxBitLen {
		panic("NewIntFromBigInt() out of bound")
	}
	return BigInt{i}
}

// NewIntFromString constructs BigInt from string
func NewIntFromString(s string) (res BigInt, ok bool) {
	i, ok := newIntegerFromString(s)
	if !ok {
		return
	}
	// Check overflow
	if i.BitLen() > maxBitLen {
		ok = false
		return
	}
	return BigInt{i}, true
}

// NewIntWithDecimal constructs BigInt with decimal
// Result value is n*10^dec
func NewIntWithDecimal(n int64, dec int) BigInt {
	if dec < 0 {
		panic("NewIntWithDecimal() decimal is negative")
	}
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(dec)), nil)
	i := new(big.Int)
	i.Mul(big.NewInt(n), exp)

	// Check overflow
	if i.BitLen() > maxBitLen {
		panic("NewIntWithDecimal() out of bound")
	}
	return BigInt{i}
}

// ZeroInt returns BigInt value with zero
func ZeroInt() BigInt { return BigInt{big.NewInt(0)} }

// OneInt returns BigInt value with one
func OneInt() BigInt { return BigInt{big.NewInt(1)} }

// ToDec converts BigInt to BigDec
func (i BigInt) ToDec() BigDec {
	return NewDecFromInt(i)
}

// Int64 converts BigInt to int64
// Panics if the value is out of range
func (i BigInt) Int64() int64 {
	if !i.i.IsInt64() {
		panic("Int64() out of bound")
	}
	return i.i.Int64()
}

// IsInt64 returns true if Int64() not panics
func (i BigInt) IsInt64() bool {
	return i.i.IsInt64()
}

// Uint64 converts BigInt to uint64
// Panics if the value is out of range
func (i BigInt) Uint64() uint64 {
	if !i.i.IsUint64() {
		panic("Uint64() out of bounds")
	}
	return i.i.Uint64()
}

// IsUint64 returns true if Uint64() not panics
func (i BigInt) IsUint64() bool {
	return i.i.IsUint64()
}

// IsZero returns true if BigInt is zero
func (i BigInt) IsZero() bool {
	if i.i == nil {
		return true
	}
	return i.i.Sign() == 0
}

// IsNegative returns true if BigInt is negative
func (i BigInt) IsNegative() bool {
	return i.i.Sign() == -1
}

// IsPositive returns true if BigInt is positive
func (i BigInt) IsPositive() bool {
	return i.i.Sign() == 1
}

// Sign returns sign of BigInt
func (i BigInt) Sign() int {
	return i.i.Sign()
}

// Equal compares two Ints
func (i BigInt) Equal(i2 BigInt) bool {
	return equal(i.i, i2.i)
}

// GT returns true if first BigInt is greater than second
func (i BigInt) GT(i2 BigInt) bool {
	return gt(i.i, i2.i)
}

// GTE returns true if receiver BigInt is greater than or equal to the parameter
// BigInt.
func (i BigInt) GTE(i2 BigInt) bool {
	return gte(i.i, i2.i)
}

// LT returns true if first BigInt is lesser than second
func (i BigInt) LT(i2 BigInt) bool {
	return lt(i.i, i2.i)
}

// LTE returns true if first BigInt is less than or equal to second
func (i BigInt) LTE(i2 BigInt) bool {
	return lte(i.i, i2.i)
}

// Add adds BigInt from another
func (i BigInt) Add(i2 BigInt) (res BigInt) {
	res = BigInt{add(i.i, i2.i)}
	// Check overflow
	if res.i.BitLen() > maxBitLen {
		panic("BigInt overflow")
	}
	return
}

// AddRaw adds int64 to BigInt
func (i BigInt) AddRaw(i2 int64) BigInt {
	return i.Add(NewInt(i2))
}

// Sub subtracts BigInt from another
func (i BigInt) Sub(i2 BigInt) (res BigInt) {
	res = BigInt{sub(i.i, i2.i)}
	// Check overflow
	if res.i.BitLen() > maxBitLen {
		panic("BigInt overflow")
	}
	return
}

// SubRaw subtracts int64 from BigInt
func (i BigInt) SubRaw(i2 int64) BigInt {
	return i.Sub(NewInt(i2))
}

// Mul multiples two Ints
func (i BigInt) Mul(i2 BigInt) (res BigInt) {
	// Check overflow
	if i.i.BitLen()+i2.i.BitLen()-1 > maxBitLen {
		panic("BigInt overflow")
	}
	res = BigInt{mul(i.i, i2.i)}
	// Check overflow if sign of both are same
	if res.i.BitLen() > maxBitLen {
		panic("BigInt overflow")
	}
	return
}

// MulRaw multipies BigInt and int64
func (i BigInt) MulRaw(i2 int64) BigInt {
	return i.Mul(NewInt(i2))
}

// Quo divides BigInt with BigInt
func (i BigInt) Quo(i2 BigInt) (res BigInt) {
	// Check division-by-zero
	if i2.i.Sign() == 0 {
		panic("Division by zero")
	}
	return BigInt{div(i.i, i2.i)}
}

// QuoRaw divides BigInt with int64
func (i BigInt) QuoRaw(i2 int64) BigInt {
	return i.Quo(NewInt(i2))
}

// Mod returns remainder after dividing with BigInt
func (i BigInt) Mod(i2 BigInt) BigInt {
	if i2.Sign() == 0 {
		panic("division-by-zero")
	}
	return BigInt{mod(i.i, i2.i)}
}

// ModRaw returns remainder after dividing with int64
func (i BigInt) ModRaw(i2 int64) BigInt {
	return i.Mod(NewInt(i2))
}

// Neg negates BigInt
func (i BigInt) Neg() (res BigInt) {
	return BigInt{neg(i.i)}
}

// return the minimum of the ints
func MinInt(i1, i2 BigInt) BigInt {
	return BigInt{min(i1.BigInt(), i2.BigInt())}
}

// MaxInt returns the maximum between two integers.
func MaxInt(i, i2 BigInt) BigInt {
	return BigInt{max(i.BigInt(), i2.BigInt())}
}

// Human readable string
func (i BigInt) String() string {
	return i.i.String()
}

func (i *BigInt) Reset() {
	*i = BigInt{}
}

func (i BigInt) ProtoMessage() {}

// autogenerated
func (i BigInt) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	a := len(dAtA)
	_ = a
	var l int
	_ = l
	{
		size := i.Size()
		a -= size
		if _, err := i.MarshalTo(dAtA[a:]); err != nil {
			return 0, err
		}
		a = encodeVarintCoin(dAtA, a, uint64(size))
	}
	a--
	dAtA[a] = 0xa
	return len(dAtA) - a, nil
}

// BigInt converts BigInt to big.BigInt
func (i BigInt) BigInt() *big.Int {
	return new(big.Int).Set(i.i)
}

// Size implements the gogo proto custom type interface.
func (i *BigInt) Size() int {
	bz, _ := i.Marshal()
	return len(bz)
}

// MarshalObject implements the gogo proto custom type interface.
func (i BigInt) Marshal() ([]byte, error) {
	if i.i == nil {
		i.i = new(big.Int)
	}
	t, err := i.i.MarshalText()
	return t, err
}

// MarshalTo implements the gogo proto custom type interface.
func (i *BigInt) MarshalTo(data []byte) (n int, err error) {
	if i.i == nil {
		i.i = new(big.Int)
	}
	if len(i.i.Bytes()) == 0 {
		copy(data, []byte{0x30})
		return 1, nil
	}

	bz, err := i.Marshal()
	if err != nil {
		return 0, err
	}

	copy(data, bz)
	return len(bz), nil
}

// UnmarshalObject implements the gogo proto custom type interface.
func (i *BigInt) Unmarshal(data []byte) error {
	if len(data) == 0 {
		i = nil
		return nil
	}

	if i.i == nil {
		i.i = new(big.Int)
	}

	if err := i.i.UnmarshalText(data); err != nil {
		return err
	}

	if i.i.BitLen() > maxBitLen {
		return fmt.Errorf("integer out of range; got: %d, max: %d", i.i.BitLen(), maxBitLen)
	}

	return nil
}

// MarshalJSON defines custom encoding scheme
func (i BigInt) MarshalJSON() ([]byte, error) {
	if i.i == nil { // Necessary since default Uint initialization has i.i as nil
		i.i = new(big.Int)
	}
	return marshalJSON(i.i)
}

// UnmarshalJSON defines custom decoding scheme
func (i *BigInt) UnmarshalJSON(bz []byte) error {
	if i.i == nil { // Necessary since default BigInt initialization has i.i as nil
		i.i = new(big.Int)
	}
	return unmarshalJSON(i.i, bz)
}

// MarshalJSON for custom encoding scheme
// Must be encoded as a string for JSON precision
func marshalJSON(i encoding.TextMarshaler) ([]byte, error) {
	text, err := i.MarshalText()
	if err != nil {
		return nil, err
	}

	return json.Marshal(string(text))
}

// UnmarshalJSON for custom decoding scheme
// Must be encoded as a string for JSON precision
func unmarshalJSON(i *big.Int, bz []byte) error {
	var text string
	if err := json.Unmarshal(bz, &text); err != nil {
		return err
	}

	return unmarshalText(i, text)
}

// MarshalYAML returns the YAML representation.
func (i BigInt) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// Override Amino binary serialization by proxying to protobuf.
func (i BigInt) MarshalAmino() ([]byte, error)   { return i.Marshal() }
func (i *BigInt) UnmarshalAmino(bz []byte) error { return i.Unmarshal(bz) }
