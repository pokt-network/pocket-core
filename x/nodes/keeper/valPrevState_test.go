package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"reflect"
	"testing"
)

func TestKeeper_DeletePrevStateValPower(t *testing.T) {
	type fields struct {
		keeper Keeper
	}

	type args struct {
		ctx  sdk.Context
		addr sdk.ValAddress
	}

	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test DeletePrevStateValPower", fields{keeper: keeper}, args{
			ctx:  context,
			addr: getRandomValidatorAddress(),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper

			k.DeletePrevStateValPower(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_IterateAndExecuteOverPrevStateVals(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx sdk.Context
		fn  func(index int64, validator exported.ValidatorI) (stop bool)
	}

	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test IterateAndExecuteOverPrevStateVals", fields{keeper: keeper}, args{
			ctx: context,
			fn: func(index int64, validator exported.ValidatorI) (stop bool) {
				return true
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			k.IterateAndExecuteOverPrevStateVals(tt.args.ctx, tt.args.fn)
		})
	}
}

func TestKeeper_IterateAndExecuteOverPrevStateValsByPower(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx     sdk.Context
		handler func(address sdk.ValAddress, power int64) (stop bool)
	}
	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test IterateAndExecuteOverPrevStateValsByPower", fields{keeper: keeper}, args{
			ctx: context,
			handler: func(address sdk.ValAddress, power int64) (stop bool) {
				return true
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			k.IterateAndExecuteOverPrevStateValsByPower(tt.args.ctx, tt.args.handler)
		})
	}
}

func TestKeeper_PrevStateValidatorPower(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx  sdk.Context
		addr sdk.ValAddress
	}

	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantPower int64
	}{
		{"Test PrevStateValidatorPower", fields{keeper: keeper}, args{
			ctx:  context,
			addr: getRandomValidatorAddress(),
		}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			if gotPower := k.PrevStateValidatorPower(tt.args.ctx, tt.args.addr); gotPower != tt.wantPower {
				t.Errorf("PrevStateValidatorPower() = %v, want %v", gotPower, tt.wantPower)
			}
		})
	}
}

func TestKeeper_PrevStateValidatorsPower(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx sdk.Context
	}

	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantPower sdk.Int
	}{
		{"Test PrevStateValidatorsPower", fields{keeper: keeper}, args{context}, sdk.ZeroInt()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			if gotPower := k.PrevStateValidatorsPower(tt.args.ctx); !reflect.DeepEqual(gotPower, tt.wantPower) {
				t.Errorf("PrevStateValidatorsPower() = %v, want %v", gotPower, tt.wantPower)
			}
		})
	}
}

func TestKeeper_SetPrevStateValPower(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx   sdk.Context
		addr  sdk.ValAddress
		power int64
	}

	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test SetPrevStateValPower", fields{keeper: keeper}, args{
			ctx:   context,
			addr:  getRandomValidatorAddress(),
			power: 0,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			k.SetPrevStateValPower(tt.args.ctx, tt.args.addr, tt.args.power)
		})
	}
}

func TestKeeper_SetPrevStateValidatorsPower(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx   sdk.Context
		power sdk.Int
	}
	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test SetPrevStateValidatorsPower", fields{keeper: keeper}, args{
			ctx:   context,
			power: sdk.ZeroInt(),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			k.SetPrevStateValidatorsPower(tt.args.ctx, tt.args.power)
		})
	}
}

func TestKeeper_getValsFromPrevState(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx sdk.Context
	}

	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantValidators []types.Validator
	}{
		{"Test getValsFromPrevState", fields{keeper: keeper}, args{ctx: context},
			[]types.Validator{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			if gotValidators := k.getValsFromPrevState(tt.args.ctx); !reflect.DeepEqual(gotValidators, tt.wantValidators) {
				t.Errorf("getValsFromPrevState() = %v, want %v", gotValidators, tt.wantValidators)
			}
		})
	}
}

func TestKeeper_prevStateValidatorsIterator(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx sdk.Context
	}
	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name         string
		fields       fields
		args         args
		wantIterator sdk.Iterator
	}{
		{"Test prevStateValidatorsIterator", fields{keeper: keeper}, args{ctx: context},
			sdk.KVStorePrefixIterator(context.KVStore(keeper.storeKey), types.PrevStateValidatorsPowerKey),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			gotIterator := k.prevStateValidatorsIterator(tt.args.ctx)
			gotIterator.Valid()
		})
	}
}
