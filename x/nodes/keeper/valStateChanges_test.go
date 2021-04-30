package keeper

import (
	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"reflect"
	"testing"
)

func TestKeeper_FinishUnstakingValidator(t *testing.T) {
	type fields struct {
		keeper Keeper
	}

	type args struct {
		ctx       sdk.Context
		validator types.Validator
	}

	validator := getStakedValidator()
	validator.StakedTokens = sdk.NewInt(0)
	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   sdk.Error
	}{
		{"Test FinishUnstakingValidator", fields{keeper: keeper}, args{
			ctx:       context,
			validator: validator,
		}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			// todo: add more tests scenarios
			k.FinishUnstakingValidator(tt.args.ctx, tt.args.validator)
		})
	}
}

func TestValidatorStateChange_EditAndValidateStakeValidator(t *testing.T) {
	stakeAmount := sdk.NewInt(100000000000)
	accountAmount := sdk.NewInt(1000000000000).Add(stakeAmount)
	bumpStakeAmount := sdk.NewInt(1000000000000)
	newChains := []string{"0021"}
	val := getUnstakedValidator()
	val.StakedTokens = sdk.ZeroInt()
	// updatedStakeAmount
	updateStakeAmountApp := val
	updateStakeAmountApp.StakedTokens = bumpStakeAmount
	// updatedStakeAmountFail
	updateStakeAmountAppFail := val
	updateStakeAmountAppFail.StakedTokens = stakeAmount.Sub(sdk.OneInt())
	// updatedStakeAmountNotEnoughCoins
	notEnoughCoinsAccount := stakeAmount
	// updateChains
	updateChainsVal := val
	updateChainsVal.StakedTokens = stakeAmount
	updateChainsVal.Chains = newChains
	// updateServiceURL
	updateServiceURL := val
	updateServiceURL.StakedTokens = stakeAmount
	updateServiceURL.Chains = newChains
	updateServiceURL.ServiceURL = "https://newServiceUrl.com"
	//same app no change no fail
	updateNothingval := val
	updateNothingval.StakedTokens = stakeAmount
	tests := []struct {
		name          string
		accountAmount sdk.BigInt
		origApp       types.Validator
		amount        sdk.BigInt
		want          types.Validator
		err           sdk.Error
	}{
		{
			name:          "edit stake amount of existing validator",
			accountAmount: accountAmount,
			origApp:       val,
			amount:        stakeAmount,
			want:          updateStakeAmountApp,
		},
		{
			name:          "FAIL edit stake amount of existing validator",
			accountAmount: accountAmount,
			origApp:       val,
			amount:        stakeAmount,
			want:          updateStakeAmountAppFail,
			err:           types.ErrMinimumEditStake("pos"),
		},
		{
			name:          "edit stake the chains of the validator",
			accountAmount: accountAmount,
			origApp:       val,
			amount:        stakeAmount,
			want:          updateChainsVal,
		},
		{
			name:          "edit stake the serviceurl of the validator",
			accountAmount: accountAmount,
			origApp:       val,
			amount:        stakeAmount,
			want:          updateChainsVal,
		},
		{
			name:          "FAIL not enough coins to bump stake amount of existing validator",
			accountAmount: notEnoughCoinsAccount,
			origApp:       val,
			amount:        stakeAmount,
			want:          updateStakeAmountApp,
			err:           types.ErrNotEnoughCoins("pos"),
		},
		{
			name:          "update nothing for the validator",
			accountAmount: accountAmount,
			origApp:       val,
			amount:        stakeAmount,
			want:          updateNothingval,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test setup
			codec.UpgradeHeight = -1
			context, _, keeper := createTestInput(t, true)
			coins := sdk.NewCoins(sdk.NewCoin(keeper.StakeDenom(context), tt.accountAmount))
			err := keeper.AccountKeeper.MintCoins(context, types.StakedPoolName, coins)
			if err != nil {
				t.Fail()
			}
			err = keeper.AccountKeeper.SendCoinsFromModuleToAccount(context, types.StakedPoolName, tt.origApp.Address, coins)
			if err != nil {
				t.Fail()
			}
			err = keeper.StakeValidator(context, tt.origApp, tt.amount)
			if err != nil {
				t.Fail()
			}
			// test begins here
			err = keeper.ValidateValidatorStaking(context, tt.want, tt.want.StakedTokens)
			if err != nil {
				if tt.err.Error() != err.Error() {
					t.Fatalf("Got error %s wanted error %s", err, tt.err)
				}
				return
			}
			// edit stake
			_ = keeper.StakeValidator(context, tt.want, tt.want.StakedTokens)
			tt.want.Status = sdk.Staked
			// see if the changes stuck
			got, _ := keeper.GetValidator(context, tt.origApp.Address)
			if !got.Equals(tt.want) {
				t.Fatalf("Got app %s\nWanted app %s", got.String(), tt.want.String())
			}
		})

	}
}

