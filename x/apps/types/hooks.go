package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

type MultiAppsHooks []AppHooks

func NewMultiStakingHooks(hooks ...AppHooks) MultiAppsHooks {
	return hooks
}

func (h MultiAppsHooks) BeforeApplicationRegistered(ctx sdk.Context, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].BeforeApplicationRegistered(ctx, valAddr)
	}
}

// nolint
func (h MultiAppsHooks) AfterApplicationRegistered(ctx sdk.Context, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterApplicationRegistered(ctx, valAddr)
	}
}

func (h MultiAppsHooks) BeforeApplicationRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, address sdk.ValAddress) {
	for i := range h {
		h[i].BeforeApplicationRemoved(ctx, consAddr, address)
	}
}

func (h MultiAppsHooks) AfterApplicationRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterApplicationRemoved(ctx, consAddr, valAddr)
	}
}

func (h MultiAppsHooks) BeforeApplicationStaked(ctx sdk.Context, consAddr sdk.ConsAddress, address sdk.ValAddress) {
	for i := range h {
		h[i].BeforeApplicationStaked(ctx, consAddr, address)
	}
}

func (h MultiAppsHooks) AfterApplicationStaked(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterApplicationStaked(ctx, consAddr, valAddr)
	}
}

func (h MultiAppsHooks) BeforeApplicationBeginUnstaking(ctx sdk.Context, consAddr sdk.ConsAddress, address sdk.ValAddress) {
	for i := range h {
		h[i].BeforeApplicationBeginUnstaking(ctx, consAddr, address)
	}
}
func (h MultiAppsHooks) AfterApplicationBeginUnstaking(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterApplicationBeginUnstaking(ctx, consAddr, valAddr)
	}
}

func (h MultiAppsHooks) BeforeApplicationBeginUnstaked(ctx sdk.Context, consAddr sdk.ConsAddress, address sdk.ValAddress) {
	for i := range h {
		h[i].BeforeApplicationUnstaked(ctx, consAddr, address)
	}
}
func (h MultiAppsHooks) AfterApplicationBeginUnstaked(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterApplicationUnstaked(ctx, consAddr, valAddr)
	}
}
func (h MultiAppsHooks) BeforeApplicationSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	for i := range h {
		h[i].BeforeApplicationSlashed(ctx, valAddr, fraction)
	}
}

func (h MultiAppsHooks) AfterApplicationSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	for i := range h {
		h[i].AfterApplicationSlashed(ctx, valAddr, fraction)
	}
}
