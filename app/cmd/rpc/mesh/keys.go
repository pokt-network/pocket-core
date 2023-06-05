package mesh

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/puzpuzpuz/xsync"
	"github.com/xeipuuv/gojsonschema"
	log2 "log"
	"os"
)

type nodeFileItem struct {
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Keys []string `json:"keys"`
}

type fallbackNodeFileItem struct {
	PrivateKey  string `json:"priv_key"`
	ServicerUrl string `json:"servicer_url"`
}

// getServicersFilePath - return servicers file path resolved by config.json
func getServicersFilePath() string {
	return app.GlobalMeshConfig.DataDir + app.FS + app.GlobalMeshConfig.ServicerPrivateKeyFile
}

// loadServicersFromFile return a sync.Map of nodes/servicers that could be used to start working or calculate a reload
func loadServicersFromFile() (nodes *xsync.MapOf[string, *fullNode], servicers *xsync.MapOf[string, *servicer]) {
	nodes = xsync.NewMapOf[*fullNode]()
	servicers = xsync.NewMapOf[*servicer]()

	path := getServicersFilePath()

	fallbackSchemaLoader := gojsonschema.NewSchemaLoader()
	fallbackSchemaStringLoader := gojsonschema.NewStringLoader(fallbackNodeFileSchema)
	fallbackSchema, fallbackSchemaError := fallbackSchemaLoader.Compile(fallbackSchemaStringLoader)
	if fallbackSchemaError != nil {
		log2.Fatal(fmt.Errorf("an error occurred loading fallback json schema: %s", fallbackSchemaError.Error()))
	}

	currentSchemaLoader := gojsonschema.NewSchemaLoader()
	currentSchemaStringLoader := gojsonschema.NewStringLoader(nodeFileSchema)
	currentSchema, currentSchemaError := currentSchemaLoader.Compile(currentSchemaStringLoader)
	if currentSchemaError != nil {
		log2.Fatal(fmt.Errorf("an error occurred loading json schema: %s", currentSchemaError.Error()))
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log2.Fatal(fmt.Errorf("an error occurred attempting to read the servicer key file: %s", err.Error()))
	}

	strData := gojsonschema.NewStringLoader(string(data[:]))

	if r, e := fallbackSchema.Validate(strData); e != nil || len(r.Errors()) > 0 {
		if r2, e2 := currentSchema.Validate(strData); e2 != nil || len(r2.Errors()) > 0 {
			log2.Fatal(fmt.Errorf("unable to parse file %s to any of the supported key schemas", path))
		} else {
			var readServicers []nodeFileItem
			// load servicers with new format
			if err := json.Unmarshal(data, &readServicers); err != nil {
				log2.Fatal(fmt.Errorf("an error occurred attempting to parse the servicer key file: %s", err.Error()))
			}

			for _, n := range readServicers {
				var node *fullNode

				if v, ok := nodes.Load(n.URL); !ok {
					node = createNode(n.URL, n.Name)
					nodes.Store(n.URL, node)
				} else {
					node = v
				}

				for index, pkStr := range n.Keys {
					pk, err := crypto.NewPrivateKey(pkStr)
					if err != nil {
						log2.Fatal(fmt.Errorf("error parsing private key on node=%s index=%d of the file %s", n.URL, index, path))
					}

					address, err := sdk.AddressFromHex(pk.PubKey().Address().String())
					if err != nil {
						log2.Fatal(fmt.Errorf("error getting address from private key on node=%s index=%d of the file %s", n.URL, index, path))
					}

					addressStr := address.String()

					if s, ok := servicers.Load(addressStr); ok {
						node.Servicers.Store(addressStr, s)
					} else {
						newServicer := &servicer{
							PrivateKey: pk,
							Address:    address,
							Node:       node,
						}
						servicers.Store(addressStr, newServicer)
						node.Servicers.Store(addressStr, newServicer)
					}
				}
			}
		}
	} else {
		// load servicer with fallback one.
		var readServicers []fallbackNodeFileItem
		// load servicers with new format
		if err := json.Unmarshal(data, &readServicers); err != nil {
			log2.Fatal(fmt.Errorf("an error occurred attempting to parse the servicer key file: %s", err.Error()))
		}

		for index, n := range readServicers {
			var node *fullNode

			if v, ok := nodes.Load(n.ServicerUrl); !ok {
				node = createNode(n.ServicerUrl, "")
				nodes.Store(n.ServicerUrl, node)
			} else {
				node = v
			}

			pk, err := crypto.NewPrivateKey(n.PrivateKey)
			if err != nil {
				log2.Fatal(fmt.Errorf("error parsing private key on node=%s index=%d of the file %s", n.ServicerUrl, index, path))
			}

			address, err := sdk.AddressFromHex(pk.PubKey().Address().String())
			if err != nil {
				log2.Fatal(fmt.Errorf("error getting address from private key on node=%s index=%d of the file %s", n.ServicerUrl, index, path))
			}

			addressStr := address.String()

			if s, ok := servicers.Load(addressStr); ok {
				node.Servicers.Store(addressStr, s)
			} else {
				newServicer := &servicer{
					PrivateKey: pk,
					Address:    address,
					Node:       node,
				}
				servicers.Store(addressStr, newServicer)
				node.Servicers.Store(addressStr, newServicer)
			}
		}
	}

	return
}

// loadServicerNodes - read servicer address and cast to sdk.Address
func loadServicerNodes() (totalNodes, totalServicers int) {
	nodes, servicers := loadServicersFromFile()

	nodes.Range(func(key string, value *fullNode) bool {
		nodesMap.Store(key, value)
		return true
	})

	servicers.Range(func(key string, value *servicer) bool {
		servicerMap.Store(key, value)
		return true
	})

	loadedServicerList := make([]string, 0)

	servicerMap.Range(func(key string, value *servicer) bool {
		loadedServicerList = append(loadedServicerList, value.Address.String())
		return true
	})

	totalNodes = nodes.Size()
	totalServicers = servicers.Size()
	mutex.Lock()
	servicerList = loadedServicerList
	mutex.Unlock()

	return
}
