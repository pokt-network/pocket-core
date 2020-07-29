package types

import (
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewExpiredProofsSubmissionError(t *testing.T) {
	assert.Equal(t, NewExpiredProofsSubmissionError(ModuleName), sdk.NewError(ModuleName, CodeExpiredProofsSubmissionError, ExpiredProofsSubmissionError.Error()))
}

func TestNewMerkleNodeNotFoundError(t *testing.T) {
	assert.Equal(t, NewMerkleNodeNotFoundError(ModuleName), sdk.NewError(ModuleName, CodeMerkleNodeNotFoundError, MerkleNodeNotFoundError.Error()))
}

func TestNewEmptyMerkleTreeError(t *testing.T) {
	assert.Equal(t, NewEmptyMerkleTreeError(ModuleName), sdk.NewError(ModuleName, CodeEmptyMerkleTreeError, EmptyMerkleTreeError.Error()))
}

func TestNewInvalidMerkleVerifyError(t *testing.T) {
	assert.Equal(t, NewInvalidMerkleVerifyError(ModuleName), sdk.NewError(ModuleName, CodeInvalidMerkleVerifyError, InvalidMerkleVerifyError.Error()))
}

func TestClaimNotFoundError(t *testing.T) {
	assert.Equal(t, NewClaimNotFoundError(ModuleName), sdk.NewError(ModuleName, CodeClaimNotFoundError, ClaimNotFoundError.Error()))
}

func TestCousinLeafEquivalentError(t *testing.T) {
	assert.Equal(t, NewCousinLeafEquivalentError(ModuleName), sdk.NewError(ModuleName, CodeCousinLeafEquivalentError, CousinLeafEquivalentError.Error()))
}

func TestInvalidLeafCousinProofsComboError(t *testing.T) {
	assert.Equal(t, NewInvalidLeafCousinProofsComboError(ModuleName), sdk.NewError(ModuleName, CodeInvalidLeafCousinProofsCombo, InvalidLeafCousinProofsCombo.Error()))
}

func TestInvalidRootError(t *testing.T) {
	assert.Equal(t, NewInvalidRootError(ModuleName), sdk.NewError(ModuleName, CodeInvalidRootError, InvalidRootError.Error()))
}

func TestInvalidHashLengthError(t *testing.T) {
	assert.Equal(t, NewInvalidHashLengthError(ModuleName), sdk.NewError(ModuleName, CodeInvalidHashLengthError, InvalidHashLengthError.Error()))
}

func TestInvalidAppPubKeyError(t *testing.T) {
	assert.Equal(t, NewInvalidAppPubKeyError(ModuleName), sdk.NewError(ModuleName, CodeInvalidAppPubKeyError, InvalidAppPubKeyError.Error()))
}
