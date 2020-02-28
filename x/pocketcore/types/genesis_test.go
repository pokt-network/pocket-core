package types

import (
	"encoding/binary"
	"github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestValidateGenesis(t *testing.T) {
	appPubKeyProof := getRandomPubKey().RawString()
	appPubKeyClaim := getRandomPubKey().RawString()
	pk := getRandomPubKey()
	servicerAddr := pk.Address()
	nn, err := NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	rootHash := Hash([]byte("fakeRoot"))
	rootSum := binary.LittleEndian.Uint64(rootHash)
	root := HashSum{
		Hash: rootHash,
		Sum:  rootSum,
	}
	invalidParams := GenesisState{
		Params: Params{
			SessionNodeCount:      0,
			ClaimSubmissionWindow: 0,
			SupportedBlockchains:  nil,
			ClaimExpiration:       0,
		},
		Proofs: []Receipt{{
			SessionHeader: SessionHeader{
				ApplicationPubKey:  appPubKeyProof,
				Chain:              nn,
				SessionBlockHeight: 1,
			},
			ServicerAddress: servicerAddr.String(),
			TotalRelays:     100,
		}},
		Claims: []MsgClaim{{
			SessionHeader: SessionHeader{
				ApplicationPubKey:  appPubKeyClaim,
				Chain:              nn,
				SessionBlockHeight: 1,
			},
			MerkleRoot:  root,
			TotalRelays: 1000,
			FromAddress: types.Address(servicerAddr),
		}},
	}
	invalidProofs := GenesisState{
		Params: Params{
			SessionNodeCount:      1,
			ClaimSubmissionWindow: 5,
			SupportedBlockchains:  []string{nn},
			ClaimExpiration:       50,
		},
		Proofs: []Receipt{{
			SessionHeader: SessionHeader{
				ApplicationPubKey:  appPubKeyProof,
				Chain:              nn,
				SessionBlockHeight: 1,
			},
			ServicerAddress: servicerAddr.String(),
			TotalRelays:     -1,
		}},
		Claims: []MsgClaim{{
			SessionHeader: SessionHeader{
				ApplicationPubKey:  appPubKeyClaim,
				Chain:              nn,
				SessionBlockHeight: 1,
			},
			MerkleRoot:  root,
			TotalRelays: 1000,
			FromAddress: types.Address(servicerAddr),
		}},
	}
	invalidClaims := GenesisState{
		Params: Params{
			SessionNodeCount:      1,
			ClaimSubmissionWindow: 5,
			SupportedBlockchains:  []string{nn},
			ClaimExpiration:       50,
		},
		Proofs: []Receipt{{
			SessionHeader: SessionHeader{
				ApplicationPubKey:  appPubKeyProof,
				Chain:              nn,
				SessionBlockHeight: 1,
			},
			ServicerAddress: servicerAddr.String(),
			TotalRelays:     100,
		}},
		Claims: []MsgClaim{{
			SessionHeader: SessionHeader{
				ApplicationPubKey:  appPubKeyClaim,
				Chain:              nn,
				SessionBlockHeight: 1,
			},
			MerkleRoot:  root,
			TotalRelays: -1000,
			FromAddress: types.Address(servicerAddr),
		}},
	}
	validGenesisState := GenesisState{
		Params: Params{
			SessionNodeCount:      1,
			ClaimSubmissionWindow: 5,
			SupportedBlockchains:  []string{nn},
			ClaimExpiration:       50,
		},
		Proofs: []Receipt{{
			SessionHeader: SessionHeader{
				ApplicationPubKey:  appPubKeyProof,
				Chain:              nn,
				SessionBlockHeight: 1,
			},
			ServicerAddress: servicerAddr.String(),
			TotalRelays:     100,
		}},
		Claims: []MsgClaim{{
			SessionHeader: SessionHeader{
				ApplicationPubKey:  appPubKeyClaim,
				Chain:              nn,
				SessionBlockHeight: 1,
			},
			MerkleRoot:  root,
			TotalRelays: 1000,
			FromAddress: types.Address(servicerAddr),
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
			name:         "Bad proofs",
			genesisState: invalidProofs,
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
	appPubKeyProof := getRandomPubKey().RawString()
	appPubKeyClaim := getRandomPubKey().RawString()
	pk := getRandomPubKey()
	servicerAddr := pk.Address()
	nn, err := NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	rootHash := Hash([]byte("fakeRoot"))
	rootSum := binary.LittleEndian.Uint64(rootHash)
	root := HashSum{
		Hash: rootHash,
		Sum:  rootSum,
	}
	validGenesisState := GenesisState{
		Params: Params{
			SessionNodeCount:      1,
			ClaimSubmissionWindow: 5,
			SupportedBlockchains:  []string{nn},
			ClaimExpiration:       50,
		},
		Proofs: []Receipt{{
			SessionHeader: SessionHeader{
				ApplicationPubKey:  appPubKeyProof,
				Chain:              nn,
				SessionBlockHeight: 1,
			},
			ServicerAddress: servicerAddr.String(),
			TotalRelays:     100,
		}},
		Claims: []MsgClaim{{
			SessionHeader: SessionHeader{
				ApplicationPubKey:  appPubKeyClaim,
				Chain:              nn,
				SessionBlockHeight: 1,
			},
			MerkleRoot:  root,
			TotalRelays: 1000,
			FromAddress: types.Address(servicerAddr),
		}},
	}
	DefaultGenState := GenesisState{Params: Params{
		SessionNodeCount:      DefaultSessionNodeCount,
		ClaimSubmissionWindow: DefaultClaimSubmissionWindow,
		SupportedBlockchains:  DefaultSupportedBlockchains,
		ClaimExpiration:       DefaultClaimExpiration,
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
