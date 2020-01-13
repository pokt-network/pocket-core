package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"reflect"
	"testing"
)

func TestKeeper_SetHooks(t *testing.T) {
	type fields struct {
		Keeper *Keeper
	}
	type args struct {
		sh types.POSHooks
	}

	_, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Keeper
	}{
		{"Test Set Hooks", fields{Keeper: &keeper}, args{
			sh: nil,
		}, &keeper},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.Keeper

			if got := k.SetHooks(tt.args.sh); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetHooks() = %v, want %v", got, tt.want)
			}
		})
	}
}
