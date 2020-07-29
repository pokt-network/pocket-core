package keeper

import (
	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-amino"
	"reflect"
	"testing"
)

func Test_queryDAO(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
	}
	context, keeper := createTestKeeperAndContext(t, false)
	jsonresponse, _ := amino.MarshalJSONIndent(sdk.ZeroInt(), "", "  ")
	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{"Test QueryDao", args{
			ctx: context,
			k:   keeper,
		}, jsonresponse, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := queryDAO(tt.args.ctx, tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryDAO() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("queryDAO() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_queryACL(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
	}
	ctx, k := createTestKeeperAndContext(t, false)
	acl := k.GetACL(ctx)
	jsonresponse, err := codec.MarshalJSONIndent(types.ModuleCdc, acl)
	if err != nil {
		t.Fatalf("failed to JSON marshal result: " + err.Error())
	}
	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{"Test QueryACL", args{
			ctx: ctx,
			k:   k,
		}, jsonresponse, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := queryACL(tt.args.ctx, tt.args.k)
			var acl1 types.ACL
			err := k.cdc.UnmarshalJSON(got, &acl1)
			assert.Nil(t, err)
			assert.Equal(t, acl.String(), acl1.String())
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("queryACL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_queryUpgrade(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
	}
	ctx, k := createTestKeeperAndContext(t, false)
	upgrade := k.GetUpgrade(ctx)
	jsonresponse, err := codec.MarshalJSONIndent(types.ModuleCdc, upgrade)
	if err != nil {
		t.Fatalf("failed to JSON marshal result: " + err.Error())
	}
	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{"Test QueryUpgrade", args{
			ctx: ctx,
			k:   k,
		}, jsonresponse, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u types.Upgrade
			got, got1 := queryUpgrade(tt.args.ctx, tt.args.k)
			err := k.cdc.UnmarshalJSON(got, &u)
			assert.Nil(t, err)
			assert.Equal(t, upgrade, u)
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("queryUpgrade() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
