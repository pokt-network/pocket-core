package keeper

import (
	"fmt"
	log2 "log"
	"sync"

	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// Keeper of the staking store
type Keeper struct {
	storeKey      sdk.StoreKey
	Cdc           *codec.Codec
	AccountKeeper types.AuthKeeper
	PocketKeeper  types.PocketKeeper // todo combine all modules
	Paramstore    sdk.Subspace
	// codespace
	codespace sdk.CodespaceType
	// Cache
	validatorCache *sdk.Cache
	valPowerCache  *lockedCache
	stakedValAddrs *lockedCache
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, accountKeeper types.AuthKeeper,
	paramstore sdk.Subspace, codespace sdk.CodespaceType) Keeper {
	// ensure staked module accounts are set
	if addr := accountKeeper.GetModuleAddress(types.StakedPoolName); addr == nil {
		log2.Fatal(fmt.Errorf("%s module account has not been set", types.StakedPoolName))
	}
	cache := sdk.NewCache(int(types.ValidatorCacheSize))

	return Keeper{
		storeKey:       key,
		AccountKeeper:  accountKeeper,
		Paramstore:     paramstore.WithKeyTable(ParamKeyTable()),
		codespace:      codespace,
		validatorCache: cache,
		valPowerCache:  &lockedCache{nil, &sync.Mutex{}},
		stakedValAddrs: &lockedCache{nil, &sync.Mutex{}},
		Cdc:            cdc,
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

func (k Keeper) UpgradeCodec(ctx sdk.Ctx) {
	if ctx.IsOnUpgradeHeight() {
		k.ConvertState(ctx)
	}
}

func (k Keeper) ConvertState(ctx sdk.Ctx) {
	k.Cdc.SetUpgradeOverride(false)
	_ = k.InitSigningInfosCache(ctx)
	k.valPowerCache = &lockedCache{nil, &sync.Mutex{}}
	k.stakedValAddrs = &lockedCache{nil, &sync.Mutex{}}
	params := k.GetParams(ctx)
	prevStateTotalPower := k.PrevStateValidatorsPower(ctx)
	validators := k.GetAllValidators(ctx)
	waitingValidators := k.GetWaitingValidators(ctx)
	prevProposer := k.GetPreviousProposer(ctx)
	var prevStateValidatorPowers []types.PrevStatePowerMapping
	k.IterateAndExecuteOverPrevStateValsByPower(ctx, func(addr sdk.Address, power int64) (stop bool) {
		prevStateValidatorPowers = append(prevStateValidatorPowers, types.PrevStatePowerMapping{Address: addr, Power: power})
		return false
	})
	signingInfos := make([]types.ValidatorSigningInfo, 0)
	k.IterateAndExecuteOverValSigningInfo(ctx, func(address sdk.Address, info types.ValidatorSigningInfo) (stop bool) {
		signingInfos = append(signingInfos, info)
		return false
	})
	err := k.UpgradeMissedBlocksArray(ctx, validators) // TODO might be able to remove missed array code
	if err != nil {
		panic(err)
	}
	k.Cdc.SetUpgradeOverride(true)
	// custom logic for minSignedPerWindow
	params.MinSignedPerWindow = params.MinSignedPerWindow.QuoInt64(params.SignedBlocksWindow)
	k.SetParams(ctx, params)
	k.SetPrevStateValidatorsPower(ctx, prevStateTotalPower)
	k.SetWaitingValidators(ctx, waitingValidators)
	k.SetValidators(ctx, validators)
	k.SetPreviousProposer(ctx, prevProposer)
	k.SetValidatorSigningInfos(ctx, signingInfos)
	k.Cdc.DisableUpgradeOverride()
}

type lockedCache struct {
	store interface{}
	l     *sync.Mutex
}

func (lc *lockedCache) Set(ctx sdk.Ctx, v interface{}) {
	lc.l.Lock()
	defer lc.l.Unlock()
	if ctx.IsPrevCtx() {
		return
	}
	lc.store = v
}

func (lc *lockedCache) Get(ctx sdk.Ctx) (interface{}, bool) {
	lc.l.Lock()
	defer lc.l.Unlock()
	if ctx.IsPrevCtx() {
		return nil, false
	}
	s := lc.store
	return s, true
}

func (lc *lockedCache) Peek() bool {
	return lc.store != nil
}
