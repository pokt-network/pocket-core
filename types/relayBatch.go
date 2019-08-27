package types

// TODO not complete -> security problems of double spend / replay attack (consider nonce)

type RelayBatchHeader struct {
	DevID       []byte
	Token       Token
	BlockNumber int // when the session started
}

type RelaySummary struct {
	RelayHash  []byte
	DigitalSig []byte
}

type RelayBatch struct {
	RelayBatchHeader RelayBatchHeader
	RelaySummary     []RelaySummary
}

func (rb *RelayBatch) AddDevID(devid []byte) {
	rb.RelayBatchHeader.DevID = devid
}

func (rb *RelayBatch) AddClientToken(token Token) {
	rb.RelayBatchHeader.Token = token
}

func (rb *RelayBatch) AddBlockNumber(blkNum int) {
	rb.RelayBatchHeader.BlockNumber = blkNum
}

func (rb *RelayBatch) AddRelaySummary(rs RelaySummary) {
	rb.RelaySummary = append(rb.RelaySummary, rs)
}

func (rb *RelayBatch) ErrorCheck() error {
	if rb.RelayBatchHeader.DevID == nil || len(rb.RelayBatchHeader.DevID) == 0 {
		return MissingDevidError
	}
	if rb.RelayBatchHeader.BlockNumber == 0 {
		return ZeroBlockError
	}
	if rb.RelayBatchHeader.Token.ExpDate == nil || len(rb.RelayBatchHeader.Token.ExpDate) == 0 { // TODO edit with token formalization
		return InvalidTokenError
	}
	if len(rb.RelaySummary) == 0 {
		return ZeroRelaySummaryError
	}
	return nil
}

func (rs *RelaySummary) AddRelayHash(relayHash []byte) {
	rs.RelayHash = relayHash
}

func (rs *RelaySummary) AddDigitalSignature(signature []byte) {
	rs.DigitalSig = signature
}

func (rs *RelaySummary) Validate(devID []byte) bool {
	// compares client public key to digital signature + relay hash
	return VerifySignature(devID, rs.RelayHash, rs.DigitalSig)
}

// TODO what is the best implementation of keeping track of the relay batches?
// Tied to a session?
// In a map?
// TODO what happens if validation of relay summary fails?
