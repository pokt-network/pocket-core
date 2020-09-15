package types

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	"strings"
)

const (
	ACLKeySep = "/"
)

func NewACLKey(subspaceName, paramName string) string {
	return subspaceName + ACLKeySep + paramName
}

func SplitACLKey(aclKey string) (subspaceName, paramName string) {
	s := strings.Split(aclKey, ACLKeySep)
	subspaceName = s[0]
	paramName = s[1]
	return
}

type ACL []ACLPair // cant use map cause of amino concrete marshal in tx

func (a ACL) Validate(adjacencyMap map[string]bool) error {
	for _, aclPair := range a {
		key := aclPair.Key
		val := aclPair.Addr
		_, ok := adjacencyMap[key]
		if !ok {
			return ErrInvalidACL(ModuleName, fmt.Errorf("the key: %s is not a recognized parameter", key))
		}
		adjacencyMap[key] = true
		if val == nil {
			return ErrInvalidACL(ModuleName, fmt.Errorf("the address provided for: %s is nil", key))
		}
	}
	var unOwnedParams []string
	for key, val := range adjacencyMap {
		if !val {
			unOwnedParams = append(unOwnedParams, key)
		}
	}
	if len(unOwnedParams) != 0 {
		return ErrInvalidACL(ModuleName, fmt.Errorf("the following params have no owner: %v aclPair", unOwnedParams))
	}
	return nil
}

func (a ACL) GetOwner(permKey string) sdk.Address {
	for _, aclPair := range a {
		if aclPair.Key == permKey {
			return aclPair.Addr
		}
	}
	return nil
}

func (a *ACL) SetOwner(permKey string, ownerValue sdk.Address) {
	for i, aclPair := range *a {
		if aclPair.Key == permKey {
			aclPair.Addr = ownerValue
			(*a)[i] = aclPair
			return
		}
	}
	temp := append(*a, ACLPair{
		Key:  permKey,
		Addr: ownerValue,
	})
	*a = temp
}

func (a ACL) GetAll() map[string]sdk.Address {
	m := make(map[string]sdk.Address)
	for _, aclPair := range a {
		m[aclPair.Key] = aclPair.Addr
	}
	return m
}

func (a ACL) String() string {
	return fmt.Sprintf(`ACL:
%v`, a.GetAll())
}

//type ACLPair struct {
//	Key  string      `json:"acl_key"`
//	Addr sdk.Address `json:"address"`
//}
