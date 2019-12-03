package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
)

// Implements AppHooks interface
var _ types.AppHooks = Keeper{}

func (k Keeper) BeforeApplicationRegistered(ctx sdk.Context, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeApplicationRegistered(ctx, valAddr)
	}
}

func (k Keeper) AfterApplicationRegistered(ctx sdk.Context, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterApplicationRegistered(ctx, valAddr)
	}
}

func (k Keeper) BeforeApplicationRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeApplicationRemoved(ctx, consAddr, valAddr)
	}
}

func (k Keeper) AfterApplicationRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterApplicationRemoved(ctx, consAddr, valAddr)
	}
}

func (k Keeper) BeforeApplicationStaked(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeApplicationStaked(ctx, consAddr, valAddr)
	}
}

func (k Keeper) AfterApplicationStaked(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterApplicationStaked(ctx, consAddr, valAddr)
	}
}

func (k Keeper) BeforeApplicationBeginUnstaking(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeApplicationBeginUnstaking(ctx, consAddr, valAddr)
	}
}

func (k Keeper) AfterApplicationBeginUnstaking(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterApplicationBeginUnstaking(ctx, consAddr, valAddr)
	}
}

func (k Keeper) BeforeApplicationUnstaked(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeApplicationUnstaked(ctx, consAddr, valAddr)
	}
}

func (k Keeper) AfterApplicationUnstaked(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterApplicationUnstaked(ctx, consAddr, valAddr)
	}
}

func (k Keeper) AfterApplicationSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	if k.hooks != nil {
		k.hooks.AfterApplicationSlashed(ctx, valAddr, fraction)
	}
}

func (k Keeper) BeforeApplicationSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	if k.hooks != nil {
		k.hooks.BeforeApplicationSlashed(ctx, valAddr, fraction)
	}
}
