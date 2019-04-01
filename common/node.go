package common

import (
	"crypto/sha256"
	"strings"
)

type NodeWorldState struct {
	Enode  string       `json:"enode"`
	Stake  int          `json:"stake"`
	Active bool         `json:"status"`
	IsVal  bool         `json:"isval"`
	Chains []Blockchain `json:"chains"`
}

type Blockchain struct {
	Name    string `json:"name"`
	NetID   string `json:"netid"`
	Version string `json:"string"`
}

func (nws NodeWorldState) EnodeSplit() (gid string, ip string, port string, discport string) {
	var url []string
	e := nws.Enode
	enodeSplit := strings.Split(e, "@")
	if strings.Contains(enodeSplit[1], "?") {
		contact := strings.Split(enodeSplit[1], "?")
		url = strings.Split(contact[0], ":")
		discport = strings.Split(contact[1], "=")[1]
	}
	url = strings.Split(enodeSplit[1], ":")
	ip = url[0]
	port = url[1]
	hash := strings.TrimPrefix(enodeSplit[0], "enode://")
	gid = hash
	return
}

func SHA256FromString(s string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(s))
	return hasher.Sum(nil)
}

func SHA256FromBytes(b []byte) []byte {
	hasher := sha256.New()
	hasher.Write(b)
	return hasher.Sum(nil)
}
