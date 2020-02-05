package types

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/go-amino"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestValidators_JSON(t *testing.T) {
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	testvalidator := Validator{
		Address:                 sdk.Address(pub.Address()),
		PublicKey:               pub,
		Jailed:                  false,
		Status:                  sdk.Staked,
		StakedTokens:            sdk.ZeroInt(),
		Chains:                  []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"},
		ServiceURL:              "google.com",
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
	rand.Read(pub[:])

	v := Validators{
		Validator{
			Address:                 sdk.Address(pub.Address()),
			PublicKey:               pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			Chains:                  []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"},
			ServiceURL:              "google.com",
			UnstakingCompletionTime: time.Unix(0, 0).UTC(),
		},
	}
	tests := []struct {
		name    string
		v       Validators
		wantOut string
	}{
		{"String Test", v, fmt.Sprintf(`Validator
  Address:           		  %s
  Validator Cons Pubkey:      %s
  Jailed:                     %v
  Status:                     %s
  Tokens:               	  %s
  ServiceURL:                 %s
  Chains:                     %v
  Unstaking Completion Time:  %v`,
			sdk.Address(pub.Address()), pub.RawString(), false, sdk.Staked, sdk.ZeroInt(), "google.com", []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}, time.Unix(0, 0).UTC(),
		)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := tt.v.String(); gotOut != tt.wantOut {
				t.Errorf("String() = %v, want %v", gotOut, tt.wantOut)
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
	rand.Read(pub[:])

	want, _ := amino.MarshalJSON(Validator{
		Address:                 sdk.Address(pub.Address()),
		PublicKey:               pub,
		Jailed:                  false,
		Status:                  sdk.Staked,
		Chains:                  []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"},
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
			Chains:                  []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"},
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
	rand.Read(pub[:])

	marshal, _ := amino.MarshalJSON(Validator{
		Address:                 sdk.Address(pub.Address()),
		PublicKey:               pub,
		Jailed:                  false,
		Status:                  sdk.Staked,
		Chains:                  []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"},
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
			Chains:                  []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"},
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
