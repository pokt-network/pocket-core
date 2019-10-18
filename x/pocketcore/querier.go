package pocketcore

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pokt-network/pocket-core/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QueryResResolve Queries Result Payload for a resolve query
type QueryResResolve struct {
	Value string `json:"value"`
}

// implement fmt.Stringer
func (r QueryResResolve) String() string {
	return r.Value
}

// QueryResNames Queries Result Payload for a names query
type QueryResNodes []types.Node

// query endpoints supported by the blockchain Querier
const (
	QueryNode         = "node"
	QueryNodes        = "nodes"
	QueryApplication  = "application"
	QueryApplications = "applications"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper PocketCoreKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryNode:
			return queryNode(ctx, path[1:], req, keeper.NodeKeeper)
		case QueryNodes:
			return queryNodes(ctx, req, keeper.NodeKeeper)
		case QueryApplication:
			return queryApplication(ctx, path[1:], req, keeper.ApplicationKeeper)
		case QueryApplications:
			return queryApplications(ctx, req, keeper.ApplicationKeeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

// nolint: unparam
func queryNode(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper nodeKeeper) ([]byte, sdk.Error) {
	whois := keeper.GetNode(ctx, path[0])

	res, err := codec.MarshalJSONIndent(keeper.cdc, whois)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryNodes(ctx sdk.Context, _ abci.RequestQuery, keeper nodeKeeper) ([]byte, sdk.Error) {
	appList, sdkError := keeper.GetAllNodes(ctx)

	if sdkError != nil {
		return nil, sdkError
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, appList)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryApplication(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper applicationKeeper) ([]byte, sdk.Error) {
	application := keeper.GetApplication(ctx, path[0])

	res, err := codec.MarshalJSONIndent(keeper.cdc, application)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryApplications(ctx sdk.Context, _ abci.RequestQuery, keeper applicationKeeper) ([]byte, sdk.Error) {
	appList, sdkError := keeper.GetAllApplications(ctx)

	if sdkError != nil {
		return nil, sdkError
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, appList)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
