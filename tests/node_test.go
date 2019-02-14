package tests

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/node"
)

func writeSampleConfigFiles() error {
	datadirectory := config.GlobalConfig().DD
	cFile, err := filepath.Abs("./fixtures/chains.json")
	if err != nil {
		return err
	}
	pFile, err := filepath.Abs("./fixtures/peers.json")
	if err != nil {
		return err
	}
	dwl, err := filepath.Abs("./fixtures/developer_whitelist.json")
	if err != nil {
		return err
	}
	swl, err := filepath.Abs("./fixtures/service_whitelist.json")
	if err != nil {
		return err
	}
	err = copyFile(cFile, datadirectory+"/chains.json")
	if err != nil {
		return err
	}
	err = copyFile(pFile, datadirectory+"/peers.json")
	if err != nil {
		return err
	}
	copyFile(dwl, datadirectory+"/developer_whitelist.json")
	if err != nil {
		return err
	}
	copyFile(swl, datadirectory+"/service_whitelist.json")
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

func TestDispatchPeers(t *testing.T) {
	n := dummyNode()
	chain := n.Blockchains[0]
	dp := node.DispatchPeers()
	dp.Add(n)
	result := dp.PeersByChain(chain)
	if len(result) != 1 {
		t.Fatalf("DispatchPeers.add(Node) returned a result other than length of 1")
	}
}

func TestPeerList(t *testing.T) {
	n := dummyNode()
	pl := node.PeerList()
	pl.Clear()
	dp := node.DispatchPeers()
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
	pl.CopyToDP()
	if len(dp.Map) != 1 {
		t.Fatalf("PeerList.CopyToDispatch(Node) resulted in len(DispatchPeers) != 1")
	}
	pl.Remove(n)
	if pl.Count() != 0 {
		t.Fatalf("PeerList.Remove(Node) did not appear to remove the node")
	}
}

func TestSelf(t *testing.T) {
	if _, err := node.Self(); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestWhiteList(t *testing.T) {
	if err := writeSampleConfigFiles(); err != nil {
		t.Fatalf(err.Error())
	}
	node.DWLFile()
	dwl := node.DWL()
	if dwl.Count() != len(dwl.M) {
		t.Fatalf("WhiteList.Count() does not equal len(dwl)")
	}
	if dwl.Count() != 1 {
		t.Fatal("After calling WhiteListFile() the number of developers in the white list was not 1")
	}
	dwl.Add("TEST")
	if len(dwl.M) != 2 {
		t.Fatalf("WhiteList.Add(ID) did not result in adding one more developer to the white list")
	}
	if !dwl.Contains("TEST") {
		t.Fatalf("WhiteList.Contains(ID) did not find the added develeoper")
	}
	dwl.Remove("TEST")
	if dwl.Contains("TEST") {
		t.Fatalf("After calling WhiteList.Remove(ID) the ID still exists")
	}
}
