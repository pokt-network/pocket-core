package app

import (
	"encoding/json"
	pocket "github.com/pokt-network/pocket-core/x/pocketcore"
)

func GenerateChain(ticker, netid, version, client, inter string) (chain string, err error) {
	chain, err = pocket.GenerateChain(ticker, netid, version, client, inter)
	return
}

func GenerateAAT(appPubKey, clientPubKey, passphrase string) (aatjson []byte, err error) {
	aat, err := pocket.GenerateAAT(MustGetKeybase(), appPubKey, clientPubKey, passphrase)
	return json.MarshalIndent(aat, "", "  ")
}
