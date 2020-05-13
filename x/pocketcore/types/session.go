package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	appexported "github.com/pokt-network/pocket-core/x/apps/exported"
	nodeexported "github.com/pokt-network/pocket-core/x/nodes/exported"
	sdk "github.com/pokt-network/posmint/types"
	"log"
)

// "Session" - The relationship between an application and the pocket network
type Session struct {
	SessionHeader `json:"header"`
	SessionKey    `json:"key"`
	SessionNodes  `json:"nodes"`
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
func (s Session) Validate(node nodeexported.ValidatorI, app appexported.ApplicationI, sessionNodeCount int) sdk.Error {
	// validate chain
	if len(s.Chain) == 0 {
		return NewEmptyNonNativeChainError(ModuleName)
	}
	// validate sessionBlockHeight
	if s.SessionBlockHeight < 1 {
		return NewInvalidBlockHeightError(ModuleName)
	}
	// validate the app public key
	if err := PubKeyVerification(s.ApplicationPubKey); err != nil {
		return err
	}
	// validate app corresponds to appPubKey
	if app.GetPublicKey().RawString() != s.ApplicationPubKey {
		return NewInvalidAppPubKeyError(ModuleName)
	}
	// validate app chains
	chains := app.GetChains()
	found := false
	for _, c := range chains {
		if c == s.Chain {
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

// "SessionNodes" - Service nodes in a session
type SessionNodes []nodeexported.ValidatorI

// "NewSessionNodes" - Generates nodes for the session
func NewSessionNodes(sessionCtx, ctx sdk.Ctx, keeper PosKeeper, chain string, sessionKey SessionKey, sessionNodesCount int) (sessionNodes SessionNodes, err sdk.Error) {
	// validate chain
	if len(chain) == 0 {
		return nil, NewEmptyNonNativeChainError(ModuleName)
	}
	// validate sessionKey
	if err = sessionKey.Validate(); err != nil {
		return nil, NewInvalidSessionKeyError(ModuleName, err)
	}
	// all nodes at session genesis
	nodes := keeper.GetValidatorsByChain(sessionCtx, chain)
	// get the total number of nodes
	totalNodes := int64(len(nodes))
	// validate nodes
	if len(nodes) < sessionNodesCount {
		return nil, NewInsufficientNodesError(ModuleName)
	}
	sessionNodes = make(SessionNodes, sessionNodesCount)
	// only select the nodes if not jailed
	for i, numOfNodes := 0, 0; ; i++ {
		// generate the random index
		index, err := PseudoRandomGeneration(totalNodes, sessionKey)
		if err != nil {
			return nil, sdk.ErrInternal("error with NewSessionNodes generation: " + err.Error())
		}
		// hash the session key to provide new entropy
		sessionKey = Hash(sessionKey)
		// get the node from the array
		n := nodes[index]
		// cross check the node from the `new` or `end` world state
		res := keeper.Validator(ctx, n.GetAddress())
		// if not found or jailed, don't add to session and continue
		if res == nil || res.IsJailed() {
			continue
		}
		if sessionNodes.ContainsAddress(res.GetAddress()) {
			continue
		}
		// else add the node to the session
		sessionNodes[numOfNodes] = n
		// increment the number of nodes in the sessionNodes slice
		numOfNodes++
		// if maxing out the session count end loop
		if numOfNodes == sessionNodesCount {
			break
		}
	}
	// return the nodes
	return sessionNodes, nil
}

// "Validate" - Validates the session node object
func (sn SessionNodes) Validate(sessionNodesCount int) sdk.Error {
	if len(sn) < sessionNodesCount || sn[0] == nil {
		return NewInsufficientNodesError(ModuleName)
	}
	return nil
}

// "Contains" - Verifies if the session nodes contain the node object
func (sn SessionNodes) Contains(nodeVerify nodeexported.ValidatorI) bool {
	// if nil return
	if nodeVerify == nil {
		return false
	}
	// loop over the nodes
	for _, node := range sn {
		if node.GetPublicKey().Equals(nodeVerify.GetPublicKey()) {
			return true
		}
	}
	return false
}

// "Contains" - Verifies if the session nodes contains the node using the address
func (sn SessionNodes) ContainsAddress(addr sdk.Address) bool {
	// if nil return
	if addr == nil {
		return false
	}
	// loop over the nodes
	for _, node := range sn {
		if node == nil {
			continue
		}
		if node.GetAddress().String() == addr.String() {
			return true
		}
	}
	return false
}

// "SessionKey" - the hash identifier of the session
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

// "Sessionheader" - The header of the session
type SessionHeader struct {
	ApplicationPubKey  string `json:"app_public_key"` // the application public key
	Chain              string `json:"chain"`          // the nonnative chain in the session
	SessionBlockHeight int64  `json:"session_height"` // the session block height
}

// "ValidateHeader" - Validates the header of the session
func (sh SessionHeader) ValidateHeader() sdk.Error {
	// check the app public key for validity
	if err := PubKeyVerification(sh.ApplicationPubKey); err != nil {
		return err
	}
	// verify the chain hash
	if err := NetworkIdentifierVerification(sh.Chain); err != nil {
		return err
	}
	// verify the block height
	if sh.SessionBlockHeight < 1 {
		return NewInvalidBlockHeightError(ModuleName)
	}
	return nil
}

// "Hash" - The cryptographic hash representation of the session header
func (sh SessionHeader) Hash() []byte {
	res := sh.Bytes()
	return Hash(res)
}

// "HashString" - The hex string representation of the hash
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

// "BlockHash" - Returns the hash from the ctx block header
func BlockHash(ctx sdk.Context) string {
	return hex.EncodeToString(ctx.BlockHeader().LastBlockId.Hash)
}
