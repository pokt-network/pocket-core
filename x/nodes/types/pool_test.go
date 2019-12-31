package types

import (
	sdk "github.com/pokt-network/posmint/types"
	"reflect"
	"testing"
)

func TestDAOPool_String(t *testing.T) {
	tests := []struct {
		name string
		dao  DAOPool
		want string
	}{
		{"EmptyPool", DAOPool(NewPool(sdk.ZeroInt())), "Tokens:Tokens In DAO Tokens: 0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dao.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPool(t *testing.T) {
	type args struct {
		tokens sdk.Int
	}
	tests := []struct {
		name string
		args args
		want Pool
	}{
		{"EmptyPool", args{tokens: sdk.ZeroInt()}, Pool{Tokens: sdk.ZeroInt()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPool(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStakingPool_String(t *testing.T) {
	tests := []struct {
		name string
		bp   StakingPool
		want string
	}{
		{"EmptyPool", StakingPool(NewPool(sdk.ZeroInt())), "Staked Tokens: 0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bp.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
