package types

import (
	appexported "github.com/pokt-network/pocket-core/x/apps/exported"
	posexported "github.com/pokt-network/pocket-core/x/nodes/exported"
	sdk "github.com/pokt-network/posmint/types"
)

type PosKeeper interface {
	AwardCoinsTo(ctx sdk.Context, amount sdk.Int, address sdk.ValAddress)
	GetStakedTokens(ctx sdk.Context) sdk.Int
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator posexported.ValidatorI, found bool)
	TotalTokens(ctx sdk.Context) sdk.Int
	BurnValidator(ctx sdk.Context, address sdk.ValAddress, severityPercentage sdk.Dec)
	JailValidator(ctx sdk.Context, addr sdk.ConsAddress)
	GetAllValidators(ctx sdk.Context) (validators []posexported.ValidatorI)
	SessionBlockFrequency(ctx sdk.Context) (res int64)
	StakeDenom(ctx sdk.Context) (res string)
}

type AppsKeeper interface {
	GetStakedTokens(ctx sdk.Context) sdk.Int
	GetApplication(ctx sdk.Context, addr sdk.ValAddress) (application appexported.ApplicationI, found bool)
	GetAllApplications(ctx sdk.Context) (applications []appexported.ApplicationI)
	TotalTokens(ctx sdk.Context) sdk.Int
	BurnApplication(ctx sdk.Context, address sdk.ValAddress, severityPercentage sdk.Dec)
	JailApplication(ctx sdk.Context, addr sdk.ConsAddress)
	SetRelayCoefficient(ctx sdk.Context, newCoefficient int)
}
