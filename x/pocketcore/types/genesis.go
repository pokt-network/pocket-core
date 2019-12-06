package types

type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
}

func ValidateGenesis(data GenesisState) error {
	return data.Params.Validate()
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}
