package session

type PeerList struct {
	List map[string]SessionPeer `json:"Peers"`
}

func (sPL *PeerList) Get(gid string) SessionPeer {
	return sPL.List[gid]
}

func (sPL *PeerList) Set(gid string, sP SessionPeer) {
	sPL.List[gid] = sP
}

func (sPL *PeerList) Count() int { // TODO probably need MUX for thread safety
	return len(sPL.List)
}
