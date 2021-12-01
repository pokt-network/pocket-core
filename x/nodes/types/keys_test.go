package types

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
)

func TestAddressFromPrevStateValidatorPowerKey(t *testing.T) {
	type args struct {
		key []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{{"sampleByteArray", args{key: []byte{0x51, 0x41, 0x33}}, []byte{0x41, 0x33}}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddressFromKey(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddressFromKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetValMissedBlockKey(t *testing.T) {
	type args struct {
		v types.Address
		i int64
	}
	ca, _ := types.AddressFromHex("29f0a60104f3218a2cb51e6a269182d5dc271447114e342086d9c922a106a3c0")

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(1))

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca, int64(1)}, append(append([]byte{0x12}, ca.Bytes()...), []byte{1, 0, 0, 0, 0, 0, 0, 0}...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Print(ca.String())
			if got := GetValMissedBlockKey(tt.args.v, tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValMissedBlockKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetValMissedBlockPrefixKey(t *testing.T) {
	type args struct {
		v types.Address
	}
	ca, _ := types.AddressFromHex("29f0a60104f3218a2cb51e6a269182d5dc271447114e342086d9c922a106a3c0")

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca}, append([]byte{0x12}, ca.Bytes()...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetValMissedBlockPrefixKey(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValMissedBlockPrefixKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetValidatorSigningInfoAddress(t *testing.T) {
	type args struct {
		key []byte
	}
	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		t.Fatalf(err.Error())
	}
	ca := types.Address(pub.Address())

	tests := []struct {
		name  string
		args  args
		wantV types.Address
	}{
		{"sampleByteArray", args{append([]byte{0x11}, ca.Bytes()...)}, ca.Bytes()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotV, _ := GetValidatorSigningInfoAddress(tt.args.key); !reflect.DeepEqual(gotV, tt.wantV) {
				t.Errorf("GetValidatorSigningInfoAddress() = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

func TestGetValidatorSigningInfoKey(t *testing.T) {
	type args struct {
		v types.Address
	}
	ca, _ := types.AddressFromHex("29f0a60104f3218a2cb51e6a269182d5dc271447114e342086d9c922a106a3c0")

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca}, append([]byte{0x11}, ca.Bytes()...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForValidatorSigningInfo(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForValidatorSigningInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyForUnstakingValidators(t *testing.T) {
	type args struct {
		unstakingTime time.Time
	}
	ut := time.Now()

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ut}, append([]byte{0x41}, types.FormatTimeBytes(ut)...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForUnstakingValidators(tt.args.unstakingTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForUnstakingValidators() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyForValByAllVals(t *testing.T) {
	type args struct {
		addr types.Address
	}
	ca, _ := types.AddressFromHex("29f0a60104f3218a2cb51e6a269182d5dc271447114e342086d9c922a106a3c0")

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca}, append([]byte{0x21}, ca.Bytes()...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForValByAllVals(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForValByAllVals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyForValidatorAward(t *testing.T) {
	type args struct {
		address types.Address
	}
	ca, _ := types.AddressFromHex("29f0a60104f3218a2cb51e6a269182d5dc271447114e342086d9c922a106a3c0")

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca}, append([]byte{0x51}, ca.Bytes()...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForValidatorAward(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForValidatorAward() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyForValidatorBurn(t *testing.T) {
	type args struct {
		address types.Address
	}
	ca, _ := types.AddressFromHex("29f0a60104f3218a2cb51e6a269182d5dc271447114e342086d9c922a106a3c0")

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca}, append([]byte{0x52}, ca.Bytes()...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForValidatorBurn(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForValidatorBurn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyForValidatorInStakingSet(t *testing.T) {
	type args struct {
		validator Validator
	}
	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		t.Fatalf(err.Error())
	}
	operAddrInvr := types.CopyBytes(pub.Address())
	for i, b := range operAddrInvr {
		operAddrInvr[i] = ^b
	}

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"NewValidator", args{validator: NewValidator(types.Address(pub.Address()), pub, []string{"0001"}, "https://www.google.com:443", types.ZeroInt(), types.Address(pub.Address()))}, append([]byte{0x23, 0, 0, 0, 0, 0, 0, 0, 0}, operAddrInvr...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForValidatorInStakingSet(tt.args.validator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForValidatorInStakingSet() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestKeyForValidatorPrevStateStateByPower(t *testing.T) {
	type args struct {
		address types.Address
	}
	ca, _ := types.AddressFromHex("29f0a60104f3218a2cb51e6a269182d5dc271447114e342086d9c922a106a3c0")

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca}, append([]byte{0x31}, ca.Bytes()...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForValidatorPrevStateStateByPower(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForValidatorPrevStateStateByPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStakedValPowerRankKey(t *testing.T) {
	type args struct {
		validator Validator
	}
	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		t.Fatalf(err.Error())
	}
	operAddrInvr := types.CopyBytes(pub.Address())
	for i, b := range operAddrInvr {
		operAddrInvr[i] = ^b
	}

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"NewValidator", args{validator: NewValidator(types.Address(pub.Address()), pub, []string{"0001"}, "https://www.google.com:443", types.ZeroInt(), types.Address(pub.Address()))}, append([]byte{0x23, 0, 0, 0, 0, 0, 0, 0, 0}, operAddrInvr...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStakedValPowerRankKey(tt.args.validator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getStakedValPowerRankKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyForValWaitingToBeginUnstaking(t *testing.T) {
	type args struct {
		addr types.Address
	}

	ca, _ := types.AddressFromHex("29f0a60104f3218a2cb51e6a269182d5dc271447114e342086d9c922a106a3c0")

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca}, append([]byte{0x43}, ca.Bytes()...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForValWaitingToBeginUnstaking(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForValWaitingToBeginUnstaking() = %v, want %v", got, tt.want)
			}
		})
	}
}
