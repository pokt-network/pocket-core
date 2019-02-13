package node

import (
	"fmt"
	"time"
	
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/util"
)

func Register() {
	c := config.GlobalConfig()
	resp, err := util.RPCRequest("http://"+c.DisIP+":"+c.DisCPort+"/v1/register", Self(), util.POST)
	if err != nil {
		util.ExitGracefully(err.Error())
	}
	fmt.Println(resp)
}

func UnRegister(count int) {
	c := config.GlobalConfig()
	if _, err := util.RPCRequest("http://"+c.DisIP+":"+c.DisCPort+"/v1/unregister", Self(), util.POST); err != nil {
		fmt.Println("Error, unable to unregister node at Pocket Incorporated's Dispatcher, trying again!")
		time.Sleep(2)
		if count > 5 {
			util.ExitGracefully("Please contact Pocket Incorporated with this error! As your node was unable to be unregistered")
		}
		UnRegister(count + 1)
	}
	util.ExitGracefully("you have been unregistered! Thank you for using Pocket Core MVP! Goodbye")
}