func TestKeeper_JailValidator(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx  sdk.Context
		addr sdk.Address
	}

	validator := getStakedValidator()
	context, _, keeper := createTestInput(t, true)
	keeper.SetValidator(context, validator)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test JailValidator", fields{keeper: keeper}, args{
			ctx:  context,
			addr: validator.GetAddress(),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			k.JailValidator(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_ReleaseWaitingValidators(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx sdk.Context
	}

	validator := getUnstakingValidator()
	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test ReleaseWaitingValidators", fields{keeper: keeper}, args{ctx: context}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			k.SetWaitingValidator(tt.args.ctx, validator)
			k.ReleaseWaitingValidators(tt.args.ctx)
		})
	}
}

func TestKeeper_StakeValidator(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx       sdk.Context
		validator types.Validator
		amount    sdk.BigInt
	}

	validator := getStakedValidator()
	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   sdk.Error
	}{
		{"Test StakeValidator", fields{keeper: keeper}, args{
			ctx:       context,
			validator: validator,
			amount:    sdk.ZeroInt(),
		}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			if got := k.StakeValidator(tt.args.ctx, tt.args.validator, tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StakeValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_UnjailValidator(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx  sdk.Context
		addr sdk.Address
	}
	validator := getStakedValidator()
	context, _, keeper := createTestInput(t, true)
	validator.Jailed = true
	keeper.SetValidator(context, validator)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test UnjailValidator", fields{keeper: keeper}, args{
			ctx:  context,
			addr: validator.GetAddress(),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			k.UnjailValidator(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_UpdateTendermintValidators(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx sdk.Context
	}

	//validator := getStakedValidator()
	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name        string
		fields      fields
		args        args
		wantUpdates []abci.ValidatorUpdate
	}{
		{"Test UpdateTenderMintValidators", fields{keeper: keeper}, args{ctx: context},
			[]abci.ValidatorUpdate{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			if gotUpdates := k.UpdateTendermintValidators(tt.args.ctx); !assert.True(t, len(gotUpdates) == len(tt.wantUpdates)) {
				t.Errorf("UpdateTendermintValidators() = %v, want %v", gotUpdates, tt.wantUpdates)
			}
		})
	}
}

func TestKeeper_ValidateValidatorBeginUnstaking(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx       sdk.Context
		validator types.Validator
	}

	validator := getStakedValidator()
	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   sdk.Error
	}{
		{"Test ValidateValidatorBeginUnstaking", fields{keeper: keeper}, args{
			ctx:       context,
			validator: validator,
		}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			if got := k.ValidateValidatorBeginUnstaking(tt.args.ctx, tt.args.validator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateValidatorBeginUnstaking() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_ValidateValidatorFinishUnstaking(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx       sdk.Context
		validator types.Validator
	}

	validator := getUnstakingValidator()
	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   sdk.Error
	}{
		{"Test ValidateValidatorFinishUnstaking", fields{keeper: keeper}, args{
			ctx:       context,
			validator: validator,
		}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			if got := k.ValidateValidatorFinishUnstaking(tt.args.ctx, tt.args.validator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateValidatorFinishUnstaking() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_ValidateValidatorStaking(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx       sdk.Context
		validator types.Validator
		amount    sdk.BigInt
	}

	validator := getUnstakedValidator()
	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   sdk.Error
	}{
		{"Test ValidateValidatorStaking", fields{keeper: keeper}, args{
			ctx:       context,
			validator: validator,
			amount:    sdk.NewInt(1000000),
		}, types.ErrNotEnoughCoins(types.ModuleName)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			if got := k.ValidateValidatorStaking(tt.args.ctx, tt.args.validator, tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateValidatorStaking() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_WaitToBeginUnstakingValidator(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx       sdk.Context
		validator types.Validator
	}

	validator := getStakedValidator()
	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   sdk.Error
	}{
		{"Test WaitToBeginUnstakingValidator", fields{keeper: keeper}, args{
			ctx:       context,
			validator: validator,
		}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			if got := k.WaitToBeginUnstakingValidator(tt.args.ctx, tt.args.validator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WaitToBeginUnstakingValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}
