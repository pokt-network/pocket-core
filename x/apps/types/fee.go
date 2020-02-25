package types

const (
	StakeFee   = 100000
	UnstakeFee = 100000
	UnjailFee  = 100000
)

var (
	AppFeeMap = map[string]int64{
		MsgAppStakeName:   StakeFee,
		MsgAppUnstakeName: UnstakeFee,
		MsgAppUnjailName:  UnjailFee,
	}
)
