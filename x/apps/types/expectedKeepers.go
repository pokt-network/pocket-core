package types

import (
	sdk "github.com/pokt-network/pocket-core/types"
	appexported "github.com/pokt-network/pocket-core/x/apps/exported"
	authexported "github.com/pokt-network/pocket-core/x/auth/exported"
)

type PosKeeper interface {
	StakeDenom(ctx sdk.Ctx) (res string)
	// GetStakedTokens total staking tokens supply which is staked
	GetStakedTokens(ctx sdk.Ctx) sdk.BigInt
}

type PocketKeeper interface {
	// clear the cache of validators for sessions and relays
	ClearSessionCache()
}

// AuthKeeper defines the expected supply Keeper (noalias)
type AuthKeeper interface {
	// get total supply of tokens
	GetSupply(ctx sdk.Ctx) authexported.SupplyI
	// set total supply of tokens
	SetSupply(ctx sdk.Ctx, supply authexported.SupplyI)
	// get the address of a module account
	GetModuleAddress(name string) sdk.Address
	// get the module account structure
	GetModuleAccount(ctx sdk.Ctx, moduleName string) authexported.ModuleAccountI
	// set module account structure
	SetModuleAccount(sdk.Ctx, authexported.ModuleAccountI)
	// send coins to/from module accounts
	SendCoinsFromModuleToModule(ctx sdk.Ctx, senderModule, recipientModule string, amt sdk.Coins) sdk.Error
	// send coins from module to validator
	SendCoinsFromModuleToAccount(ctx sdk.Ctx, senderModule string, recipientAddr sdk.Address, amt sdk.Coins) sdk.Error
	// send coins from validator to module
	SendCoinsFromAccountToModule(ctx sdk.Ctx, senderAddr sdk.Address, recipientModule string, amt sdk.Coins) sdk.Error
	// mint coins
	MintCoins(ctx sdk.Ctx, moduleName string, amt sdk.Coins) sdk.Error
	// burn coins
	BurnCoins(ctx sdk.Ctx, name string, amt sdk.Coins) sdk.Error
	// iterate accounts
	IterateAccounts(ctx sdk.Ctx, process func(authexported.Account) (stop bool))
	// get coins
	GetCoins(ctx sdk.Ctx, addr sdk.Address) sdk.Coins
	// set coins
	SetCoins(ctx sdk.Ctx, addr sdk.Address, amt sdk.Coins) sdk.Error
	// has coins
	HasCoins(ctx sdk.Ctx, addr sdk.Address, amt sdk.Coins) bool
	// send coins
	SendCoins(ctx sdk.Ctx, fromAddr sdk.Address, toAddr sdk.Address, amt sdk.Coins) sdk.Error
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
	TotalTokens(sdk.Ctx) sdk.BigInt
	// jail a application
	JailApplication(sdk.Ctx, sdk.Address)
	// unjail a application
	UnjailApplication(sdk.Ctx, sdk.Address)
	// MaxApplications returns the maximum amount of staked applications
	MaxApplications(sdk.Ctx) int64
}
