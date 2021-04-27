package types

import (
	sdk "github.com/pokt-network/pocket-core/types"
)

// ensure ProtoMsg interface compliance at compile time
var (
	_ sdk.ProtoMsg = &MsgChangeParam{}
	_ sdk.ProtoMsg = &MsgDAOTransfer{}
	_ sdk.ProtoMsg = &MsgUpgrade{}
)

const (
	MsgDAOTransferName = "dao_tranfer"
	MsgChangeParamName = "change_param"
	MsgUpgradeName     = "upgrade"
)

//----------------------------------------------------------------------------------------------------------------------
// MsgChangeParam structure for changing governance parameters
// type MsgChangeParam struct {
// 	FromAddress sdk.Address `json:"address"`
// 	ParamKey    string      `json:"param_key"`
// 	ParamVal    []byte      `json:"param_value"`
// }

// Route provides router key for msg
func (msg MsgChangeParam) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgChangeParam) Type() string { return MsgChangeParamName }

// GetFee get fee for msg
func (msg MsgChangeParam) GetFee() sdk.BigInt {
	return sdk.NewInt(GovFeeMap[msg.Type()])
}

// GetSigner return address(es) that must sign over msg.GetSignBytes()
func (msg MsgChangeParam) GetSigner() sdk.Address {
	return msg.FromAddress
}

// GetSigner return address(es) that must sign over msg.GetSignBytes()
func (msg MsgChangeParam) GetRecipient() sdk.Address {
	return nil
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgChangeParam) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check
func (msg MsgChangeParam) ValidateBasic() sdk.Error {
	if msg.FromAddress == nil {
		return sdk.ErrInvalidAddress("nil address")
	}
	if msg.ParamKey == "" {
		return ErrEmptyKey(ModuleName)
	}
	if msg.ParamVal == nil {
		return ErrEmptyValue(ModuleName)
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------

// MsgDAOTransfer structure for changing governance parameters
// type MsgDAOTransfer struct {
// 	FromAddress sdk.Address `json:"from_address"`
// 	ToAddress   sdk.Address `json:"to_address"`
// 	Amount      sdk.BigInt     `json:"amount"`
// 	Action      string      `json:"action"`
// }

// Route provides router key for msg
func (msg MsgDAOTransfer) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgDAOTransfer) Type() string { return MsgDAOTransferName }

// GetFee get fee for msg
func (msg MsgDAOTransfer) GetFee() sdk.BigInt {
	return sdk.NewInt(GovFeeMap[msg.Type()])
}

// GetSigner return address(es) that must sign over msg.GetSignBytes()
func (msg MsgDAOTransfer) GetSigner() sdk.Address {
	return msg.FromAddress
}

// GetSigner return address(es) that must sign over msg.GetSignBytes()
func (msg MsgDAOTransfer) GetRecipient() sdk.Address {
	return nil
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgDAOTransfer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check
func (msg MsgDAOTransfer) ValidateBasic() sdk.Error {
	if msg.FromAddress == nil {
		return sdk.ErrInvalidAddress("nil from address")
	}
	if msg.Amount.Int64() == 0 {
		return ErrZeroValueDAOAction(ModuleName)
	}
	daoAction, err := DAOActionFromString(msg.Action)
	if err != nil {
		return err
	}
	if daoAction == DAOTransfer && msg.ToAddress == nil {
		return sdk.ErrInvalidAddress("nil to address")
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------

// MsgUpgrade structure for changing governance parameters
// type MsgUpgrade struct {
// 	Address sdk.Address `json:"address"`
// 	Upgrade Upgrade     `json:"upgrade"`
// }

// Route provides router key for msg
func (msg MsgUpgrade) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgUpgrade) Type() string { return MsgUpgradeName }

// GetFee get fee for msg
func (msg MsgUpgrade) GetFee() sdk.BigInt {
	return sdk.NewInt(GovFeeMap[msg.Type()])
}

// GetSigner return address(es) that must sign over msg.GetSignBytes()
func (msg MsgUpgrade) GetSigner() sdk.Address {
	return msg.Address
}

// GetSigner return address(es) that must sign over msg.GetSignBytes()
func (msg MsgUpgrade) GetRecipient() sdk.Address {
	return nil
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgUpgrade) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check
func (msg MsgUpgrade) ValidateBasic() sdk.Error {
	if msg.Address == nil {
		return sdk.ErrInvalidAddress("nil from address")
	}
	if msg.Upgrade.UpgradeHeight() == 0 {
		return ErrZeroHeightUpgrade(ModuleName)
	}
	if msg.Upgrade.UpgradeVersion() == "" {
		return ErrZeroHeightUpgrade(ModuleName)
	}
	return nil
}
