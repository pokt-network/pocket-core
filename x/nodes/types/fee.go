package types

const (
	StakeFee   = 100000
	UnstakeFee = 100000
	UnjailFee  = 100000
	SendFee    = 100000
)

var (
	NodeFeeMap = map[string]int64{
		MsgStakeName:   StakeFee,
		MsgUnstakeName: UnstakeFee,
		MsgUnjailName:  UnjailFee,
		MsgSendName:    SendFee,
	}
)
