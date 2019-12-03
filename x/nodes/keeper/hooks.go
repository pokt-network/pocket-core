package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
)

// Implements POSHooks interface
var _ types.POSHooks = Keeper{}

func (k Keeper) BeforeValidatorRegistered(ctx sdk.Context, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeValidatorRegistered(ctx, valAddr)
	}
}

func (k Keeper) AfterValidatorRegistered(ctx sdk.Context, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterValidatorRegistered(ctx, valAddr)
	}
}

func (k Keeper) BeforeValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeValidatorRemoved(ctx, consAddr, valAddr)
	}
}

func (k Keeper) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterValidatorRemoved(ctx, consAddr, valAddr)
	}
}

func (k Keeper) BeforeValidatorStaked(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeValidatorStaked(ctx, consAddr, valAddr)
	}
}

func (k Keeper) AfterValidatorStaked(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterValidatorStaked(ctx, consAddr, valAddr)
	}
}

func (k Keeper) BeforeValidatorBeginUnstaking(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeValidatorBeginUnstaking(ctx, consAddr, valAddr)
	}
}

func (k Keeper) AfterValidatorBeginUnstaking(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterValidatorBeginUnstaking(ctx, consAddr, valAddr)
	}
}

func (k Keeper) BeforeValidatorUnstaked(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeValidatorUnstaked(ctx, consAddr, valAddr)
	}
}

func (k Keeper) AfterValidatorUnstaked(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterValidatorUnstaked(ctx, consAddr, valAddr)
	}
}

func (k Keeper) AfterValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	if k.hooks != nil {
		k.hooks.AfterValidatorSlashed(ctx, valAddr, fraction)
	}
}

func (k Keeper) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	if k.hooks != nil {
		k.hooks.BeforeValidatorSlashed(ctx, valAddr, fraction)
	}
}
