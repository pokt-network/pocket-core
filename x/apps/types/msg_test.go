package types

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"math/rand"
	"reflect"
	"testing"
)

var msgAppStake MsgAppStake
var msgBeginAppUnstake MsgBeginAppUnstake
var msgAppUnjail MsgAppUnjail

func init() {
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	moduleCdc = codec.New()
	RegisterCodec(moduleCdc)
	codec.RegisterCrypto(moduleCdc)
	moduleCdc.Seal()

	msgAppStake = MsgAppStake{
		Address: sdk.ValAddress(pub.Address()),
		PubKey:  pub,
		Chains:  []string{"886ba5bcb77e1064530052fed1a3f145"},
		Value:   sdk.NewInt(10),
	}
	msgAppUnjail = MsgAppUnjail{sdk.ValAddress(pub.Address())}
	msgBeginAppUnstake = MsgBeginAppUnstake{sdk.ValAddress(pub.Address())}
}

func TestMsgApp_GetSigners(t *testing.T) {
	type args struct {
		msgAppStake MsgAppStake
	}
	tests := []struct {
		name string
		args
		want []sdk.AccAddress
	}{
		{
			name: "return signers",
			args: args{msgAppStake},
			want: []sdk.AccAddress{sdk.AccAddress(msgAppStake.Address)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.msgAppStake.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsgApp_GetSignBytes(t *testing.T) {
	type args struct {
		msgAppStake MsgAppStake
	}
	tests := []struct {
		name string
		args
		want []byte
	}{
		{
			name: "return signers",
			args: args{msgAppStake},
			want: sdk.MustSortJSON(moduleCdc.MustMarshalJSON(msgAppStake)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.msgAppStake.GetSignBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsgApp_Route(t *testing.T) {
	type args struct {
		msgAppStake MsgAppStake
	}
	tests := []struct {
		name string
		args
		want string
	}{
		{
			name: "return signers",
			args: args{msgAppStake},
			want: RouterKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.msgAppStake.Route(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsgApp_Type(t *testing.T) {
	type args struct {
		msgAppStake MsgAppStake
	}
	tests := []struct {
		name string
		args
		want string
	}{
		{
			name: "return signers",
			args: args{msgAppStake},
			want: "stake_application",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.msgAppStake.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsgApp_ValidateBasic(t *testing.T) {
	type args struct {
		msgAppStake MsgAppStake
	}
	tests := []struct {
		name string
		args
		want sdk.Error
		msg  string
	}{
		{
			name: "errs if no Address",
			args: args{MsgAppStake{}},
			want: ErrNilApplicationAddr(DefaultCodespace),
		},
		{
			name: "errs if no stake lower than zero",
			args: args{MsgAppStake{Address: msgAppStake.Address, Value: sdk.NewInt(-1)}},
			want: ErrBadStakeAmount(DefaultCodespace),
		},
		{
			name: "errs if no native chains supported",
			args: args{MsgAppStake{Address: msgAppStake.Address, Value: sdk.NewInt(1), Chains: []string{}}},
			want: ErrNoChains(DefaultCodespace),
		},
		{
			name: "returns err",
			args: args{MsgAppStake{Address: msgAppStake.Address, Value: msgAppStake.Value, Chains: []string{"a"}}},
			want: types.NewInvalidHashLengthError("pocketcore"),
		},
		{
			name: "returns nil if valid address",
			args: args{msgAppStake},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.msgAppStake.ValidateBasic(); got != nil {
				if !reflect.DeepEqual(got.Error(), tt.want.Error()) {
					t.Errorf("ValidatorBasic() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestMsgBeginAppUnstake_GetSigners(t *testing.T) {
	type args struct {
		msgBeginAppUnstake MsgBeginAppUnstake
	}
	tests := []struct {
		name string
		args
		want []sdk.AccAddress
	}{
		{
			name: "return signers",
			args: args{msgBeginAppUnstake},
			want: []sdk.AccAddress{sdk.AccAddress(msgAppStake.Address)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.msgBeginAppUnstake.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsgBeginAppUnstake_GetSignBytes(t *testing.T) {
	type args struct {
		msgBeginAppUnstake MsgBeginAppUnstake
	}
	tests := []struct {
		name string
		args
		want []byte
	}{
		{
			name: "return signers",
			args: args{msgBeginAppUnstake},
			want: sdk.MustSortJSON(moduleCdc.MustMarshalJSON(msgBeginAppUnstake)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.msgBeginAppUnstake.GetSignBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSignBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsgBeginAppUnstake_Route(t *testing.T) {
	type args struct {
		msgBeginAppUnstake MsgBeginAppUnstake
	}
	tests := []struct {
		name string
		args
		want string
	}{
		{
			name: "return signers",
			args: args{msgBeginAppUnstake},
			want: RouterKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.msgBeginAppUnstake.Route(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsgBeginAppUnstake_Type(t *testing.T) {
	type args struct {
		msgBeginAppUnstake MsgBeginAppUnstake
	}
	tests := []struct {
		name string
		args
		want string
	}{
		{
			name: "return signers",
			args: args{msgBeginAppUnstake},
			want: "begin_unstaking_application",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.msgBeginAppUnstake.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsgBeginAppUnstake_ValidateBasic(t *testing.T) {
	type args struct {
		msgBeginAppUnstake MsgBeginAppUnstake
	}
	tests := []struct {
		name string
		args
		want sdk.Error
		msg  string
	}{
		{
			name: "errs if no Address",
			args: args{MsgBeginAppUnstake{}},
			want: ErrNilApplicationAddr(DefaultCodespace),
		},
		{
			name: "returns nil if valid address",
			args: args{msgBeginAppUnstake},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.msgBeginAppUnstake.ValidateBasic(); got != nil {
				if !reflect.DeepEqual(got.Error(), tt.want.Error()) {
					t.Errorf("ValidatorBasic() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestMsgAppUnjail_Route(t *testing.T) {
	type args struct {
		msgAppUnjail MsgAppUnjail
	}
	tests := []struct {
		name string
		args
		want string
	}{
		{
			name: "return signers",
			args: args{msgAppUnjail},
			want: RouterKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.msgAppUnjail.Route(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsgAppUnjail_Type(t *testing.T) {
	type args struct {
		msgAppUnjail MsgAppUnjail
	}
	tests := []struct {
		name string
		args
		want string
	}{
		{
			name: "return signers",
			args: args{msgAppUnjail},
			want: "unjail",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.msgAppUnjail.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsgAppUnjail_GetSigners(t *testing.T) {
	type args struct {
		msgAppUnjail MsgAppUnjail
	}
	tests := []struct {
		name string
		args
		want []sdk.AccAddress
	}{
		{
			name: "return signers",
			args: args{msgAppUnjail},
			want: []sdk.AccAddress{sdk.AccAddress(msgAppUnjail.AppAddr)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.msgAppUnjail.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsgAppUnjail_GetSignBytes(t *testing.T) {
	type args struct {
		msgAppUnjail MsgAppUnjail
	}
	tests := []struct {
		name string
		args
		want []byte
	}{
		{
			name: "return signers",
			args: args{msgAppUnjail},
			want: sdk.MustSortJSON(moduleCdc.MustMarshalJSON(msgAppUnjail)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.msgAppUnjail.GetSignBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsgAppUnjail_ValidateBasic(t *testing.T) {
	type args struct {
		msgAppUnjail MsgAppUnjail
	}
	tests := []struct {
		name string
		args
		want sdk.Error
	}{
		{
			name: "errs if no Address",
			args: args{MsgAppUnjail{}},
			want: ErrBadApplicationAddr(DefaultCodespace),
		},
		{
			name: "returns nil if valid address",
			args: args{msgAppUnjail},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.msgAppUnjail.ValidateBasic(); got != nil {
				if !reflect.DeepEqual(got.Error(), tt.want.Error()) {
					t.Errorf("GetSigners() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
