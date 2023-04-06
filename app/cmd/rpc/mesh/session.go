package mesh

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	sdk "github.com/pokt-network/pocket-core/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/robfig/cron/v3"
	"io"
	log2 "log"
	"net/http"
	"time"
)

// DispatchSessionNode - app session node structure
type DispatchSessionNode struct {
	Address       string          `json:"address"`
	Chains        []string        `json:"chains"`
	Jailed        bool            `json:"jailed"`
	OutputAddress string          `json:"output_address"`
	PublicKey     string          `json:"public_key"`
	ServiceUrl    string          `json:"service_url"`
	Status        sdk.StakeStatus `json:"status"`
	Tokens        string          `json:"tokens"`
	UnstakingTime time.Time       `json:"unstaking_time"`
}

// DispatchSession - app session structure
type DispatchSession struct {
	Header pocketTypes.SessionHeader `json:"header"`
	Key    string                    `json:"key"`
	Nodes  []DispatchSessionNode     `json:"nodes"`
}

// DispatchResponse handle /v1/client/dispatch response due to was unable to inflate it using pocket core struct
// it was throwing an error about Nodes unmarshalling
type DispatchResponse struct {
	BlockHeight int64           `json:"block_height"`
	Session     DispatchSession `json:"session"`
}

// Contains - evaluate if the dispatch response contains passed address in their node list
func (sn DispatchResponse) Contains(addr sdk.Address) bool {
	// if nil return
	if addr == nil {
		return false
	}
	// loop over the nodes
	for _, node := range sn.Session.Nodes {
		// There is reference to node address so that way we don't have to recreate address twice for pre-leanpokt
		address, err := sdk.AddressFromHex(node.Address)
		if err != nil {
			log2.Fatal(err)
		}
		if _, ok := servicerMap.Load(address.String()); ok {
			return true
		}
	}
	return false
}

// ShouldKeep - evaluate if this dispatch response is one that we need to keep for the running mesh node.
func (sn DispatchResponse) ShouldKeep() bool {
	// loop over the nodes
	for _, node := range sn.Session.Nodes {
		if _, ok := servicerMap.Load(node.Address); ok {
			return true
		}
	}
	// if hit here, no one of in the map match the dispatch response nodes.
	return false
}

// GetSupportedNodes - return a list of the supported nodes of running mesh node from the DispatchResponse payload.
func (sn DispatchResponse) GetSupportedNodes() []string {
	nodes := make([]string, 0)
	// loop over the nodes
	for _, node := range sn.Session.Nodes {
		// There is reference to node address so that way we don't have to recreate address twice for pre-leanpokt
		if _, ok := servicerMap.Load(node.Address); ok {
			nodes = append(nodes, node.Address)
		}
	}
	// if hit here, no one of in the map match the dispatch response nodes.
	return nodes
}

type AppSessionCache struct {
	PublicKey       string
	Chain           string
	Dispatch        *DispatchResponse
	RemainingRelays int64
	IsValid         bool
	Error           *SdkErrorResponse
}

// getAppSession - call ServicerURL to get an application session using retrieve header
func getAppSession(relay *pocketTypes.Relay, model interface{}) *SdkErrorResponse {
	servicerNode := getServicerFromPubKey(relay.Proof.ServicerPubKey)
	payload := pocketTypes.MeshSession{
		SessionHeader: pocketTypes.SessionHeader{
			ApplicationPubKey:  relay.Proof.Token.ApplicationPublicKey,
			Chain:              relay.Proof.Blockchain,
			SessionBlockHeight: relay.Proof.SessionBlockHeight,
		},
		Meta:               relay.Meta,
		ServicerPubKey:     relay.Proof.ServicerPubKey,
		Blockchain:         relay.Proof.Blockchain,
		SessionBlockHeight: relay.Proof.SessionBlockHeight,
	}
	logger.Debug(fmt.Sprintf("reading session from servicer %s", servicerNode.Address.String()))
	jsonData, e := json.Marshal(payload)
	if e != nil {
		return NewSdkErrorFromPocketSdkError(sdk.ErrInternal(e.Error()))
	}

	requestURL := fmt.Sprintf(
		"%s%s?authtoken=%s",
		servicerNode.Node.URL,
		ServicerSessionEndpoint,
		servicerAuthToken.Value,
	)
	req, e := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if e != nil {
		return NewSdkErrorFromPocketSdkError(sdk.ErrInternal(e.Error()))
	}

	req.Header.Set("Content-Type", "application/json")
	if app.GlobalMeshConfig.UserAgent != "" {
		req.Header.Set("User-Agent", app.GlobalMeshConfig.UserAgent)
	}

	resp, e := servicerClient.Do(req)
	if e != nil {
		return NewSdkErrorFromPocketSdkError(sdk.ErrInternal(e.Error()))
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return // add log here
		}
	}(resp.Body)

	if resp.StatusCode == 401 {
		return NewSdkErrorFromPocketSdkError(
			sdk.ErrUnauthorized(
				fmt.Sprintf("wrong auth form %s", ServicerSessionEndpoint),
			),
		)
	}

	isSuccess := resp.StatusCode == 200

	if !isSuccess {
		result := RPCSessionResult{}
		e = json.NewDecoder(resp.Body).Decode(&result)
		if e != nil {
			return NewSdkErrorFromPocketSdkError(sdk.ErrInternal(e.Error()))
		}
		return nil
	} else {
		e = json.NewDecoder(resp.Body).Decode(model)
		if e != nil {
			return NewSdkErrorFromPocketSdkError(sdk.ErrInternal(e.Error()))
		}
		return nil
	}
}

// getSessionHashFromRelay - calculate the session header and late the hash of it
func getSessionHashFromRelay(r *pocketTypes.Relay) []byte {
	header := pocketTypes.SessionHeader{
		ApplicationPubKey:  r.Proof.Token.ApplicationPublicKey,
		Chain:              r.Proof.Blockchain,
		SessionBlockHeight: r.Proof.SessionBlockHeight,
	}

	return header.Hash()
}

// cleanOldSessions - clean up sessions that are longer than 50 blocks (just to be sure they are not needed)
func cleanOldSessions(c *cron.Cron) {
	_, err := c.AddFunc(fmt.Sprintf("@every %ds", app.GlobalMeshConfig.SessionCacheCleanUpInterval), func() {
		servicerMap.Range(func(_ string, servicerNode *servicer) bool {
			servicerNode.SessionCache.Range(func(key string, appSession *AppSessionCache) bool {
				hash, err := hex.DecodeString(key)
				if err != nil {
					logger.Error("error decoding session hash to delete from cache " + err.Error())
					return true
				}

				if appSession.Dispatch == nil {
					servicerNode.DeleteAppSession(hash)
				} else if appSession.Dispatch.Session.Header.SessionBlockHeight < (servicerNode.Node.Status.Height - 6) {
					servicerNode.DeleteAppSession(hash)
				}

				return true
			})
			return true
		})
	})

	if err != nil {
		log2.Fatal(err)
	}
}
