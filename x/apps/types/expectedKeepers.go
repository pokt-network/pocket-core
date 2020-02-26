package types

import (
	appexported "github.com/pokt-network/pocket-core/x/apps/exported"
	sdk "github.com/pokt-network/posmint/types"
	supplyexported "github.com/pokt-network/posmint/x/supply/exported"
)

type PosKeeper interface {
	StakeDenom(ctx sdk.Ctx) (res string)
	// GetStakedTokens total staking tokens supply which is staked
	GetStakedTokens(ctx sdk.Ctx) sdk.Int
}

// SupplyKeeper defines the expected supply Keeper (noalias)
type SupplyKeeper interface {
	// get total supply of tokens
	GetSupply(ctx sdk.Ctx) supplyexported.SupplyI
	// get the address of a module account
	GetModuleAddress(name string) sdk.Address
	// get the module account structure
	GetModuleAccount(ctx sdk.Ctx, moduleName string) supplyexported.ModuleAccountI
	// set module account structure
	SetModuleAccount(sdk.Ctx, supplyexported.ModuleAccountI)
	// send coins to/from module accounts
	SendCoinsFromModuleToModule(ctx sdk.Ctx, senderModule, recipientModule string, amt sdk.Coins) sdk.Error
	// send coins from module to application
	SendCoinsFromModuleToAccount(ctx sdk.Ctx, senderModule string, recipientAddr sdk.Address, amt sdk.Coins) sdk.Error
	// send coins from application to module
	SendCoinsFromAccountToModule(ctx sdk.Ctx, senderAddr sdk.Address, recipientModule string, amt sdk.Coins) sdk.Error
	// burn coins
	BurnCoins(ctx sdk.Ctx, name string, amt sdk.Coins) sdk.Error
	// mint coins for testing
	MintCoins(ctx sdk.Ctx, moduleName string, amt sdk.Coins) sdk.Error
}

// ApplicationSet expected properties for the set of all applications (noalias)
type ApplicationSet interface {
	// iterate through applications by address, execute func for each application
	IterateAndExecuteOverApps(sdk.Ctx, func(index int64, application appexported.ApplicationI) (stop bool))
	// iterate through staked applications by address, execute func for each application
	IterateAndExecuteOverStakedApps(sdk.Ctx, func(index int64, application appexported.ApplicationI) (stop bool))
	// get a particular application by address
	Application(sdk.Ctx, sdk.Address) appexported.ApplicationI
	// total staked tokens within the application set
	TotalTokens(sdk.Ctx) sdk.Int
	// jail a application
	JailApplication(sdk.Ctx, sdk.Address)
	// unjail a application
	UnjailApplication(sdk.Ctx, sdk.Address)
	// MaxApplications returns the maximum amount of staked applications
	MaxApplications(sdk.Ctx) uint64
}
