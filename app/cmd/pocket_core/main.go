package main

import (
	"github.com/pokt-network/pocket-core/app/cmd/cli"
	"github.com/pokt-network/pocket-core/app/cmd/rpc"
)

func main() {
	go rpc.StartRPC("8081")
	cli.Execute()
}
