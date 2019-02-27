package _const

import (
	"crypto"
)

// TODO will change SHA to algorithm decided on by core team
const (
	// defines the session hashing algorithm
	SESSIONHASHINGALGORITHM  = crypto.SHA3_256
	VALIDATEHASHINGALGORITHM = crypto.SHA3_256
)
