package keeper

import (
	"reflect"
	"testing"
	"time"

	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
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
	val.OutputAddress = val.Address
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
	// nil output addresss
	nilOutputAddress := val
	nilOutputAddress.OutputAddress = nil
	nilOutputAddress.StakedTokens = stakeAmount
	//same app no change no fail
	updateNothingval := val
	updateNothingval.StakedTokens = stakeAmount
	//new staked amount doesn't push into the next bin
	failPip22 := val
	failPip22.StakedTokens = sdk.NewInt(29999000000)
	//New staked amount does push into the next bin
	passPip22NextBin := val
	passPip22NextBin.StakedTokens = sdk.NewInt(30001000000)
	//All updates should pass above the ceiling
	passPip22AboveCeil := val
	passPip22AboveCeil.StakedTokens = sdk.NewInt(60000000000).Add(sdk.OneInt())

	tests := []struct {
		name          string
		accountAmount sdk.BigInt
		origApp       types.Validator
		amount        sdk.BigInt
		want          types.Validator
		err           sdk.Error
		PIP22Edit     bool
	}{
		{
			name:          "edit stake amount of existing validator",
			accountAmount: accountAmount,
			origApp:       val,
			amount:        stakeAmount,
			want:          updateStakeAmountApp,
			PIP22Edit:     true,
		},
		{
			name:          "FAIL edit stake amount of existing validator",
			accountAmount: accountAmount,
			origApp:       val,
			amount:        stakeAmount,
			want:          updateStakeAmountAppFail,
			err:           types.ErrMinimumEditStake("pos"),
			PIP22Edit:     false,
		},
		{
			name:          "edit stake the chains of the validator",
			accountAmount: accountAmount,
			origApp:       val,
			amount:        stakeAmount,
			want:          updateChainsVal,
			PIP22Edit:     false,
		},
		{
			name:          "edit stake the serviceurl of the validator",
			accountAmount: accountAmount,
			origApp:       val,
			amount:        stakeAmount,
			want:          updateChainsVal,
			PIP22Edit:     false,
		},
		{
			name:          "FAIL not enough coins to bump stake amount of existing validator",
			accountAmount: notEnoughCoinsAccount,
			origApp:       val,
			amount:        stakeAmount,
			want:          updateStakeAmountApp,
			err:           types.ErrNotEnoughCoins("pos"),
			PIP22Edit:     false,
		},
		{
			name:          "update nothing for the validator",
			accountAmount: accountAmount,
			origApp:       val,
			amount:        stakeAmount,
			want:          updateNothingval,
			PIP22Edit:     false,
		},
		{
			name:          "PIP22 not enough to bump bin",
			accountAmount: accountAmount,
			origApp:       val,
			amount:        sdk.NewInt(15001000000),
			want:          failPip22,
			err:           types.ErrSameBinEditStake("pos"),
			PIP22Edit:     true,
		},
		{
			name:          "PIP22 update to next bin",
			accountAmount: accountAmount,
			origApp:       val,
			amount:        sdk.NewInt(15001000000),
			want:          passPip22NextBin,
			PIP22Edit:     true,
		},
		{
			name:          "PIP22 above ceil",
			accountAmount: accountAmount,
			origApp:       val,
			amount:        sdk.NewInt(60000000000),
			want:          passPip22AboveCeil,
			PIP22Edit:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test setup
			codec.UpgradeHeight = -1
			if tt.PIP22Edit {
				codec.UpgradeFeatureMap[codec.RSCALKey] = -1
				codec.UpgradeFeatureMap[codec.VEDITKey] = -1
			} else {
				codec.UpgradeFeatureMap[codec.RSCALKey] = 0
				codec.UpgradeFeatureMap[codec.VEDITKey] = 0
			}
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
			err = keeper.StakeValidator(context, tt.origApp, tt.amount, tt.origApp.PublicKey)
			if err != nil {
				t.Fail()
			}
			// test begins here
			err = keeper.ValidateValidatorStaking(context, tt.want, tt.want.StakedTokens, sdk.Address(tt.origApp.PublicKey.Address()))
			if err != nil {
				if tt.err.Error() != err.Error() {
					t.Fatalf("Got error %s wanted error %s", err, tt.err)
				}
				return
			}
			// edit stake
			_ = keeper.StakeValidator(context, tt.want, tt.want.StakedTokens, tt.want.PublicKey)
			tt.want.Status = sdk.Staked
			// see if the changes stuck
			got, _ := keeper.GetValidator(context, tt.origApp.Address)
			assert.Nil(t, got.OutputAddress, "OutputAddress was set before NCUST update")
			// Manually updated `got` to account for post NCUST updates
			got.OutputAddress = tt.want.OutputAddress
			if !got.Equals(tt.want) {
				t.Fatalf("Got app %s\nWanted app %s", got.String(), tt.want.String())
			}
		})

	}
}
func TestValidatorStateChange_EditAndValidateStakeValidatorAfterNonCustodialUpgrade(t *testing.T) {
	originalUpgradeHeight := codec.UpgradeHeight
	originalTestMode := codec.TestMode
	t.Cleanup(func() {
		codec.UpgradeHeight = originalUpgradeHeight
		codec.TestMode = originalTestMode
	})

	stakeAmount := sdk.NewInt(100000000000)
	accountAmount := sdk.NewInt(1000000000000).Add(stakeAmount)
	bumpStakeAmount := sdk.NewInt(1000000000000)
	newChains := []string{"0021"}
	val := getUnstakedValidator()
	val.StakedTokens = sdk.ZeroInt()
	val.OutputAddress = val.Address
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
	// nil output addresss
	nilOutputAddress := val
	nilOutputAddress.OutputAddress = nil
	nilOutputAddress.StakedTokens = stakeAmount
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
			name:          "FAIL nil output address",
			accountAmount: notEnoughCoinsAccount,
			origApp:       val,
			amount:        stakeAmount,
			want:          nilOutputAddress,
			err:           types.ErrNilOutputAddr("pos"),
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
			codec.TestMode = -3
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
			err = keeper.StakeValidator(context, tt.origApp, tt.amount, tt.origApp.PublicKey)
			if err != nil {
				t.Fail()
			}
			// test begins here
			err = keeper.ValidateValidatorStaking(context, tt.want, tt.want.StakedTokens, sdk.Address(tt.origApp.PublicKey.Address()))
			if err != nil {
				if tt.err.Error() != err.Error() {
					t.Fatalf("Got error %s wanted error %s", err, tt.err)
				}
				return
			}
			// edit stake
			_ = keeper.StakeValidator(context, tt.want, tt.want.StakedTokens, tt.want.PublicKey)
			tt.want.Status = sdk.Staked
			// see if the changes stuck
			got, _ := keeper.GetValidator(context, tt.origApp.Address)
			if !got.Equals(tt.want) {
				t.Fatalf("Got app %s\nWanted app %s", got.String(), tt.want.String())
			}
		})
	}
}

