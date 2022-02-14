package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	appexported "github.com/pokt-network/pocket-core/x/apps/exported"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"log"
)

// "Session" - The relationship between an application and the pocket network

func (s Session) IsSealed() bool {
	return false
}

func (s Session) Seal() CacheObject {
	return s
}

// "NewSession" - create a new session from seed data
func NewSession(sessionCtx, ctx sdk.Ctx, keeper PosKeeper, sessionHeader SessionHeader, blockHash string, sessionNodesCount int) (Session, sdk.Error) {
	// first generate session key
	sessionKey, err := NewSessionKey(sessionHeader.ApplicationPubKey, sessionHeader.Chain, blockHash)
	if err != nil {
		return Session{}, err
	}
	// then generate the service nodes for that session
	sessionNodes, err := NewSessionNodes(sessionCtx, ctx, keeper, sessionHeader.Chain, sessionKey, sessionNodesCount)
	if err != nil {
		return Session{}, err
	}
	// then populate the structure and return
	return Session{
		SessionKey:    sessionKey,
		SessionHeader: sessionHeader,
		SessionNodes:  sessionNodes,
	}, nil
}

// "Validate" - Validates a session object
func (s Session) Validate(node sdk.Address, app appexported.ApplicationI, sessionNodeCount int) sdk.Error {
	// validate chain
	if len(s.SessionHeader.Chain) == 0 {
		return NewEmptyNonNativeChainError(ModuleName)
	}
	// validate sessionBlockHeight
	if s.SessionHeader.SessionBlockHeight < 1 {
		return NewInvalidBlockHeightError(ModuleName)
	}
	// validate the app public key
	if err := PubKeyVerification(s.SessionHeader.ApplicationPubKey); err != nil {
		return err
	}
	// validate app corresponds to appPubKey
	if app.GetPublicKey().RawString() != s.SessionHeader.ApplicationPubKey {
		return NewInvalidAppPubKeyError(ModuleName)
	}
	// validate app chains
	chains := app.GetChains()
	found := false
	for _, c := range chains {
		if c == s.SessionHeader.Chain {
			found = true
			break
		}
	}
	if !found {
		return NewUnsupportedBlockchainAppError(ModuleName)
	}
	// validate sessionNodes
	err := s.SessionNodes.Validate(sessionNodeCount)
	if err != nil {
		return err
	}
	// validate node is of the session
	if !s.SessionNodes.Contains(node) {
		return NewInvalidSessionError(ModuleName)
	}
	return nil
}

var _ CacheObject = Session{} // satisfies the cache object interface

func (s Session) MarshalObject() ([]byte, error) {
	return ModuleCdc.ProtoMarshalBinaryBare(&s)
}

func (s Session) UnmarshalObject(b []byte) (CacheObject, error) {
	err := ModuleCdc.ProtoUnmarshalBinaryBare(b, &s)
	return s, err
}

func (s Session) Key() ([]byte, error) {
	return s.SessionHeader.Hash(), nil
}

// "SessionNodes" - Service nodes in a session
type SessionNodes []sdk.Address

// "NewSessionNodes" - Generates nodes for the session
func NewSessionNodes(sessionCtx, ctx sdk.Ctx, keeper PosKeeper, chain string, sessionKey SessionKey, sessionNodesCount int) (sessionNodes SessionNodes, err sdk.Error) {
	// all nodesAddrs at session genesis
	nodesAddrs, totalNodes := keeper.GetValidatorsByChain(sessionCtx, chain)
	// validate nodesAddrs
	if totalNodes < sessionNodesCount {
		return nil, NewInsufficientNodesError(ModuleName)
	}
	sessionNodes = make(SessionNodes, sessionNodesCount)
	var node exported.ValidatorI
	//unique address map to avoid re-checking a pseudorandomly selected servicer
	m := make(map[string]struct{})
	// only select the nodesAddrs if not jailed
	for i, numOfNodes := 0, 0; ; i++ {
		//if this is true we already checked all nodes we got on getValidatorsBychain
		if len(m) >= totalNodes {
			return nil, NewInsufficientNodesError(ModuleName)
		}
		// generate the random index
		index := PseudorandomSelection(sdk.NewInt(int64(totalNodes)), sessionKey)
		// merkleHash the session key to provide new entropy
		sessionKey = Hash(sessionKey)
		// get the node from the array
		n := nodesAddrs[index.Int64()]
		//if we already have seen this address we continue as it's either on the list or discarded
		if _, ok := m[n.String()]; ok {
			continue
		}
		//add the node address to the map
		m[n.String()] = struct{}{}

		// cross check the node from the `new` or `end` world state
		node = keeper.Validator(ctx, n)
		// if not found or jailed, don't add to session and continue
		if node == nil || node.IsJailed() || !NodeHasChain(chain, node) || sessionNodes.Contains(node.GetAddress()) {
			continue
		}
		// else add the node to the session
		sessionNodes[numOfNodes] = n
		// increment the number of nodesAddrs in the sessionNodes slice
		numOfNodes++
		// if maxing out the session count end loop
		if numOfNodes == sessionNodesCount {
			break
		}
	}
	// return the nodesAddrs
	return sessionNodes, nil
}

