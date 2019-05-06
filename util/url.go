package util

import (
	"errors"
	"github.com/pokt-network/pocket-core/config"
	"net/http"
	"strings"
	"time"
)

const (
	HTTPS = "https://"
	HTTP  = "http://"
)

func URLProto(query string) (string, error) {
	if strings.Contains(query, "//") {
		query = strings.TrimLeft(query, "//")
	}
	_, err := Ping(HTTPS + query)
	if err != nil {
		_, err := Ping(HTTP + query)
		if err != nil {
			return "", errors.New(err.Error())
		}
		return HTTP + query, nil
	}
	return HTTPS + query, nil
}

func Ping(url string) (int, error) {
	client := http.Client{}
	client.Timeout = time.Duration(config.GlobalConfig().RequestTimeout) * time.Millisecond
	req, err := http.NewRequest("HEAD", url, http.NoBody)
	if err != nil {
		return 0, err
	}
	defer req.Body.Close()
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}
