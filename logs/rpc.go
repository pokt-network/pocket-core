package logs

import (
	"encoding/json"
	"github.com/pokt-network/pocket-core/config"
	"io/ioutil"
	"path/filepath"
	"time"
)

const (
	REQUEST  = "REQUEST"
	RESPONSE = "RESPONSE"
)

type RPCLog struct {
	Type      string `json:"type"`
	ClientIP  string `json:"clientip"`
	Payload   string `json:"payload"`
	Timestamp string `json:"timestamp"`
}

func newRPCLog(isRequest bool, ip string, body string) RPCLog {
	if isRequest {
		return RPCLog{REQUEST, ip, body, time.Now().UTC().String()}
	}
	return RPCLog{RESPONSE, ip, body, time.Now().UTC().String()}
}

func NewRPCLog(isRequest bool, ip string, body string) error {
	if config.GlobalConfig().RPCLogs {
		log := newRPCLog(isRequest, ip, body)
		f, err := json.MarshalIndent(log, "", "  ")
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(config.GlobalConfig().DD+string(filepath.Separator)+"rpc.log", f, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
