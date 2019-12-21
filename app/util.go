package app

import (
	"encoding/json"
	pocket "github.com/pokt-network/pocket-core/x/pocketcore"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
)

func GenerateChain(ticker, netid, version, client, inter string) (chain string, err error) {
	chain, err = (*app.mm.GetModule(pocketTypes.ModuleName)).(pocket.AppModule).GenerateChain(ticker, netid, version, client, inter)
	return
}

func GenerateAAT(appPubKey, clientPubKey, passphrase string) (aatjson []byte, err error) {
	aat, err := (*app.mm.GetModule(pocketTypes.ModuleName)).(pocket.AppModule).GenerateAAT(appPubKey, clientPubKey, passphrase)
	return json.MarshalIndent(aat, "", "  ")
}
