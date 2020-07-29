package types

import "github.com/pokt-network/pocket-core/types"

type FeeMultiplier struct {
	Key        string `json:"key"`
	Multiplier int64  `json:"multiplier"`
}

type FeeMultipliers struct {
	FeeMultis []FeeMultiplier `json:"fee_multiplier"`
	Default   int64           `json:"default"`
}

func (fm FeeMultipliers) GetFee(msg types.Msg) types.Int {
	for _, feeMultiplier := range fm.FeeMultis {
		if feeMultiplier.Key == msg.Type() {
			return msg.GetFee().Mul(types.NewInt(feeMultiplier.Multiplier))
		}
	}
	return msg.GetFee().Mul(types.NewInt(fm.Default))
}
