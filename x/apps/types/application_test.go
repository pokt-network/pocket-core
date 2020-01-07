package types

import (
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestNewApplication(t *testing.T) {
	type args struct {
		addr          sdk.ValAddress
		consPubKey    crypto.PubKey
		tokensToStake sdk.Int
		chains        []string
		serviceURL    string
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name string
		args args
		want Application
	}{
		{"defaultApplication", args{sdk.ValAddress(pub.Address()), pub, sdk.ZeroInt(), []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}, "google.com"},
			Application{
				Address:                 sdk.ValAddress(pub.Address()),
				ConsPubKey:              pub,
				Jailed:                  false,
				Status:                  sdk.Bonded,
				StakedTokens:            sdk.ZeroInt(),
				Chains:                  []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"},
				UnstakingCompletionTime: time.Unix(0, 0).UTC(), // zero out because status: bonded
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewApplication(tt.args.addr, tt.args.consPubKey, tt.args.chains, tt.args.tokensToStake); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApplication() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_AddStakedTokens(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	type args struct {
		tokens sdk.Int
	}
	tests := []struct {
		name   string
		panics bool
		fields fields
		args   args
		want  interface{}
	}{
		{
			"Default Add Token Test",
			false,
			fields{
				Address:                 sdk.ValAddress(pub.Address()),
				ConsPubKey:              pub,
				Jailed:                  false,
				Status:                  sdk.Bonded,
				StakedTokens:            sdk.ZeroInt(),
				UnstakingCompletionTime: time.Time{},
			},
			args{
					tokens: sdk.NewInt(100),
			},
			Application{
				Address:                 sdk.ValAddress(pub.Address()),
				ConsPubKey:              pub,
				Jailed:                  false,
				Status:                  sdk.Bonded,
				StakedTokens:            sdk.NewInt(100),
				UnstakingCompletionTime: time.Time{},
			},
		},
		{
			" panics Add negative amount",
			true,
			fields{
				Address:                 sdk.ValAddress(pub.Address()),
				ConsPubKey:              pub,
				Jailed:                  false,
				Status:                  sdk.Bonded,
				StakedTokens:            sdk.ZeroInt(),
				UnstakingCompletionTime: time.Time{},
			},
			args{
				tokens: sdk.NewInt(-1),
			},
			fmt.Sprint("should not happen: trying to add negative tokens -1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			switch tt.panics{
			case true:
				defer func(){
					err := recover()
					if !reflect.DeepEqual(fmt.Sprintf("%v", err), tt.want) {
						t.Errorf("AddStakedTokens() = %v, want %v", err, tt.want)
					}
				}()
				_ = v.AddStakedTokens(tt.args.tokens)
			default:
				if got := v.AddStakedTokens(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("AddStakedTokens() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestApplication_ConsAddress(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   sdk.ConsAddress
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, sdk.ConsAddress(pub.Address())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.ConsAddress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConsAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_ConsensusPower(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{"Default Test / 0 power", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, 0},
		{"Default Test / 1 power", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.NewInt(1000000),
			UnstakingCompletionTime: time.Time{},
		}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.ConsensusPower(); got != tt.want {
				t.Errorf("ConsensusPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_Equals(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	type args struct {
		v2 Application
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"Default Test Equal", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{Application{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}}, true},
		{"Default Test Not Equal", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{Application{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.Equals(tt.args.v2); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_GetAddress(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   sdk.ValAddress
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, sdk.ValAddress(pub.Address())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetAddress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_GetConsAddr(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   sdk.ConsAddress
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, sdk.ConsAddress(pub.Address())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetConsAddr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConsAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_GetConsPubKey(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   crypto.PubKey
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, pub},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetConsPubKey(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConsPubKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_GetConsensusPower(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetConsensusPower(); got != tt.want {
				t.Errorf("GetConsensusPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_GetStatus(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   sdk.BondStatus
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, sdk.Bonded},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetStatus(); got != tt.want {
				t.Errorf("GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_GetTokens(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   sdk.Int
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, sdk.ZeroInt()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetTokens(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_IsJailed(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.IsJailed(); got != tt.want {
				t.Errorf("IsJailed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_IsStaked(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Default Test / bonded true", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, true},
		{"Default Test / Unbonding false", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonding,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
		{"Default Test / Unbonded false", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.IsStaked(); got != tt.want {
				t.Errorf("IsStaked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_IsUnstaked(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Default Test / bonded false", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
		{"Default Test / Unbonding false", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonding,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
		{"Default Test / Unbonded true", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.IsUnstaked(); got != tt.want {
				t.Errorf("IsUnstaked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_IsUnstaking(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Default Test / bonded false", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
		{"Default Test / Unbonding true", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonding,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, true},
		{"Default Test / Unbonded false", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.IsUnstaking(); got != tt.want {
				t.Errorf("IsUnstaking() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_PotentialConsensusPower(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{"Default Test / potential power 0", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.PotentialConsensusPower(); got != tt.want {
				t.Errorf("PotentialConsensusPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_RemoveStakedTokens(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	type args struct {
		tokens sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Application
	}{
		{"Remove 0 tokens having 0 tokens ", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{tokens: sdk.ZeroInt()}, Application{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}},
		{"Remove 99 tokens having 100 tokens ", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.NewInt(100),
			UnstakingCompletionTime: time.Time{},
		}, args{tokens: sdk.NewInt(99)}, Application{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.OneInt(),
			UnstakingCompletionTime: time.Time{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.RemoveStakedTokens(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveStakedTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_UpdateStatus(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	type args struct {
		newStatus sdk.BondStatus
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Application
	}{
		{"Test Bonded -> Unbonding", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{newStatus: sdk.Unbonding}, Application{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonding,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}},
		{"Test Unbonding -> Unbonded", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonding,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{newStatus: sdk.Unbonded}, Application{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}},
		{"Test Unbonded -> Bonded", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{newStatus: sdk.Bonded}, Application{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Application{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.UpdateStatus(tt.args.newStatus); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_GetChains(t *testing.T) {
	type args struct {
		addr          sdk.ValAddress
		consPubKey    crypto.PubKey
		tokensToStake sdk.Int
		chains        []string
		serviceURL    string
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"defaultApplication",
			args{sdk.ValAddress(pub.Address()), pub, sdk.ZeroInt(), []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}, "google.com"},
			[]string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApplication(tt.args.addr, tt.args.consPubKey, tt.args.chains, tt.args.tokensToStake)
			if got := app.GetChains(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_GetMaxRelays(t *testing.T) {
	type args struct {
		addr          sdk.ValAddress
		consPubKey    crypto.PubKey
		tokensToStake sdk.Int
		chains        []string
		serviceURL    string
		maxRelays     sdk.Int
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name string
		args args
		want sdk.Int
	}{
		{
			"defaultApplication",
			args{sdk.ValAddress(pub.Address()), pub, sdk.ZeroInt(), []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}, "google.com", sdk.NewInt(1)},
			sdk.NewInt(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := Application{
				Address:    tt.args.addr,
				ConsPubKey: tt.args.consPubKey,
				Chains:     tt.args.chains,
				MaxRelays:  tt.args.maxRelays,
			}
			if got := app.GetMaxRelays(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMaxRelays() = %v, want %v", got, tt.want)
			}
		})
	}
}
