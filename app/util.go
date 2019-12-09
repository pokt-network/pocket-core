package app

import "encoding/json"

func GenerateChain(ticker, netid, version, client, inter string) (chain string, err error) {
	chain, err = pocketModule.GenerateChain(ticker, netid, version, client, inter)
	return
}

func GenerateAAT(appPubKey, clientPubKey, passphrase string) (aatjson []byte, err error) {
	aat, err := pocketModule.GenerateAAT(appPubKey, clientPubKey, passphrase)
	return json.MarshalIndent(aat, "", "  ")
}
