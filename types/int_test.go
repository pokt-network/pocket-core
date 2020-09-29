package types

import (
	"math/big"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromInt64(t *testing.T) {
	for n := 0; n < 20; n++ {
		r := rand.Int63()
		require.Equal(t, r, NewInt(r).Int64())
	}
}

func TestIntPanic(t *testing.T) {
	// Max BigInt = 2^255-1 = 5.789e+76
	// Min BigInt = -(2^255-1) = -5.789e+76
	require.NotPanics(t, func() { NewIntWithDecimal(1, 76) })
	i1 := NewIntWithDecimal(1, 76)
	require.NotPanics(t, func() { NewIntWithDecimal(2, 76) })
	i2 := NewIntWithDecimal(2, 76)
	require.NotPanics(t, func() { NewIntWithDecimal(3, 76) })
	i3 := NewIntWithDecimal(3, 76)

	require.Panics(t, func() { NewIntWithDecimal(6, 76) })
	require.Panics(t, func() { NewIntWithDecimal(9, 80) })

	// Overflow check
	require.NotPanics(t, func() { i1.Add(i1) })
	require.NotPanics(t, func() { i2.Add(i2) })
	require.Panics(t, func() { i3.Add(i3) })

	require.NotPanics(t, func() { i1.Sub(i1.Neg()) })
	require.NotPanics(t, func() { i2.Sub(i2.Neg()) })
	require.Panics(t, func() { i3.Sub(i3.Neg()) })

	require.Panics(t, func() { i1.Mul(i1) })
	require.Panics(t, func() { i2.Mul(i2) })
	require.Panics(t, func() { i3.Mul(i3) })

	require.Panics(t, func() { i1.Neg().Mul(i1.Neg()) })
	require.Panics(t, func() { i2.Neg().Mul(i2.Neg()) })
	require.Panics(t, func() { i3.Neg().Mul(i3.Neg()) })

	// Underflow check
	i3n := i3.Neg()
	require.NotPanics(t, func() { i3n.Sub(i1) })
	require.NotPanics(t, func() { i3n.Sub(i2) })
	require.Panics(t, func() { i3n.Sub(i3) })

	require.NotPanics(t, func() { i3n.Add(i1.Neg()) })
	require.NotPanics(t, func() { i3n.Add(i2.Neg()) })
	require.Panics(t, func() { i3n.Add(i3.Neg()) })

	require.Panics(t, func() { i1.Mul(i1.Neg()) })
	require.Panics(t, func() { i2.Mul(i2.Neg()) })
	require.Panics(t, func() { i3.Mul(i3.Neg()) })

	// Bound check
	intmax := NewIntFromBigInt(new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(255), nil), big.NewInt(1)))
	intmin := intmax.Neg()
	require.NotPanics(t, func() { intmax.Add(ZeroInt()) })
	require.NotPanics(t, func() { intmin.Sub(ZeroInt()) })
	require.Panics(t, func() { intmax.Add(OneInt()) })
	require.Panics(t, func() { intmin.Sub(OneInt()) })

	// Division-by-zero check
	require.Panics(t, func() { i1.Quo(NewInt(0)) })
}

// Tests below uses randomness
// Since we are using *big.BigInt as underlying value
// and (U/)BigInt is immutable value(see TestImmutability(U/)BigInt)
// it is safe to use randomness in the tests
func TestIdentInt(t *testing.T) {
	for d := 0; d < 1000; d++ {
		n := rand.Int63()
		i := NewInt(n)

		ifromstr, ok := NewIntFromString(strconv.FormatInt(n, 10))
		require.True(t, ok)

		cases := []int64{
			i.Int64(),
			i.BigInt().Int64(),
			ifromstr.Int64(),
			NewIntFromBigInt(big.NewInt(n)).Int64(),
			NewIntWithDecimal(n, 0).Int64(),
		}

		for tcnum, tc := range cases {
			require.Equal(t, n, tc, "BigInt is modified during conversion. tc #%d", tcnum)
		}
	}
}

func minint(i1, i2 int64) int64 {
	if i1 < i2 {
		return i1
	}
	return i2
}

