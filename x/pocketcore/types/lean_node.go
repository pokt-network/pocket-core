package types

import (
	"fmt"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"sync"
)

var GlobalNodesLean map[string]*LeanNode

type LeanNode struct {
	PrivateKey        crypto.PrivateKey
	EvidenceCache     *CacheStorage
	SessionCache      *CacheStorage
	EvidenceSealedMap *sync.Map
}

func InitNodeWithCacheLean(pk crypto.PrivateKey) {

	if GlobalNodesLean == nil {
		GlobalNodesLean = make(map[string]*LeanNode)
	}

	leanNode := LeanNode{PrivateKey: pk, EvidenceCache: new(CacheStorage), SessionCache: new(CacheStorage), EvidenceSealedMap: &sync.Map{}}
	key := sdk.GetAddress(pk.PublicKey()).String()
	_, exists := GlobalNodesLean[key]
	if exists {
		fmt.Println(key + " already added as a lean node")
		return
	}

	leanNode.EvidenceCache.Init(GlobalPocketConfig.DataDir, GlobalPocketConfig.EvidenceDBName+"_"+key, GlobalTenderMintConfig.LevelDBOptions, GlobalPocketConfig.MaxEvidenceCacheEntires, false)
	leanNode.SessionCache.Init(GlobalPocketConfig.DataDir, "", GlobalTenderMintConfig.LevelDBOptions, GlobalPocketConfig.MaxSessionCacheEntries, true)
	GlobalNodesLean[key] = &leanNode
}

func GetNodeLean(address *sdk.Address) (*LeanNode, error) {
	node, ok := GlobalNodesLean[address.String()]
	if !ok {
		return nil, fmt.Errorf("failed to find private key for %s", address.String())
	}
	return node, nil
}
