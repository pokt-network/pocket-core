package types

import (
	"github.com/pokt-network/pocket-core/crypto"
	"sync"
)

var GlobalNodesLean map[string]*LeanNode

type LeanNode struct {
	PrivateKey 							crypto.PrivateKey
	GlobalEvidenceCacheLean            *CacheStorage
	GlobalSessionCacheLean             *CacheStorage
	GlobalEvidenceSealedMapLean        *sync.Map
}
