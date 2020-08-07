package types

// "GenesisState" - The state of the module from the beginning
type GenesisState struct {
	Params Params     `json:"params" yaml:"params"` // governance params
	Claims []MsgClaim `json:"claims"`               // outstanding claims
}

// "ValidateGenesis" - Returns an error on an invalid genesis object
func ValidateGenesis(gs GenesisState) error {
	// validate the params
	err := gs.Params.Validate()
	if err != nil {
		return err
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
