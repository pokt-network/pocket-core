package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/gov/types"
	"os"
)

// Allocate subspace used for keepers
func (k Keeper) Subspace(s string) sdk.Subspace {
	_, ok := k.spaces[s]
	if ok {
		fmt.Println(fmt.Errorf("subspace %s already occupied", s))
		os.Exit(1)
	}
	if s == "" {
		fmt.Println(fmt.Errorf("cannot use empty stirng for subspace"))
		os.Exit(1)
	}
	space := sdk.NewSubspace(s)
	space.SetCodec(k.cdc)
	k.spaces[s] = space
	return space
}

func (k Keeper) AddSubspaces(subspaces ...sdk.Subspace) {
	for _, space := range subspaces {
		_, ok := k.spaces[space.Name()]
		if ok {
			fmt.Println(fmt.Errorf("subspace %s already occupied", space.Name()))
			os.Exit(1)
		}
		if space.Name() == "" {
			fmt.Println(fmt.Errorf("cannot use empty stirng for subspace"))
			os.Exit(1)
		}
		space.SetCodec(k.cdc)
		k.spaces[space.Name()] = space
	}
}

// Get existing substore from keeper
func (k Keeper) GetSubspace(s string) (sdk.Subspace, bool) {
	space, ok := k.spaces[s]
	if !ok {
		return sdk.Subspace{}, false
	}
	return space, ok
}

func (k Keeper) GetAllParamNames(ctx sdk.Ctx) (paramNames map[string]bool) {
	paramNames = make(map[string]bool)
	for _, space := range k.spaces {
		keys := space.GetAllParamKeys(ctx)
		for _, key := range keys {
			paramNames[space.Name()+"/"+key] = false // set to false for adjacency matrix
		}
	}
	return
}

func (k Keeper) GetAllParamNameValue(ctx sdk.Ctx) (paramNames map[string]string) {
	paramNames = make(map[string]string)
	for _, space := range k.spaces {
		keys := space.GetAllParamKeys(ctx)
		for _, key := range keys {
			paramNames[space.Name()+"/"+key] = string(space.GetIfExistsRaw(ctx, []byte(key)))
		}
	}
	return
}

func (k Keeper) HandleUpgrade(ctx sdk.Ctx, aclKey string, paramValue interface{}, owner sdk.Address) sdk.Result {
	if err := k.VerifyACL(ctx, aclKey, owner); err != nil {
		return err.Result()
	}
	subspaceName, paramKey := types.SplitACLKey(aclKey)
	space, ok := k.spaces[subspaceName]
	if !ok {
		k.Logger(ctx).Error(types.ErrSubspaceNotFound(types.ModuleName, subspaceName).Error())
		os.Exit(1)
	}
	space.Set(ctx, []byte(paramKey), paramValue)
	k.spaces[subspaceName] = space
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventParamChange,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, fmt.Sprintf("modified: %s to: %v", aclKey, paramValue)),
			sdk.NewAttribute(sdk.AttributeKeySender, owner.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, owner.String()),
		),
	})
	// if upgrade, emit separate upgrade event
	if aclKey == types.NewACLKey(types.ModuleName, string(types.UpgradeKey)) {
		u, ok := paramValue.(types.Upgrade)
		if !ok {
			ctx.Logger().Error(fmt.Sprintf("unable to convert %v to upgrade, can't emit event about upgrade, at height: %d", paramValue, ctx.BlockHeight()))
			return sdk.Result{Events: ctx.EventManager().Events()}
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventUpgrade,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, fmt.Sprintf("UPGRADE CONFIRMED: %s at height %v", u.UpgradeVersion(), u.UpgradeHeight())),
			sdk.NewAttribute(sdk.AttributeKeySender, owner.String()),
		))
	}
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func (k Keeper) ModifyParam(ctx sdk.Ctx, aclKey string, paramValue []byte, owner sdk.Address) sdk.Result {
	if err := k.VerifyACL(ctx, aclKey, owner); err != nil {
		return err.Result()
	}
	subspaceName, paramKey := types.SplitACLKey(aclKey)
	space, ok := k.spaces[subspaceName]
	if !ok {
		k.Logger(ctx).Error(types.ErrSubspaceNotFound(types.ModuleName, subspaceName).Error())
		os.Exit(1)
	}
	space.Update(ctx, []byte(paramKey), paramValue)
	k.spaces[subspaceName] = space
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventParamChange,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, fmt.Sprintf("modified: %s to: %v", aclKey, paramValue)),
			sdk.NewAttribute(sdk.AttributeKeySender, owner.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, owner.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
