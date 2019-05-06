package util

import (
	"io/ioutil"
	"net/http"
	
	"github.com/pokt-network/pocket-core/logs"
)

// "IP" returns the public ip of the user
func IP() (string, error) {
	url := "https://api.ipify.org?format=text"
	resp, err := http.Get(url)
	logs.NewLog("Getting IP address from ipify.org...", logs.InfoLevel, logs.JSONLogFormat)
	if err != nil {
		return "", err
<<<<<<< 57ceb161d287776fc08ba212726bb3bf39a278c6
	}
	if resp != nil {
		if resp.Body!=nil{
			defer resp.Body.Close()
		}
=======
>>>>>>> fixed nil pointer error
	}
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ",", err
	}
	logs.NewLog("IP returned: "+string(ip)+" from ipify!", logs.InfoLevel, logs.JSONLogFormat)
	return string(ip), nil
}
