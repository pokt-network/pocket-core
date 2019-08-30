package session

import (
	"github.com/pokt-network/pocket-core/legacy"
	types "github.com/pokt-network/pocket-core/types"
)

type Session struct {
	SessionKey      SessionKey        `json:"sessionkey"`
	Developer       types.Developer   `json:"developer"`
	NonNativeChain  legacy.Blockchain `json:"nonnativechain"`
	LatestBlockID SessionBlock     `json:"latestBlock"`
	Nodes           SessionNodes      `json:"sessionNodes"`
}

func NewSession() {

}
