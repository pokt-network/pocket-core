package keeper

import (
	"github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_VerifyACL(t *testing.T) {
	ctx, keeper := createTestKeeperAndContext(t, false)
	posACLKey := `pos/foo`
	posACLKey2 := `pos/bar`
	addr := getRandomValidatorAddress()
	addr2 := getRandomValidatorAddress()
	acl := types.ACL(make([]types.ACLPair, 0))
	acl.SetOwner(posACLKey, addr)
	acl.SetOwner(posACLKey2, addr2)
	keeper.SetParams(ctx, types.Params{
		ACL:      acl,
		DAOOwner: addr,
		Upgrade:  types.Upgrade{},
	})
	assert.Nil(t, keeper.VerifyACL(ctx, posACLKey, addr))
	assert.NotNil(t, keeper.VerifyACL(ctx, posACLKey, addr2))
}
