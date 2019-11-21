package types

const (
	SUPPORTEDTOKENVERSIONS = "0.0.1" // todo
)

type AAT struct {
	Version    AATVersion `json:"version"`
	AATMessage AATMessage `json:"aatMessage"`
	Signature  string     `json:"signature"`
}

type AATMessage struct {
	ApplicationPublicKey string `json:"ApplicaitonAddress"`
	ClientPublicKey      string `json:"ClientPublicKey"`
}

type AATVersion string

func (av AATVersion) IsIncluded() bool {
	if av == "" {
		return false
	}
	return true
}

func (ac AATVersion) IsSupported() bool {
	if ac == SUPPORTEDTOKENVERSIONS {
		return true
	}
	return false
}
