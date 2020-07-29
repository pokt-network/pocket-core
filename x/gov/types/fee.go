package types

const (
	DAOTransferFee    = 10000
	MsgChangeParamFee = 10000
	MsgUpgradeFee     = 10000
)

var (
	GovFeeMap = map[string]int64{
		MsgDAOTransferName: DAOTransferFee,
		MsgChangeParamName: MsgChangeParamFee,
		MsgUpgradeName:     MsgUpgradeFee,
	}
)
