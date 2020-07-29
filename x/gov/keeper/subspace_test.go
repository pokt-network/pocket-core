package keeper

import (
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestModifyParam(t *testing.T) {
	addr := getRandomValidatorAddress()
	var aclKey = types.NewACLKey(types.ModuleName, string(types.DAOOwnerKey))
	ctx, k := createTestKeeperAndContext(t, false)
	jbyte, _ := amino.MarshalJSON(addr)
	res := k.ModifyParam(ctx, aclKey, jbyte, k.GetACL(ctx).GetOwner(aclKey))
	assert.Zero(t, res.Code)
	s, ok := k.GetSubspace(types.DefaultParamspace)
	assert.True(t, ok)
	var b sdk.Address
	s.Get(ctx, []byte("daoOwner"), &b)
	assert.Equal(t, addr, b)
	// Test "message.sender" event emission
	assert.Equal(
		t,
		true,
		ContainsEvent(
			res.Events,
			abci.Event(
				sdk.NewEvent(
					sdk.EventTypeMessage,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
					sdk.NewAttribute(sdk.AttributeKeySender, k.GetACL(ctx).GetOwner(aclKey).String()),
				),
			),
		),
	)
}
