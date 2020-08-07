package types

import (
	"fmt"
	"github.com/pokt-network/pocket-core/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmtypes "github.com/tendermint/tendermint/types"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

var codespace = types.CodespaceType(ModuleName)

func TestErrBadDelegationAmount(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Bad Delegation Amount", args{codespace: codespace}, types.NewError(codespace, CodeInvalidDelegation, "amount must be > 0")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrBadDelegationAmount(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBadDelegationAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBadDenom(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Bad coin Denom", args{codespace: codespace}, types.NewError(codespace, CodeInvalidDelegation, "invalid coin denomination")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrBadDenom(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBadDenom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBadSendAmount(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Bad Send Amount", args{codespace: codespace}, types.NewError(codespace, CodeBadSend, "the amount to send must be positive")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrBadSendAmount(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrBadSendAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrBadValidatorAddr(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Bad Validator Address", args{codespace: codespace}, types.NewError(codespace, CodeInvalidValidator, "validator does not exist for that address")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrNoValidatorFound(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrNoValidatorFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrCantHandleEvidence(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Can't Store Evidence", args{codespace: codespace}, types.NewError(codespace, CodeCantHandleEvidence, "Warning: the DS evidence is unable to be handled")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrCantHandleEvidence(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrCantHandleEvidence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrMinimumStake(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Minimun Stake", args{codespace: codespace}, types.NewError(codespace, CodeMinimumStake, "validator isn't staking above the minimum")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrMinimumStake(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrMinimumStake() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrMissingSelfDelegation(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Missing Self Delegation", args{codespace: codespace}, types.NewError(codespace, CodeMissingSelfDelegation, "validator has no self-delegation; cannot be unjailed")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrMissingSelfDelegation(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrMissingSelfDelegation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrNilValidatorAddr(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Nil Validator Address", args{codespace: codespace}, types.NewError(codespace, CodeInvalidInput, "validator address is nil")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrNilValidatorAddr(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrNilValidatorAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrNoChains(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"No Chains", args{codespace: codespace}, types.NewError(codespace, CodeNoChains, "validator must stake with hosted blockchains")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrNoChains(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrNoChains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrNoServiceURL(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"No Service URL", args{codespace: codespace}, types.NewError(codespace, CodeNoServiceURL, "validator must stake with a serviceurl")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrNoServiceURL(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrNoServiceURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrNoSigningInfoFound(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
		consAddr  types.Address
	}
	var pub ed25519.PubKeyEd25519
	_, err := rand.Read(pub[:])
	if err != nil {
		t.Fatalf(err.Error())
	}
	ca := types.Address(pub.Address())

	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"No Signing Info Found", args{codespace: codespace, consAddr: ca}, types.NewError(codespace, CodeMissingSigningInfo, fmt.Sprintf("no signing info found for address: %s", ca))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrNoSigningInfoFound(tt.args.codespace, tt.args.consAddr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrNoSigningInfoFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrNoValidatorForAddress(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"No Validator For Address", args{codespace: codespace}, types.NewError(codespace, CodeInvalidValidator, "that address is not associated with any known validator")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrNoValidatorForAddress(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrNoValidatorForAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrNoValidatorFound(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"No Validator Found", args{codespace: codespace}, types.NewError(codespace, CodeInvalidValidator, "validator does not exist for that address")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrNoValidatorFound(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrNoValidatorFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrNotEnoughCoins(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Not Enough Coins", args{codespace: codespace}, types.NewError(codespace, CodeNotEnoughCoins, "validator does not have enough coins in their account")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrNotEnoughCoins(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrNotEnoughCoins() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrSelfDelegationTooLowToUnjail(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Self Delegation Too Low To Unjail", args{codespace: codespace}, types.NewError(codespace, CodeValidatorNotJailed, "validator's self delegation less than MinSelfDelegation, cannot be unjailed")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrSelfDelegationTooLowToUnjail(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrSelfDelegationTooLowToUnjail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrValidatorJailed(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Validator Jailed", args{codespace: codespace}, types.NewError(codespace, CodeValidatorJailed, "validator jailed")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrValidatorJailed(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrValidatorJailed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrValidatorNotJailed(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Validator Not Jailed", args{codespace: codespace}, types.NewError(codespace, CodeValidatorNotJailed, "validator not jailed, cannot be unjailed")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrValidatorNotJailed(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrValidatorNotJailed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrValidatorPubKeyExists(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Validator Public Key Exist", args{codespace: codespace}, types.NewError(codespace, CodeInvalidValidator, "validator already exist for this pubkey, must use new validator pubkey")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrValidatorPubKeyExists(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrValidatorPubKeyExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrValidatorPubKeyTypeNotSupported(t *testing.T) {
	type args struct {
		codespace      types.CodespaceType
		keyType        string
		supportedTypes []string
	}

	keyType := "secp256k1"
	keyTypes := []string{tmtypes.ABCIPubKeyTypeEd25519}
	msg := fmt.Sprintf("validator pubkey type %s is not supported, must use %s", keyType, strings.Join(keyTypes, ","))

	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Validator Public Key type Not Supported", args{codespace: codespace, keyType: keyType, supportedTypes: keyTypes}, types.NewError(codespace, CodeInvalidValidator, msg)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrValidatorPubKeyTypeNotSupported(tt.args.codespace, tt.args.keyType, tt.args.supportedTypes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrValidatorPubKeyTypeNotSupported() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrValidatorStatus(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Validator Status", args{codespace: codespace}, types.NewError(codespace, CodeInvalidStatus, "validator status is not valid")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrValidatorStatus(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrValidatorStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrValidatorTombstoned(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Validator Already Tombstoned", args{codespace: codespace}, types.NewError(codespace, CodeValidatorTombstoned, "Warning: validator is already tombstoned")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrValidatorTombstoned(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrValidatorTombstoned() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrValidatorWaitingToUnstake(t *testing.T) {
	type args struct {
		codespace types.CodespaceType
	}
	tests := []struct {
		name string
		args args
		want types.Error
	}{
		{"Validator Waiting To Unstake", args{codespace: codespace}, types.NewError(codespace, CodeWaitingValidator, "validator is currently waiting to unstake")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrValidatorWaitingToUnstake(tt.args.codespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrValidatorWaitingToUnstake() = %v, want %v", got, tt.want)
			}
		})
	}
}
