package types

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
)

// Param module codespace constants
const (
	CodeUnknownSubspace         sdk.CodeType = 1
	CodeSettingParameter                     = 2
	CodeEmptyData                            = 3
	CodeInvalidACL                           = 4
	CodeUnauthorizedParamChange              = 5
	CodeSubspaceNotFound                     = 6
	CodeUnrecognizedDAOAction                = 7
	CodeZeroValueDAOAction                   = 8
	CodeZeroHeightUpgrade                    = 9
	CodeEmptyVersionUpgrade                  = 10
)

func ErrZeroHeightUpgrade(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeZeroHeightUpgrade, "the upgrade Height must not be zero")
}

func ErrEmptyVersionUpgrade(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyVersionUpgrade, "the upgrade Version string must not be empty")
}

func ErrZeroValueDAOAction(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeZeroValueDAOAction, "dao action value must not be zero: ")
}

func ErrUnrecognizedDAOAction(codespace sdk.CodespaceType, daoActionString string) sdk.Error {
	return sdk.NewError(codespace, CodeUnrecognizedDAOAction, "unrecognized dao action: "+daoActionString)
}

func ErrInvalidACL(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidACL, "invalid ACL: "+err.Error())
}

func ErrSubspaceNotFound(codespace sdk.CodespaceType, subspaceName string) sdk.Error {
	return sdk.NewError(codespace, CodeSubspaceNotFound, fmt.Sprintf("the subspace %s cannot be found", subspaceName))
}

func ErrUnauthorizedParamChange(codespace sdk.CodespaceType, owner sdk.Address, param string) sdk.Error {
	return sdk.NewError(codespace, CodeUnauthorizedParamChange,
		fmt.Sprintf("the param change is unathorized: Account %s does not have the permission to change param %s", owner, param))
}

// ErrUnknownSubspace returns an unknown subspace error.
func ErrUnknownSubspace(codespace sdk.CodespaceType, space string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownSubspace, fmt.Sprintf("unknown subspace %s", space))
}

// ErrSettingParameter returns an error for failing to set a parameter.
func ErrSettingParameter(codespace sdk.CodespaceType, key, subkey, value, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeSettingParameter, fmt.Sprintf("error setting parameter %s on %s (%s): %s", value, key, subkey, msg))
}

// ErrEmptyChanges returns an error for empty parameter changes.
func ErrEmptyChanges(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyData, "submitted parameter changes are empty")
}

// ErrEmptySubspace returns an error for an empty subspace.
func ErrEmptySubspace(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyData, "parameter subspace is empty")
}

// ErrEmptyKey returns an error for when an empty key is given.
func ErrEmptyKey(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyData, "parameter key is empty")
}

// ErrEmptyValue returns an error for when an empty key is given.
func ErrEmptyValue(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyData, "parameter value is empty")
}
