package keeper

//
//import (
//	"github.com/pokt-network/posmint/codec"
//	sdk "github.com/pokt-network/posmint/types"
//	"github.com/pokt-network/posmint/x/crisis"
//	"github.com/pokt-network/posmint/x/params"
//	"reflect"
//	"testing"
//
//)
//
//func TestModuleAccountInvariants(t *testing.T) {
//	type args struct {
//		k Keeper
//	}
//
//	_, _, keeper := createTestInput(t, true)
//
//	tests := []struct {
//		name string
//		args args
//		want sdk.Invariant
//	}{
//		{"Test ModuleAccountInvariants",args{k:keeper},
//			func(context sdk.Context) (string, bool){
//				return "", true
//			}},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := ModuleAccountInvariants(tt.args.k); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("ModuleAccountInvariants() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNonNegativePowerInvariant(t *testing.T) {
//	type args struct {
//		k Keeper
//	}
//
//	_, _, keeper := createTestInput(t, true)
//
//	tests := []struct {
//		name string
//		args args
//		want sdk.Invariant
//	}{
//		{"Test NonNegativePowerInvariant",args{k:keeper},func(context sdk.Context) (string, bool){
//			return "", true
//		}},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NonNegativePowerInvariant(tt.args.k); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NonNegativePowerInvariant() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestRegisterInvariants(t *testing.T) {
//	type args struct {
//		ir sdk.InvariantRegistry
//		k  Keeper
//	}
//	_, _, keeper := createTestInput(t, true)
//
//	cdc := codec.New()
//	paramsKeeper := params.NewKeeper(
//		cdc, sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey), params.DefaultCodespace,
//	)
//
//	cKeeper := crisis.NewKeeper(paramsKeeper.Subspace("crisis"), 1, nil, "test")
//
//	tests := []struct {
//		name string
//		args args
//	}{
//		{"Test RegisterInvariants",args{
//			ir: cKeeper,
//			k:  keeper,
//		}},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//		})
//	}
//}
