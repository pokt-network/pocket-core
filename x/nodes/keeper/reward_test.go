package keeper

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	pcTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
)

type args struct {
	consAddress sdk.Address
}

func TestSetAndGetProposer(t *testing.T) {
	validator := getStakedValidator()
	consAddress := validator.GetAddress()

	tests := []struct {
		name            string
		args            args
		expectedAddress sdk.Address
	}{
		{
			name:            "can set the preivous proposer",
			args:            args{consAddress: consAddress},
			expectedAddress: consAddress,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			keeper.SetPreviousProposer(context, test.args.consAddress)
			receivedAddress := keeper.GetPreviousProposer(context)
			assert.True(t, test.expectedAddress.Equals(receivedAddress), "addresses do not match ")
		})
	}
}

func TestMint(t *testing.T) {
	validator := getStakedValidator()
	validatorAddress := validator.Address

	tests := []struct {
		name     string
		amount   sdk.BigInt
		expected string
		address  sdk.Address
		panics   bool
	}{
		{
			name:     "mints a coin",
			amount:   sdk.NewInt(90),
			expected: "a reward of ",
			address:  validatorAddress,
			panics:   false,
		},
		{
			name:     "errors invalid ammount of coins",
			amount:   sdk.NewInt(-1),
			expected: "negative coin amount: -1",
			address:  validatorAddress,
			panics:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			switch test.panics {
			case true:
				defer func() {
					err := recover().(error)
					assert.Contains(t, err.Error(), test.expected, "error does not match")
				}()
				_ = keeper.mint(context, test.amount, test.address)
			default:
				result := keeper.mint(context, test.amount, test.address)
				assert.Contains(t, result.Log, test.expected, "does not contain message")
				coins := keeper.AccountKeeper.GetCoins(context, sdk.Address(test.address))
				assert.True(t, sdk.NewCoins(sdk.NewCoin(keeper.StakeDenom(context), test.amount)).IsEqual(coins), "coins should match")
			}
		})
	}
}

func verifyAccountBalance(
	t *testing.T,
	k Keeper,
	ctx sdk.Context,
	address sdk.Address,
	expected sdk.BigInt,
) {
	acc := k.GetAccount(ctx, address)
	expectedCoins := sdk.NewCoins(sdk.NewCoin("upokt", expected))
	assert.True(t, acc.Coins.IsEqual(expectedCoins))
	if !acc.Coins.IsEqual(expectedCoins) {
		fmt.Println(
			"Balance mismatch",
			address,
			"actual=",
			acc.Coins,
			"expected=",
			expectedCoins,
		)
	}
}

func TestKeeper_rewardFromFees(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx              sdk.Context
		previousProposer sdk.Address
		Output           sdk.Address
		Amount           sdk.BigInt
	}

	originalTestMode := codec.TestMode
	originalRSCAL := codec.UpgradeFeatureMap[codec.RSCALKey]
	t.Cleanup(func() {
		codec.TestMode = originalTestMode
		codec.UpgradeFeatureMap[codec.RSCALKey] = originalRSCAL
	})

	stakedValidator := getStakedValidator()
	stakedValidator.OutputAddress = getRandomValidatorAddress()
	codec.UpgradeFeatureMap[codec.RSCALKey] = 0
	codec.TestMode = -3
	amount := sdk.NewInt(10000)
	fees := sdk.NewCoins(sdk.NewCoin("upokt", amount))
	context, _, keeper := createTestInput(t, true)
	fp := keeper.getFeePool(context)
	keeper.AccountKeeper.SetCoins(context, fp.GetAddress(), fees)
	fp = keeper.getFeePool(context)
	keeper.SetValidator(context, stakedValidator)
	assert.Equal(t, fees, fp.GetCoins())

	_, proposerCut := keeper.splitFeesCollected(context, amount)

	totalSupplyPrev := keeper.AccountKeeper.GetSupply(context).
		GetTotal().
		AmountOf("upokt")

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test blockReward", fields{keeper: keeper},
			args{
				ctx:              context,
				previousProposer: stakedValidator.GetAddress(),
				Output:           stakedValidator.OutputAddress,
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			ctx := tt.args.ctx
			k.blockReward(tt.args.ctx, tt.args.previousProposer)

			verifyAccountBalance(t, k, ctx, tt.args.Output, proposerCut)
			verifyAccountBalance(t, k, ctx, tt.args.previousProposer, sdk.ZeroInt())

			totalSupply := k.AccountKeeper.GetSupply(ctx).
				GetTotal().
				AmountOf("upokt")
			assert.True(t, totalSupply.Equal(totalSupplyPrev))
		})
	}
}

