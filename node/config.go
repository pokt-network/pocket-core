package node

import (
	"fmt"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/util"
)

const chainsFileExample = "[{\"blockchain\": {\"name\": \"ethereum\",\"netid\": \"1\",\"version\": \"1.0\"},\"port\": \"8545\",\"medium\": \"rpc\"},{\"blockchain\": {\"name\": \"bitcoin\",\"netid\": \"1\",\"version\": \"1.0\"},\"port\": \"8333\",\"medium\": \"rpc\"}]"

type FileName int

const (
	ChainFile FileName = iota + 1
	PeerFile
	DeveWL
	SerWL
)

func fileErrorMessage(fn FileName) {
	var path string
	var filename string
	switch fn {
	case ChainFile:
		path = config.GlobalConfig().CFile
		filename = "chains"
	case PeerFile:
		path = config.GlobalConfig().PFile
		filename = "peer"
	case DeveWL:
		path = config.GlobalConfig().DWL
		filename = "developer white list file"
	case SerWL:
		path = config.GlobalConfig().SNWL
		filename = "service white list file"
	}
	fmt.Println("There seems to be something wrong with your" + filename + " file @ " + path)
	fmt.Println("Please ensure that it is in the proper format:")
	res, err := util.StringToPrettyJSON(chainsFileExample)
	if err == nil {
		fmt.Println(string(res))
	}
}

func ConfigFiles() error {
	// Map.json
	if err := ManualPeersFile(config.GlobalConfig().PFile); err != nil { // add Map from file
		logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
		fileErrorMessage(PeerFile)
		return err
	}
	// chains.json
	if err := CFile(config.GlobalConfig().CFile); err != nil {
		logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
		fileErrorMessage(ChainFile)
		util.ExitGracefully(err.Error() + " " + config.GlobalConfig().CFile) // currently just exit
	}
	// whitelists for centralized dispatcher
	WhiteListInit()
	if err := SWLFile(); err != nil {
		logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
		if config.GlobalConfig().Dispatch {
			fileErrorMessage(SerWL)
		}
		return err
	}
	if err := DWLFile(); err != nil {
		logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
		if config.GlobalConfig().Dispatch {
			fileErrorMessage(DeveWL)
		}
		return err
	}
	return nil
}
