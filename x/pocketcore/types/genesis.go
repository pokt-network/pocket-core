package types

import (
	"errors"
)

// "GenesisState" - The state of the module from the beginning
type GenesisState struct {
	Params   Params     `json:"params" yaml:"params"` // governance params
	Receipts []Receipt  `json:"receipts"`             // verified proofs
	Claims   []MsgClaim `json:"claims"`               // outstanding claims
}

// "ValidateGenesis" - Returns an error on an invalid genesis object
func ValidateGenesis(gs GenesisState) error {
	// validate the params
	err := gs.Params.Validate()
	if err != nil {
		return err
	}
	// validate receipts
	for _, reciept := range gs.Receipts {
		// validate the servicers address for the receipt
		if err := AddressVerification(reciept.ServicerAddress); err != nil {
			return err
		}
		// validate the header for the receipt
		if err := reciept.ValidateHeader(); err != nil {
			return err
		}
		// validate the total for the receipt
		if reciept.Total <= 0 {
			return errors.New("total relays for receipt is not positive")
		}
		// test byte conversion of evidence
		if reciept.EvidenceType != RelayEvidence && reciept.EvidenceType != ChallengeEvidence {
			return NewInvalidEvidenceErr(ModuleName)
		}
	}
	// validate each claim
	for _, claim := range gs.Claims {
		if err := claim.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}

// "DefaultGenesisState" - Returns the default genesis state for pocketcore module
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}
