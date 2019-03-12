package node

import (
	"errors"
	"fmt"
	"time"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/util"
)

// "Register" marks a service node 'ready for work' in the database.
func Register() {
	c := config.GlobalConfig()
	s, err := Self()
	if err != nil {
		ExitGracefully("error registering node " + err.Error())
	}
	resp, err := util.StructRPCReq("http://"+c.DisIP+":"+c.DisRPort+"/v1/register", s, util.POST)
	if err != nil {
		ExitGracefully("error registering node " + err.Error())
	}
	fmt.Println(resp)
}

// "Unregister" removes a service node from the database
func UnRegister(count int) error {
	c := config.GlobalConfig()
	s, err := Self()
	if err != nil {
		return err
	}
	if _, err := util.StructRPCReq("http://"+c.DisIP+":"+c.DisRPort+"/v1/unregister", s, util.POST); err != nil {
		fmt.Println("Error, unable to unregister node at Pocket Incorporated's Dispatcher, trying again!")
		time.Sleep(2)
		if count > 5 {
			return errors.New("please contact Pocket Incorporated with this error! As your node was unable to be unregistered")
		}
		UnRegister(count + 1)
	}
	return nil
}
