package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

type MultiPOSHooks []POSHooks

func NewMultiStakingHooks(hooks ...POSHooks) MultiPOSHooks {
	return hooks
}

func (h MultiPOSHooks) BeforeValidatorRegistered(ctx sdk.Context, valAddr sdk.Address) {
	for i := range h {
		h[i].BeforeValidatorRegistered(ctx, valAddr)
	}
}

// nolint
func (h MultiPOSHooks) AfterValidatorRegistered(ctx sdk.Context, valAddr sdk.Address) {
	for i := range h {
		h[i].AfterValidatorRegistered(ctx, valAddr)
	}
}

func (h MultiPOSHooks) BeforeValidtorRemoved(ctx sdk.Context, consAddr sdk.Address, address sdk.Address) {
	for i := range h {
		h[i].BeforeValidatorRemoved(ctx, address)
	}
}

func (h MultiPOSHooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	for i := range h {
		h[i].AfterValidatorRemoved(ctx, valAddr)
	}
}

func (h MultiPOSHooks) BeforeValidatorStaked(ctx sdk.Context, consAddr sdk.Address, address sdk.Address) {
	for i := range h {
		h[i].BeforeValidatorStaked(ctx, address)
	}
}

func (h MultiPOSHooks) AfterValidatorStaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	for i := range h {
		h[i].AfterValidatorStaked(ctx, valAddr)
	}
}

func (h MultiPOSHooks) BeforeValidatorBeginUnstaking(ctx sdk.Context, consAddr sdk.Address, address sdk.Address) {
	for i := range h {
		h[i].BeforeValidatorBeginUnstaking(ctx, address)
	}
}
func (h MultiPOSHooks) AfterValidatorBeginUnstaking(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	for i := range h {
		h[i].AfterValidatorBeginUnstaking(ctx, valAddr)
	}
}

func (h MultiPOSHooks) BeforeValidatorBeginUnstaked(ctx sdk.Context, consAddr sdk.Address, address sdk.Address) {
	for i := range h {
		h[i].BeforeValidatorUnstaked(ctx, consAddr, address)
	}
}
func (h MultiPOSHooks) AfterValidatorBeginUnstaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address) {
	for i := range h {
		h[i].AfterValidatorUnstaked(ctx, consAddr, valAddr)
	}
}
func (h MultiPOSHooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.Address, fraction sdk.Dec) {
	for i := range h {
		h[i].BeforeValidatorSlashed(ctx, valAddr, fraction)
	}
}

func (h MultiPOSHooks) AfterValidatorSlashed(ctx sdk.Context, valAddr sdk.Address, fraction sdk.Dec) {
	for i := range h {
		h[i].AfterValidatorSlashed(ctx, valAddr, fraction)
	}
}
