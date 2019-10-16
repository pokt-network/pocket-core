package blockchain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pokt-network/pocket-core/tests/fixtures"
	"github.com/pokt-network/pocket-core/types"
)

func GetLatestBlockID(ctx sdk.Context) types.BlockID {
	// return fixtures.GenerateBlockHash()
	header := ctx.BlockHeader()
	return types.BlockID(header.GetLastBlockId())
}

func GetLatestSessionBlockID(ctx sdk.Context) types.BlockID {
	//return fixtures.GenerateBlockHash()
	latestsessionBlockHeight := GetLatestSessionBlockHeight(ctx)
	ctxAtHeight := ctx.WithBlockHeight(latestsessionBlockHeight)
	return GetLatestBlockID(ctxAtHeight)
}

func GetLatestSessionBlockHeight(ctx sdk.Context) int64 {
	//return fixtures.GenerateBlockHash()
	blkHeight := ctx.BlockHeight()
	return (blkHeight / SESSIONBLOCKFREQUENCY) * SESSIONBLOCKFREQUENCY
}

func GetNodes() (*types.Nodes, error) { // this is essentially -> dispatchPeers()
	// todo
	return fixtures.GetNodes()
}

func GetApplications() (*types.Applications, error) {
	// todo
	return fixtures.GetApplications()
}

func GetMaxNumberOfRelaysForApp(applicationPubKey string) int {
	// todo
	return 5000
}
