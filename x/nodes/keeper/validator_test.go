package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"reflect"
	"testing"
)

func TestKeeper_GetValidators(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx         sdk.Context
		maxRetrieve uint16
	}

	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantValidators []types.Validator
	}{
		{"Test GetValidators 0", fields{keeper: keeper}, args{
			ctx:         context,
			maxRetrieve: 0,
		}, []types.Validator{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper

			if gotValidators := k.GetValidators(tt.args.ctx, tt.args.maxRetrieve); !reflect.DeepEqual(gotValidators, tt.wantValidators) {
				t.Errorf("GetValidators() = %v, want %v", gotValidators, tt.wantValidators)
			}
		})
	}
}
