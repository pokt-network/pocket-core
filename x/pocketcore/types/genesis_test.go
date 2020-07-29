package types

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/pokt-network/pocket-core/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateGenesis(t *testing.T) {
	appPubKeyClaim := getRandomPubKey().RawString()
	pk := getRandomPubKey()
	servicerAddr := pk.Address()
	nn := hex.EncodeToString([]byte{01})
	rootHash := Hash([]byte("fakeRoot"))
	root := HashRange{
		Hash:  rootHash,
		Range: Range{0, 100},
	}
	invalidParams := GenesisState{
		Params: Params{
			SessionNodeCount:      0,
			ClaimSubmissionWindow: 0,
			SupportedBlockchains:  nil,
			ClaimExpiration:       0,
		},
		Claims: []MsgClaim{{
			SessionHeader: SessionHeader{
				ApplicationPubKey:  appPubKeyClaim,
				Chain:              nn,
				SessionBlockHeight: 1,
			},
			MerkleRoot:   root,
			TotalProofs:  1000,
			FromAddress:  types.Address(servicerAddr),
			EvidenceType: RelayEvidence,
		}},
	}
	invalidClaims := GenesisState{
		Params: Params{
			SessionNodeCount:      1,
			ClaimSubmissionWindow: 5,
			SupportedBlockchains:  []string{nn},
			ClaimExpiration:       50,
		},
		Claims: []MsgClaim{{
			SessionHeader: SessionHeader{
				ApplicationPubKey:  appPubKeyClaim,
				Chain:              nn,
				SessionBlockHeight: 1,
			},
			MerkleRoot:   root,
			TotalProofs:  -1000,
			FromAddress:  types.Address(servicerAddr),
			EvidenceType: RelayEvidence,
		}},
	}
	validGenesisState := GenesisState{
		Params: Params{
			SessionNodeCount:      1,
			ClaimSubmissionWindow: 5,
			SupportedBlockchains:  []string{nn},
			ClaimExpiration:       50,
		},
		Claims: []MsgClaim{{
			SessionHeader: SessionHeader{
				ApplicationPubKey:  appPubKeyClaim,
				Chain:              nn,
				SessionBlockHeight: 1,
			},
			MerkleRoot:   root,
			TotalProofs:  1000,
			FromAddress:  types.Address(servicerAddr),
			EvidenceType: RelayEvidence,
		}},
	}
	tests := []struct {
		name         string
		genesisState GenesisState
		hasError     bool
	}{
		{
			name:         "Bad params",
			genesisState: invalidParams,
			hasError:     true,
		},
		{
			name:         "Bad claims",
			genesisState: invalidClaims,
			hasError:     true,
		},
		{
			name:         "Valid genesis state",
			genesisState: validGenesisState,
			hasError:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, ValidateGenesis(tt.genesisState) != nil, tt.hasError)
		})
	}
}

func TestDefaultGenesisState(t *testing.T) {
	appPubKeyClaim := getRandomPubKey().RawString()
	pk := getRandomPubKey()
	servicerAddr := pk.Address()
	nn := hex.EncodeToString([]byte{01})
	rootHash := Hash([]byte("fakeRoot"))
	root := HashRange{
		Hash:  rootHash,
		Range: Range{0, 100},
	}
	validGenesisState := GenesisState{
		Params: Params{
			SessionNodeCount:      1,
			ClaimSubmissionWindow: 5,
			SupportedBlockchains:  []string{nn},
			ClaimExpiration:       50,
		},
		Claims: []MsgClaim{{
			SessionHeader: SessionHeader{
				ApplicationPubKey:  appPubKeyClaim,
				Chain:              nn,
				SessionBlockHeight: 1,
			},
			MerkleRoot:  root,
			TotalProofs: 1000,
			FromAddress: types.Address(servicerAddr),
		}},
	}
	DefaultGenState := GenesisState{Params: Params{
		SessionNodeCount:           DefaultSessionNodeCount,
		ClaimSubmissionWindow:      DefaultClaimSubmissionWindow,
		SupportedBlockchains:       DefaultSupportedBlockchains,
		ClaimExpiration:            DefaultClaimExpiration,
		ReplayAttackBurnMultiplier: DefaultReplayAttackBurnMultiplier,
		MinimumNumberOfProofs:      DefaultMinimumNumberOfProofs,
	}}
	tests := []struct {
		name         string
		genesisState GenesisState
		isEqual      bool
	}{
		{
			name:         "Valid genesis state, but not default",
			genesisState: validGenesisState,
			isEqual:      false,
		},
		{
			name:         "DefaultGenesisState",
			genesisState: DefaultGenState,
			isEqual:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, reflect.DeepEqual(DefaultGenesisState(), tt.genesisState), tt.isEqual)
		})
	}
}
