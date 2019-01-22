package node

import (
	"fmt"
	"time"
	
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/util"
)

func Register() {
	if _, err := util.RPCRequest("http://"+_const.DISPATCHIP+":"+_const.DISPATCHCLIENTPORT+"/v1/register", GetSelf(), util.POST); err != nil {
		util.ExitGracefully("Error, unable to register node at Pocket Incorporated's Dispatcher")
	}
	fmt.Println("Node registered successfully " , GetPeerList().Map)
}

func UnRegister(count int) {
	if _, err := util.RPCRequest("http://"+_const.DISPATCHIP+":"+_const.DISPATCHCLIENTPORT+"/v1/register", GetSelf(), util.POST); err != nil {
		fmt.Println("Error, unable to unregister node at Pocket Incorporated's Dispatcher, trying again!")
		time.Sleep(2)
		if count > 5 {
			util.ExitGracefully("Please contact Pocket Incorporated with this error! As your node was unable to be unregistered")
		}
		UnRegister(count + 1)
	}
	util.ExitGracefully("you have been unregistered! Thank you for using Pocket Core MVP! Goodbye")
}
