package keeper

import (
	"container/list"
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/pokt-network/posmint/x/params"
	"github.com/tendermint/tendermint/libs/log"
)

const aminoCacheSize = 500

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// keeper of the staking store
type Keeper struct {
	storeKey           sdk.StoreKey
	cdc                *codec.Codec
	accountKeeper      auth.AccountKeeper
	coinKeeper         bank.Keeper
	supplyKeeper       types.SupplyKeeper
	Paramstore         params.Subspace
	validatorCache     map[string]cachedValidator
	validatorCacheList *list.List

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, accountKeeper auth.AccountKeeper, coinKeeper bank.Keeper, supplyKeeper types.SupplyKeeper,
	paramstore params.Subspace, codespace sdk.CodespaceType) Keeper {
	// ensure staked module accounts are set
	if addr := supplyKeeper.GetModuleAddress(types.StakedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.StakedPoolName))
	}
	return Keeper{
		storeKey:           key,
		cdc:                cdc,
		accountKeeper:      accountKeeper,
		coinKeeper:         coinKeeper,
		supplyKeeper:       supplyKeeper,
		Paramstore:         paramstore.WithKeyTable(ParamKeyTable()),
		validatorCache:     make(map[string]cachedValidator, aminoCacheSize),
		validatorCacheList: list.New(),
		codespace:          codespace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}
