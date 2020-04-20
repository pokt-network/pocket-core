package types

const (
	ClaimFee = 100000 // fee for claim message (in uPOKT)
	ProofFee = 100000 // fee for proof message (in uPOKT)
)

var (
	// map of message name to fee value
	PocketFeeMap = map[string]int64{
		MsgClaimName: ClaimFee,
		MsgProofName: ProofFee,
	}
)