func maxint(i1, i2 int64) int64 {
	if i1 > i2 {
		return i1
	}
	return i2
}

func TestArithInt(t *testing.T) {
	for d := 0; d < 1000; d++ {
		n1 := int64(rand.Int31())
		i1 := NewInt(n1)
		n2 := int64(rand.Int31())
		i2 := NewInt(n2)

		cases := []struct {
			ires BigInt
			nres int64
		}{
			{i1.Add(i2), n1 + n2},
			{i1.Sub(i2), n1 - n2},
			{i1.Mul(i2), n1 * n2},
			{i1.Quo(i2), n1 / n2},
			{i1.AddRaw(n2), n1 + n2},
			{i1.SubRaw(n2), n1 - n2},
			{i1.MulRaw(n2), n1 * n2},
			{i1.QuoRaw(n2), n1 / n2},
			{MinInt(i1, i2), minint(n1, n2)},
			{MaxInt(i1, i2), maxint(n1, n2)},
			{i1.Neg(), -n1},
		}

		for tcnum, tc := range cases {
			require.Equal(t, tc.nres, tc.ires.Int64(), "BigInt arithmetic operation does not match with int64 operation. tc #%d", tcnum)
		}
	}

}

func TestCompInt(t *testing.T) {
	for d := 0; d < 1000; d++ {
		n1 := int64(rand.Int31())
		i1 := NewInt(n1)
		n2 := int64(rand.Int31())
		i2 := NewInt(n2)

		cases := []struct {
			ires bool
			nres bool
		}{
			{i1.Equal(i2), n1 == n2},
			{i1.GT(i2), n1 > n2},
			{i1.LT(i2), n1 < n2},
		}

		for tcnum, tc := range cases {
			require.Equal(t, tc.nres, tc.ires, "BigInt comparison operation does not match with int64 operation. tc #%d", tcnum)
		}
	}
}

func minuint(i1, i2 uint64) uint64 {
	if i1 < i2 {
		return i1
	}
	return i2
}

func maxuint(i1, i2 uint64) uint64 {
	if i1 > i2 {
		return i1
	}
	return i2
}

func randint() BigInt {
	return NewInt(rand.Int63())
}

func TestImmutabilityAllInt(t *testing.T) {
	ops := []func(*BigInt){
		func(i *BigInt) { _ = i.Add(randint()) },
		func(i *BigInt) { _ = i.Sub(randint()) },
		func(i *BigInt) { _ = i.Mul(randint()) },
		func(i *BigInt) { _ = i.Quo(randint()) },
		func(i *BigInt) { _ = i.AddRaw(rand.Int63()) },
		func(i *BigInt) { _ = i.SubRaw(rand.Int63()) },
		func(i *BigInt) { _ = i.MulRaw(rand.Int63()) },
		func(i *BigInt) { _ = i.QuoRaw(rand.Int63()) },
		func(i *BigInt) { _ = i.Neg() },
		func(i *BigInt) { _ = i.IsZero() },
		func(i *BigInt) { _ = i.Sign() },
		func(i *BigInt) { _ = i.Equal(randint()) },
		func(i *BigInt) { _ = i.GT(randint()) },
		func(i *BigInt) { _ = i.LT(randint()) },
		func(i *BigInt) { _ = i.String() },
	}

	for i := 0; i < 1000; i++ {
		n := rand.Int63()
		ni := NewInt(n)

		for opnum, op := range ops {
			op(&ni)

			require.Equal(t, n, ni.Int64(), "BigInt is modified by operation. tc #%d", opnum)
			require.Equal(t, NewInt(n), ni, "BigInt is modified by operation. tc #%d", opnum)
		}
	}
}

type intop func(BigInt, *big.Int) (BigInt, *big.Int)

func intarith(uifn func(BigInt, BigInt) BigInt, bifn func(*big.Int, *big.Int, *big.Int) *big.Int) intop {
	return func(ui BigInt, bi *big.Int) (BigInt, *big.Int) {
		r := rand.Int63()
		br := new(big.Int).SetInt64(r)
		return uifn(ui, NewInt(r)), bifn(new(big.Int), bi, br)
	}
}

