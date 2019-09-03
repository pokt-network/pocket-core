package fixtures

import (
	"encoding/json"
	"github.com/pokt-network/pocket-core/types"
	"os"
)

var (
	gp     = os.Getenv("GOPATH")
	source = gp + string(os.PathSeparator) + "src" + string(os.PathSeparator) + "github.com" + string(os.PathSeparator) +
		"pokt-network" + string(os.PathSeparator) + "pocket-core" + string(os.PathSeparator) + "tests" +
		string(os.PathSeparator) + "fixtures" + string(os.PathSeparator) + "JSON" + string(os.PathSeparator)
)

func GetDevelopers() (*types.Developers, error) {
	result := &types.Developers{}
	jsonFile, err := os.Open(source + "randomDeveloperPool.json")
	if err != nil {
		return nil, err
	}
	jsonParser := json.NewDecoder(jsonFile)
	if err = jsonParser.Decode(result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetNodes() (*types.Nodes, error) {
	result := &types.Nodes{}
	jsonFile, err := os.Open(source + "randomNodePool.json")
	if err != nil {
		return nil, err
	}
	jsonParser := json.NewDecoder(jsonFile)
	if err = jsonParser.Decode(result); err != nil {
		return nil, err
	}
	return result, nil
}
