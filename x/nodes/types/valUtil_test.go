package types

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-amino"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestValidators_JSON(t *testing.T) {
	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}

	testvalidator := Validator{
		Address:                 sdk.Address(pub.Address()),
		PublicKey:               pub,
		Jailed:                  false,
		Status:                  sdk.Staked,
		StakedTokens:            sdk.ZeroInt(),
		Chains:                  []string{"0001"},
		ServiceURL:              "https://www.google.com:443",
		UnstakingCompletionTime: time.Unix(0, 0).UTC(),
	}

	vals := Validators{testvalidator}

	result, _ := json.Marshal(vals)

	tests := []struct {
		name    string
		v       Validators
		wantOut []byte
		wantErr bool
	}{
		{"JSON Validators Test", vals, result, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := tt.v.JSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("JSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("JSON() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestValidators_String(t *testing.T) {

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}

	v := Validators{
		Validator{
			Address:                 sdk.Address(pub.Address()),
			PublicKey:               pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			Chains:                  []string{"0001"},
			ServiceURL:              "https://www.google.com:443",
			UnstakingCompletionTime: time.Unix(0, 0).UTC(),
		},
	}
	tests := []struct {
		name    string
		v       Validators
		wantOut string
	}{
		{"String Test", v, fmt.Sprintf("Address:\t\t%s\nPublic Key:\t\t%s\nJailed:\t\t\t%v\nStatus:\t\t\t%s\nTokens:\t\t\t%s\n"+
			"ServiceUrl:\t\t%s\nChains:\t\t\t%v\nUnstaking Completion Time:\t\t%v"+
			"\n----",
			sdk.Address(pub.Address()), pub.RawString(), false, sdk.Staked, sdk.ZeroInt(), "https://www.google.com:443", []string{"0001"}, time.Unix(0, 0).UTC(),
		)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if gotOut := tt.v.String(); gotOut != tt.wantOut {
				t.Errorf("String() = \n%v \nwant \b%v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestValidator_MarshalJSON(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		Chains                  []string
		ServiceURL              string
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}

	want, _ := amino.MarshalJSON(Validator{
		Address:                 sdk.Address(pub.Address()),
		PublicKey:               pub,
		Jailed:                  false,
		Status:                  sdk.Staked,
		Chains:                  []string{"0001"},
		ServiceURL:              "www.pokt.network",
		StakedTokens:            sdk.ZeroInt(),
		UnstakingCompletionTime: time.Time{},
	})

	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{"Marshall JSON Test", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			Chains:                  []string{"0001"},
			ServiceURL:              "www.pokt.network",
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, want, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				Chains:                  tt.fields.Chains,
				ServiceURL:              tt.fields.ServiceURL,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			got, err := v.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		Chains                  []string
		ServiceURL              string
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}

	marshal, _ := amino.MarshalJSON(Validator{
		Address:                 sdk.Address(pub.Address()),
		PublicKey:               pub,
		Jailed:                  false,
		Status:                  sdk.Staked,
		Chains:                  []string{"0001"},
		ServiceURL:              "www.pokt.network",
		StakedTokens:            sdk.ZeroInt(),
		UnstakingCompletionTime: time.Time{},
	})

	//amino.UnmarshalJSON(marshal,Validator{})

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Unmarshal JSON Test", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			Chains:                  []string{"0001"},
			ServiceURL:              "www.pokt.network",
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{data: marshal}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				Chains:                  tt.fields.Chains,
				ServiceURL:              tt.fields.ServiceURL,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if err := v.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateServiceURL(t *testing.T) {
	validURL := "https://foo.bar:8080"
	// missing prefix
	invalidURLNoPrefix := "foo.bar:8080"
	// wrong prefix
	invalidURLWrongPrefix := "ws://foo.bar:8080"
	// no port
	invalidURLNoPort := "ws://foo.bar"
	// bad port
	invalidURLBadPort := "ws://foo.bar:66666"
	// bad url
	invalidURLBad := "https://foobar:8080"
	assert.Nil(t, ValidateServiceURL(validURL))
	assert.NotNil(t, ValidateServiceURL(invalidURLNoPrefix), "invalid no prefix")
	assert.NotNil(t, ValidateServiceURL(invalidURLWrongPrefix), "invalid wrong prefix")
	assert.NotNil(t, ValidateServiceURL(invalidURLNoPort), "invalid no port")
	assert.NotNil(t, ValidateServiceURL(invalidURLBadPort), "invalid bad port")
	assert.NotNil(t, ValidateServiceURL(invalidURLBad), "invalid bad url")
}