func intarithraw(uifn func(BigInt, int64) BigInt, bifn func(*big.Int, *big.Int, *big.Int) *big.Int) intop {
	return func(ui BigInt, bi *big.Int) (BigInt, *big.Int) {
		r := rand.Int63()
		br := new(big.Int).SetInt64(r)
		return uifn(ui, r), bifn(new(big.Int), bi, br)
	}
}

func TestImmutabilityArithInt(t *testing.T) {
	size := 500

	ops := []intop{
		intarith(BigInt.Add, (*big.Int).Add),
		intarith(BigInt.Sub, (*big.Int).Sub),
		intarith(BigInt.Mul, (*big.Int).Mul),
		intarith(BigInt.Quo, (*big.Int).Quo),
		intarithraw(BigInt.AddRaw, (*big.Int).Add),
		intarithraw(BigInt.SubRaw, (*big.Int).Sub),
		intarithraw(BigInt.MulRaw, (*big.Int).Mul),
		intarithraw(BigInt.QuoRaw, (*big.Int).Quo),
	}

	for i := 0; i < 100; i++ {
		uis := make([]BigInt, size)
		bis := make([]*big.Int, size)

		n := rand.Int63()
		ui := NewInt(n)
		bi := new(big.Int).SetInt64(n)

		for j := 0; j < size; j++ {
			op := ops[rand.Intn(len(ops))]
			uis[j], bis[j] = op(ui, bi)
		}

		for j := 0; j < size; j++ {
			require.Equal(t, 0, bis[j].Cmp(uis[j].BigInt()), "BigInt is different from *big.BigInt. tc #%d, BigInt %s, *big.BigInt %s", j, uis[j].String(), bis[j].String())
			require.Equal(t, NewIntFromBigInt(bis[j]), uis[j], "BigInt is different from *big.BigInt. tc #%d, BigInt %s, *big.BigInt %s", j, uis[j].String(), bis[j].String())
			require.True(t, uis[j].i != bis[j], "Pointer Addresses are equal. tc #%d, BigInt %s, *big.BigInt %s", j, uis[j].String(), bis[j].String())
		}
	}
}

func TestEncodingRandom(t *testing.T) {
	for i := 0; i < 1000; i++ {
		n := rand.Int63()
		ni := NewInt(n)
		var ri BigInt

		str, err := ni.MarshalAmino()
		require.Nil(t, err)
		err = (&ri).UnmarshalAmino(str)
		require.Nil(t, err)

		require.Equal(t, ni, ri, "MarshalAmino * UnmarshalAmino is not identity. tc #%d, Expected %s, Actual %s", i, ni.String(), ri.String())
		require.True(t, ni.i != ri.i, "Pointer Addresses are equal. tc #%d", i)

		bz, err := ni.MarshalJSON()
		require.Nil(t, err)
		err = (&ri).UnmarshalJSON(bz)
		require.Nil(t, err)

		require.Equal(t, ni, ri, "MarshalJSON * UnmarshalJSON is not identity. tc #%d, Expected %s, Actual %s", i, ni.String(), ri.String())
		require.True(t, ni.i != ri.i, "Pointer Addresses are equal. tc #%d", i)
	}

	for i := 0; i < 1000; i++ {
		n := rand.Uint64()
		ni := NewUint(n)
		var ri Uint

		str, err := ni.MarshalAmino()
		require.Nil(t, err)
		err = (&ri).UnmarshalAmino(str)
		require.Nil(t, err)

		require.Equal(t, ni, ri, "MarshalAmino * UnmarshalAmino is not identity. tc #%d, Expected %s, Actual %s", i, ni.String(), ri.String())
		require.True(t, ni.i != ri.i, "Pointer Addresses are equal. tc #%d", i)

		bz, err := ni.MarshalJSON()
		require.Nil(t, err)
		err = (&ri).UnmarshalJSON(bz)
		require.Nil(t, err)

		require.Equal(t, ni, ri, "MarshalJSON * UnmarshalJSON is not identity. tc #%d, Expected %s, Actual %s", i, ni.String(), ri.String())
		require.True(t, ni.i != ri.i, "Pointer Addresses are equal. tc #%d", i)
	}
}