func TestKeeper_rewardFromRelays(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx               sdk.Context
		validator         sdk.Address
		Output            sdk.Address
		validatorNoOutput sdk.Address
		OutputNoOutput    sdk.Address
	}

	originalTestMode := codec.TestMode
	originalRSCAL := codec.UpgradeFeatureMap[codec.RSCALKey]
	t.Cleanup(func() {
		codec.TestMode = originalTestMode
		codec.UpgradeFeatureMap[codec.RSCALKey] = originalRSCAL
	})

	stakedValidator := getStakedValidator()
	stakedValidatorNoOutput := getStakedValidator()
	stakedValidatorNoOutput.OutputAddress = nil
	stakedValidator.OutputAddress = getRandomValidatorAddress()
	codec.TestMode = -3
	codec.UpgradeFeatureMap[codec.RSCALKey] = 0
	context, _, keeper := createTestInput(t, true)
	context = context.WithBlockHeight(codec.NonCustodial2AllowanceHeight)
	keeper.SetValidator(context, stakedValidator)
	keeper.SetValidator(context, stakedValidatorNoOutput)
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test RelayReward", fields{keeper: keeper},
			args{
				ctx:               context,
				validator:         stakedValidator.GetAddress(),
				Output:            stakedValidator.OutputAddress,
				validatorNoOutput: stakedValidatorNoOutput.GetAddress(),
			}},
	}

	totalSupplyPrev := keeper.AccountKeeper.GetSupply(context).
		GetTotal().
		AmountOf("upokt")

	relays := sdk.NewInt(10000)
	rewardCost := keeper.AccountKeeper.GetFee(context, pcTypes.MsgClaim{}).
		Add(keeper.AccountKeeper.GetFee(context, pcTypes.MsgProof{}))
	totalReward := relays.Mul(keeper.RelaysToTokensMultiplier(context))
	nodeReward, _ := keeper.splitRewards(context, totalReward)
	outputReward := nodeReward.Sub(rewardCost)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			ctx := tt.args.ctx

			k.RewardForRelays(tt.args.ctx, relays, tt.args.validator)
			verifyAccountBalance(t, k, ctx, tt.args.Output, outputReward)
			verifyAccountBalance(t, k, ctx, tt.args.validator, rewardCost)

			totalSupply := k.AccountKeeper.GetSupply(ctx).
				GetTotal().
				AmountOf("upokt")
			assert.True(t, totalSupply.Equal(totalSupplyPrev.Add(totalReward)))

			// no output now
			k.RewardForRelays(tt.args.ctx, relays, tt.args.validatorNoOutput)
			verifyAccountBalance(t, k, ctx, tt.args.validatorNoOutput, nodeReward)

			totalSupply = k.AccountKeeper.GetSupply(ctx).
				GetTotal().
				AmountOf("upokt")
			assert.True(t, totalSupply.Equal(totalSupplyPrev.Add(totalReward.MulRaw(2))))
		})
	}
}

