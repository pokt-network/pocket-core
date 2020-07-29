package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth/exported"
	"github.com/pokt-network/pocket-core/x/auth/types"
	"os"
)

// GetModuleAddress returns an address based on the module name
func (k Keeper) GetModuleAddress(moduleName string) sdk.Address {
	permAddr, ok := k.permAddrs[moduleName]
	if !ok {
		return nil
	}
	return permAddr.GetAddress()
}

// GetModuleAddressAndPermissions returns an address and permissions based on the module name
func (k Keeper) GetModuleAddressAndPermissions(moduleName string) (addr sdk.Address, permissions []string) {
	permAddr, ok := k.permAddrs[moduleName]
	if !ok {
		return addr, permissions
	}
	return permAddr.GetAddress(), permAddr.GetPermissions()
}

// GetModuleAccountAndPermissions gets the module account from the auth account store and its
// registered permissions
func (k Keeper) GetModuleAccountAndPermissions(ctx sdk.Ctx, moduleName string) (exported.ModuleAccountI, []string) {
	addr, perms := k.GetModuleAddressAndPermissions(moduleName)
	if addr == nil {
		return nil, []string{}
	}
	acc := k.GetAccount(ctx, addr)
	if acc != nil {
		macc, ok := acc.(exported.ModuleAccountI)
		if !ok {
			fmt.Println("account that is retrieved is not a module account")
			return types.ModuleAccount{}, []string{}
		}
		return macc, perms
	}

	// create a new module account
	macc := types.NewEmptyModuleAccount(moduleName, perms...)
	maccI := (k.NewAccount(ctx, macc)).(exported.ModuleAccountI) // set the account number
	k.SetModuleAccount(ctx, maccI)

	return maccI, perms
}

// GetModuleAccount gets the module account from the auth account store
func (k Keeper) GetModuleAccount(ctx sdk.Ctx, moduleName string) exported.ModuleAccountI {
	acc, _ := k.GetModuleAccountAndPermissions(ctx, moduleName)
	return acc
}

// SetModuleAccount sets the module account to the auth account store
func (k Keeper) SetModuleAccount(ctx sdk.Ctx, macc exported.ModuleAccountI) {
	k.SetAccount(ctx, macc)
}

// ValidatePermissions validates that the module account has been granted
// permissions within its set of allowed permissions.
func (k Keeper) ValidatePermissions(macc exported.ModuleAccountI) error {
	permAddr := k.permAddrs[macc.GetName()]
	for _, perm := range macc.GetPermissions() {
		if !permAddr.HasPermission(perm) {
			return fmt.Errorf("invalid module permission %s", perm)
		}
	}
	return nil
}

// NewAccount creates a new account
func (k Keeper) NewAccount(ctx sdk.Ctx, acc exported.Account) exported.Account {
	return acc
}

// GetAccount implements sdk.Keeper.
func (k Keeper) GetAccount(ctx sdk.Ctx, addr sdk.Address) exported.Account {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AddressStoreKey(addr))
	if bz == nil {
		return nil
	}
	acc, err := k.decodeAccount(bz)
	if err != nil {
		return nil // Could not decode account
	}
	return acc
}

// GetAllAccounts returns all accounts in the accountKeeper.
func (k Keeper) GetAllAccounts(ctx sdk.Ctx) []exported.Account {
	var accounts []exported.Account
	appendAccount := func(acc exported.Account) (stop bool) {
		accounts = append(accounts, acc)
		return false
	}
	k.IterateAccounts(ctx, appendAccount)
	return accounts
}

// GetAllAccounts returns all accounts in the accountKeeper.
func (k Keeper) GetAllAccountsExport(ctx sdk.Ctx) []exported.Account {
	var accounts []exported.Account
	appendAccount := func(acc exported.Account) (stop bool) {
		//not get empty coins accounts
		if !acc.GetCoins().Empty() {
			//sanity check here
			if acc.GetAddress() != nil && acc.GetPubKey() != nil {
				accounts = append(accounts, acc)
			}
		}
		return false
	}
	k.IterateAccounts(ctx, appendAccount)
	return accounts
}

// SetAccount implements sdk.Keeper.
func (k Keeper) SetAccount(ctx sdk.Ctx, acc exported.Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalBinaryBare(acc)
	if err != nil {
		ctx.Logger().Error(fmt.Errorf("error marshalling account %v err: %s", acc, err.Error()).Error())
		os.Exit(1)
	}
	store.Set(types.AddressStoreKey(addr), bz)
}

// RemoveAccount removes an account for the account mapper store.
// NOTE: this will cause supply invariant violation if called
func (k Keeper) RemoveAccount(ctx sdk.Ctx, acc exported.Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AddressStoreKey(addr))
}

// IterateAccounts implements sdk.Keeper.
func (k Keeper) IterateAccounts(ctx sdk.Ctx, process func(exported.Account) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AddressStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		acc, err := k.decodeAccount(val)
		if err != nil {
			ctx.Logger().Error(fmt.Errorf("error while iterating accounts: unmarshalling account %v err: %s", val, err.Error()).Error())
			continue
		}
		if process(acc) {
			return
		}
		iter.Next()
	}
}

func (k Keeper) decodeAccount(bz []byte) (acc exported.Account, err error) {
	err = k.cdc.UnmarshalBinaryBare(bz, &acc)
	return
}

// NewAccountWithAddress implements sdk.AuthKeeper.
func (k Keeper) NewAccountWithAddress(ctx sdk.Ctx, addr sdk.Address) (exported.Account, error) {
	acc := types.BaseAccount{}
	err := acc.SetAddress(addr)
	if err != nil {
		return nil, fmt.Errorf("unable to create a new account with address %s", addr)
	}
	return &acc, nil
}

// GetPubKey Returns the PublicKey of the account at address
func (k Keeper) GetPubKey(ctx sdk.Ctx, addr sdk.Address) (crypto.PublicKey, sdk.Error) {
	acc := k.GetAccount(ctx, addr)
	if acc == nil {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}
	return acc.GetPubKey(), nil
}
