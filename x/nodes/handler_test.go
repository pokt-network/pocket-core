package nodes

import (
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/keeper"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	type args struct {
		k keeper.Keeper
	}
	tests := []struct {
		name string
		args args
		want sdk.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleMsgBeginUnstake(t *testing.T) {
	type args struct {
		ctx sdk.Context
		msg types.MsgBeginUnstake
		k   keeper.Keeper
	}
	tests := []struct {
		name string
		args args
		want sdk.Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleMsgBeginUnstake(tt.args.ctx, tt.args.msg, tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgBeginUnstake() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleMsgSend(t *testing.T) {
	type args struct {
		ctx sdk.Context
		msg types.MsgSend
		k   keeper.Keeper
	}
	tests := []struct {
		name string
		args args
		want sdk.Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleMsgSend(tt.args.ctx, tt.args.msg, tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgSend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleMsgUnjail(t *testing.T) {
	type args struct {
		ctx sdk.Context
		msg types.MsgUnjail
		k   keeper.Keeper
	}
	tests := []struct {
		name string
		args args
		want sdk.Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleMsgUnjail(tt.args.ctx, tt.args.msg, tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgUnjail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_ValidateBeginUnstakeSigner(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   keeper.Keeper
		v   types.Validator
		msg types.MsgBeginUnstake
	}
	unauthSigner := getRandomValidatorAddress()
	validator := getStakedValidator()
	validatorNoOuptut := validator
	validatorNoOuptut.OutputAddress = nil
	context, _, keeper := createTestInput(t, true)
	msgAuthorizedByValidator := types.MsgBeginUnstake{
		Address: validator.Address,
		Signer:  validator.Address,
	}
	msgAuthorizedByOutput := types.MsgBeginUnstake{
		Address: validator.Address,
		Signer:  validator.OutputAddress,
	}
	msgUnauthorizedSigner := types.MsgBeginUnstake{
		Address: validator.Address,
		Signer:  unauthSigner,
	}
	tests := []struct {
		name string
		args args
		want sdk.CodeType
	}{
		{"Test ValidateBeginUnstake With Output Address & AuthorizedByValidator", args{
			ctx: context,
			k:   keeper,
			v:   validator,
			msg: msgAuthorizedByValidator,
		}, 0},
		{"Test ValidateBeginUnstake With Output Address & AuthorizedByOutput", args{
			ctx: context,
			k:   keeper,
			v:   validator,
			msg: msgAuthorizedByOutput,
		}, 0},
		{"Test ValidateBeginUnstake Without Output Address & AuthorizedByValidator", args{
			ctx: context,
			k:   keeper,
			v:   validatorNoOuptut,
			msg: msgAuthorizedByValidator,
		}, 0},
		{"Test ValidateBeginUnstake Without Output Address & AuthroizedByOutput", args{
			ctx: context,
			k:   keeper,
			v:   validatorNoOuptut,
			msg: msgAuthorizedByOutput,
		}, types.CodeUnauthorizedSigner},
		{"Test ValidateBeginUnstake Without Output Address & Unauthorized", args{
			ctx: context,
			k:   keeper,
			v:   validatorNoOuptut,
			msg: msgUnauthorizedSigner,
		}, types.CodeUnauthorizedSigner},

		{"Test ValidateBeginUnstake With Output Address & Unauthorized", args{
			ctx: context,
			k:   keeper,
			v:   validator,
			msg: msgUnauthorizedSigner,
		}, types.CodeUnauthorizedSigner},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keeper.SetValidator(tt.args.ctx, tt.args.v)
			keeper.SetValidatorSigningInfo(tt.args.ctx, tt.args.v.Address, types.ValidatorSigningInfo{
				Address:             tt.args.v.Address,
				StartHeight:         0,
				Index:               0,
				JailedUntil:         time.Time{},
				MissedBlocksCounter: 0,
				JailedBlocksCounter: 0,
			})
			res := handleMsgBeginUnstake(tt.args.ctx, tt.args.msg, tt.args.k)
			assert.Equal(t, tt.want, res.Code)
		})
	}
}