func TestKeeper_rewardFromRelaysPIP22NoEXP(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx        sdk.Context
		baseReward sdk.BigInt
		relays     int64
		validator1 types.Validator
		validator2 types.Validator
		validator3 types.Validator
		validator4 types.Validator
	}

	codec.UpgradeFeatureMap[codec.RSCALKey] = 3
	context, _, keeper := createTestInput(t, true)
	context = context.WithBlockHeight(3)
	p := keeper.GetParams(context)
	p.ServicerStakeFloorMultiplier = types.DefaultServicerStakeFloorMultiplier
	p.ServicerStakeWeightMultiplier = types.DefaultServicerStakeWeightMultiplier
	p.ServicerStakeFloorMultiplierExponent = sdk.NewDec(1)
	p.ServicerStakeWeightCeiling = 60000000000
	keeper.SetParams(context, p)

	stakedValidatorBin1 := getStakedValidator()
	stakedValidatorBin1.StakedTokens = keeper.ServicerStakeFloorMultiplier(context)
	stakedValidatorBin2 := getStakedValidator()
	stakedValidatorBin2.StakedTokens = keeper.ServicerStakeFloorMultiplier(context).Mul(sdk.NewInt(2))
	stakedValidatorBin3 := getStakedValidator()
	stakedValidatorBin3.StakedTokens = keeper.ServicerStakeFloorMultiplier(context).Mul(sdk.NewInt(3))
	stakedValidatorBin4 := getStakedValidator()
	stakedValidatorBin4.StakedTokens = keeper.ServicerStakeFloorMultiplier(context).Mul(sdk.NewInt(4))

	numRelays := int64(10000)
	base := sdk.NewDec(1).Quo(keeper.ServicerStakeWeightMultiplier(context)).Mul(sdk.NewDec(numRelays)).Mul(sdk.NewDecWithPrec(89, 2)).TruncateInt().Mul(keeper.RelaysToTokensMultiplier(context))

	keeper.SetValidator(context, stakedValidatorBin1)
	keeper.SetValidator(context, stakedValidatorBin2)
	keeper.SetValidator(context, stakedValidatorBin3)
	keeper.SetValidator(context, stakedValidatorBin4)
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test RelayReward", fields{keeper: keeper},
			args{
				ctx:        context,
				baseReward: base,
				relays:     numRelays,
				validator1: stakedValidatorBin1,
				validator2: stakedValidatorBin2,
				validator3: stakedValidatorBin3,
				validator4: stakedValidatorBin4,
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			ctx := tt.args.ctx
			k.RewardForRelays(tt.args.ctx, sdk.NewInt(tt.args.relays), tt.args.validator1.GetAddress())
			acc := k.GetAccount(ctx, tt.args.validator1.GetAddress())
			assert.False(t, acc.Coins.IsZero())
			assert.True(t, acc.Coins.IsEqual(sdk.NewCoins(sdk.NewCoin("upokt", tt.args.baseReward))))
			k.RewardForRelays(tt.args.ctx, sdk.NewInt(tt.args.relays), tt.args.validator2.GetAddress())
			acc = k.GetAccount(ctx, tt.args.validator2.GetAddress())
			assert.False(t, acc.Coins.IsZero())
			assert.True(t, acc.Coins.IsEqual(sdk.NewCoins(sdk.NewCoin("upokt", tt.args.baseReward.Mul(sdk.NewInt(2))))))
			k.RewardForRelays(tt.args.ctx, sdk.NewInt(tt.args.relays), tt.args.validator3.GetAddress())
			acc = k.GetAccount(ctx, tt.args.validator3.GetAddress())
			assert.False(t, acc.Coins.IsZero())
			assert.True(t, acc.Coins.IsEqual(sdk.NewCoins(sdk.NewCoin("upokt", tt.args.baseReward.Mul(sdk.NewInt(3))))))
			k.RewardForRelays(tt.args.ctx, sdk.NewInt(tt.args.relays), tt.args.validator4.GetAddress())
			acc = k.GetAccount(ctx, tt.args.validator4.GetAddress())
			assert.False(t, acc.Coins.IsZero())
			assert.True(t, acc.Coins.IsEqual(sdk.NewCoins(sdk.NewCoin("upokt", tt.args.baseReward.Mul(sdk.NewInt(4))))))
		})
	}
}

