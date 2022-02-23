package types

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
)

func TestMsgBeginUnstake_GetSignBytes(t *testing.T) {
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
		Signer:  va,
	}

	encodedmsg, _ := ModuleCdc.MarshalJSON(&mesg)
	encodedmsg = sdk.MustSortJSON(encodedmsg)
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{"Test GetSignBytes", fields{va}, encodedmsg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgBeginUnstake{
				Address: tt.fields.Address,
				Signer:  tt.fields.Address,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMsgBeginUnstake_GetSigners(t *testing.T) {
	type fields struct {
		Address sdk.Address
		Signer  sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	var pub2 crypto.Ed25519PublicKey
	_, err = rand.Read(pub2[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())
	va2 := sdk.Address(pub2.Address())
	mesg := MsgBeginUnstake{
		Address: va,
		Signer:  va2,
	}

	tests := []struct {
		name   string
		fields fields
		want   []sdk.Address
	}{
		{"Test GetSigners", fields{va, va2}, []sdk.Address{mesg.Signer, mesg.Address}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgBeginUnstake{
				Address: tt.fields.Address,
				Signer:  tt.fields.Signer,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgBeginUnstake_Route(t *testing.T) {
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
			msg := MsgBeginUnstake{
				Address: tt.fields.Address,
				Signer:  tt.fields.Address,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgBeginUnstake_Type(t *testing.T) {
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
			msg := MsgBeginUnstake{
				Address: tt.fields.Address,
				Signer:  tt.fields.Address,
			}
			if got := msg.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgBeginUnstake_ValidateBasic(t *testing.T) {
	type fields struct {
		Address sdk.Address
		Signer  sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	var pub2 crypto.Ed25519PublicKey
	_, err = rand.Read(pub2[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())
	va2 := sdk.Address(pub2.Address())
	tests := []struct {
		name   string
		fields fields
		want   sdk.Error
	}{
		{"Test Validate Basic error", fields{nil, va2}, sdk.NewError(codespace, CodeInvalidInput, "validator address is nil")},
		{"Test Validate Basic error", fields{va, nil}, sdk.NewError(codespace, CodeNilSigner, "signer address is nil")},
		{"Test Validate Basic pass", fields{va, va2}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgBeginUnstake{
				Address: tt.fields.Address,
				Signer:  tt.fields.Signer,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgSend_GetSignBytes(t *testing.T) {
	type fields struct {
		FromAddress sdk.Address
		ToAddress   sdk.Address
		Amount      sdk.BigInt
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())
	_, err = rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va2 := sdk.Address(pub.Address())

	mesg := MsgSend{
		FromAddress: va,
		ToAddress:   va2,
		Amount:      sdk.OneInt(),
	}

	encmesg, _ := ModuleCdc.MarshalJSON(&mesg)
	encmesg = sdk.MustSortJSON(encmesg)

	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{"Test GetSignBytes", fields{
			FromAddress: va,
			ToAddress:   va2,
			Amount:      sdk.OneInt(),
		}, encmesg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgSend{
				FromAddress: tt.fields.FromAddress,
				ToAddress:   tt.fields.ToAddress,
				Amount:      tt.fields.Amount,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSignBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgSend_GetSigners(t *testing.T) {
	type fields struct {
		FromAddress sdk.Address
		ToAddress   sdk.Address
		Amount      sdk.BigInt
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())
	_, err = rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va2 := sdk.Address(pub.Address())

	tests := []struct {
		name   string
		fields fields
		want   []sdk.Address
	}{
		{"Test GetSigners", fields{
			FromAddress: va,
			ToAddress:   va2,
			Amount:      sdk.OneInt(),
		}, []sdk.Address{sdk.Address(va)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgSend{
				FromAddress: tt.fields.FromAddress,
				ToAddress:   tt.fields.ToAddress,
				Amount:      tt.fields.Amount,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgSend_Route(t *testing.T) {
	type fields struct {
		FromAddress sdk.Address
		ToAddress   sdk.Address
		Amount      sdk.BigInt
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())
	_, err = rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va2 := sdk.Address(pub.Address())

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Test Route", fields{
			FromAddress: va,
			ToAddress:   va2,
			Amount:      sdk.OneInt(),
		}, ModuleName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgSend{
				FromAddress: tt.fields.FromAddress,
				ToAddress:   tt.fields.ToAddress,
				Amount:      tt.fields.Amount,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgSend_Type(t *testing.T) {
	type fields struct {
		FromAddress sdk.Address
		ToAddress   sdk.Address
		Amount      sdk.BigInt
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())
	_, err = rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va2 := sdk.Address(pub.Address())

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Test Type", fields{
			FromAddress: va,
			ToAddress:   va2,
			Amount:      sdk.OneInt(),
		}, "send"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgSend{
				FromAddress: tt.fields.FromAddress,
				ToAddress:   tt.fields.ToAddress,
				Amount:      tt.fields.Amount,
			}
			if got := msg.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgSend_ValidateBasic(t *testing.T) {
	type fields struct {
		FromAddress sdk.Address
		ToAddress   sdk.Address
		Amount      sdk.BigInt
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())
	_, err = rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va2 := sdk.Address(pub.Address())

	tests := []struct {
		name   string
		fields fields
		want   sdk.Error
	}{
		{"Test ValidateBasic ok", fields{
			FromAddress: va,
			ToAddress:   va2,
			Amount:      sdk.OneInt(),
		}, nil},
		{"Test ValidateBasic empty from", fields{
			FromAddress: nil,
			ToAddress:   va2,
			Amount:      sdk.OneInt(),
		}, ErrNoValidatorFound(DefaultCodespace)},
		{"Test ValidateBasic empty to", fields{
			FromAddress: va,
			ToAddress:   nil,
			Amount:      sdk.OneInt(),
		}, ErrNoValidatorFound(DefaultCodespace)},
		{"Test ValidateBasic bad amount", fields{
			FromAddress: va,
			ToAddress:   va2,
			Amount:      sdk.NewInt(-1),
		}, ErrBadSendAmount(DefaultCodespace)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgSend{
				FromAddress: tt.fields.FromAddress,
				ToAddress:   tt.fields.ToAddress,
				Amount:      tt.fields.Amount,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgStake_GetSigners(t *testing.T) {
	type fields struct {
		Address       sdk.Address
		PubKey        crypto.PublicKey
		Chains        []string
		Value         sdk.BigInt
		ServiceURL    string
		OutputAddress sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	var pub2 crypto.Ed25519PublicKey
	_, err = rand.Read(pub2[:])
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
			PubKey:        pub,
			Chains:        chains,
			Value:         value,
			ServiceURL:    surl,
			OutputAddress: sdk.Address(pub2.Address()),
		}, []sdk.Address{sdk.Address(pub.Address()), sdk.Address(pub2.Address())}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgStake{
				PublicKey:  tt.fields.PubKey,
				Chains:     tt.fields.Chains,
				Value:      tt.fields.Value,
				ServiceUrl: tt.fields.ServiceURL,
				Output:     tt.fields.OutputAddress,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgStake_Route(t *testing.T) {
	type fields struct {
		Address       sdk.Address
		PubKey        crypto.PublicKey
		Chains        []string
		Value         sdk.BigInt
		ServiceURL    string
		OutputAddress sdk.Address
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
			PubKey:        pub,
			Chains:        chains,
			Value:         value,
			ServiceURL:    surl,
			OutputAddress: sdk.Address(pub.Address()),
		}, RouterKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgStake{
				PublicKey:  tt.fields.PubKey,
				Chains:     tt.fields.Chains,
				Value:      tt.fields.Value,
				ServiceUrl: tt.fields.ServiceURL,
				Output:     tt.fields.OutputAddress,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgStake_Type(t *testing.T) {
	type fields struct {
		Address       sdk.Address
		PubKey        crypto.PublicKey
		Chains        []string
		Value         sdk.BigInt
		ServiceURL    string
		OutputAddress sdk.Address
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
			PubKey:        pub,
			Chains:        chains,
			Value:         value,
			ServiceURL:    surl,
			OutputAddress: sdk.Address(pub.Address()),
		}, "stake_validator"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgStake{
				PublicKey:  tt.fields.PubKey,
				Chains:     tt.fields.Chains,
				Value:      tt.fields.Value,
				ServiceUrl: tt.fields.ServiceURL,
				Output:     tt.fields.OutputAddress,
			}
			if got := msg.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgStake_ValidateBasic(t *testing.T) {
	type fields struct {
		Address    sdk.Address
		PubKey     crypto.PublicKey
		Chains     []string
		Value      sdk.BigInt
		ServiceURL string
		Output     sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	var pub2 crypto.Ed25519PublicKey
	_, err = rand.Read(pub2[:])
	if err != nil {
		_ = err
	}
	outputAddr := sdk.Address(pub2.Address())
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
			Output:     outputAddr,
		}, nil},
		{"Test Validate Basic bad value", fields{
			PubKey:     pub,
			Chains:     chains,
			Value:      sdk.NewInt(-1),
			ServiceURL: surl,
			Output:     outputAddr,
		}, ErrBadDelegationAmount(DefaultCodespace)},
		{"Test Validate Basic bad Chains", fields{
			PubKey:     pub,
			Chains:     []string{},
			Value:      value,
			ServiceURL: surl,
			Output:     outputAddr,
		}, ErrNoChains(DefaultCodespace)},
		{"Test Validate Basic bad chain in Chains", fields{
			PubKey:     pub,
			Chains:     []string{""},
			Value:      value,
			ServiceURL: surl,
			Output:     outputAddr,
		}, ErrInvalidNetworkIdentifier(DefaultCodespace, fmt.Errorf("net id is empty"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgStake{
				PublicKey:  tt.fields.PubKey,
				Chains:     tt.fields.Chains,
				Value:      tt.fields.Value,
				ServiceUrl: tt.fields.ServiceURL,
				Output:     tt.fields.Output,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnjail_GetSignBytes(t *testing.T) {
	type fields struct {
		ValidatorAddr sdk.Address
		OutputAddr    sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())
	var pub2 crypto.Ed25519PublicKey
	_, err = rand.Read(pub2[:])
	if err != nil {
		_ = err
	}
	va2 := sdk.Address(pub2.Address())

	mesg := MsgUnjail{
		ValidatorAddr: va,
		Signer:        va2,
	}

	encmesg, _ := ModuleCdc.MarshalJSON(&mesg)

	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{"Test GetSignBytes", fields{va, va2}, encmesg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUnjail{
				ValidatorAddr: tt.fields.ValidatorAddr,
				Signer:        va2,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSignBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnjail_GetSigners(t *testing.T) {
	type fields struct {
		ValidatorAddr sdk.Address
		OutputAddr    sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())
	var pub2 crypto.Ed25519PublicKey
	_, err = rand.Read(pub2[:])
	if err != nil {
		_ = err
	}
	va2 := sdk.Address(pub2.Address())
	tests := []struct {
		name   string
		fields fields
		want   []sdk.Address
	}{
		{"Test GetSigners", fields{va, va2}, []sdk.Address{va2, va}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUnjail{
				ValidatorAddr: tt.fields.ValidatorAddr,
				Signer:        va2,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnjail_Route(t *testing.T) {
	type fields struct {
		ValidatorAddr sdk.Address
		Output        sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())
	var pub2 crypto.Ed25519PublicKey
	_, err = rand.Read(pub2[:])
	if err != nil {
		_ = err
	}
	va2 := sdk.Address(pub2.Address())
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Test Route", fields{va, va2}, ModuleName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUnjail{
				ValidatorAddr: tt.fields.ValidatorAddr,
				Signer:        tt.fields.Output,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnjail_Type(t *testing.T) {
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
			msg := MsgUnjail{
				ValidatorAddr: tt.fields.ValidatorAddr,
				Signer:        tt.fields.ValidatorAddr,
			}
			if got := msg.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnjail_ValidateBasic(t *testing.T) {
	type fields struct {
		ValidatorAddr sdk.Address
		Output        sdk.Address
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := sdk.Address(pub.Address())
	var pub2 crypto.Ed25519PublicKey
	_, err = rand.Read(pub2[:])
	if err != nil {
		_ = err
	}
	va2 := sdk.Address(pub2.Address())
	tests := []struct {
		name   string
		fields fields
		want   sdk.Error
	}{
		{"Test ValidateBasic OK", fields{va, va2}, nil},
		{"Test ValidateBasic bad address", fields{nil, va2}, ErrNoValidatorFound(DefaultCodespace)},
		{"Test ValidateBasic bad output", fields{va, nil}, ErrNilSignerAddr(DefaultCodespace)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUnjail{
				ValidatorAddr: tt.fields.ValidatorAddr,
				Signer:        tt.fields.Output,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}