func TestEncodingTableInt(t *testing.T) {

	cases := []struct {
		i   BigInt
		bz  []byte
		str string
	}{
		{NewInt(0), []byte("\"0\""), "0"},
		{NewInt(100), []byte("\"100\""), "100"},
		{NewInt(51842), []byte("\"51842\""), "51842"},
		{NewInt(19513368), []byte("\"19513368\""), "19513368"},
		{NewInt(999999999999), []byte("\"999999999999\""), "999999999999"},
	}

	for tcnum, tc := range cases {
		bz, err := tc.i.MarshalJSON()
		require.Nil(t, err, "Error marshaling BigInt. tc #%d, err %s", tcnum, err)
		require.Equal(t, tc.bz, bz, "Marshaled value is different from exported. tc #%d", tcnum)
		var i BigInt
		err = (&i).UnmarshalJSON(bz)
		require.Nil(t, err, "Error unmarshaling BigInt. tc #%d, err %s", tcnum, err)
		require.Equal(t, tc.i, i, "Unmarshaled value is different from exported. tc #%d", tcnum)

		str, err := tc.i.MarshalAmino()
		require.Nil(t, err, "Error marshaling BigInt. tc #%d, err %s", tcnum, err)
		require.Equal(t, tc.str, string(str), "Marshaled value is different from exported. tc #%d", tcnum)
		err = (&i).UnmarshalAmino(str)
		require.Nil(t, err, "Error unmarshaling BigInt. tc #%d, err %s", tcnum, err)
		require.Equal(t, tc.i, i, "Unmarshaled value is different from exported. tc #%d", tcnum)
	}
}

func TestEncodingTableUint(t *testing.T) {
	var i Uint

	cases := []struct {
		i   Uint
		bz  []byte
		str string
	}{
		{NewUint(0), []byte("\"0\""), "0"},
		{NewUint(100), []byte("\"100\""), "100"},
		{NewUint(51842), []byte("\"51842\""), "51842"},
		{NewUint(19513368), []byte("\"19513368\""), "19513368"},
		{NewUint(999999999999), []byte("\"999999999999\""), "999999999999"},
	}

	for tcnum, tc := range cases {
		bz, err := tc.i.MarshalJSON()
		require.Nil(t, err, "Error marshaling BigInt. tc #%d, err %s", tcnum, err)
		require.Equal(t, tc.bz, bz, "Marshaled value is different from exported. tc #%d", tcnum)
		err = (&i).UnmarshalJSON(bz)
		require.Nil(t, err, "Error unmarshaling BigInt. tc #%d, err %s", tcnum, err)
		require.Equal(t, tc.i, i, "Unmarshaled value is different from exported. tc #%d", tcnum)

		str, err := tc.i.MarshalAmino()
		require.Nil(t, err, "Error marshaling BigInt. tc #%d, err %s", tcnum, err)
		require.Equal(t, tc.str, string(str), "Marshaled value is different from exported. tc #%d", tcnum)
		err = (&i).UnmarshalAmino(str)
		require.Nil(t, err, "Error unmarshaling BigInt. tc #%d, err %s", tcnum, err)
		require.Equal(t, tc.i, i, "Unmarshaled value is different from exported. tc #%d", tcnum)
	}
}

func TestSerializationOverflow(t *testing.T) {
	bx, _ := new(big.Int).SetString("91888242871839275229946405745257275988696311157297823662689937894645226298583", 10)
	x := BigInt{bx}
	y := new(BigInt)

	// require amino deserialization to fail due to overflow
	xStr, err := x.MarshalAmino()
	require.NoError(t, err)

	err = y.UnmarshalAmino(xStr)
	require.Error(t, err)

	// require JSON deserialization to fail due to overflow
	bz, err := x.MarshalJSON()
	require.NoError(t, err)

	err = y.UnmarshalJSON(bz)
	require.Error(t, err)
}