func TestKeeper_checkPIP22CheckCeiling(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx        sdk.Context
		baseReward sdk.BigInt
		relays     int64
		validator1 types.Validator
		validator2 types.Validator
	}

	codec.UpgradeFeatureMap[codec.RSCALKey] = 3
	context, _, keeper := createTestInput(t, true)
	context = context.WithBlockHeight(3)
	p := keeper.GetParams(context)
	p.ServicerStakeFloorMultiplier = types.DefaultServicerStakeFloorMultiplier
	p.ServicerStakeWeightMultiplier = types.DefaultServicerStakeWeightMultiplier
	p.ServicerStakeFloorMultiplierExponent = sdk.NewDec(1)
	p.ServicerStakeWeightCeiling = 15000000000
	keeper.SetParams(context, p)

	stakedValidatorBin1 := getStakedValidator()
	stakedValidatorBin1.StakedTokens = keeper.ServicerStakeFloorMultiplier(context)
	stakedValidatorBin2 := getStakedValidator()
	stakedValidatorBin2.StakedTokens = keeper.ServicerStakeFloorMultiplier(context).Mul(sdk.NewInt(2))

	numRelays := int64(10000)
	base := sdk.NewDec(1).Quo(keeper.ServicerStakeWeightMultiplier(context)).Mul(sdk.NewDec(numRelays)).Mul(sdk.NewDecWithPrec(89, 2)).TruncateInt().Mul(keeper.RelaysToTokensMultiplier(context))

	keeper.SetValidator(context, stakedValidatorBin1)
	keeper.SetValidator(context, stakedValidatorBin2)
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test RelayReward", fields{keeper: keeper},
			args{
				ctx:        context,
				baseReward: base,
				relays:     numRelays,
				validator1: stakedValidatorBin1,
				validator2: stakedValidatorBin2,
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			ctx := tt.args.ctx
			k.RewardForRelays(tt.args.ctx, sdk.NewInt(tt.args.relays), tt.args.validator1.GetAddress())
			acc := k.GetAccount(ctx, tt.args.validator1.GetAddress())
			assert.False(t, acc.Coins.IsZero())
			assert.True(t, acc.Coins.IsEqual(sdk.NewCoins(sdk.NewCoin("upokt", tt.args.baseReward))))
			k.RewardForRelays(tt.args.ctx, sdk.NewInt(tt.args.relays), tt.args.validator2.GetAddress())
			acc = k.GetAccount(ctx, tt.args.validator2.GetAddress())
			assert.False(t, acc.Coins.IsZero())
			assert.True(t, acc.Coins.IsEqual(sdk.NewCoins(sdk.NewCoin("upokt", tt.args.baseReward))))
		})
	}
}

func TestKeeper_rewardFromRelaysPIP22EXP(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx        sdk.Context
		validator1 types.Validator
		validator2 types.Validator
		validator3 types.Validator
		validator4 types.Validator
	}

	codec.UpgradeFeatureMap[codec.RSCALKey] = 3
	context, _, keeper := createTestInput(t, true)
	context = context.WithBlockHeight(3)
	p := keeper.GetParams(context)
	p.ServicerStakeFloorMultiplier = types.DefaultServicerStakeFloorMultiplier
	p.ServicerStakeWeightMultiplier = types.DefaultServicerStakeWeightMultiplier
	p.ServicerStakeFloorMultiplierExponent = sdk.NewDecWithPrec(50, 2)
	p.ServicerStakeWeightMultiplier = sdk.NewDec(1)
	p.ServicerStakeWeightCeiling = 60000000000
	keeper.SetParams(context, p)

	stakedValidatorBin1 := getStakedValidator()
	stakedValidatorBin1.StakedTokens = keeper.ServicerStakeFloorMultiplier(context)
	stakedValidatorBin2 := getStakedValidator()
	stakedValidatorBin2.StakedTokens = keeper.ServicerStakeFloorMultiplier(context).Mul(sdk.NewInt(2))
	stakedValidatorBin3 := getStakedValidator()
	stakedValidatorBin3.StakedTokens = keeper.ServicerStakeFloorMultiplier(context).Mul(sdk.NewInt(3))
	stakedValidatorBin4 := getStakedValidator()
	stakedValidatorBin4.StakedTokens = keeper.ServicerStakeFloorMultiplier(context).Mul(sdk.NewInt(4))

	keeper.SetValidator(context, stakedValidatorBin1)
	keeper.SetValidator(context, stakedValidatorBin2)
	keeper.SetValidator(context, stakedValidatorBin3)
	keeper.SetValidator(context, stakedValidatorBin4)
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test RelayReward", fields{keeper: keeper},
			args{
				ctx:        context,
				validator1: stakedValidatorBin1,
				validator2: stakedValidatorBin2,
				validator3: stakedValidatorBin3,
				validator4: stakedValidatorBin4,
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			ctx := tt.args.ctx
			k.RewardForRelays(tt.args.ctx, sdk.NewInt(1000), tt.args.validator1.GetAddress())
			acc := k.GetAccount(ctx, tt.args.validator1.GetAddress())
			assert.False(t, acc.Coins.IsZero())
			assert.True(t, acc.Coins.IsEqual(sdk.NewCoins(sdk.NewCoin("upokt", sdk.NewInt(890000)))))
			k.RewardForRelays(tt.args.ctx, sdk.NewInt(1000), tt.args.validator2.GetAddress())
			acc = k.GetAccount(ctx, tt.args.validator2.GetAddress())
			assert.False(t, acc.Coins.IsZero())
			assert.True(t, acc.Coins.IsEqual(sdk.NewCoins(sdk.NewCoin("upokt", sdk.NewInt(1258650)))))
			k.RewardForRelays(tt.args.ctx, sdk.NewInt(1000), tt.args.validator3.GetAddress())
			acc = k.GetAccount(ctx, tt.args.validator3.GetAddress())
			assert.False(t, acc.Coins.IsZero())
			assert.True(t, acc.Coins.IsEqual(sdk.NewCoins(sdk.NewCoin("upokt", sdk.NewInt(1541525)))))
			k.RewardForRelays(tt.args.ctx, sdk.NewInt(1000), tt.args.validator4.GetAddress())
			acc = k.GetAccount(ctx, tt.args.validator4.GetAddress())
			assert.False(t, acc.Coins.IsZero())
			assert.True(t, acc.Coins.IsEqual(sdk.NewCoins(sdk.NewCoin("upokt", sdk.NewInt(1780000)))))
		})
	}
}

