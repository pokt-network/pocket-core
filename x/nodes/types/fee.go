package types

const (
	StakeFee   = 10000
	UnstakeFee = 10000
	UnjailFee  = 10000
	SendFee    = 10000
)

var (
	NodeFeeMap = map[string]int64{
		MsgStakeName:   StakeFee,
		MsgUnstakeName: UnstakeFee,
		MsgUnjailName:  UnjailFee,
		MsgSendName:    SendFee,
	}
)
