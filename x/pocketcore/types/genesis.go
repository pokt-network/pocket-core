package types

import (
	"errors"
)

type GenesisState struct {
	Params Params     `json:"params" yaml:"params"`
	Proofs []Receipt  `json:"proofs"`
	Claims []MsgClaim `json:"claims"`
}

func ValidateGenesis(data GenesisState) error {
	err := data.Params.Validate()
	if err != nil {
		return err
	}
	// validate proofsMap
	for _, proof := range data.Proofs {
		if err := AddressVerification(proof.ServicerAddress); err != nil {
			return err
		}
		if err := proof.ValidateHeader(); err != nil {
			return err
		}
		if proof.Total <= 0 {
			return errors.New("total relays for Proof is negative")
		}
	}
	// validate each claim
	for _, claim := range data.Claims {
		if err := claim.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}
