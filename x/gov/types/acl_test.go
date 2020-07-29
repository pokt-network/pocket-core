package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestACLGetSetOwner(t *testing.T) {
	acl := ACL(make([]ACLPair, 0))
	a := getRandomValidatorAddress()
	acl.SetOwner("gov/acl", a)
	assert.Equal(t, acl.GetOwner("gov/acl").String(), a.String())
}

func TestValidateACL(t *testing.T) {
	acl := createTestACL()
	adjMap := createTestAdjacencyMap()
	assert.Nil(t, acl.Validate(adjMap))
	acl.SetOwner("gov/acl2", getRandomValidatorAddress())
	assert.NotNil(t, acl.Validate(adjMap))
}
