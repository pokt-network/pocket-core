package util

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

const (
	HTTPS   = "https://"
	HTTP    = "http://"
	TIMEOUT = time.Duration(200 * time.Millisecond)
)

func URLProto(query string) (string, error) {
	if strings.Contains(query, "//") {
		query = strings.TrimLeft(query, "//")
	}
	_, err := Ping(HTTPS + query)
	if err != nil {
		_, err := Ping(HTTP + query)
		if err != nil {
			return "", errors.New("unreachable site error: " + err.Error())
		}
		return HTTP + query, nil
	}
	return HTTPS + query, nil
}

func Ping(url string) (int, error) {
	client := http.Client{}
	client.Timeout = TIMEOUT
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
	return resp.StatusCode, nil
}
