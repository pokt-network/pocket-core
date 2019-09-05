package types

import "github.com/pokt-network/pocket-core/crypto"

type AAT struct {
	Version    string     `json:"version"`
	AATMessage AATMessage `json:"aatMessage"`
	Signature  crypto.Signature     `json:"signature"`
}

type AATMessage struct {
	ApplicationPublicKey string `json:"ApplicaitonAddress"`
	ClientAddress      string `json:"ClientAddress"`
}
