package keeper

import (
	"fmt"
	"reflect"
	"testing"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_SendCoins(t *testing.T) {
	type fields struct {
		Keeper Keeper
	}
	type args struct {
		ctx         sdk.Context
		fromAddress sdk.Address
		toAddress   sdk.Address
		amount      sdk.BigInt
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

func TestKeeper_GetAccount(t *testing.T) {
	ctx, accs, keeper := createTestInput(t, false)
	acc := keeper.GetAccount(ctx, accs[0].GetAddress())
	assert.NotNil(t, acc)
	assert.Equal(t, accs[0], acc)
}
