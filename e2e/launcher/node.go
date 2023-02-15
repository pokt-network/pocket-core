package launcher

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
)

type runtimeStatus int

const (
	NotStarted = runtimeStatus(iota)
	Running
	Finished
)

type Node struct {
	PocketServer PocketServer
	Address      string
	DataDir      string

	configuration types.Config
	privateKey    crypto.PrivateKey
}

func newNode(nodeConfig *NodeConfiguration, networkRootDir, genesisPath, executablePath string) (*Node, error) {
	pkBytes, err := hex.DecodeString(nodeConfig.PrivateKey)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.NewPrivateKeyBz(pkBytes)
	if err != nil {
		return nil, err
	}

	address := strings.ToLower(privateKey.PublicKey().Address().String())
	nodeDataDir, err := os.MkdirTemp(networkRootDir, address+"-")
	if err != nil {
		return nil, err
	}

	pocketCoreConfig, err := loadNodePocketCoreConfiguration(nodeConfig.ConfigPath)
	if err != nil {
		return nil, err
	}

	pocketCoreConfig.TendermintConfig.RootDir = nodeDataDir
	pocketCoreConfig.TendermintConfig.RPC.RootDir = nodeDataDir
	pocketCoreConfig.TendermintConfig.P2P.RootDir = nodeDataDir
	pocketCoreConfig.TendermintConfig.Mempool.RootDir = nodeDataDir
	pocketCoreConfig.TendermintConfig.Consensus.RootDir = nodeDataDir
	pocketCoreConfig.PocketConfig.DataDir = nodeDataDir

	// write config.json & genesis.json
	if err = os.Mkdir(filepath.Join(nodeDataDir, "config"), os.ModePerm); err != nil {
		return nil, err
	}

	if configBz, err := json.Marshal(pocketCoreConfig); err != nil {
		return nil, err
	} else if err = writeBytesToFile(filepath.Join(nodeDataDir, "config", "config.json"), configBz); err != nil {
		return nil, err
	}

	if err = copyFile(genesisPath, filepath.Join(nodeDataDir, "config", "genesis.json")); err != nil {
		return nil, err
	}

	// set identity files
	if err = writeNodePrivValFile(privateKey, nodeDataDir, pocketCoreConfig.TendermintConfig.PrivValidatorKey); err != nil {
		return nil, err
	}

	if err = writeNodeKeyFile(privateKey, nodeDataDir, pocketCoreConfig.TendermintConfig.NodeKey); err != nil {
		return nil, err
	}

	if err = writePrivValState(nodeDataDir, pocketCoreConfig.TendermintConfig.PrivValidatorState); err != nil {
		return nil, err
	}

	node := &Node{
		Address:       address,
		DataDir:       nodeDataDir,
		configuration: pocketCoreConfig,
		privateKey:    privateKey,
		PocketServer:  NewPocketServer(executablePath),
	}
	return node, nil
}

func writePrivValState(datadir, privValState string) error {
	pvkBz, err := app.Codec().MarshalJSONIndent(privval.FilePVLastSignState{}, "", "  ")
	if err != nil {
		return err
	}
	return writeBytesToFile(filepath.Join(datadir, privValState), pvkBz)
}

func writeNodeKeyFile(privateKey crypto.PrivateKey, datadir string, tmNodeKey string) error {
	cdc := app.Codec()
	nodeKey := p2p.NodeKey{
		PrivKey: privateKey.PrivKey(),
	}
	pvkBz, err := cdc.MarshalJSONIndent(nodeKey, "", "  ")
	if err != nil {
		return err
	}

	return writeBytesToFile(filepath.Join(datadir, tmNodeKey), pvkBz)
}

func writeNodePrivValFile(privateKey crypto.PrivateKey, datadir string, privValidatorKey string) error {
	cdc := app.Codec()
	privValKey := privval.FilePVKey{
		Address: privateKey.PubKey().Address(),
		PubKey:  privateKey.PubKey(),
		PrivKey: privateKey.PrivKey(),
	}
	pvkBz, err := cdc.MarshalJSONIndent(privValKey, "", "  ")
	if err != nil {
		return err
	}

	return writeBytesToFile(filepath.Join(datadir, privValidatorKey), pvkBz)
}

func loadNodePocketCoreConfiguration(path string) (pcc types.Config, err error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(contents, &pcc)
	if err != nil {
		return
	}

	return
}
