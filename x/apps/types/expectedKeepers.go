package types

import (
	appexported "github.com/pokt-network/pocket-core/x/apps/exported"
	sdk "github.com/pokt-network/posmint/types"
	supplyexported "github.com/pokt-network/posmint/x/supply/exported"
)

type PosKeeper interface {
	StakeDenom(ctx sdk.Context) (res string)
	// GetStakedTokens total staking tokens supply which is staked
	GetStakedTokens(ctx sdk.Context) sdk.Int
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
	// send coins from module to application
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.Address, amt sdk.Coins) sdk.Error
	// send coins from application to module
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.Address, recipientModule string, amt sdk.Coins) sdk.Error
	// burn coins
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error
	// mint coins for testing
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
}

// ApplicationSet expected properties for the set of all applications (noalias)
type ApplicationSet interface {
	// iterate through applications by address, execute func for each application
	IterateAndExecuteOverApps(sdk.Context, func(index int64, application appexported.ApplicationI) (stop bool))
	// iterate through staked applications by address, execute func for each application
	IterateAndExecuteOverStakedApps(sdk.Context, func(index int64, application appexported.ApplicationI) (stop bool))
	// get a particular application by address
	Application(sdk.Context, sdk.Address) appexported.ApplicationI
	// total staked tokens within the application set
	TotalTokens(sdk.Context) sdk.Int
	// jail a application
	JailApplication(sdk.Context, sdk.Address)
	// unjail a application
	UnjailApplication(sdk.Context, sdk.Address)
	// MaxApplications returns the maximum amount of staked applications
	MaxApplications(sdk.Context) uint64
}

//_______________________________________________________________________________
// Event Hooks
// These can be utilized to communicate between the pos keeper and another
// keeper which must take particular actions when applications change
// state. The second keeper must implement this interface, which then the
// staking keeper can call.

// AppHooks event hooks for staking application object (noalias)
type AppHooks interface {
	BeforeApplicationRegistered(ctx sdk.Context, valAddr sdk.Address)
	AfterApplicationRegistered(ctx sdk.Context, valAddr sdk.Address)
	BeforeApplicationRemoved(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address)
	AfterApplicationRemoved(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address)
	BeforeApplicationStaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address)
	AfterApplicationStaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address)
	BeforeApplicationBeginUnstaking(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address)
	AfterApplicationBeginUnstaking(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address)
	BeforeApplicationUnstaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address)
	AfterApplicationUnstaked(ctx sdk.Context, consAddr sdk.Address, valAddr sdk.Address)
	BeforeApplicationSlashed(ctx sdk.Context, valAddr sdk.Address, fraction sdk.Dec)
	AfterApplicationSlashed(ctx sdk.Context, valAddr sdk.Address, fraction sdk.Dec)
}
