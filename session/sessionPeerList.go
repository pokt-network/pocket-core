package session

type SessionPeerList struct {
	List map[string]SessionPeer `json:"Peers"`
}

func (sPL *SessionPeerList) Get(gid string) SessionPeer{
	return sPL.List[gid]
}

func (sPL *SessionPeerList) Set(gid string, sP SessionPeer){
	sPL.List[gid]=sP
}

func (sPL *SessionPeerList) Count() int {
	return len(sPL.List)
}
