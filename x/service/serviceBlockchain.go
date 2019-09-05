package service

import "github.com/pokt-network/pocket-core/types"

type ServiceBlockchain types.AminoBuffer

type ServiceBlockchains map[string]ServiceBlockchain
