// This package is the starting point of Pocket Core.
package main

import (
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/posmint/config"
	"github.com/pokt-network/posmint/types"
	"github.com/tendermint/go-amino"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// "init" is a built in function that is automatically called before main.
func init() {
	// generates seed for randomization
	rand.Seed(time.Now().UTC().UnixNano())
}

// "main" is the starting function of the client.
func main() {
	startClient()
}

// "startClient" Starts the client with the given initial configuration.
func startClient() {
	rootDir := "app/cmd/pocket_core/test/"
	vKFP := "val_priv_key.json"
	sFP := "val_state_fp.json"
	nk := "node_key.json"
	cdc := app.MakeCodec()
	app.Cdc = cdc
	err := config.LoadOrGenerateNodeKeyFile(cdc, rootDir+string(os.PathSeparator)+nk)
	if err != nil {
		panic(err)
	}
	app.Cdc.RegisterConcrete(types.ValAddress{}, "validatorAddr", &amino.ConcreteOptions{})
	config.LoadOrGenFilePV(rootDir+string(os.PathSeparator)+vKFP, rootDir+string(os.PathSeparator)+sFP)
	app.InitKeybase(rootDir, "")
	keybase, err := app.GetKeybase()
	if err != nil {
		panic(err)
	}
	kp, err := keybase.List()
	k := kp[0]
	fmt.Println(types.Bech32ifyConsPub(k.PubKey))
	fmt.Println(types.ValAddress(k.PubKey.Address()).String())
	app.InitHostedChains(rootDir + "config/chains.json") // todo
	app.SetCoinbasePassphrase("")                        // todo
	app.SetGenesisFilepath(rootDir + "config/genesis.json")
	tmNode := app.InitTendermintNode(rootDir, "", nk, vKFP, sFP, "84153F412E8148C8545FAD7173CB0BC2D87102C2@localhost:26656", "84153F412E8148C8545FAD7173CB0BC2D87102C2@localhost:26656", "localhost:26656") // todo
	err := tmNode.Start()
	if err != nil {
		panic(err)
	}
	// We trap kill signals (2,3,15,9)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		os.Kill,
		os.Interrupt)

	defer func() {
		sig := <-signalChannel
		tmNode, _ := app.GetTendermintNode()
		tmNode.Stop()
		message := fmt.Sprintf("Exit signal %s received\n", sig)
		fmt.Println(message)
		os.Exit(3)
	}()
}
