package unit

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/node"
)

func writeSampleConfigFiles() error {
	datadirectory := config.GlobalConfig().DD + _const.FILESEPARATOR
	cFile, err := filepath.Abs("fixtures" + _const.FILESEPARATOR + "chains.json")
	if err != nil {
		return err
	}
	pFile, err := filepath.Abs("fixtures" + _const.FILESEPARATOR + "peers.json")
	if err != nil {
		return err
	}
	err = copyFile(cFile, datadirectory+"chains.json")
	if err != nil {
		return err
	}
	err = copyFile(pFile, datadirectory+"peers.json")
	if err != nil {
		return err
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func dummyNode() node.Node {
	chains := []node.Blockchain{{Name: "ethereum", NetID: "1", Version: "1"}}
	n := node.Node{
		GID:         "test",
		IP:          "123",
		RelayPort:   "0",
		ClientPort:  "0",
		ClientID:    "0",
		CliVersion:  "0",
		Blockchains: chains,
	}
	return n
}

func TestConfigFiles(t *testing.T) {
	if err := writeSampleConfigFiles(); err != nil {
		t.Fatalf(err.Error())
	}
	if err := node.ConfigFiles(); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestChains(t *testing.T) {
	if err := writeSampleConfigFiles(); err != nil {
		t.Fatalf(err.Error())
	}
	if err := node.CFile(config.GlobalConfig().CFile); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestPeerList(t *testing.T) {
	n := dummyNode()
	pl := node.PeerList()
	pl.Clear()
	pl.Add(n)
	if pl.Count() != len(pl.M) {
		t.Fatalf("PeerList.Count() does not equal len(PeerList)")
	}
	if pl.Count() != 1 {
		t.Fatalf("PeerList.add(Node) returned a result other than length of 1")
	}
	if !pl.Contains(n.GID) {
		t.Fatalf("Peerlist.Contains(Node) returned false when it should contain the node")
	}
	pl.Remove(n)
	if pl.Count() != 0 {
		t.Fatalf("PeerList.Remove(Node) did not appear to remove the node")
	}
}
