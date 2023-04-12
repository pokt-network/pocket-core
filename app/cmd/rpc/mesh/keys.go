package mesh

import "github.com/pokt-network/pocket-core/app"

type nodeFileItem struct {
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Keys []string `json:"keys"`
}

type fallbackNodeFileItem struct {
	PrivateKey  string `json:"priv_key"`
	ServicerUrl string `json:"servicer_url"`
}

// getServicersFilePath - return servicers file path resolved by config.json
func getServicersFilePath() string {
	return app.GlobalMeshConfig.DataDir + app.FS + app.GlobalMeshConfig.ServicerPrivateKeyFile
}
