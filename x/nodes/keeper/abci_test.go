package keeper

import (
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
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

func TestKeeper_ConvertValidatorsState(t *testing.T) {
	v1 := getValidator()
	v1.OutputAddress = nil
	v2 := getValidator()
	v2.OutputAddress = nil
	v3 := getValidator()
	v3.OutputAddress = nil
	v := []types.Validator{v1, v2, v3}
	lv1 := v1.ToLegacy()
	lv2 := v2.ToLegacy()
	lv3 := v3.ToLegacy()
	lvs := []types.LegacyValidator{lv1, lv2, lv3}
	ctx, _, k := createTestInput(t, true)
	// manually set the validators
	store := ctx.KVStore(k.storeKey)
	for i, lv := range lvs {
		bz, err := k.Cdc.MarshalBinaryLengthPrefixed(&lv, ctx.BlockHeight())
		if err != nil {
			ctx.Logger().Error("could not marshal validator: " + err.Error())
		}
		err = store.Set(types.KeyForValByAllVals(lv.Address), bz)
		// convert the state, can be commented out as not needed,
		//intentionally left here as a reminder that state convert for this was planned but not needed and can be removed next version
		//k.ConvertValidatorsState(ctx)
		// manually get validators using new structure
		value, err := store.Get(types.KeyForValByAllVals(lv.Address))
		assert.Nil(t, err)
		var val types.Validator
		err = k.Cdc.UnmarshalBinaryLengthPrefixed(value, &val, ctx.BlockHeight())
		assert.Nil(t, err)
		assert.Equal(t, v[i], val)
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
