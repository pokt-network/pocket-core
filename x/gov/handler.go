package gov

import (
	"fmt"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/gov/keeper"
	"github.com/pokt-network/pocket-core/x/gov/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Ctx, msg sdk.LegacyMsg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgChangeParam:
			return handleMsgChangeParam(ctx, *msg, k)
		case *types.MsgDAOTransfer:
			return handleMsgDaoTransfer(ctx, *msg, k)
		case *types.MsgUpgrade:
			return handleMsgUpgrade(ctx, *msg, k)
		case types.MsgChangeParam:
			return handleMsgChangeParam(ctx, msg, k)
		case types.MsgDAOTransfer:
			return handleMsgDaoTransfer(ctx, msg, k)
		case types.MsgUpgrade:
			return handleMsgUpgrade(ctx, msg, k)
		default:
			errMsg := fmt.Sprintf("unrecognized gov message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgChangeParam(ctx sdk.Ctx, msg types.MsgChangeParam, k keeper.Keeper) sdk.Result {
	return k.ModifyParam(ctx, msg.ParamKey, msg.ParamVal, msg.FromAddress)
}

func handleMsgDaoTransfer(ctx sdk.Ctx, msg types.MsgDAOTransfer, k keeper.Keeper) sdk.Result {
	da, err := types.DAOActionFromString(msg.Action)
	if err != nil {
		return err.Result()
	}
	switch da {
	case types.DAOTransfer:
		return k.DAOTransferFrom(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	case types.DAOBurn:
		return k.DAOBurn(ctx, msg.FromAddress, msg.Amount)
	}
	return sdk.Result{}
}

func handleMsgUpgrade(ctx sdk.Ctx, msg types.MsgUpgrade, k keeper.Keeper) sdk.Result {
	return k.HandleUpgrade(ctx, types.NewACLKey(ModuleName, string(types.UpgradeKey)), msg.Upgrade, msg.Address)
}
