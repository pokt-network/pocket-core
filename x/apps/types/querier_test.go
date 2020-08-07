package types

import (
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	"math/rand"
	"reflect"
	"testing"
)

func TestNewQueryApplicationParams(t *testing.T) {
	type args struct {
		applicationAddr types.Address
	}
	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	va := types.Address(pub.Address())

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