// handleStakeForTesting is a helper fnction to stake a new node with
// a given MsgStake in the same way as handleStake does.
func handleStakeForTesting(
	ctx sdk.Context,
	k Keeper,
	msg types.MsgStake,
	signer crypto.PublicKey,
) sdk.Error {
	validator := types.NewValidatorFromMsg(msg)
	validator.StakedTokens = sdk.ZeroInt()
	if err := k.ValidateValidatorStaking(
		ctx, validator, msg.Value, sdk.Address(signer.Address())); err != nil {
		return err
	}
	return k.StakeValidator(ctx, validator, msg.Value, signer)
}

func TestValidatorStateChange_OutputAddressEdit(t *testing.T) {
	ctx, _, k := createTestInput(t, true)

	originalUpgradeHeight := codec.UpgradeHeight
	originalTestMode := codec.TestMode
	originalNCUST := codec.UpgradeFeatureMap[codec.NonCustodialUpdateKey]
	originalOEDIT := codec.UpgradeFeatureMap[codec.OutputAddressEditKey]
	t.Cleanup(func() {
		codec.UpgradeHeight = originalUpgradeHeight
		codec.TestMode = originalTestMode
		codec.UpgradeFeatureMap[codec.NonCustodialUpdateKey] = originalNCUST
		codec.UpgradeFeatureMap[codec.OutputAddressEditKey] = originalOEDIT
	})

	// Enable EditStake
	codec.UpgradeHeight = -1

	// Make sure NCUST is disabled
	codec.TestMode = 0

	stakeAmount := sdk.NewCoin(k.StakeDenom(ctx), sdk.NewInt(k.MinimumStake(ctx)))
	outputPubKey := getRandomPubKey()
	outputAddress := sdk.Address(outputPubKey.Address())

	runStake := func(
		operatorPubKey crypto.PublicKey,
		outputAddress sdk.Address,
		signer crypto.PublicKey,
	) sdk.Error {
		msgStake := types.MsgStake{
			Chains:     []string{"0021", "0040"},
			ServiceUrl: "https://www.pokt.network:443",
			Value:      stakeAmount.Amount,
			PublicKey:  operatorPubKey,
			Output:     outputAddress,
		}
		return handleStakeForTesting(ctx, k, msgStake, signer)
	}

	// Create and fund accounts
	operatorPubKey1 := getRandomPubKey()
	operatorPubKey2 := getRandomPubKey()
	operatorPubKey3 := getRandomPubKey()
	operatorAddr1 := sdk.Address(operatorPubKey1.Address())
	operatorAddr2 := sdk.Address(operatorPubKey2.Address())
	operatorAddr3 := sdk.Address(operatorPubKey3.Address())
	assert.Nil(t, fundAccount(ctx, k, operatorAddr1, stakeAmount))
	assert.Nil(t, fundAccount(ctx, k, operatorAddr2, stakeAmount))
	assert.Nil(t, fundAccount(ctx, k, outputAddress, stakeAmount))

	// Stake two nodes before NCUST
	assert.Nil(t, runStake(operatorPubKey1, outputAddress, operatorPubKey1))
	assert.Nil(t, runStake(operatorPubKey2, outputAddress, operatorPubKey2))

	// Verify staked nodes having nil output address
	validatorCur, found := k.GetValidator(ctx, operatorAddr1)
	assert.True(t, found)
	assert.Nil(t, validatorCur.OutputAddress)
	validatorCur, found = k.GetValidator(ctx, operatorAddr2)
	assert.True(t, found)
	assert.Nil(t, validatorCur.OutputAddress)

	// Enable NCUST
	codec.UpgradeFeatureMap[codec.NonCustodialUpdateKey] = -1

	// Set an output address on Addr1 after NCUST
	assert.Nil(t, runStake(operatorPubKey1, outputAddress, operatorPubKey1))
	validatorCur, found = k.GetValidator(ctx, operatorAddr1)
	assert.True(t, found)
	assert.Equal(t, validatorCur.OutputAddress, outputAddress)

	// Attempt to change the output address with operator's signature --> Fail
	err := runStake(operatorPubKey1, operatorAddr1, operatorPubKey1)
	assert.NotNil(t, err)
	assert.Equal(t, k.codespace, err.Codespace())
	assert.Equal(t, types.CodeUnequalOutputAddr, err.Code())

	// Attempt to change the output address with output's signature --> Fail
	err = runStake(operatorPubKey1, operatorAddr1, outputPubKey)
	assert.NotNil(t, err)
	assert.Equal(t, k.codespace, err.Codespace())
	assert.Equal(t, types.CodeUnauthorizedSigner, err.Code())

	// Enable OEDIT
	codec.UpgradeFeatureMap[codec.OutputAddressEditKey] = -1

	// Attempt to change the output address with operator's signature --> Fail
	err = runStake(operatorPubKey1, operatorAddr1, operatorPubKey1)
	assert.NotNil(t, err)
	assert.Equal(t, k.codespace, err.Codespace())
	assert.Equal(t, types.CodeDisallowedOutputAddressEdit, err.Code())

	// Attempt to change the output address with output's signature --> Success
	err = runStake(operatorPubKey1, operatorAddr1, outputPubKey)
	assert.Nil(t, err)
	validatorCur, found = k.GetValidator(ctx, operatorAddr1)
	assert.True(t, found)
	assert.Equal(t, validatorCur.OutputAddress, operatorAddr1)

	// Attempt to change the output address from nil
	// with operator's signature --> Success
	err = runStake(operatorPubKey2, outputAddress, operatorPubKey2)
	assert.Nil(t, err)
	validatorCur, found = k.GetValidator(ctx, operatorAddr2)
	assert.True(t, found)
	assert.Equal(t, validatorCur.OutputAddress, outputAddress)

	// New non-custodial stake with output's signature --> Success
	err = runStake(operatorPubKey3, outputAddress, outputPubKey)
	assert.Nil(t, err)
	validatorCur, found = k.GetValidator(ctx, operatorAddr3)
	assert.True(t, found)
	assert.Equal(t, validatorCur.OutputAddress, outputAddress)
}

