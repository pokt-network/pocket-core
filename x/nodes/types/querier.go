package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

// query endpoints supported by the staking Querier
const (
	QueryValidators          = "validators"
	QueryValidator           = "validator"
	QueryUnstakingValidators = "unstaking_validators"
	QueryStakedValidators    = "staked_validators"
	QueryUnstakedValidators  = "unstaked_validators"
	QueryStakedPool          = "stakedPool"
	QueryUnstakedPool        = "unstakedPool"
	QueryDAO                 = "dao"
	QueryParameters          = "parameters"
	QuerySigningInfo         = "signingInfo"
	QuerySigningInfos        = "signingInfos"
	QueryAccountBalance      = "account_balance"
)

type QueryValidatorParams struct {
	Address sdk.Address
}

func NewQueryValidatorParams(validatorAddr sdk.Address) QueryValidatorParams {
	return QueryValidatorParams{
		Address: validatorAddr,
	}
}

type QueryValidatorsParams struct {
	Page, Limit int
}

func NewQueryValidatorsParams(page, limit int) QueryValidatorsParams {
	return QueryValidatorsParams{page, limit}
}

type QueryAccountBalanceParams struct {
	sdk.Address
}

type QueryUnstakingValidatorsParams struct {
	Page, Limit int
}

func NewQueryUnstakingValidatorsParams(page, limit int) QueryUnstakingValidatorsParams {
	return QueryUnstakingValidatorsParams{page, limit}
}

type QueryStakedValidatorsParams struct {
	Page, Limit int
}

func NewQueryStakedValidatorsParams(page, limit int) QueryStakedValidatorsParams {
	return QueryStakedValidatorsParams{page, limit}
}

type QueryUnstakedValidatorsParams struct {
	Page, Limit int
}

func NewQueryUnstakedValidatorsParams(page, limit int) QueryUnstakedValidatorsParams {
	return QueryUnstakedValidatorsParams{page, limit}
}

// QuerySigningInfoParams defines the params for the following queries:
// - 'custom/slashing/signingInfo'
type QuerySigningInfoParams struct {
	Address sdk.Address
}

func NewQuerySigningInfoParams(consAddr sdk.Address) QuerySigningInfoParams {
	return QuerySigningInfoParams{consAddr}
}

// QuerySigningInfosParams defines the params for the following queries:
// - 'custom/slashing/signingInfos'
type QuerySigningInfosParams struct {
	Page, Limit int
}

func NewQuerySigningInfosParams(page, limit int) QuerySigningInfosParams {
	return QuerySigningInfosParams{page, limit}
}
