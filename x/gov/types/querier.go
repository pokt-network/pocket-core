package types

import (
	sdk "github.com/pokt-network/pocket-core/types"
)

// query endpoints supported by the staking Querier
const (
	ModuleName                         = "gov"           // ModuleKey defines the name of the module
	RouterKey                          = ModuleName      // RouterKey defines the routing key for a Parameter Change
	StoreKey                           = "gov"           // StoreKey is the string store key for the param store
	TStoreKey                          = "transient_gov" // TStoreKey is the string store key for the param transient store
	DefaultCodespace sdk.CodespaceType = ModuleName      // default codespace for governance errors
	QuerierRoute                       = ModuleName      // QuerierRoute is the querier route for the staking module
	QueryACL                           = "acl"
	QueryDAO                           = "dao"
	QueryUpgrade                       = "upgrade"
	QueryDAOOwner                      = "daoOwner"
)

type QueryACLParams struct{}

type QueryDAOParams struct{}

type QueryUpgradeParams struct{}
