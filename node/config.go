package node

import (
	"fmt"
	"github.com/pokt-network/pocket-core/const"
	"time"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/util"
)

const chainsFileExample = "[{\"blockchain\": {\"name\": \"ethereum\",\"netid\": \"1\",\"host\":\"localhost\",\"port\": \"8545\",\"medium\": \"rpc\"},{\"blockchain\": {\"name\": \"bitcoin\",\"netid\": \"1\",\"host\":\"localhost\",\"port\": \"8333\",\"medium\": \"rpc\"}]"
const devFileExample = "[\"DEVID1\"]"
const serFileExample = "[\"GID1\"]"

type FileName int

const (
	ChainFile FileName = iota + 1
	DeveWL
	SerWL
)

func fileErrorMessage(fn FileName) {
	var path, filename, example string
	switch fn {
	case ChainFile:
		path = config.GlobalConfig().CFile
		filename = "chains"
		example = chainsFileExample
	case DeveWL:
		path = config.GlobalConfig().DWL
		filename = "developer white list file"
		example = devFileExample
	case SerWL:
		path = config.GlobalConfig().SNWL
		filename = "service white list file"
		example = serFileExample
	}
	fmt.Println("There seems to be something wrong with your " + filename + " file @ " + path)
	fmt.Println("Please ensure that it is in the proper format:")
	res, err := util.StringToPrettyJSON(example)
	if err == nil {
		fmt.Println(string(res))
	}
}

func swlConfigFile() error {
	if err := SWLFile(); err != nil {
		logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
		if config.GlobalConfig().Dispatch {
			fileErrorMessage(SerWL)
			return err
		}
	}
	return nil
}

func chainsConfigFile() {
	// chains.json
	c := config.GlobalConfig().CFile
	if c == _const.CHAINFILEPLACEHOLDER {
		c = config.GlobalConfig().DD+_const.FILESEPARATOR+"chains.json"
	}
	if err := CFile(c); err != nil {
		logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
		fileErrorMessage(ChainFile)
		ExitGracefully(err.Error() + " " + c) // currently just exit
	}
}

func dwlConfigFile() error {
	if err := DWLFile(); err != nil {
		logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
		if config.GlobalConfig().Dispatch {
			fileErrorMessage(DeveWL)
			return err
		}
	}
	return nil
}

// "ConfigFiles" configure the client based off of files in the data directory.
func ConfigFiles() error {
	chainsConfigFile()
	WhiteListInit()
	err := dwlConfigFile()
	err2 := swlConfigFile()
	if err != nil {
		return err
	}
	if err2 != nil {
		return err2
	}
	go WLRefresh()
	return nil
}

// "WLRefresh" updates data structure in memory from file for both whitelists after a certain amount of time.
func WLRefresh() {
	for {
		var err error
		err = dwlConfigFile()
		if err != nil {
			fmt.Println("Error with Developers WL " + err.Error())
		}
		err = swlConfigFile()
		if err != nil {
			fmt.Println("Error with Developers WL " + err.Error())
		}
		if !config.GlobalConfig().Dispatch {
			err := UpdateWhiteList()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		time.Sleep(time.Duration(config.GlobalConfig().PRefresh) * time.Second)
	}
}
