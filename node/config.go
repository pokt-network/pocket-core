package node

import (
	"fmt"
	"time"
	
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/util"
)

const chainsFileExample = "[{\"blockchain\": {\"name\": \"ethereum\",\"netid\": \"1\",\"version\": \"1.0\"},\"host\":\"localhost\",\"port\": \"8545\",\"medium\": \"rpc\"},{\"blockchain\": {\"name\": \"bitcoin\",\"netid\": \"1\",\"version\": \"1.0\"},\"host\":\"localhost\",\"port\": \"8333\",\"medium\": \"rpc\"}]"
const peerFileExample = "[{\"gid\":\"gid1\",\"ip\":\"localhost\",\"relayport\":\"8080\",\"clientport\":\"8081\",\"clientid\":\"pocket_core\",\"cliversion\":\"0.0.1\",\"blockchains\":[{\"name\":\"ethereum\", \"version\":\"0\",\"netid\":\"0\"}]}]"
const devFileExample = "[\"DEVID1\"]"
const serFileExample = "[\"GID1\"]"

type FileName int

const (
	ChainFile FileName = iota + 1
	PeerFile
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
	case PeerFile:
		path = config.GlobalConfig().PFile
		filename = "peer"
		example = peerFileExample
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

func peerConfigFile() error {
	// Map.json
	if err := ManualPeersFile(config.GlobalConfig().PFile); err != nil { // add Map from file
		logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
		fileErrorMessage(PeerFile)
		return err
	}
	return nil
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
	if err := CFile(config.GlobalConfig().CFile); err != nil {
		logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
		fileErrorMessage(ChainFile)
		ExitGracefully(err.Error() + " " + config.GlobalConfig().CFile) // currently just exit
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

func ConfigFiles() error {
	chainsConfigFile()
	err1 := peerConfigFile()
	WhiteListInit()
	err2 := dwlConfigFile()
	err3 := swlConfigFile()
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	if err3 != nil {
		return err3
	}
	if config.GlobalConfig().Dispatch {
		go WLRefresh()
	}
	return nil
}

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
		time.Sleep(time.Duration(config.GlobalConfig().PRefresh) * time.Second)
	}
}
