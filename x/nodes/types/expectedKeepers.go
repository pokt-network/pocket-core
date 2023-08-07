package types

import (
	sdk "github.com/pokt-network/pocket-core/types"
	authexported "github.com/pokt-network/pocket-core/x/auth/exported"
	posexported "github.com/pokt-network/pocket-core/x/nodes/exported"
)

// AuthKeeper defines the expected supply Keeper (noalias)
type AuthKeeper interface {
	GetSupply(ctx sdk.Ctx) authexported.SupplyI
	SetSupply(ctx sdk.Ctx, supply authexported.SupplyI)
	GetModuleAddress(name string) sdk.Address
	GetModuleAccount(ctx sdk.Ctx, moduleName string) authexported.ModuleAccountI
	SetModuleAccount(sdk.Ctx, authexported.ModuleAccountI)
	SendCoinsFromModuleToModule(ctx sdk.Ctx, senderModule, recipientModule string, amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToAccount(ctx sdk.Ctx, senderModule string, recipientAddr sdk.Address, amt sdk.Coins) sdk.Error
	SendCoinsFromAccountToModule(ctx sdk.Ctx, senderAddr sdk.Address, recipientModule string, amt sdk.Coins) sdk.Error
	MintCoins(ctx sdk.Ctx, moduleName string, amt sdk.Coins) sdk.Error
	BurnCoins(ctx sdk.Ctx, name string, amt sdk.Coins) sdk.Error
	IterateAccounts(ctx sdk.Ctx, process func(authexported.Account) (stop bool))
	GetCoins(ctx sdk.Ctx, addr sdk.Address) sdk.Coins
	SetCoins(ctx sdk.Ctx, addr sdk.Address, amt sdk.Coins) sdk.Error
	HasCoins(ctx sdk.Ctx, addr sdk.Address, amt sdk.Coins) bool
	SendCoins(ctx sdk.Ctx, fromAddr sdk.Address, toAddr sdk.Address, amt sdk.Coins) sdk.Error
	GetAccount(ctx sdk.Ctx, addr sdk.Address) authexported.Account
	GetFee(ctx sdk.Ctx, msg sdk.Msg) sdk.BigInt
}

type PocketKeeper interface {
	// clear the cache of validators for sessions and relays
	ClearSessionCache()
}

// ValidatorSet expected properties for the set of all validators (noalias)
type ValidatorSet interface {
	// iterate through validators by address, execute func for each validator
	IterateAndExecuteOverVals(sdk.Ctx, func(index int64, validator posexported.ValidatorI) (stop bool))
	// iterate through staked validators by address, execute func for each validator
	IterateAndExecuteOverStakedVals(sdk.Ctx, func(index int64, validator posexported.ValidatorI) (stop bool))
	// iterate through the validator set of the prevState block by address, execute func for each validator
	IterateAndExecuteOverPrevStateVals(sdk.Ctx, func(index int64, validator posexported.ValidatorI) (stop bool))
	// get a particular validator by address
	Validator(sdk.Ctx, sdk.Address) posexported.ValidatorI
	// total staked tokens within the validator set
	TotalTokens(sdk.Ctx) sdk.BigInt
	// jail a validator
	JailValidator(sdk.Ctx, sdk.Address)
	// unjail a validator
	UnjailValidator(sdk.Ctx, sdk.Address)
	// MaxValidators returns the maximum amount of staked validators
	MaxValidators(sdk.Ctx) int64
}
