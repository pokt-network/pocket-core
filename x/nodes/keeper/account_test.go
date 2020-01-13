package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"reflect"
	"testing"
)

func TestKeeper_SendCoins(t *testing.T) {
	type fields struct {
		Keeper Keeper
	}
	type args struct {
		ctx         sdk.Context
		fromAddress sdk.ValAddress
		toAddress   sdk.ValAddress
		amount      sdk.Int
	}

	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   sdk.Error
	}{
		{"Test Send 0 Coins", fields{Keeper: keeper}, args{
			ctx:         context,
			fromAddress: getRandomValidatorAddress(),
			toAddress:   getRandomValidatorAddress(),
			amount:      sdk.ZeroInt(),
		}, nil},
		{"Test Send 100 Coins", fields{Keeper: keeper}, args{
			ctx:         context,
			fromAddress: getRandomValidatorAddress(),
			toAddress:   getRandomValidatorAddress(),
			amount:      sdk.NewInt(100),
		}, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", sdk.NewCoins(), sdk.NewCoins(sdk.NewCoin(keeper.StakeDenom(context), sdk.NewInt(100)))))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.Keeper

			if got := k.SendCoins(tt.args.ctx, tt.args.fromAddress, tt.args.toAddress, tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SendCoins() = %v, want %v", got, tt.want)
			}
		})
	}
}
