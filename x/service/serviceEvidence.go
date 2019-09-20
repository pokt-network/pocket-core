package service

// Evidence is a slice of type `ServiceCertificate`
// which are individual proofs of work completed
type Evidence []ServiceCertificate

// the header of the relay batch that is used
// to identify the relay batch in the global map
type EvidenceHeader struct {
	SessionHash       string
	ApplicationPubKey string
}

// add proof of work completed (type service certificate) to the evidence structure
func (e Evidence) AddEvidence(sc ServiceCertificate) error {
	if e == nil || len(e) == 0 {
		return EmptyEvidenceError
	}
	// if the increment counter is less than the evidence slice
	if len(e) < sc.Counter {
		return InvalidEvidenceSizeError
	}
	// if the evidence at index[increment counter] is not empty
	if e[sc.Counter].Signature != "" {
		return DuplicateEvidenceError
	}
	// set evidence at index[service certificate] = proof of work completed (Service Certificate)
	e[sc.Counter] = sc
	return nil
}
