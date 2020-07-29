package types

import (
	sdk "github.com/pokt-network/pocket-core/types"
)

const (
	DAOAccountName              = "dao"
	DAOTransferString           = "dao_transfer"
	DAOBurnString               = "dao_burn"
	DAOTransfer       DAOAction = iota + 1
	DAOBurn
)

type DAOAction int

func (da DAOAction) String() string {
	switch da {
	case DAOTransfer:
		return DAOTransferString
	case DAOBurn:
		return DAOBurnString
	}
	return ""
}

func DAOActionFromString(s string) (DAOAction, sdk.Error) {
	switch s {
	case DAOTransferString:
		return DAOTransfer, nil
	case DAOBurnString:
		return DAOBurn, nil
	default:
		return 0, ErrUnrecognizedDAOAction(ModuleName, s)
	}
}