func TestValidatorStateChange_Delegators(t *testing.T) {
	ctx, _, k := createTestInput(t, true)

	originalUpgradeHeight := codec.UpgradeHeight
	originalTestMode := codec.TestMode
	originalNCUST := codec.UpgradeFeatureMap[codec.NonCustodialUpdateKey]
	originalOEDIT := codec.UpgradeFeatureMap[codec.OutputAddressEditKey]
	originalReward := codec.UpgradeFeatureMap[codec.RewardDelegatorsKey]
	t.Cleanup(func() {
		codec.UpgradeHeight = originalUpgradeHeight
		codec.TestMode = originalTestMode
		codec.UpgradeFeatureMap[codec.NonCustodialUpdateKey] = originalNCUST
		codec.UpgradeFeatureMap[codec.OutputAddressEditKey] = originalOEDIT
		codec.UpgradeFeatureMap[codec.RewardDelegatorsKey] = originalReward
	})

	// Enable EditStake, NCUST, and OEDIT
	codec.TestMode = 0
	codec.UpgradeHeight = -1
	codec.UpgradeFeatureMap[codec.NonCustodialUpdateKey] = -1
	codec.UpgradeFeatureMap[codec.OutputAddressEditKey] = -1

	// Prepare accounts
	outputPubKey := getRandomPubKey()
	operatorPubKey1 := getRandomPubKey()
	operatorPubKey2 := getRandomPubKey()
	operatorAddr1 := sdk.Address(operatorPubKey1.Address())
	outputAddress := sdk.Address(outputPubKey.Address())
	operatorAddr2 := sdk.Address(operatorPubKey2.Address())

	// Fund output address for two nodes
	stakeAmount := sdk.NewCoin(k.StakeDenom(ctx), sdk.NewInt(k.MinimumStake(ctx)))
	assert.Nil(t, fundAccount(ctx, k, outputAddress, stakeAmount))
	assert.Nil(t, fundAccount(ctx, k, outputAddress, stakeAmount))

	runStake := func(
		operatorPubkey crypto.PublicKey,
		delegators map[string]uint32,
		signer crypto.PublicKey,
	) sdk.Error {
		msgStake := types.MsgStake{
			Chains:           []string{"0021", "0040"},
			ServiceUrl:       "https://www.pokt.network:443",
			Value:            stakeAmount.Amount,
			PublicKey:        operatorPubkey,
			Output:           outputAddress,
			RewardDelegators: delegators,
		}
		return handleStakeForTesting(ctx, k, msgStake, signer)
	}

	singleDelegator := map[string]uint32{}
	singleDelegator[getRandomValidatorAddress().String()] = 1

	// Attempt to set a delegators before the upgrade --> The field is ignored
	assert.Nil(t, runStake(operatorPubKey1, singleDelegator, outputPubKey))
	validatorCur, found := k.GetValidator(ctx, operatorAddr1)
	assert.True(t, found)
	assert.Nil(t, validatorCur.RewardDelegators)

	// Enable RewardDelegators
	codec.UpgradeFeatureMap[codec.RewardDelegatorsKey] = -1

	// Attempt to change the delegators with output's signature --> Fail
	err := runStake(operatorPubKey1, singleDelegator, outputPubKey)
	assert.NotNil(t, err)
	assert.Equal(t, k.codespace, err.Codespace())
	assert.Equal(t, types.CodeDisallowedRewardDelegatorEdit, err.Code())

	// Attempt to set the delegators with operator's signature --> Success
	err = runStake(operatorPubKey1, singleDelegator, operatorPubKey1)
	assert.Nil(t, err)
	validatorCur, found = k.GetValidator(ctx, operatorAddr1)
	assert.True(t, found)
	assert.True(
		t,
		sdk.CompareStringMaps(validatorCur.RewardDelegators, singleDelegator),
	)

	// Attempt to reset the delegators with operator's signature --> Success
	err = runStake(operatorPubKey1, nil, operatorPubKey1)
	assert.Nil(t, err)
	validatorCur, found = k.GetValidator(ctx, operatorAddr1)
	assert.True(t, found)
	assert.Nil(t, validatorCur.RewardDelegators)

	// New stake with delegators can be signed by the output --> Success
	err = runStake(operatorPubKey2, singleDelegator, outputPubKey)
	assert.Nil(t, err)
	validatorCur, found = k.GetValidator(ctx, operatorAddr2)
	assert.True(t, found)
	assert.True(
		t,
		sdk.CompareStringMaps(validatorCur.RewardDelegators, singleDelegator),
	)
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
			if got := k.StakeValidator(tt.args.ctx, tt.args.validator, tt.args.amount, tt.args.validator.PublicKey); !reflect.DeepEqual(got, tt.want) {
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

	originalTestMode := codec.TestMode
	t.Cleanup(func() {
		codec.TestMode = originalTestMode
	})

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
			codec.TestMode = -2
			if got := k.ValidateValidatorStaking(tt.args.ctx, tt.args.validator, tt.args.amount, sdk.Address(tt.args.validator.PublicKey.Address())); !reflect.DeepEqual(got, tt.want) {
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

func TestKeeper_ValidateUnjailMessage(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
		v   types.Validator
		msg types.MsgUnjail
	}

	originalTestMode := codec.TestMode
	t.Cleanup(func() {
		codec.TestMode = originalTestMode
	})
	codec.TestMode = -3

	unauthSigner := getRandomValidatorAddress()
	validator := getStakedValidator()
	validator.Jailed = true
	validator.OutputAddress = getRandomValidatorAddress()
	validatorNoOuptut := validator
	validatorNoOuptut.OutputAddress = nil
	context, _, keeper := createTestInput(t, true)
	msgUnjailAuthorizedByValidator := types.MsgUnjail{
		ValidatorAddr: validator.Address,
		Signer:        validator.Address,
	}
	msgUnjailAuthorizedByOutput := types.MsgUnjail{
		ValidatorAddr: validator.Address,
		Signer:        validator.OutputAddress,
	}
	msgUnjailUnauthorizedSigner := types.MsgUnjail{
		ValidatorAddr: validator.Address,
		Signer:        unauthSigner,
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{"Test ValidateUnjailMessage With Output Address & AuthorizedByValidator", args{
			ctx: context,
			k:   keeper,
			v:   validator,
			msg: msgUnjailAuthorizedByValidator,
		}, nil},
		{"Test ValidateUnjailMessage With Output Address & AuthorizedByOutput", args{
			ctx: context,
			k:   keeper,
			v:   validator,
			msg: msgUnjailAuthorizedByOutput,
		}, nil},
		{"Test ValidateUnjailMessage Without Output Address & AuthorizedByValidator", args{
			ctx: context,
			k:   keeper,
			v:   validatorNoOuptut,
			msg: msgUnjailAuthorizedByValidator,
		}, nil},
		{"Test ValidateUnjailMessage Without Output Address & AuthroizedByOutput", args{
			ctx: context,
			k:   keeper,
			v:   validatorNoOuptut,
			msg: msgUnjailAuthorizedByOutput,
		}, types.ErrUnauthorizedSigner("pos")},
		{"Test ValidateUnjailMessage Without Output Address & Unauthorized", args{
			ctx: context,
			k:   keeper,
			v:   validatorNoOuptut,
			msg: msgUnjailUnauthorizedSigner,
		}, types.ErrUnauthorizedSigner("pos")},

		{"Test ValidateUnjailMessage With Output Address & Unauthorized", args{
			ctx: context,
			k:   keeper,
			v:   validator,
			msg: msgUnjailUnauthorizedSigner,
		}, types.ErrUnauthorizedSigner("pos")},
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
			_, err := tt.args.k.ValidateUnjailMessage(tt.args.ctx, tt.args.msg)
			assert.Equal(t, tt.want, err)
		})
	}
}