func TestKeeper_RewardForRelaysPerChain(t *testing.T) {
	Height_PIP22 := int64(3)
	Height_PerChainRTTM := int64(10)
	Chain_Normal := "0001"
	Chain_HighProfit := "0002"
	RTTM_Default := int64(10000)
	RTTM_High := int64(15000)
	MinStake := int64(15000000000)
	RewardWeight := int64(4)
	ServicerAllocation := int64(85)
	NumOfRelays := sdk.NewInt(10)

	ExpectedRewards := func(multiplier int64) sdk.BigInt {
		return NumOfRelays.
			MulRaw(multiplier).
			MulRaw(RewardWeight).
			MulRaw(ServicerAllocation).
			QuoRaw(100)
	}

	codec.UpgradeFeatureMap[codec.RSCALKey] = Height_PIP22
	codec.UpgradeFeatureMap[codec.PerChainRTTM] = Height_PerChainRTTM

	ctx, _, keeper := createTestInput(t, true)

	// Add a validator
	validator := getStakedValidator()
	validator.Chains = []string{Chain_Normal, Chain_HighProfit}
	validator.StakedTokens = sdk.NewInt(MinStake * RewardWeight * 2)
	keeper.SetValidator(ctx, validator)

	p := keeper.GetParams(ctx)
	p.RelaysToTokensMultiplier = RTTM_Default
	p.DAOAllocation = 10
	p.ProposerAllocation = int64(100) - p.DAOAllocation - ServicerAllocation

	// Activate PIP-22
	ctx = ctx.WithBlockHeight(Height_PIP22)
	p.ServicerStakeFloorMultiplier = MinStake
	p.ServicerStakeWeightCeiling = MinStake * RewardWeight
	p.ServicerStakeWeightMultiplier = sdk.NewDec(1)
	p.ServicerStakeFloorMultiplierExponent = sdk.NewDec(1)
	keeper.SetParams(ctx, p)

	// Make sure PerChainRTTM is empty
	assert.NotNil(t, p.RelaysToTokensMultiplierMap)
	assert.Zero(t, len(p.RelaysToTokensMultiplierMap))

	// Set PerChainRTTM
	ctx = ctx.WithBlockHeight(Height_PerChainRTTM)
	p.RelaysToTokensMultiplierMap[Chain_HighProfit] = RTTM_High
	keeper.SetParams(ctx, p)

	// Make sure the default RTTM and PerChainRTTM
	p = keeper.GetParams(ctx)
	assert.Equal(t, len(p.RelaysToTokensMultiplierMap), 1)
	assert.Equal(t, p.RelaysToTokensMultiplierMap[Chain_HighProfit], RTTM_High)
	assert.Equal(t, p.RelaysToTokensMultiplier, RTTM_Default)

	// Verify the default multiplier
	rewardsDefault := keeper.RewardForRelays(
		ctx,
		NumOfRelays,
		validator.Address,
	)
	assert.True(t, rewardsDefault.Equal(ExpectedRewards(RTTM_Default)))

	// Verify the default multiplier with the chain ID
	rewardsNormalChain := keeper.RewardForRelaysPerChain(
		ctx,
		Chain_Normal,
		NumOfRelays,
		validator.Address,
	)
	assert.True(t, rewardsDefault.Equal(rewardsNormalChain))

	// Verify rewards with a non-default multiplier
	rewardsHighProfit := keeper.RewardForRelaysPerChain(
		ctx,
		Chain_HighProfit,
		NumOfRelays,
		validator.Address,
	)
	assert.True(t, rewardsHighProfit.Equal(ExpectedRewards(RTTM_High)))
	assert.True(t, rewardsDefault.LT(rewardsHighProfit))
}

