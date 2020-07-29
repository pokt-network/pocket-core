package types

import (
	appexported "github.com/pokt-network/pocket-core/x/apps/exported"
	nodesexported "github.com/pokt-network/pocket-core/x/nodes/exported"
	sdk "github.com/pokt-network/pocket-core/types"
)

type PosKeeper interface {
	RewardForRelays(ctx sdk.Ctx, relays sdk.Int, address sdk.Address)
	GetStakedTokens(ctx sdk.Ctx) sdk.Int
	Validator(ctx sdk.Ctx, addr sdk.Address) nodesexported.ValidatorI
	TotalTokens(ctx sdk.Ctx) sdk.Int
	BurnForChallenge(ctx sdk.Ctx, challenges sdk.Int, address sdk.Address)
	JailValidator(ctx sdk.Ctx, addr sdk.Address)
	AllValidators(ctx sdk.Ctx) (validators []nodesexported.ValidatorI)
	GetStakedValidators(ctx sdk.Ctx) (validators []nodesexported.ValidatorI)
	BlocksPerSession(ctx sdk.Ctx) (res int64)
	StakeDenom(ctx sdk.Ctx) (res string)
	GetValidatorsByChain(ctx sdk.Ctx, networkID string) (validators []nodesexported.ValidatorI)
}

type AppsKeeper interface {
	GetStakedTokens(ctx sdk.Ctx) sdk.Int
	Application(ctx sdk.Ctx, addr sdk.Address) appexported.ApplicationI
	AllApplications(ctx sdk.Ctx) (applications []appexported.ApplicationI)
	TotalTokens(ctx sdk.Ctx) sdk.Int
	JailApplication(ctx sdk.Ctx, addr sdk.Address)
}

type AuthKeeper interface {
	GetFee(ctx sdk.Ctx, msg sdk.Msg) sdk.Int
}
