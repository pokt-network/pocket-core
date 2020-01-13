package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"reflect"
	"testing"
	"time"
)

func TestKeeper_GetBlockHeight(t *testing.T) {
	type fields struct {
		Keeper Keeper
	}

	context, _, keeper := createTestInput(t, true)

	type args struct {
		ctx sdk.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{"Test GetBlockHeight", fields{Keeper: keeper}, args{ctx: context}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.Keeper
			if got := k.GetBlockHeight(tt.args.ctx); got != tt.want {
				t.Errorf("GetBlockHeight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_GetBlockTime(t *testing.T) {
	type fields struct {
		Keeper Keeper
	}
	type args struct {
		ctx sdk.Context
	}
	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		{"Test GetBlockTime default", fields{Keeper: keeper}, args{ctx: context}, time.Time{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.Keeper
			if got := k.GetBlockTime(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBlockTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_GetLatestBlockID(t *testing.T) {
	type fields struct {
		Keeper Keeper
	}
	type args struct {
		ctx sdk.Context
	}
	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   abci.BlockID
	}{
		{"GetLatestBlockID", fields{Keeper: keeper}, args{ctx: context}, abci.BlockID{
			Hash:        nil,
			PartsHeader: abci.PartSetHeader{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.Keeper
			if got := k.GetLatestBlockID(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLatestBlockID() = %v, want %v", got, tt.want)
			}
		})
	}
}
