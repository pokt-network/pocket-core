package types

import (
	sdk "github.com/pokt-network/pocket-core/types"
)

// query endpoints supported by the staking Querier
const (
	QueryApplications    = "applications"
	QueryApplication     = "application"
	QueryAppStakedPool   = "appStakedPool"
	QueryAppUnstakedPool = "appUnstakedPool"
	QueryParameters      = "parameters"
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

type QueryApplicationsWithOpts struct {
	Page          int             `json:"page"`
	Limit         int             `json:"per_page"`
	StakingStatus sdk.StakeStatus `json:"staking_status"`
	Blockchain    string          `json:"blockchain"`
}

func (opts QueryApplicationsWithOpts) IsValid(app Application) bool {
	if opts.StakingStatus != 0 {
		if opts.StakingStatus != app.Status {
			return false
		}
	}
	if opts.Blockchain != "" {
		var contains bool
		for _, chain := range app.Chains {
			if chain == opts.Blockchain {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}
	return true
}

type QueryStakedApplicationsParams struct {
	Page, Limit int
}
