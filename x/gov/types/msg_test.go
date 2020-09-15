package types

import (
	"testing"

	"github.com/pokt-network/pocket-core/types"
	"github.com/stretchr/testify/assert"
)

func TestMsgChangeParam_ValidateBasic(t *testing.T) {
	cdc := makeTestCodec()
	bytes, _ := cdc.MarshalJSON(false)
	m := MsgChangeParam{
		FromAddress: getRandomValidatorAddress(),
		ParamKey:    "bank/sendenabled",
		ParamVal:    bytes,
	}
	assert.Nil(t, m.ValidateBasic())
	m = MsgChangeParam{
		FromAddress: getRandomValidatorAddress(),
		ParamKey:    "",
		ParamVal:    bytes,
	}
	assert.NotNil(t, m.ValidateBasic())
	m = MsgChangeParam{
		ParamKey: "bank/sendenabled",
		ParamVal: bytes,
	}
	assert.NotNil(t, m.ValidateBasic())
	m = MsgChangeParam{
		FromAddress: getRandomValidatorAddress(),
		ParamKey:    "bank/sendenabled",
	}
	assert.NotNil(t, m.ValidateBasic())
}

func TestAminoPrimitive(t *testing.T) {
	cdc := makeTestCodec()
	bytesbool, _ := cdc.MarshalJSON(false)
	bytesint, _ := cdc.MarshalJSON(int64(23))
	assert.NotNil(t, bytesbool)
	assert.NotNil(t, bytesint)
	var b bool
	var i int64
	err := cdc.UnmarshalJSON(bytesbool, &b)
	assert.Nil(t, err)

	err = cdc.UnmarshalJSON(bytesint, &i)
	assert.Nil(t, err)
}

func TestMsgDAOTransfer_ValidateBasic(t *testing.T) {
	m := MsgDAOTransfer{
		FromAddress: getRandomValidatorAddress(),
		ToAddress:   getRandomValidatorAddress(),
		Amount:      types.OneInt(),
		Action:      DAOTransferString,
	}
	assert.Nil(t, m.ValidateBasic())
	m = MsgDAOTransfer{
		FromAddress: getRandomValidatorAddress(),
		ToAddress:   getRandomValidatorAddress(),
		Amount:      types.OneInt(),
	}
	assert.NotNil(t, m.ValidateBasic())
	m = MsgDAOTransfer{
		FromAddress: getRandomValidatorAddress(),
		ToAddress:   getRandomValidatorAddress(),
		Amount:      types.ZeroInt(),
		Action:      DAOTransferString,
	}
	assert.NotNil(t, m.ValidateBasic())
	m = MsgDAOTransfer{
		FromAddress: getRandomValidatorAddress(),
		Amount:      types.ZeroInt(),
		Action:      DAOTransferString,
	}
	assert.NotNil(t, m.ValidateBasic())
	m = MsgDAOTransfer{
		ToAddress: getRandomValidatorAddress(),
		Amount:    types.ZeroInt(),
		Action:    DAOTransferString,
	}
	assert.NotNil(t, m.ValidateBasic())
}

func TestMsgUpgrade_ValidateBasic(t *testing.T) {
	m := MsgUpgrade{
		Address: getRandomValidatorAddress(),
		Upgrade: Upgrade{
			Height:  100,
			Version: "2.0.0",
		},
	}
	assert.Nil(t, m.ValidateBasic())
	m = MsgUpgrade{
		Address: getRandomValidatorAddress(),
		Upgrade: Upgrade{
			Height:  0,
			Version: "2.0.0",
		},
	}
	assert.NotNil(t, m.ValidateBasic())
	m = MsgUpgrade{
		Address: getRandomValidatorAddress(),
		Upgrade: Upgrade{
			Height:  100,
			Version: "",
		},
	}
	assert.NotNil(t, m.ValidateBasic())
}
