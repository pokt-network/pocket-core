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
	CodeInvalidValidator      CodeType          = 101
	CodeInvalidDelegation     CodeType          = 102
	CodeInvalidInput          CodeType          = 103
	CodeValidatorJailed       CodeType          = 104
	CodeValidatorNotJailed    CodeType          = 105
	CodeMissingSelfDelegation CodeType          = 106
	CodeMissingSigningInfo    CodeType          = 108
	CodeBadSend               CodeType          = 109
	CodeInvalidStatus         CodeType          = 110
	CodeMinimumStake          CodeType          = 111
	CodeNotEnoughCoins        CodeType          = 112
	CodeValidatorTombstoned   CodeType          = 113
	CodeCantHandleEvidence    CodeType          = 114
	CodeNoChains              CodeType          = 115
	CodeNoServiceURL          CodeType          = 116
	CodeWaitingValidator      CodeType          = 117
)

func ErrValidatorWaitingToUnstake(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeWaitingValidator, "validator is currently waiting to unstake")
}

func ErrNoServiceURL(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNoServiceURL, "validator must stake with a serviceurl")
}
func ErrNoChains(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNoChains, "validator must stake with hosted blockchains")
}
func ErrNilValidatorAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "validator address is nil")
}
func ErrValidatorStatus(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidStatus, "validator status is not valid")
}
func ErrNoValidatorFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "validator does not exist for that address")
}

func ErrNotEnoughCoins(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNotEnoughCoins, "validator does not have enough coins in their account")
}

func ErrValidatorTombstoned(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeValidatorTombstoned, "Warning: validator is already tombstoned")
}

func ErrCantHandleEvidence(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeCantHandleEvidence, "Warning: the DS evidence is unable to be handled")
}

func ErrMinimumStake(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMinimumStake, "validator isn't staking above the minimum")
}

func ErrValidatorPubKeyExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "validator already exist for this pubkey, must use new validator pubkey")
}

func ErrValidatorPubKeyTypeNotSupported(codespace sdk.CodespaceType, keyType string, supportedTypes []string) sdk.Error {
	msg := fmt.Sprintf("validator pubkey type %s is not supported, must use %s", keyType, strings.Join(supportedTypes, ","))
	return sdk.NewError(codespace, CodeInvalidValidator, msg)
}

func ErrBadSendAmount(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeBadSend, "the amount to send must be positive")
}

func ErrBadDenom(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation, "invalid coin denomination")
}

func ErrBadDelegationAmount(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation, "amount must be > 0")
}

func ErrNoValidatorForAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "that address is not associated with any known validator")
}

func ErrValidatorJailed(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeValidatorJailed, "validator still jailed, cannot yet be unjailed")
}

func ErrValidatorNotJailed(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeValidatorNotJailed, "validator not jailed, cannot be unjailed")
}

func ErrMissingSelfDelegation(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMissingSelfDelegation, "validator has no self-delegation; cannot be unjailed")
}

func ErrSelfDelegationTooLowToUnjail(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeValidatorNotJailed, "validator's self delegation less than MinSelfDelegation, cannot be unjailed")
}

func ErrNoSigningInfoFound(codespace sdk.CodespaceType, consAddr sdk.Address) sdk.Error {
	return sdk.NewError(codespace, CodeMissingSigningInfo, fmt.Sprintf("no signing info found for address: %s", consAddr))
}
