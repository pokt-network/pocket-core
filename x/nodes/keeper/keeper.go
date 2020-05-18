package keeper

import (
	"container/list"
	"fmt"
	log2 "log"

	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// Keeper of the staking store
type Keeper struct {
	storeKey           sdk.StoreKey
	cdc                *codec.Codec
	AccountKeeper      types.AuthKeeper
	PocketKeeper       types.PocketKeeper // todo combine all modules
	Paramstore         sdk.Subspace
	validatorCache     map[string]cachedValidator
	validatorCacheList *list.List

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, accountKeeper types.AuthKeeper,
	paramstore sdk.Subspace, codespace sdk.CodespaceType) Keeper {
	// ensure staked module accounts are set
	if addr := accountKeeper.GetModuleAddress(types.StakedPoolName); addr == nil {
		log2.Fatal(fmt.Errorf("%s module account has not been set", types.StakedPoolName))
	}
	return Keeper{
		storeKey:           key,
		cdc:                cdc,
		AccountKeeper:      accountKeeper,
		Paramstore:         paramstore.WithKeyTable(ParamKeyTable()),
		validatorCache:     make(map[string]cachedValidator, types.ValidatorCacheSize),
		validatorCacheList: list.New(),
		codespace:          codespace,
	}
}

// Logger - returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Ctx) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Codespace - Retrieve the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}
