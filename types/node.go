package types

import (
	"github.com/pokt-network/pocket-core/legacy"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto"
	"time"
)

// "Node" is the base structure for a Pocket Network Node"
type Node struct {
	Account         `json:"routing"`
	URL             []byte              `json:"url"`
	SupportedChains NodeSupportedChains `json:"supportedChains"`
	IsAlive         bool
}

type Nodes []Node

type NodeSupportedChain struct {
	legacy.Blockchain `json:"blockchain"`
}

type NodeSupportedChains map[string]NodeSupportedChain // [hex]->Blockchain

func (nsc NodeSupportedChains) Add(hexBlockchainhash string, blockchain NodeSupportedChain) {
	nsc[hexBlockchainhash] = blockchain
}

func (nsc NodeSupportedChains) Contains(hexBlockchainHash string) bool {
	_, contains := nsc[hexBlockchainHash]
	return contains
}

// posmint compatible node

// ValidatorI expected validator functions
type PosmintNode struct {
	Address                 sdk.ValAddress      `json:"address" yaml:"address"`               // address of the validator; bech encoded in JSON
	ConsPubKey              crypto.PubKey       `json:"cons_pubkey" yaml:"cons_pubkey"`       // the consensus public key of the validator; bech encoded in JSON
	Jailed                  bool                `json:"jailed" yaml:"jailed"`                 // has the validator been jailed from bonded status?
	Status                  sdk.BondStatus      `json:"status" yaml:"status"`                 // validator status (bonded/unbonding/unbonded)
	Chains                  map[string]struct{} `json:"chains" yaml:"chains"`                 // validator non native blockchains
	ServiceURL              string              `json:"serviceurl" yaml:"serviceurl"`         // url where the pocket service api is hosted
	StakedTokens            sdk.Int             `json:"Tokens" yaml:"Tokens"`                 // tokens staked in the network
	UnstakingCompletionTime time.Time           `json:"unstaking_time" yaml:"unstaking_time"` // if unstaking, min time for the va
}

var _ exported.ValidatorI = PosmintNode{}

// return the TM validator address
func (pn PosmintNode) ConsAddress() sdk.ConsAddress   { return sdk.ConsAddress(pn.ConsPubKey.Address()) }
func (pn PosmintNode) GetChains() map[string]struct{} { return pn.Chains }
func (pn PosmintNode) GetServiceURL() string          { return pn.ServiceURL }
func (pn PosmintNode) IsStaked() bool                 { return pn.GetStatus().Equal(sdk.Bonded) }
func (pn PosmintNode) IsUnstaked() bool               { return pn.GetStatus().Equal(sdk.Unbonded) }
func (pn PosmintNode) IsUnstaking() bool              { return pn.GetStatus().Equal(sdk.Unbonding) }
func (pn PosmintNode) IsJailed() bool                 { return pn.Jailed }
func (pn PosmintNode) GetStatus() sdk.BondStatus      { return pn.Status }
func (pn PosmintNode) GetAddress() sdk.ValAddress     { return pn.Address }
func (pn PosmintNode) GetConsPubKey() crypto.PubKey   { return pn.ConsPubKey }
func (pn PosmintNode) GetConsAddr() sdk.ConsAddress   { return sdk.ConsAddress(pn.ConsPubKey.Address()) }
func (pn PosmintNode) GetTokens() sdk.Int             { return pn.StakedTokens }
func (pn PosmintNode) GetConsensusPower() int64       { return 0 } // not needed for test right now
