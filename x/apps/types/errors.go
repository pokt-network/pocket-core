// nolint
package types

import (
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"strings"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace          sdk.CodespaceType = ModuleName
	CodeInvalidApplication    CodeType          = 101
	CodeInvalidInput          CodeType          = 103
	CodeApplicationJailed     CodeType          = 104
	CodeApplicationNotJailed  CodeType          = 105
	CodeMissingSelfDelegation CodeType          = 106
	CodeInvalidStatus         CodeType          = 110
	CodeMinimumStake          CodeType          = 111
	CodeNotEnoughCoins        CodeType          = 112
	CodeInvalidStakeAmount    CodeType          = 115
	CodeNoChains              CodeType          = 116
)

func ErrNoChains(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNoChains, "validator must stake with hosted blockchains")
}
func ErrNilApplicationAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "application address is nil")
}
func ErrApplicationStatus(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidStatus, "application status is not valid")
}
func ErrNoApplicationFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidApplication, "application does not exist for that address")
}

func ErrBadStakeAmount(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidStakeAmount, "the stake amount is invalid")
}

func ErrNotEnoughCoins(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNotEnoughCoins, "application does not have enough coins in their account")
}

func ErrMinimumStake(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMinimumStake, "application isn't staking above the minimum")
}

func ErrApplicationPubKeyExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidApplication, "application already exist for this pubkey, must use new application pubkey")
}

func ErrApplicationPubKeyTypeNotSupported(codespace sdk.CodespaceType, keyType string, supportedTypes []string) sdk.Error {
	msg := fmt.Sprintf("application pubkey type %s is not supported, must use %s", keyType, strings.Join(supportedTypes, ","))
	return sdk.NewError(codespace, CodeInvalidApplication, msg)
}

func ErrNoApplicationForAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidApplication, "that address is not associated with any known application")
}

func ErrBadApplicationAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidApplication, "application does not exist for that address")
}

func ErrApplicationJailed(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeApplicationJailed, "application still jailed, cannot yet be unjailed")
}

func ErrApplicationNotJailed(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeApplicationNotJailed, "application not jailed, cannot be unjailed")
}

func ErrMissingAppStake(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMissingSelfDelegation, "application has no stake; cannot be unjailed")
}

func ErrStakeTooLow(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeApplicationNotJailed, "application's self delegation less than min stake, cannot be unjailed")
}
