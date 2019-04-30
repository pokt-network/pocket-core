package core

import (
	"github.com/pokt-network/pocket-core/core/fbs"
	"strconv"
)

func UnmarshalBlockchain(flatBuffer []byte) Blockchain {
	res := fbs.GetRootAsBlockchain(flatBuffer, 0)
	return Blockchain{string(res.NameBytes()), strconv.Itoa(int(res.Netid())), strconv.Itoa(int(res.Version()))}
}

// TODO testing with empty URL field
func UnmarshalRelay(flatbuffer []byte) Relay {
	res := fbs.GetRootAsRelay(flatbuffer, 0)
	return Relay{res.BlockchainBytes(), res.PayloadBytes(), res.DevidBytes(), Token{res.Token(&fbs.Token{}).ExpdateBytes()}, res.MethodBytes(), res.UrlBytes()}
}

func UnmarshalRelayMessage(flatbuffer []byte) RelayMessage {
	res := fbs.GetRootAsRelayMessage(flatbuffer, 0)
	r := &fbs.Relay{}
	return RelayMessage{
		Relay{
			res.Relay(r).BlockchainBytes(),
			res.Relay(r).PayloadBytes(), res.Relay(r).DevidBytes(),
			Token{res.Relay(r).Token(&fbs.Token{}).ExpdateBytes()},
			res.Relay(r).MethodBytes(),
			res.Relay(r).UrlBytes()},
		res.SignatureBytes()}
}
