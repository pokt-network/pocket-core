package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

// query endpoints supported by the staking Querier
const (
	QueryApplications          = "applications"
	QueryApplication           = "application"
	QueryUnstakingApplications = "unstaking_applications"
	QueryStakedApplications    = "staked_applications"
	QueryUnstakedApplications  = "unstaked_applications"
	QueryAppStakedPool         = "appStakedPool"
	QueryAppUnstakedPool       = "appUnstakedPool"
	QueryParameters            = "parameters"
)

type QueryAppParams struct {
	Address sdk.Address
}

func NewQueryAppParams(applicationAddr sdk.Address) QueryAppParams {
	return QueryAppParams{
		Address: applicationAddr,
	}
}

type QueryAppsParams struct {
	Page, Limit int
}

func NewQueryApplicationsParams(page, limit int) QueryAppsParams {
	return QueryAppsParams{page, limit}
}

type QueryUnstakingApplicationsParams struct {
	Page, Limit int
}

func NewQueryUnstakingApplicationsParams(page, limit int) QueryUnstakingApplicationsParams {
	return QueryUnstakingApplicationsParams{page, limit}
}

type QueryStakedApplicationsParams struct {
	Page, Limit int
}

func NewQueryStakedApplicationsParams(page, limit int) QueryStakedApplicationsParams {
	return QueryStakedApplicationsParams{page, limit}
}

type QueryUnstakedApplicationsParams struct {
	Page, Limit int
}

func NewQueryUnstakedApplicationsParams(page, limit int) QueryUnstakedApplicationsParams {
	return QueryUnstakedApplicationsParams{page, limit}
}
