package logs

import (
	"encoding/json"
	"github.com/pokt-network/pocket-core/config"
	"os"
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
	URL       string `json:"path"`
	Payload   string `json:"payload"`
	Timestamp string `json:"timestamp"`
}

func newRPCLog(isRequest bool, url, ip, body string) RPCLog {
	if isRequest {
		return RPCLog{REQUEST, ip, url, body, time.Now().UTC().String()}
	}
	return RPCLog{RESPONSE, ip, url, body, time.Now().UTC().String()}
}

func NewRPCLog(isRequest bool, url, ip, body string) error {
	if config.GlobalConfig().RPCLogs {
		log := newRPCLog(isRequest, url, ip, body)
		j, err := json.MarshalIndent(log, "", "  ")
		if err != nil {
			return err
		}
		f, err := os.OpenFile(config.GlobalConfig().DD+string(filepath.Separator)+"access.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err = f.Write(j); err != nil {
			return err
		}
	}
	return nil
}
