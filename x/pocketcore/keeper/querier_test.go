package keeper

import (
	"testing"

	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQuerySupportedBlockchains(t *testing.T) {
	ctx, _, _, _, k, _, _ := createTestInput(t, false)
	p := types.Params{
		SupportedBlockchains: []string{"ethereum"},
	}
	k.SetParams(ctx, p)
	sbbz, err := querySupportedBlockchains(ctx, abci.RequestQuery{}, k)
	assert.Nil(t, err)
	var sb []string
	er := makeTestCodec().UnmarshalJSON(sbbz, &sb)
	assert.Nil(t, er)
	assert.Equal(t, sb, []string{"ethereum"})
}

func TestQueryParameters(t *testing.T) {
	ctx, _, _, _, k, _, _ := createTestInput(t, false)
	p := types.Params{
		SupportedBlockchains: []string{"ethereum"},
	}
	k.SetParams(ctx, p)
	sbbz, err := queryParameters(ctx, k)
	assert.Nil(t, err)
	var params types.Params
	er := makeTestCodec().UnmarshalJSON(sbbz, &params)
	assert.Nil(t, er)
	assert.Equal(t, params, p)
}
