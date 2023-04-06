package mesh

import (
	"encoding/hex"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/puzpuzpuz/xsync"
	"time"
)

// Servicer - represents a staked address read from servicer_private_key_file
type servicer struct {
	SessionCache *xsync.MapOf[string, *AppSessionCache]
	PrivateKey   crypto.PrivateKey
	Address      sdk.Address
	Node         *fullNode
}

// LoadAppSession - retrieve from cache (memory or persistent) an app session cache
func (s *servicer) LoadAppSession(hash []byte) (*AppSessionCache, bool) {
	sHash := hex.EncodeToString(hash)
	if v, ok := s.SessionCache.Load(sHash); ok {
		return v, ok
	}

	return nil, false
}

// StoreAppSession - store in cache (memory and persistent) an appCache
func (s *servicer) StoreAppSession(hash []byte, appSession *AppSessionCache) {
	hashString := hex.EncodeToString(hash)
	s.SessionCache.Store(hashString, appSession)

	return
}

// DeleteAppSession - delete an app session from cache (memory and persistent)
func (s *servicer) DeleteAppSession(hash []byte) {
	sHash := hex.EncodeToString(hash)
	s.SessionCache.Delete(sHash)
}

// reloadServicers - read key file again and manipulate current state after detect differences (add/remove)
func reloadServicers() {
	logger.Debug("initializing servicer reload process")
	nodes, servicers := loadServicersFromFile()

	currentNodes := mapset.NewSet[string]()
	currentServicers := mapset.NewSet[string]()
	reloadedNodes := mapset.NewSet[string]()
	reloadedServicers := mapset.NewSet[string]()

	nodesMap.Range(func(key string, _ *fullNode) bool {
		currentNodes.Add(key)
		return true
	})
	nodes.Range(func(key string, _ *fullNode) bool {
		reloadedNodes.Add(key)
		return true
	})
	servicerMap.Range(func(key string, _ *servicer) bool {
		currentServicers.Add(key)
		return true
	})
	servicers.Range(func(key string, _ *servicer) bool {
		reloadedServicers.Add(key)
		return true
	})

	newNodes := reloadedNodes.Difference(currentNodes)
	newServicers := reloadedServicers.Difference(currentServicers)

	removedNodes := currentNodes.Difference(reloadedNodes)
	removedServicers := currentServicers.Difference(reloadedServicers)

	orphanServicers := xsync.NewMapOf[*servicer]()

	removedNodes.Each(func(s string) bool {
		if node, ok := nodesMap.LoadAndDelete(s); ok {
			node.Servicers.Range(func(key string, s *servicer) bool {
				address := s.Address.String()
				if v, ok := servicerMap.LoadAndDelete(address); ok {
					// so we can use same servicer if basically was a move from node A to node B
					orphanServicers.Store(address, v)
				}
				return true
			})
			node.stop()
		}
		return true
	})
	removedServicers.Each(func(address string) bool {
		if s, ok := servicerMap.LoadAndDelete(address); ok {
			_, loaded := s.Node.Servicers.LoadAndDelete(address)
			if loaded {
				mutex.Lock() // lock it because could be in use by rpc
				s.Node.NeedResize = true
				mutex.Unlock()
			}
		}
		return true
	})

	newNodes.Each(func(s string) bool {
		node, _ := nodes.Load(s)
		// just add it because is new
		nodesMap.Store(s, node)

		node.Servicers.Range(func(key string, s *servicer) bool {
			address := s.Address.String()
			if currentServicer, ok := servicerMap.Load(address); ok {
				node.Servicers.Store(address, currentServicer)
			} else if v, ok := orphanServicers.Load(address); ok {
				// remove it from orphans and assign to new node
				orphanServicers.Delete(address)
				node.Servicers.Store(address, v)
			}
			return true
		})
		// if cron is already running it will omit this call.
		node.start()
		return true
	})
	newServicers.Each(func(address string) bool {
		s, _ := servicers.Load(address)
		servicerMap.Store(address, s)
		if node, ok := nodesMap.Load(s.Node.URL); ok {
			// set the already existent fullNode reference instead of the new one.
			s.Node = node

			if orphan, ok := orphanServicers.Load(address); ok {
				orphan.Node = node
				s = orphan
				// reassign to reuse the orphan one
			}

			s.Node.Servicers.Store(address, s)

			mutex.Lock() // lock it because could be in use by rpc
			s.Node.NeedResize = true
			mutex.Unlock()
		}
		orphanServicers.Delete(address) // just in case it remain on orphans one.
		return true
	})

	totalModifications := newNodes.Cardinality() + newServicers.Cardinality() + removedNodes.Cardinality() + removedServicers.Cardinality()

	if totalModifications > 0 {
		logger.Info(fmt.Sprintf("Servicer reload detect a total of %d modifications:", totalModifications))
		logger.Info(fmt.Sprintf("New Nodes: %d Removed Nodes: %d", newNodes.Cardinality(), removedNodes.Cardinality()))
		logger.Info(fmt.Sprintf("New Servicer: %d Removed Servicer: %d", newServicers.Cardinality(), removedServicers.Cardinality()))

		newServicersList := make([]string, 0)
		totalNodes := nodesMap.Size()

		nodesMap.Range(func(key string, node *fullNode) bool {
			node.ResizeWorker()

			return true
		})

		servicerMap.Range(func(key string, s *servicer) bool {
			newServicersList = append(newServicersList, s.Address.String())
			return true
		})

		mutex.Lock()
		servicerList = newServicersList
		mutex.Unlock()

		logger.Debug(fmt.Sprintf("Current Node Map lenght after modification is: %d", totalNodes))
		logger.Debug(fmt.Sprintf("Current Servicer Map lenght after modification is: %d", len(servicerList)))
		// run connectivity checks only for the new nodes
		connectivityChecks(newNodes)
	}

	logger.Debug("servicers reload process done")
}

// initKeysHotReload - initialize keys read every N time base on MeshConfig.KeysHotReloadInterval
func initKeysHotReload() {
	if app.GlobalMeshConfig.KeysHotReloadInterval <= 0 {
		logger.Info("skipping hot reload due to keys_hot_reload_interval is less or equal to 0")
		return
	}

	logger.Info(fmt.Sprintf("keys hot reload set to run every %d milliseconds", app.GlobalMeshConfig.KeysHotReloadInterval))

	for {
		time.Sleep(time.Duration(app.GlobalMeshConfig.KeysHotReloadInterval) * time.Millisecond)
		reloadServicers()
	}
}

// getServicerFromPubKey - return the servicer instance base on public key.
func getServicerFromPubKey(pubKey string) *servicer {
	servicerAddress, err := GetAddressFromPubKeyAsString(pubKey)

	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"unable to decode servicer public key %s",
				pubKey,
			),
		)
		return nil
	}

	s, ok := servicerMap.Load(servicerAddress)

	if !ok {
		logger.Error(
			fmt.Sprintf(
				"unable to find servicer with address=%s",
				servicerAddress,
			),
		)
		return nil
	}

	return s
}

// GetServicerLists - return a copy of the servicers as array of strings
func GetServicerLists() (servicers []string) {
	mutex.Lock()
	for _, a := range servicerList {
		servicers = append(servicers, a)
	}
	mutex.Unlock()
	return
}

// ServicersSize - return how many servicers are registered under servicerMap
func ServicersSize() int {
	return servicerMap.Size()
}
