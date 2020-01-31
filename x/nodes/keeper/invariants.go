package keeper

import (
	"bytes"
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes/exported"

	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
)

// RegisterInvariants registers all staking invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-accounts",
		ModuleAccountInvariants(k))
	ir.RegisterRoute(types.ModuleName, "nonnegative-power",
		NonNegativePowerInvariant(k))
}

// ModuleAccountInvariants checks that the staked ModuleAccounts pools
// reflects the tokens actively staked and not staked
func ModuleAccountInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		staked := sdk.ZeroInt()
		notStaked := sdk.ZeroInt()
		stakedPool := k.GetStakedTokens(ctx)
		notStakedPool := k.GetUnstakedTokens(ctx)

		k.IterateAndExecuteOverVals(ctx, func(_ int64, validator exported.ValidatorI) bool {
			switch validator.GetStatus() {
			case sdk.Staked, sdk.Unstaking:
				staked = staked.Add(validator.GetTokens())
			case sdk.Unstaked:
				notStaked = notStaked.Add(validator.GetTokens())
			default:
				panic("invalid validator status")
			}
			return false
		})

		broken := !stakedPool.Equal(staked) || !notStakedPool.Equal(notStaked)

		// Staked tokens should equal sum of tokens with staked validators
		// Not-staked tokens should equal unstaked tokens from validators
		return sdk.FormatInvariant(types.ModuleName, "staked and not staked module account coins", fmt.Sprintf(
			"\tPool's staked tokens: %v\n"+
				"\tsum of staked tokens: %v\n"+
				"not staked token invariance:\n"+
				"\tPool's not staked tokens: %v\n"+
				"\tsum of not staked tokens: %v\n"+
				"module accounts total (staked + not staked):\n"+
				"\tModule Accounts' tokens: %v\n"+
				"\tsum tokens:              %v\n",
			stakedPool, staked, notStakedPool, notStaked, stakedPool.Add(notStakedPool), staked.Add(notStaked))), broken
	}
}

// NonNegativePowerInvariant checks that all stored validators have >= 0 power.
func NonNegativePowerInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var broken bool

		iterator := k.stakedValsIterator(ctx)

		for ; iterator.Valid(); iterator.Next() {
			validator, found := k.GetValidator(ctx, iterator.Value())
			if !found {
				panic(fmt.Sprintf("validator record not found for address: %X\n", iterator.Value()))
			}

			powerKey := types.KeyForValidatorInStakingSet(validator)

			if !bytes.Equal(iterator.Key(), powerKey) {
				broken = true
				msg += fmt.Sprintf("power store invariance:\n\tvalidator.Power: %v"+
					"\n\tkey should be: %v\n\tkey in store: %v\n",
					validator.GetConsensusPower(), powerKey, iterator.Key())
			}

			if validator.StakedTokens.IsNegative() {
				broken = true
				msg += fmt.Sprintf("\tnegative tokens for validator: %v\n", validator)
			}
		}
		iterator.Close()
		return sdk.FormatInvariant(types.ModuleName, "nonnegative power", fmt.Sprintf("found invalid validator powers\n%s", msg)), broken
	}
}
