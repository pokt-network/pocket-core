package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

type MultiPOSHooks []AppHooks

func (h MultiPOSHooks) BeforeApplicationRegistered(ctx sdk.Context, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].BeforeApplicationRegistered(ctx, valAddr)
	}
}

// nolint
func (h MultiPOSHooks) AfterApplicationRegistered(ctx sdk.Context, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterApplicationRegistered(ctx, valAddr)
	}
}

func (h MultiPOSHooks) BeforeApplicationRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, address sdk.ValAddress) {
	for i := range h {
		h[i].BeforeApplicationRemoved(ctx, consAddr, address)
	}
}

func (h MultiPOSHooks) AfterApplicationRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterApplicationRemoved(ctx, consAddr, valAddr)
	}
}

func (h MultiPOSHooks) BeforeApplicationStaked(ctx sdk.Context, consAddr sdk.ConsAddress, address sdk.ValAddress) {
	for i := range h {
		h[i].BeforeApplicationStaked(ctx, consAddr, address)
	}
}

func (h MultiPOSHooks) AfterApplicationStaked(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterApplicationStaked(ctx, consAddr, valAddr)
	}
}

func (h MultiPOSHooks) BeforeApplicationBeginUnstaking(ctx sdk.Context, consAddr sdk.ConsAddress, address sdk.ValAddress) {
	for i := range h {
		h[i].BeforeApplicationBeginUnstaking(ctx, consAddr, address)
	}
}
func (h MultiPOSHooks) AfterApplicationBeginUnstaking(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterApplicationBeginUnstaking(ctx, consAddr, valAddr)
	}
}

func (h MultiPOSHooks) BeforeApplicationBeginUnstaked(ctx sdk.Context, consAddr sdk.ConsAddress, address sdk.ValAddress) {
	for i := range h {
		h[i].BeforeApplicationUnstaked(ctx, consAddr, address)
	}
}
func (h MultiPOSHooks) AfterApplicationBeginUnstaked(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterApplicationUnstaked(ctx, consAddr, valAddr)
	}
}
func (h MultiPOSHooks) BeforeApplicationSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	for i := range h {
		h[i].BeforeApplicationSlashed(ctx, valAddr, fraction)
	}
}

func (h MultiPOSHooks) AfterApplicationSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	for i := range h {
		h[i].AfterApplicationSlashed(ctx, valAddr, fraction)
	}
}
