package types

import "github.com/pokt-network/pocket-core/types"

func (fm FeeMultipliers) GetFee(msg types.Msg) types.BigInt {
	for _, feeMultiplier := range fm.FeeMultis {
		if feeMultiplier.Key == msg.Type() {
			return msg.GetFee().Mul(types.NewInt(feeMultiplier.Multiplier))
		}
	}
	return msg.GetFee().Mul(types.NewInt(fm.Default))
}
