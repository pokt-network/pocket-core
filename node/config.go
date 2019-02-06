package node

import (
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/util"
)

func Files() {
	// Map.json
	if err := ManualPeersFile(config.GlobalConfig().PFile); err != nil { // add Map from file
		logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
	}
	// chains.json
	if err := CFIle(config.GlobalConfig().CFile); err != nil {
		logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
		util.ExitGracefully(err.Error() + config.GlobalConfig().CFile)
	}
	// whitelists for centralized dispatcher
	WhiteListInit()
	if err := SWLFile(); err != nil {
		logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
	}
	if err := DWLFile(); err != nil {
		logs.NewLog(err.Error(), logs.WaringLevel, logs.JSONLogFormat)
	}
}
