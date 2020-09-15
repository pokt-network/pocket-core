package nodes

import (
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/keeper"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"reflect"
	"testing"
)

func TestNewHandler(t *testing.T) {
	type args struct {
		k keeper.Keeper
	}
	tests := []struct {
		name string
		args args
		want sdk.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleMsgBeginUnstake(t *testing.T) {
	type args struct {
		ctx sdk.Context
		msg types.MsgBeginUnstake
		k   keeper.Keeper
	}
	tests := []struct {
		name string
		args args
		want sdk.Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleMsgBeginUnstake(tt.args.ctx, tt.args.msg, tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgBeginUnstake() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleMsgSend(t *testing.T) {
	type args struct {
		ctx sdk.Context
		msg types.MsgSend
		k   keeper.Keeper
	}
	tests := []struct {
		name string
		args args
		want sdk.Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleMsgSend(tt.args.ctx, tt.args.msg, tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgSend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleMsgUnjail(t *testing.T) {
	type args struct {
		ctx sdk.Context
		msg types.MsgUnjail
		k   keeper.Keeper
	}
	tests := []struct {
		name string
		args args
		want sdk.Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleMsgUnjail(tt.args.ctx, tt.args.msg, tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgUnjail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleStake(t *testing.T) {
	type args struct {
		ctx sdk.Context
		msg types.MsgNodeStake
		k   keeper.Keeper
	}
	tests := []struct {
		name string
		args args
		want sdk.Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleStake(tt.args.ctx, tt.args.msg, tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleStake() = %v, want %v", got, tt.want)
			}
		})
	}
}
