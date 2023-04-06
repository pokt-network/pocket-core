package mesh

import (
	"encoding/json"
	"github.com/pokt-network/pocket-core/app"
	sdk "github.com/pokt-network/pocket-core/types"
	"io/ioutil"
	log2 "log"
	"os"
)

// getAuthTokenFromFile - read from path a json that match sdk.AuthToken struct
func getAuthTokenFromFile(path string) sdk.AuthToken {
	logger.Info("reading authtoken from path=" + path)
	t := sdk.AuthToken{}

	var jsonFile *os.File
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			return
		}
	}(jsonFile)

	if _, err := os.Stat(path); err == nil {
		jsonFile, err = os.OpenFile(path, os.O_RDONLY, os.ModePerm)
		if err != nil {
			log2.Fatalf("cannot open auth token json file: " + err.Error())
		}
		b, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log2.Fatalf("cannot read auth token json file: " + err.Error())
		}
		err = json.Unmarshal(b, &t)
		if err != nil {
			log2.Fatalf("cannot read auth token json file into json: " + err.Error())
		}
	}

	return t
}

// loadAuthTokens - load mesh node authtoken and servicer authtoken
func loadAuthTokens() {
	dataDir := app.GlobalMeshConfig.DataDir
	meshNodeAuthFile := dataDir + app.FS + app.GlobalMeshConfig.AuthTokenFile
	servicerAuthFile := dataDir + app.FS + app.GlobalMeshConfig.ServicerAuthTokenFile
	// used to authenticate request to mesh node on /v1/private paths
	meshAuthToken = getAuthTokenFromFile(meshNodeAuthFile)
	// used to call servicer node on private path to notify about relays
	servicerAuthToken = getAuthTokenFromFile(servicerAuthFile)
}

// GetAuthToken - return mesh auth token value load from auth file
func GetAuthToken() string {
	return meshAuthToken.Value
}
