package types

import (
	appexported "github.com/pokt-network/pocket-core/x/apps/exported"
	nodesexported "github.com/pokt-network/pocket-core/x/nodes/exported"
	sdk "github.com/pokt-network/posmint/types"
)

type PosKeeper interface {
	AwardCoinsTo(ctx sdk.Context, relays sdk.Int, address sdk.Address)
	GetStakedTokens(ctx sdk.Context) sdk.Int
	Validator(ctx sdk.Context, addr sdk.Address) nodesexported.ValidatorI
	TotalTokens(ctx sdk.Context) sdk.Int
	BurnValidator(ctx sdk.Context, address sdk.Address, severityPercentage sdk.Dec)
	JailValidator(ctx sdk.Context, addr sdk.Address)
	AllValidators(ctx sdk.Context) (validators []nodesexported.ValidatorI)
	GetStakedValidators(ctx sdk.Context) (validators []nodesexported.ValidatorI)
	SessionBlockFrequency(ctx sdk.Context) (res int64)
	StakeDenom(ctx sdk.Context) (res string)
}

type AppsKeeper interface {
	GetStakedTokens(ctx sdk.Context) sdk.Int
	Application(ctx sdk.Context, addr sdk.Address) appexported.ApplicationI
	AllApplications(ctx sdk.Context) (applications []appexported.ApplicationI)
	TotalTokens(ctx sdk.Context) sdk.Int
	BurnApplication(ctx sdk.Context, address sdk.Address, severityPercentage sdk.Dec)
	JailApplication(ctx sdk.Context, addr sdk.Address)
}
