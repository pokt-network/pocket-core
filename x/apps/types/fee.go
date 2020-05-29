package types

const (
	StakeFee   = 10000
	UnstakeFee = 10000
	UnjailFee  = 10000
)

var (
	AppFeeMap = map[string]int64{
		MsgAppStakeName:   StakeFee,
		MsgAppUnstakeName: UnstakeFee,
		MsgAppUnjailName:  UnjailFee,
	}
)
