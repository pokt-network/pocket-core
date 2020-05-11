package keeper

import (
	"container/list"
	"fmt"
	"os"

	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	govKeeper "github.com/pokt-network/posmint/x/gov/keeper"
	"github.com/tendermint/tendermint/libs/log"
)

const aminoCacheSize = 500

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// Keeper of the staking store
type Keeper struct {
	storeKey           sdk.StoreKey
	cdc                *codec.Codec
	govKeeper          govKeeper.Keeper
	AccountKeeper      types.AuthKeeper
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
		fmt.Println(fmt.Errorf("%s module account has not been set", types.StakedPoolName))
		os.Exit(1)
	}
	return Keeper{
		storeKey:           key,
		cdc:                cdc,
		AccountKeeper:      accountKeeper,
		Paramstore:         paramstore.WithKeyTable(ParamKeyTable()),
		validatorCache:     make(map[string]cachedValidator, aminoCacheSize),
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
