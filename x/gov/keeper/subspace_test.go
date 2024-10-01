package keeper

import (
	"testing"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	testSubSpaceName  = "testParamSet"
	testParamFieldKey = "field"
)

type testParamSet struct {
	m map[string]string
}

func (p *testParamSet) ParamSetPairs() sdk.ParamSetPairs {
	return sdk.ParamSetPairs{
		{Key: []byte(testParamFieldKey), Value: &p.m},
	}
}

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

func Test_ModifyParam_MapValue(t *testing.T) {
	subspace := sdk.NewSubspace(testSubSpaceName).WithKeyTable(
		sdk.NewKeyTable().RegisterParamSet(&testParamSet{}),
	)
	ctx, k := createTestKeeperAndContext(t, false, subspace)

	var (
		expectedMapValue = map[string]string{"2": "2", "3": "3", "1": "1"}

		// These strings represent the same value as `expectedMapValue` above
		mapValueVariations = []string{
			"{\"1\":\"1\", \"2\":\"2\", \"3\":\"3\"}",
			"{\"1\":\"1\", \"3\":\"3\", \"2\":\"2\"}",
			"{\"2\":\"2\", \"3\":\"3\", \"1\":\"1\"}",
			"{\"2\":\"2\", \"1\":\"1\", \"3\":\"3\"}",
			"{\"3\":\"3\", \"1\":\"1\", \"2\":\"2\"}",
			"{\"3\":\"3\", \"2\":\"2\", \"1\":\"1\"}",
		}

		// These strings are incompatible with testParamFieldKey
		invalidParamValues = []string{
			"{\"3\":3, \"1\":1, \"1\":13}", // wrong type
			"[\"123\":\"123\"]",            // wrong json
		}
	)

	aclKey := types.NewACLKey(testSubSpaceName, testParamFieldKey)
	aclOwner := k.GetACL(ctx).GetOwner(aclKey)

	for _, mapValue := range mapValueVariations {
		res := k.ModifyParam(ctx, aclKey, []byte(mapValue), aclOwner)
		assert.Zero(t, res.Code)

		s, ok := k.GetSubspace(testSubSpaceName)
		assert.True(t, ok)

		var value map[string]string
		s.Get(ctx, []byte(testParamFieldKey), &value)

		// Make sure the value is expected regardless of the order of fields in
		// `mapValue`.  Ideally we should verify the merkle root of the tree to
		// make sure binary representation of a param leaf is the same, but there
		// is no easy way to retrieve the merkle root from `ctx`.
		assert.True(t, sdk.CompareStringMaps(value, expectedMapValue))
	}

	for _, mapValue := range invalidParamValues {
		// Setting an invalid value doesn't fail
		res := k.ModifyParam(ctx, aclKey, []byte(mapValue), aclOwner)
		assert.Zero(t, res.Code)

		s, ok := k.GetSubspace(testSubSpaceName)
		assert.True(t, ok)

		// Value is not updated
		var value map[string]string
		s.Get(ctx, []byte(testParamFieldKey), &value)
		assert.True(t, sdk.CompareStringMaps(value, expectedMapValue))
	}
}
