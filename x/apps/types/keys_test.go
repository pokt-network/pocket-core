package types

import (
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestAddressFromPrevStateApplicationPowerKey(t *testing.T) {
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

func TestKeyForUnstakingApplications(t *testing.T) {
	type args struct {
		unstakingTime time.Time
	}
	ut := time.Now()

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ut}, append(UnstakingAppsKey, types.FormatTimeBytes(ut)...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForUnstakingApps(tt.args.unstakingTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForUnstakingApplications() = %v, want %v", got, tt.want)
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
		{"sampleByteArray", args{ca}, append(AllApplicationsKey, ca.Bytes()...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForAppByAllApps(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForValByAllVals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyForApplicationBurn(t *testing.T) {
	type args struct {
		address types.Address
	}
	ca, _ := types.AddressFromHex("29f0a60104f3218a2cb51e6a269182d5dc271447114e342086d9c922a106a3c0")

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca}, append(BurnApplicationKey, ca.Bytes()...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForAppBurn(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForApplicationBurn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyForApplicationInStakingSet(t *testing.T) {
	type args struct {
		application Application
	}
	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
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
		{"NewApplication", args{application: NewApplication(types.Address(pub.Address()), pub, []string{"0001"}, types.ZeroInt())}, append([]byte{0x02, 0, 0, 0, 0, 0, 0, 0, 0}, operAddrInvr...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForAppInStakingSet(tt.args.application); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForApplicationInStakingSet() = %s, want %s", got, tt.want)
			}
		})
	}
}
