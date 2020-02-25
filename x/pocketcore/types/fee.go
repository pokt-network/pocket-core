package types

const (
	ClaimFee = 100000
	ProofFee = 100000
)

var (
	PocketFeeMap = map[string]int64{
		MsgClaimName: ClaimFee,
		MsgProofName: ProofFee,
	}
)