func toArray(addr sdk.Address) [sdk.AddrLen]byte {
	var arr [sdk.AddrLen]byte
	copy(arr[:], addr)
	return arr
}

func indexToAddress(index int) sdk.Address {
	addr := make([]byte, sdk.AddrLen)
	binary.BigEndian.PutUint64(addr, uint64(index))
	return addr
}

func TestKeeper_SplitNodeRewards(t *testing.T) {
	var result map[[sdk.AddrLen]byte]sdk.BigInt
	callback := func(addr sdk.Address, rewards sdk.BigInt) {
		result[toArray(addr)] = rewards
	}

	recipient, _ := sdk.AddressFromHex("ffffffffffffffffffffffffffffffffffffffff")

	verifyResult := func(
		t *testing.T,
		totalRewards sdk.BigInt,
		expectedBalances []uint64,
	) {
		// Verify each delegator receives expected rewards
		for idx, expectedBalance := range expectedBalances {
			rewardsResult, found := result[toArray(indexToAddress(idx))]
			if expectedBalance == 0 {
				assert.False(t, found, "Rewards shouldn't be dispatched")
				continue
			}
			assert.True(t, found, "Rewards not dispatched")
			assert.Equal(t, expectedBalance, rewardsResult.Uint64(), "Wrong rewards")
			totalRewards = totalRewards.Sub(rewardsResult)
		}

		assert.False(t, totalRewards.IsNegative(), "Too many rewards")

		// Verify the recipient receives the remains if any exists
		reward, found := result[toArray(recipient)]
		if totalRewards.IsZero() {
			assert.False(t, found)
			return
		}
		assert.True(t, found)
		assert.True(t, reward.Equal(totalRewards))
	}

	delegatorMap := func(shares []uint32) map[string]uint32 {
		if len(shares) == 0 {
			return nil
		}
		m := map[string]uint32{}
		for idx, share := range shares {
			m[indexToAddress(idx).String()] = share
		}
		return m
	}

	totalBig := sdk.NewInt(10004)
	totalSmall := sdk.NewInt(81)

	// All goes to the default recipient.
	result = map[[sdk.AddrLen]byte]sdk.BigInt{}
	SplitNodeRewards(totalBig, recipient, delegatorMap([]uint32{}), callback)
	verifyResult(t, totalBig, []uint64{})

	// All goes to the delegator.
	result = map[[sdk.AddrLen]byte]sdk.BigInt{}
	SplitNodeRewards(totalBig, recipient, delegatorMap([]uint32{100}), callback)
	verifyResult(t, totalBig, []uint64{totalBig.Uint64()})

	// Multiple delegators.  Remainder goes to the recipient.
	result = map[[sdk.AddrLen]byte]sdk.BigInt{}
	SplitNodeRewards(totalBig, recipient, delegatorMap([]uint32{1, 0, 2, 30, 50}), callback)
	verifyResult(t, totalBig, []uint64{100, 0, 200, 3001, 5002})

	// Share less than a single token is truncated.
	result = map[[sdk.AddrLen]byte]sdk.BigInt{}
	SplitNodeRewards(totalSmall, recipient, delegatorMap([]uint32{1, 1, 1, 1}), callback)
	verifyResult(t, totalSmall, []uint64{})
	result = map[[sdk.AddrLen]byte]sdk.BigInt{}
	SplitNodeRewards(totalSmall, recipient, delegatorMap([]uint32{1, 2}), callback)
	verifyResult(t, totalSmall, []uint64{0, 1})

	// Invalid delegator map: all goes to the recipient
	result = map[[sdk.AddrLen]byte]sdk.BigInt{}
	SplitNodeRewards(totalSmall, recipient, delegatorMap([]uint32{1, 0xffffffff}), callback)
	verifyResult(t, totalSmall, []uint64{})
}
