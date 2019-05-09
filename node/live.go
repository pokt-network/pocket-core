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
		ExitGracefully(err.Error())
	}
	u, err := util.URLProto(c.DisIP + ":" + c.DisRPort + "/v1/register")
	if err != nil {
		ExitGracefully(err.Error())
	}
	resp, err := util.StructRPCReq(u, s, util.POST)
	if err != nil {
		ExitGracefully(err.Error())
	}
}

// "Unregister" removes a service node from the database
func UnRegister(count int) error {
	c := config.GlobalConfig()
	s, err := Self()
	if err != nil {
		return err
	}
	u, err := util.URLProto(c.DisIP + ":" + c.DisRPort + "/v1/unregister")
	if err != nil {
		return errors.New(err.Error())
	}
	if _, err := util.StructRPCReq(u, s, util.POST); err != nil {
		fmt.Println("Error, unable to unregister node at Pocket Incorporated's Dispatcher, trying again!")
		time.Sleep(2)
		if count > 5 {
			return errors.New("please contact Pocket Incorporated with this error! As your node was unable to be unregistered")
		}
		UnRegister(count + 1)
	}
	return nil
}
