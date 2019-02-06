package util

import (
	"io/ioutil"
	"net/http"
)

func IP() (string, error) {
	url := "https://api.ipify.org?format=text"
	resp, err := http.Get(url)
	if err != nil {
		return ",", err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ",", err
	}
	return string(ip), nil
}
