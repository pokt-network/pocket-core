package types

import (
	"github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"math/rand"
	"reflect"
	"testing"
)

func TestNewQueryStakedApplicationsParams(t *testing.T) {
	type args struct {
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want QueryStakedApplicationsParams
	}{
		{"Default Test", args{page: 1, limit: 1}, QueryStakedApplicationsParams{Page: 1, Limit: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryStakedApplicationsParams(tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryStakedApplicationsParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryUnstakedApplicationsParams(t *testing.T) {
	type args struct {
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want QueryUnstakedApplicationsParams
	}{
		{"Default Test", args{page: 1, limit: 1}, QueryUnstakedApplicationsParams{Page: 1, Limit: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryUnstakedApplicationsParams(tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryUnstakedApplicationsParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryUnstakingApplicationsParams(t *testing.T) {
	type args struct {
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want QueryUnstakingApplicationsParams
	}{
		{"Default Test", args{page: 1, limit: 1}, QueryUnstakingApplicationsParams{Page: 1, Limit: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryUnstakingApplicationsParams(tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryUnstakingApplicationsParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryApplicationParams(t *testing.T) {
	type args struct {
		applicationAddr types.ValAddress
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])
	va := types.ValAddress(pub.Address())

	tests := []struct {
		name string
		args args
		want QueryAppParams
	}{
		{"default Test", args{va}, QueryAppParams{Address: va}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryAppParams(tt.args.applicationAddr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryApplicationParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryApplicationsParams(t *testing.T) {
	type args struct {
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want QueryAppsParams
	}{
		{"Default Test", args{page: 1, limit: 1}, QueryAppsParams{Page: 1, Limit: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryApplicationsParams(tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryApplicationsParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
