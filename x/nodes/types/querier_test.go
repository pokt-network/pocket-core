package types

import (
	"github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"math/rand"
	"reflect"
	"testing"
)

func TestNewQuerySigningInfoParams(t *testing.T) {
	type args struct {
		consAddr types.ConsAddress
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])
	ca := types.ConsAddress(pub.Address())

	tests := []struct {
		name string
		args args
		want QuerySigningInfoParams
	}{
		{"default Test", args{ca}, QuerySigningInfoParams{ConsAddress: ca}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQuerySigningInfoParams(tt.args.consAddr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQuerySigningInfoParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQuerySigningInfosParams(t *testing.T) {
	type args struct {
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want QuerySigningInfosParams
	}{
		{"Default Test", args{limit: 1, page: 1}, QuerySigningInfosParams{Page: 1, Limit: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQuerySigningInfosParams(tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQuerySigningInfosParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryStakedValidatorsParams(t *testing.T) {
	type args struct {
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want QueryStakedValidatorsParams
	}{
		{"Default Test", args{page: 1, limit: 1}, QueryStakedValidatorsParams{Page: 1, Limit: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryStakedValidatorsParams(tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryStakedValidatorsParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryUnstakedValidatorsParams(t *testing.T) {
	type args struct {
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want QueryUnstakedValidatorsParams
	}{
		{"Default Test", args{page: 1, limit: 1}, QueryUnstakedValidatorsParams{Page: 1, Limit: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryUnstakedValidatorsParams(tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryUnstakedValidatorsParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryUnstakingValidatorsParams(t *testing.T) {
	type args struct {
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want QueryUnstakingValidatorsParams
	}{
		{"Default Test", args{page: 1, limit: 1}, QueryUnstakingValidatorsParams{Page: 1, Limit: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryUnstakingValidatorsParams(tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryUnstakingValidatorsParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryValidatorParams(t *testing.T) {
	type args struct {
		validatorAddr types.ValAddress
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])
	va := types.ValAddress(pub.Address())

	tests := []struct {
		name string
		args args
		want QueryValidatorParams
	}{
		{"default Test", args{va}, QueryValidatorParams{Address: va}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryValidatorParams(tt.args.validatorAddr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryValidatorParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryValidatorsParams(t *testing.T) {
	type args struct {
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want QueryValidatorsParams
	}{
		{"Default Test", args{page: 1, limit: 1}, QueryValidatorsParams{Page: 1, Limit: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryValidatorsParams(tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryValidatorsParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