// "Validate" - Validates the session node object
func (sn SessionNodes) Validate(sessionNodesCount int) sdk.Error {
	if len(sn) < sessionNodesCount {
		return NewInsufficientNodesError(ModuleName)
	}
	for _, n := range sn {
		if n == nil {
			return NewEmptyAddressError(ModuleName)
		}
	}
	return nil
}

// "Contains" - Verifies if the session nodes contains the node using the address
func (sn SessionNodes) Contains(addr sdk.Address) bool {
	// if nil return
	if addr == nil {
		return false
	}
	// loop over the nodes
	for _, node := range sn {
		if node == nil {
			continue
		}
		if node.Equals(addr) {
			return true
		}
	}
	return false
}

// "SessionKey" - the merkleHash identifier of the session
type SessionKey []byte

// "sessionKey" - Used for custom json
type sessionKey struct {
	AppPublicKey   string `json:"app_pub_key"`
	NonNativeChain string `json:"chain"`
	BlockHash      string `json:"blockchain"`
}

// "NewSessionKey" - generates the session key from metadata
func NewSessionKey(appPubKey string, chain string, blockHash string) (SessionKey, sdk.Error) {
	// validate appPubKey
	if err := PubKeyVerification(appPubKey); err != nil {
		return nil, err
	}
	// validate chain
	if err := NetworkIdentifierVerification(chain); err != nil {
		return nil, NewEmptyChainError(ModuleName)
	}
	// validate block addr
	if err := HashVerification(blockHash); err != nil {
		return nil, err
	}
	// marshal into json
	seed, err := json.Marshal(sessionKey{
		AppPublicKey:   appPubKey,
		NonNativeChain: chain,
		BlockHash:      blockHash,
	})
	if err != nil {
		return nil, NewJSONMarshalError(ModuleName, err)
	}
	// return the addr of the result
	return Hash(seed), nil
}

// "Validate" - Validates the session key
func (sk SessionKey) Validate() sdk.Error {
	return HashVerification(hex.EncodeToString(sk))
}

// "ValidateHeader" - Validates the header of the session
func (sh SessionHeader) ValidateHeader() sdk.Error {
	// check the app public key for validity
	if err := PubKeyVerification(sh.ApplicationPubKey); err != nil {
		return err
	}
	// verify the chain merkleHash
	if err := NetworkIdentifierVerification(sh.Chain); err != nil {
		return err
	}
	// verify the block height
	if sh.SessionBlockHeight < 1 {
		return NewInvalidBlockHeightError(ModuleName)
	}
	return nil
}

// "Hash" - The cryptographic merkleHash representation of the session header
func (sh SessionHeader) Hash() []byte {
	res := sh.Bytes()
	return Hash(res)
}

// "HashString" - The hex string representation of the merkleHash
func (sh SessionHeader) HashString() string {
	return hex.EncodeToString(sh.Hash())
}

// "Bytes" - The bytes representation of the session header
func (sh SessionHeader) Bytes() []byte {
	res, err := json.Marshal(sh)
	if err != nil {
		log.Fatal(fmt.Errorf("an error occured converting the session header into bytes:\n%v", err))
	}
	return res
}

// "BlockHash" - Returns the merkleHash from the ctx block header
func BlockHash(ctx sdk.Context) string {
	return hex.EncodeToString(ctx.BlockHeader().LastBlockId.Hash)
}

// "MaxPossibleRelays" - Returns the maximum possible amount of relays for an App on a sessions
func MaxPossibleRelays(app appexported.ApplicationI, sessionNodeCount int64) sdk.BigInt {
	//GetMaxRelays Max value is bound to math.MaxUint64,
	//current worse case is 1 chain and 5 nodes per session with a result of 3689348814741910323 which can be used safely as int64
	return app.GetMaxRelays().ToDec().Quo(sdk.NewDec(int64(len(app.GetChains())))).Quo(sdk.NewDec(sessionNodeCount)).RoundInt()
}

// "NodeHashChain" - Returns whether or not the node has the relayChain
func NodeHasChain(chain string, node exported.ValidatorI) bool {
	hasChain := false
	for _, c := range node.GetChains() {
		if c == chain {
			hasChain = true
			break
		}
	}
	return hasChain
}
