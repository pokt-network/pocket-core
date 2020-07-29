package types

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	"strings"
	"testing"
)

var codespace = sdk.CodespaceType("app")

func TestError_ErrNoChains(t *testing.T) {
	type args struct {
		codespace sdk.CodespaceType
	}
	tests := []struct {
		name string
		args
		want sdk.Error
	}{
		{
			name: "returns error for stake on unhosted blockchain",
			args: args{codespace},
			want: sdk.NewError(codespace, sdk.CodeType(116), "validator must stake with hosted blockchains"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrNoChains(tt.args.codespace); got.Error() != tt.want.Error() {
				t.Errorf("ErrorNoChains(): returns %v but want %v", got, tt.want)
			}
		})
	}
}
func TestError_ErrNilApplicationAddr(t *testing.T) {
	type args struct {
		codespace sdk.CodespaceType
	}
	tests := []struct {
		name string
		args
		want sdk.Error
	}{
		{
			name: "returns error application address is nil",
			args: args{codespace},
			want: sdk.NewError(codespace, sdk.CodeType(103), "application address is nil"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrNilApplicationAddr(tt.args.codespace); got.Error() != tt.want.Error() {
				t.Errorf("ErrNilApplicationAddr(): returns %v but want %v", got, tt.want)
			}
		})
	}
}
func TestError_ErrApplicaitonStatus(t *testing.T) {
	type args struct {
		codespace sdk.CodespaceType
	}
	tests := []struct {
		name string
		args
		want sdk.Error
	}{
		{
			name: "returns error status is invalid",
			args: args{codespace},
			want: sdk.NewError(codespace, sdk.CodeType(110), "application status is not valid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrApplicationStatus(tt.args.codespace); got.Error() != tt.want.Error() {
				t.Errorf("ErrApplicationStatus(): returns %v but want %v", got, tt.want)
			}
		})
	}
}
func TestError_ErrNoApplicationFound(t *testing.T) {
	type args struct {
		codespace sdk.CodespaceType
	}
	tests := []struct {
		name string
		args
		want sdk.Error
	}{
		{
			name: "returns error application not found for address",
			args: args{codespace},
			want: sdk.NewError(codespace, sdk.CodeType(101), "application does not exist for that address"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrNoApplicationFound(tt.args.codespace); got.Error() != tt.want.Error() {
				t.Errorf("ErrNoApplicationFound(): returns %v but want %v", got, tt.want)
			}
		})
	}
}
func TestError_ErrBadStakeAmount(t *testing.T) {
	type args struct {
		codespace sdk.CodespaceType
	}
	tests := []struct {
		name string
		args
		want sdk.Error
	}{
		{
			name: "returns error stake amount is invalid",
			args: args{codespace},
			want: sdk.NewError(codespace, sdk.CodeType(115), "the stake amount is invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrBadStakeAmount(tt.args.codespace); got.Error() != tt.want.Error() {
				t.Errorf("ErrBadStakeAmount(): returns %v but want %v", got, tt.want)
			}
		})
	}
}
func TestError_ErrNotEnoughCoins(t *testing.T) {
	type args struct {
		codespace sdk.CodespaceType
	}
	tests := []struct {
		name string
		args
		want sdk.Error
	}{
		{
			name: "returns error application does not have enough coins",
			args: args{codespace},
			want: sdk.NewError(codespace, sdk.CodeType(112), "application does not have enough coins in their account"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrNotEnoughCoins(tt.args.codespace); got.Error() != tt.want.Error() {
				t.Errorf("ErrNotEnoughCoins(): returns %v but want %v", got, tt.want)
			}
		})
	}
}
func TestError_ErrApplicaitonPubKeyExists(t *testing.T) {
	type args struct {
		codespace sdk.CodespaceType
	}
	tests := []struct {
		name string
		args
		want sdk.Error
	}{
		{
			name: "returns error application already exists for public key",
			args: args{codespace},
			want: sdk.NewError(codespace, sdk.CodeType(101), "application already exist for this pubkey, must use new application pubkey"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrApplicationPubKeyExists(tt.args.codespace); got.Error() != tt.want.Error() {
				t.Errorf("ErrApplicationPubKeyExists(): returns %v but want %v", got, tt.want)
			}
		})
	}
}
func TestError_ErrMinimumStake(t *testing.T) {
	type args struct {
		codespace sdk.CodespaceType
	}
	tests := []struct {
		name string
		args
		want sdk.Error
	}{
		{
			name: "returns error application staking lower than minimum",
			args: args{codespace},
			want: sdk.NewError(codespace, sdk.CodeType(111), "application isn't staking above the minimum"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrMinimumStake(tt.args.codespace); got.Error() != tt.want.Error() {
				t.Errorf("ErrMinimumStake(): returns %v but want %v", got, tt.want)
			}
		})
	}
}
func TestError_ErrNoApplicationFoundForAddress(t *testing.T) {
	type args struct {
		codespace sdk.CodespaceType
	}
	tests := []struct {
		name string
		args
		want sdk.Error
	}{
		{
			name: "returns error applicaiton not found",
			args: args{codespace},
			want: sdk.NewError(codespace, sdk.CodeType(101), "that address is not associated with any known application"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrNoApplicationForAddress(tt.args.codespace); got.Error() != tt.want.Error() {
				t.Errorf("ErrNoApplicationForAddress(): returns %v but want %v", got, tt.want)
			}
		})
	}
}
func TestError_ErrBadApplicaitonAddr(t *testing.T) {
	type args struct {
		codespace sdk.CodespaceType
	}
	tests := []struct {
		name string
		args
		want sdk.Error
	}{
		{
			name: "returns error application does not exist for address",
			args: args{codespace},
			want: sdk.NewError(codespace, sdk.CodeType(101), "application does not exist for that address"),
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if got := ErrBadApplicationAddr(tt.args.codespace); got.Error() != tt.want.Error() {
				t.Errorf("ErrBadApplicaitonAddr(): returns %v but want %v", got, tt.want)
			}
		})
	}
}
func TestError_ErrApplicationNotJailed(t *testing.T) {
	type args struct {
		codespace sdk.CodespaceType
	}
	tests := []struct {
		name string
		args
		want sdk.Error
	}{
		{
			name: "returns error application is not jailed",
			args: args{codespace},
			want: sdk.NewError(codespace, sdk.CodeType(105), "application not jailed, cannot be unjailed"),
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if got := ErrApplicationNotJailed(tt.args.codespace); got.Error() != tt.want.Error() {
				t.Errorf("ErrApplicationNotJailed(): returns %v but want %v", got, tt.want)
			}
		})
	}
}
func TestError_ErrMissingAppStake(t *testing.T) {
	type args struct {
		codespace sdk.CodespaceType
	}
	tests := []struct {
		name string
		args
		want sdk.Error
	}{
		{
			name: "returns error application has no stake",
			args: args{codespace},
			want: sdk.NewError(codespace, sdk.CodeType(106), "application has no stake; cannot be unjailed"),
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if got := ErrMissingAppStake(tt.args.codespace); got.Error() != tt.want.Error() {
				t.Errorf("ErrMissingAppStake(): returns %v but want %v", got, tt.want)
			}
		})
	}
}
func TestError_ErrStakeTooLow(t *testing.T) {
	type args struct {
		codespace sdk.CodespaceType
	}
	tests := []struct {
		name string
		args
		want sdk.Error
	}{
		{
			name: "returns error application stke lower than delegation",
			args: args{codespace},
			want: sdk.NewError(codespace, sdk.CodeType(105), "application's self delegation less than min stake, cannot be unjailed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrStakeTooLow(tt.args.codespace); got.Error() != tt.want.Error() {
				t.Errorf("ErrStakeTooLow(): returns %v but want %v", got, tt.want)
			}
		})
	}
}
func TestError_ErrApplicationPubKeyTypeNotSupported(t *testing.T) {
	type args struct {
		codespace sdk.CodespaceType
		types     []string
		keyType   string
	}
	tests := []struct {
		name string
		args
		want sdk.Error
	}{
		{
			name: "returns error application does not exist",
			args: args{codespace, []string{"ed25519", "blake2b"}, "int"},
			want: sdk.NewError(
				codespace,
				sdk.CodeType(101),
				fmt.Sprintf("application pubkey type %s is not supported, must use %s", "int", strings.Join([]string{"ed25519", "blake2b"}, ","))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrApplicationPubKeyTypeNotSupported(tt.args.codespace, tt.args.keyType, tt.args.types); got.Error() != tt.want.Error() {
				t.Errorf("ErrApplicationPubKeyTypeNotSupported(): returns %v but want %v", got, tt.want)
			}
		})
	}
}
