package types

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
)

func TestLegacyMsgBeginUnstake_GetSignBytes(t *testing.T) {
	type fields struct {
		Address sdk.Address
	}
	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())

	mesg := LegacyMsgBeginUnstake{
		Address: va,
	}

	encodedmsg, _ := ModuleCdc.MarshalJSON(&mesg)

	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{"Test GetSignBytes", fields{va}, encodedmsg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := LegacyMsgBeginUnstake{
				Address: tt.fields.Address,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestLegacyMsgBeginUnstake_GetSigners(t *testing.T) {
	type fields struct {
		Address sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())

	mesg := MsgBeginUnstake{
		Address: va,
	}

	tests := []struct {
		name   string
		fields fields
		want   []sdk.Address
	}{
		{"Test GetSigners", fields{va}, []sdk.Address{sdk.Address(mesg.Address)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := LegacyMsgBeginUnstake{
				Address: tt.fields.Address,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyMsgBeginUnstake_Route(t *testing.T) {
	type fields struct {
		Address sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Test Route", fields{va}, ModuleName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := LegacyMsgBeginUnstake{
				Address: tt.fields.Address,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyMsgBeginUnstake_Type(t *testing.T) {
	type fields struct {
		Address sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Test Type", fields{va}, MsgUnstakeName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := LegacyMsgBeginUnstake{
				Address: tt.fields.Address,
			}
			if got := msg.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyMsgBeginUnstake_ValidateBasic(t *testing.T) {
	type fields struct {
		Address sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())

	tests := []struct {
		name   string
		fields fields
		want   sdk.Error
	}{
		{"Test Validate Basic error", fields{nil}, sdk.NewError(codespace, CodeInvalidInput, "validator address is nil")},
		{"Test Validate Basic pass", fields{va}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := LegacyMsgBeginUnstake{
				Address: tt.fields.Address,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyMsgStake_GetSigners(t *testing.T) {
	type fields struct {
		Address    sdk.Address
		PubKey     crypto.PublicKey
		Chains     []string
		Value      sdk.BigInt
		ServiceURL string
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	chains := []string{"0001"}
	value := sdk.OneInt()
	surl := "www.pokt.network"

	tests := []struct {
		name   string
		fields fields
		want   []sdk.Address
	}{
		{"Test GetSigners", fields{
			PubKey:     pub,
			Chains:     chains,
			Value:      value,
			ServiceURL: surl,
		}, []sdk.Address{sdk.Address(pub.Address())}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := LegacyMsgStake{
				PublicKey:  tt.fields.PubKey,
				Chains:     tt.fields.Chains,
				Value:      tt.fields.Value,
				ServiceUrl: tt.fields.ServiceURL,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyMsgStake_Route(t *testing.T) {
	type fields struct {
		Address    sdk.Address
		PubKey     crypto.PublicKey
		Chains     []string
		Value      sdk.BigInt
		ServiceURL string
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	chains := []string{"0001"}
	value := sdk.OneInt()
	surl := "www.pokt.network"

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Test Route", fields{
			PubKey:     pub,
			Chains:     chains,
			Value:      value,
			ServiceURL: surl,
		}, RouterKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := LegacyMsgStake{
				PublicKey:  tt.fields.PubKey,
				Chains:     tt.fields.Chains,
				Value:      tt.fields.Value,
				ServiceUrl: tt.fields.ServiceURL,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyMsgStake_Type(t *testing.T) {
	type fields struct {
		Address    sdk.Address
		PubKey     crypto.PublicKey
		Chains     []string
		Value      sdk.BigInt
		ServiceURL string
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	chains := []string{"0001"}
	value := sdk.OneInt()
	surl := "www.pokt.network"

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Test Type", fields{
			PubKey:     pub,
			Chains:     chains,
			Value:      value,
			ServiceURL: surl,
		}, "stake_validator"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := LegacyMsgStake{
				PublicKey:  tt.fields.PubKey,
				Chains:     tt.fields.Chains,
				Value:      tt.fields.Value,
				ServiceUrl: tt.fields.ServiceURL,
			}
			if got := msg.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyMsgStake_ValidateBasic(t *testing.T) {
	type fields struct {
		Address    sdk.Address
		PubKey     crypto.PublicKey
		Chains     []string
		Value      sdk.BigInt
		ServiceURL string
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	chains := []string{"0001"}
	value := sdk.OneInt()
	surl := "https://www.pokt.network:8080"

	tests := []struct {
		name   string
		fields fields
		want   sdk.Error
	}{
		{"Test Validate Basic ok", fields{
			PubKey:     pub,
			Chains:     chains,
			Value:      value,
			ServiceURL: surl,
		}, nil},
		{"Test Validate Basic bad value", fields{
			PubKey:     pub,
			Chains:     chains,
			Value:      sdk.NewInt(-1),
			ServiceURL: surl,
		}, ErrBadDelegationAmount(DefaultCodespace)},
		{"Test Validate Basic bad Chains", fields{
			PubKey:     pub,
			Chains:     []string{},
			Value:      value,
			ServiceURL: surl,
		}, ErrNoChains(DefaultCodespace)},
		{"Test Validate Basic bad chain in Chains", fields{
			PubKey:     pub,
			Chains:     []string{""},
			Value:      value,
			ServiceURL: surl,
		}, ErrInvalidNetworkIdentifier(DefaultCodespace, fmt.Errorf("net id is empty"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := LegacyMsgStake{
				PublicKey:  tt.fields.PubKey,
				Chains:     tt.fields.Chains,
				Value:      tt.fields.Value,
				ServiceUrl: tt.fields.ServiceURL,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyMsgUnjail_GetSignBytes(t *testing.T) {
	type fields struct {
		ValidatorAddr sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())

	mesg := LegacyMsgUnjail{
		ValidatorAddr: va,
	}

	encmesg, _ := ModuleCdc.MarshalJSON(&mesg)

	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{"Test GetSignBytes", fields{va}, encmesg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := LegacyMsgUnjail{
				ValidatorAddr: tt.fields.ValidatorAddr,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSignBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyMsgUnjail_GetSigners(t *testing.T) {
	type fields struct {
		ValidatorAddr sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())

	tests := []struct {
		name   string
		fields fields
		want   []sdk.Address
	}{
		{"Test GetSigners", fields{va}, []sdk.Address{sdk.Address(va)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := LegacyMsgUnjail{
				ValidatorAddr: tt.fields.ValidatorAddr,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyMsgUnjail_Route(t *testing.T) {
	type fields struct {
		ValidatorAddr sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Test Route", fields{va}, ModuleName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := LegacyMsgUnjail{
				ValidatorAddr: tt.fields.ValidatorAddr,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyMsgUnjail_Type(t *testing.T) {
	type fields struct {
		ValidatorAddr sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Test Type", fields{va}, MsgUnjailName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := LegacyMsgUnjail{
				ValidatorAddr: tt.fields.ValidatorAddr,
			}
			if got := msg.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyMsgUnjail_ValidateBasic(t *testing.T) {
	type fields struct {
		ValidatorAddr sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())

	tests := []struct {
		name   string
		fields fields
		want   sdk.Error
	}{
		{"Test ValidateBasic OK", fields{va}, nil},
		{"Test ValidateBasic bad address", fields{nil}, ErrNoValidatorFound(DefaultCodespace)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := LegacyMsgUnjail{
				ValidatorAddr: tt.fields.ValidatorAddr,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}
