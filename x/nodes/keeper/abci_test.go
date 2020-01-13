package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestBeginBlocker(t *testing.T) {
	type args struct {
		ctx sdk.Context
		req abci.RequestBeginBlock
		k   Keeper
	}
	context, _, keeper := createTestInput(t, true)
	tests := []struct {
		name string
		args args
	}{
		{"Test BeginBlocker", args{
			ctx: context,
			req: abci.RequestBeginBlock{
				Hash:                 []byte{0x51, 0x51, 0x51},
				Header:               abci.Header{},
				LastCommitInfo:       abci.LastCommitInfo{},
				ByzantineValidators:  nil,
				XXX_NoUnkeyedLiteral: struct{}{},
				XXX_unrecognized:     nil,
				XXX_sizecache:        0,
			},
			k: keeper,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BeginBlocker(tt.args.ctx, tt.args.req, tt.args.k)
		})
	}
}

func TestEndBlocker(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
	}
	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name string
		args args
		want []abci.ValidatorUpdate
	}{
		{"Test EndBlocker", args{
			ctx: context,
			k:   keeper,
		}, []abci.ValidatorUpdate{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndBlocker(tt.args.ctx, tt.args.k); !assert.True(t, len(got) == len(tt.want)) {
				t.Errorf("EndBlocker() = %v, want %v", got, tt.want)
			}
		})
	}
}
