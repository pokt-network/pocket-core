package keeper

import (
	"fmt"
	log2 "log"

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
	k.SetParams(ctx, params)
	k.SetPrevStateValidatorsPower(ctx, prevStateTotalPower)
	k.SetWaitingValidators(ctx, waitingValidators)
	k.SetValidators(ctx, validators)
	k.SetPreviousProposer(ctx, prevProposer)
	k.SetValidatorSigningInfos(ctx, signingInfos)
	k.Cdc.DisableUpgradeOverride()
}

func (k Keeper) GetMsgStakeOutputSigner(ctx sdk.Ctx, msg sdk.Msg) sdk.Address {
	stakeMsg, ok := msg.(*types.MsgStake)
	if !ok {
		return nil
	}
	operatorAddr := sdk.Address(stakeMsg.PublicKey.Address())
	outputAddr, found := k.GetValidatorOutputAddress(ctx, operatorAddr)
	if !found {
		return nil
	}
	return outputAddr
}
