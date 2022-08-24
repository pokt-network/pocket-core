package types

import (
	"fmt"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/privval"
	"sync"
)

// GlobalEvidenceCache & GlobalSessionCache is used for the first pocket node and acts as backwards-compatibility for pre-lean pocket
var GlobalEvidenceCache *CacheStorage
var GlobalSessionCache *CacheStorage

var GlobalPocketNodes = map[string]*PocketNode{}

// PocketNode represents an entity in the network that is able to handle dispatches, servicing, challenges, and submit proofs/claims.
type PocketNode struct {
	PrivateKey      crypto.PrivateKey
	EvidenceStore   *CacheStorage
	SessionStore    *CacheStorage
	DoCacheInitOnce sync.Once
}

func (n *PocketNode) GetAddress() sdk.Address {
	return sdk.GetAddress(n.PrivateKey.PublicKey())
}

func AddPocketNode(pk crypto.PrivateKey, logger log.Logger) *PocketNode {
	key := sdk.GetAddress(pk.PublicKey()).String()
	logger.Info("Adding " + key + " to list of pocket nodes")
	node, exists := GlobalPocketNodes[key]
	if exists {
		return node
	}
	node = &PocketNode{
		PrivateKey: pk,
	}
	GlobalPocketNodes[key] = node
	return node
}

func AddPocketNodeByFilePVKey(fpvKey privval.FilePVKey, logger log.Logger) {
	key, err := crypto.PrivKeyToPrivateKey(fpvKey.PrivKey)
	if err != nil {
		return
	}
	AddPocketNode(key, logger)
}

// InitPocketNodeCache adds a PocketNode with its SessionStore and EvidenceStore initialized
func InitPocketNodeCache(node *PocketNode, c types.Config, logger log.Logger) {
	node.DoCacheInitOnce.Do(func() {
		evidenceDbName := c.PocketConfig.EvidenceDBName
		address := node.GetAddress().String()
		// In LeanPocket, we create a evidence store on disk with suffix of the node's address
		if c.PocketConfig.LeanPocket {
			evidenceDbName = evidenceDbName + "_" + address
		}
		logger.Info("Initializing " + address + " session and evidence cache")
		node.EvidenceStore = &CacheStorage{}
		node.SessionStore = &CacheStorage{}
		node.EvidenceStore.Init(c.PocketConfig.DataDir, evidenceDbName, c.TendermintConfig.LevelDBOptions, c.PocketConfig.MaxEvidenceCacheEntires, false)
		node.SessionStore.Init(c.PocketConfig.DataDir, "", c.TendermintConfig.LevelDBOptions, c.PocketConfig.MaxSessionCacheEntries, true)

		// Set the GOBSession and GOBEvidence Global for backwards compatibility for pre-LeanPocket
		if GlobalSessionCache == nil {
			GlobalSessionCache = node.SessionStore
			GlobalEvidenceCache = node.EvidenceStore
		}
	})
}

func InitPocketNodeCaches(c types.Config, logger log.Logger) {
	for _, node := range GlobalPocketNodes {
		InitPocketNodeCache(node, c, logger)
	}
}

// GetPocketNodeByAddress returns a PocketNode from global map GlobalPocketNodes
func GetPocketNodeByAddress(address *sdk.Address) (*PocketNode, error) {
	node, ok := GlobalPocketNodes[address.String()]
	if !ok {
		return nil, fmt.Errorf("failed to find private key for %s", address.String())
	}
	return node, nil
}

// CleanPocketNodes sets the global pocket nodes and its caches back to original state as if the node is starting up again.
// Cleaning up pocket nodes is used for unit and integration tests where the cache is initialized in various scenarios (relays, tx, etc).
func CleanPocketNodes() {
	for _, n := range GlobalPocketNodes {
		if n == nil {
			continue
		}
		cacheToClean := []*CacheStorage{n.EvidenceStore, n.SessionStore}
		for _, r := range cacheToClean {
			if r == nil {
				continue
			}
			r.Clear()
			if r.DB == nil {
				continue
			}
			r.DB.Close()
		}
		GlobalEvidenceCache = nil
		GlobalSessionCache = nil
		GlobalPocketNodes = map[string]*PocketNode{}
	}
}

// GetPocketNode returns a PocketNode from global map GlobalPocketNodes, it does not guarantee order
func GetPocketNode() *PocketNode {
	for _, r := range GlobalPocketNodes {
		if r != nil {
			return r
		}
	}
	return nil
}
