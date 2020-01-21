package types

import (
	posexported "github.com/pokt-network/pocket-core/x/nodes/exported"
	sdk "github.com/pokt-network/posmint/types"
	authexported "github.com/pokt-network/posmint/x/auth/exported"
	supplyexported "github.com/pokt-network/posmint/x/supply/exported"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	IterateAccounts(ctx sdk.Context, process func(authexported.Account) (stop bool))
}

// SupplyKeeper defines the expected supply Keeper (noalias)
type SupplyKeeper interface {
	// get total supply of tokens
	GetSupply(ctx sdk.Context) supplyexported.SupplyI
	// get the address of a module account
	GetModuleAddress(name string) sdk.Address
	// get the module account structure
	GetModuleAccount(ctx sdk.Context, moduleName string) supplyexported.ModuleAccountI
	// set module account structure
	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)
	// send coins to/from module accounts
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) sdk.Error
	// send coins from module to validator
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.Address, amt sdk.Coins) sdk.Error
	// send coins from validator to module
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.Address, recipientModule string, amt sdk.Coins) sdk.Error
	// mint coins
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
	// burn coins
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error
}

// ValidatorSet expected properties for the set of all validators (noalias)
// todo this is here so other modules can conform to this interface
type ValidatorSet interface {
	// iterate through validators by address, execute func for each validator
	IterateAndExecuteOverVals(sdk.Context, func(index int64, validator posexported.ValidatorI) (stop bool))
	// iterate through staked validators by address, execute func for each validator
	IterateAndExecuteOverStakedVals(sdk.Context, func(index int64, validator posexported.ValidatorI) (stop bool))
	// iterate through the validator set of the prevState block by address, execute func for each validator
	IterateAndExecuteOverPrevStateVals(sdk.Context, func(index int64, validator posexported.ValidatorI) (stop bool))
	// get a particular validator by address
	Validator(sdk.Context, sdk.Address) posexported.ValidatorI
	// total staked tokens within the validator set
	TotalTokens(sdk.Context) sdk.Int
	// jail a validator
	JailValidator(sdk.Context, sdk.Address)
	// unjail a validator
	UnjailValidator(sdk.Context, sdk.Address)
	// MaxValidators returns the maximum amount of staked validators
	MaxValidators(sdk.Context) uint64
}

//_______________________________________________________________________________
// Event Hooks
// These can be utilized to communicate between the pos keeper and another
// keeper which must take particular actions when validators change
// state. The second keeper must implement this interface, which then the
// staking keeper can call.

// POSHooks event hooks for staking validator object (noalias)
type POSHooks interface {
	BeforeValidatorRegistered(ctx sdk.Context, valAddr sdk.Address)
	AfterValidatorRegistered(ctx sdk.Context, valAddr sdk.Address)
	BeforeValidatorRemoved(ctx sdk.Context, valAddr sdk.Address)
	AfterValidatorRemoved(ctx sdk.Context, valAddr sdk.Address)
	BeforeValidatorStaked(ctx sdk.Context, valAddr sdk.Address)
	AfterValidatorStaked(ctx sdk.Context, valAddr sdk.Address)
	BeforeValidatorBeginUnstaking(ctx sdk.Context, valAddr sdk.Address)
	AfterValidatorBeginUnstaking(ctx sdk.Context, valAddr sdk.Address)
	BeforeValidatorUnstaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address)
	AfterValidatorUnstaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address)
	BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.Address, fraction sdk.Dec)
	AfterValidatorSlashed(ctx sdk.Context, valAddr sdk.Address, fraction sdk.Dec)
}
