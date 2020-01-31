package keeper

import (
	"bytes"
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/exported"

	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
)

// RegisterInvariants registers all staking invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-accounts",
		ModuleAccountInvariants(k))
	ir.RegisterRoute(types.ModuleName, "nonnegative-power",
		NonNegativePowerInvariant(k))
}

// todo add max relay invariant compared to staked status / relays

// ModuleAccountInvariants checks that the staked ModuleAccounts pools
// reflects the tokens actively staked and not staked
func ModuleAccountInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		staked := sdk.ZeroInt()
		notStaked := sdk.ZeroInt()
		stakedPool := k.GetStakedTokens(ctx)
		notStakedPool := k.GetUnstakedTokens(ctx)

		k.IterateAndExecuteOverApps(ctx, func(_ int64, application exported.ApplicationI) bool {
			switch application.GetStatus() {
			case sdk.Staked, sdk.Unstaking:
				staked = staked.Add(application.GetTokens())
			case sdk.Unstaked:
				notStaked = notStaked.Add(application.GetTokens())
			default:
				panic("invalid application status")
			}
			return false
		})

		broken := !stakedPool.Equal(staked) || !notStakedPool.Equal(notStaked)

		// Staked tokens should equal sum of tokens with staked applications
		// Not-staked tokens should equal unstaked tokens from applications
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

// NonNegativePowerInvariant checks that all stored applications have >= 0 power.
func NonNegativePowerInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var broken bool

		iterator := k.stakedAppsIterator(ctx)

		for ; iterator.Valid(); iterator.Next() {
			app, found := k.GetApplication(ctx, iterator.Value())
			if !found {
				panic(fmt.Sprintf("app record not found for address: %X\n", iterator.Value()))
			}

			powerKey := types.KeyForAppInStakingSet(app)

			if !bytes.Equal(iterator.Key(), powerKey) {
				broken = true
				msg += fmt.Sprintf("power store invariance:\n\tapp.Power: %v"+
					"\n\tkey should be: %v\n\tkey in store: %v\n",
					app.GetConsensusPower(), powerKey, iterator.Key())
			}

			if app.StakedTokens.IsNegative() {
				broken = true
				msg += fmt.Sprintf("\tnegative tokens for app: %v\n", app)
			}
		}
		iterator.Close()
		return sdk.FormatInvariant(types.ModuleName, "nonnegative power", fmt.Sprintf("found invalid application powers\n%s", msg)), broken
	}
}
