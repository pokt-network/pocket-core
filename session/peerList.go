package session

import "github.com/pokt-network/pocket-core/types"

type PeerList types.List

func NewPeerList() PeerList {
	return *(*PeerList)(types.NewList())
}

func (pl *PeerList) Get(gid string) Peer {
	return (*types.List)(pl).Get(gid).(Peer)
}

func (pl *PeerList) Set(gid string, sp Peer) {
	(*types.List)(pl).Add(gid, sp)
}

func (pl *PeerList) Count() int {
	return (*types.List)(pl).Count()
}
