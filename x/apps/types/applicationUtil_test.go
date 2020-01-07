package types

import (
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"math/rand"
	"reflect"
	"testing"
	"time"
)

var application Application
var moduleCdc *codec.Codec

func init() {
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	moduleCdc = codec.New()
	RegisterCodec(moduleCdc)
	codec.RegisterCrypto(moduleCdc)
	moduleCdc.Seal()

	application = Application{
		Address:                 sdk.ValAddress(pub.Address()),
		ConsPubKey:              pub,
		Jailed:                  false,
		Status:                  sdk.Bonded,
		StakedTokens:            sdk.NewInt(100),
		MaxRelays:               sdk.NewInt(1000),
		UnstakingCompletionTime: time.Time{},
	}
}

func TestApplicationUtil_MarshalJSON(t *testing.T) {
	type args struct {
		application Application
		codec       *codec.Codec
	}
	hexApp := hexApplication{
		Address:                 application.Address,
		ConsPubKey:              sdk.HexConsPub(application.ConsPubKey),
		Jailed:                  application.Jailed,
		Status:                  application.Status,
		StakedTokens:            application.StakedTokens,
		UnstakingCompletionTime: application.UnstakingCompletionTime,
		MaxRelays:               application.MaxRelays,
	}
	bz, _ := codec.Cdc.MarshalJSON(hexApp)

	tests := []struct {
		name string
		args
		want []byte
	}{
		{
			name: "marshals application",
			args: args{application: application, codec: moduleCdc},
			want: bz,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.args.application.MarshalJSON(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MmashalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestApplicationUtil_UnmarshalJSON(t *testing.T) {
	type args struct {
		application Application
	}
	tests := []struct {
		name string
		args
		want Application
	}{
		{
			name: "marshals application",
			args: args{application: application},
			want: application,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshaled, err := tt.args.application.MarshalJSON()
			if err != nil {
				t.Error("Cannot marshal application")
			}
			if err = tt.args.application.UnmarshalJSON(marshaled); err != nil {
				t.Errorf("Unmarshal(): returns %v but want %v", err, tt.want)
			}
			// NOTE CANNOT PERFORM DEEP EQUAL
			// Unmarshalling causes StakedTokens & MaxRelays to be
			//  assigned a new memory address overwriting the previous reference to application
			// separate them and assert absolute value rather than deep equal

			gotStaked := tt.args.application.StakedTokens
			wantStaked := tt.want.StakedTokens
			gotRelays := tt.args.application.StakedTokens
			wantRelays := tt.want.StakedTokens

			tt.args.application.StakedTokens = tt.want.StakedTokens
			tt.args.application.MaxRelays = tt.want.MaxRelays

			if !reflect.DeepEqual(tt.args.application, tt.want) {
				t.Errorf("got %v but want %v", tt.args.application, tt.want)
			}
			if !gotStaked.Equal(wantStaked) {
				t.Errorf("got %v but want %v", gotStaked, wantStaked)
			}
			if !gotRelays.Equal(wantRelays) {
				t.Errorf("got %v but want %v", gotRelays, wantRelays)
			}
		})
	}
}

func TestApplicationUtil_MustMarshalApplication(t *testing.T) {
	type args struct {
		application Application
		codec       *codec.Codec
	}
	tests := []struct {
		name string
		args
		want []byte
	}{
		{
			name: "marshals application",
			args: args{application: application, codec: moduleCdc},
			want: moduleCdc.MustMarshalBinaryLengthPrefixed(application),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustMarshalApplication(tt.args.codec, tt.args.application); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustMarshalApplication()= returns %v but want %v", got, tt.want)
			}
		})
	}
}

func TestApplicationUtil_MustUnMarshalApplication(t *testing.T) {
	type args struct {
		application Application
		codec       *codec.Codec
	}
	tests := []struct {
		name string
		args
		want Application
	}{
		{
			name: "can unmarshal application",
			args: args{application: application, codec: moduleCdc},
			want: application,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bz := MustMarshalApplication(tt.args.codec, tt.args.application)
			if unmarshaledApp := MustUnmarshalApplication(tt.args.codec, bz); !reflect.DeepEqual(unmarshaledApp, tt.want) {
				t.Errorf("got %v but want %v", unmarshaledApp, tt.want)
			}
		})
	}
}

func TestApplicationUtil_UnMarshalApplication(t *testing.T) {
	type args struct {
		application Application
		codec       *codec.Codec
	}
	tests := []struct {
		name string
		args
		want Application
	}{
		{
			name: "can unmarshal application",
			args: args{application: application, codec: moduleCdc},
			want: application,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bz := MustMarshalApplication(tt.args.codec, tt.args.application)
			unmarshaledApp, err := UnmarshalApplication(tt.args.codec, bz)
			if err != nil {
				t.Error("could not unmarshal app")
			}

			if !reflect.DeepEqual(unmarshaledApp, tt.want) {
				t.Errorf("got %v but want %v", unmarshaledApp, unmarshaledApp)
			}
		})
	}
}
