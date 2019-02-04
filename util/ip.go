package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func IP() (string, error) {
	url := "https://api.ipify.org?format=text"
	fmt.Printf("Getting IP address from  ipify ...\n")
	resp, err := http.Get(url)
	if err != nil {
		return ",", err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ",", err
	}
	fmt.Println(string(ip))
	return string(ip), nil
}
