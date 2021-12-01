package types

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	"reflect"
	"testing"
)

func TestPool_NewPool(t *testing.T) {
	tests := []struct {
		name string
		args sdk.BigInt
		want sdk.BigInt
	}{
		{
			"returns pool with tokens",
			sdk.NewInt(1),
			sdk.NewInt(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPool(tt.args); !got.Tokens.Equal(tt.want) {
				t.Errorf("NewPool.Tokens = %v, want %v", got.Tokens, tt.want)
			}
		})
	}
}
func TestPool_String(t *testing.T) {
	tests := []struct {
		name string
		args StakingPool
		want string
	}{
		{
			"returns pool with tokens",
			StakingPool{sdk.NewInt(10)},
			fmt.Sprintf(`Staked Tokens:      %s`, sdk.NewInt(10)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StakingPool.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
