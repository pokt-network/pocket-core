package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDAOAction_String(t *testing.T) {
	assert.Equal(t, DAOBurn.String(), DAOBurnString)
	assert.Equal(t, DAOTransfer.String(), DAOTransferString)
}

func TestDAOActionFromString(t *testing.T) {
	res, err := DAOActionFromString(DAOBurnString)
	assert.Nil(t, err)
	assert.Equal(t, DAOBurn, res)
	res, err = DAOActionFromString(DAOTransferString)
	assert.Nil(t, err)
	assert.Equal(t, DAOTransfer, res)
	_, err = DAOActionFromString("fake")
	assert.NotNil(t, err)
}
