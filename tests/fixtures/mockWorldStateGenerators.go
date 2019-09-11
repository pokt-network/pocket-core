package fixtures

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	//"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/legacy"
	"github.com/pokt-network/pocket-core/types"
	tdmnt "github.com/tendermint/tendermint/types"
	"io/ioutil"
	"math/rand"
	"time"
)

const (
	numberOfNodes        = 50
	numberOfApplications = 50
)

var (
	tickers = []string{"eth", "btc", "ltc"}
)

// writes json nodepool for testing
func GenerateAliveNodes() {
	RegisterPOKT()
	var result types.Nodes
	fmt.Println()
	for i := 0; i < numberOfNodes; i++ {
		result = append(result, GenerateAliveNode())
	}
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = ioutil.WriteFile("tests/fixtures/JSON/randomNodePool.json", output, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GenerateAliveNode() (node types.Node) {
	randomSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randomSource)
	balance := GeneratePOKT(random.Int63())
	stakeAmount := GeneratePOKT(random.Int63())
	nsc := types.NodeSupportedChains{}
	nsc.Add(hex.EncodeToString(GenerateNonNativeBlockchain()), types.NodeSupportedChain{})
	privateKey, _ := legacy.NewPrivateKey()
	pubKey := legacy.NewPublicKey(privateKey)
	pubKeyBytes := legacy.FromECDSAPub(pubKey)
	hexPubKey := types.AccountPublicKey(hex.EncodeToString(pubKeyBytes))
	node = types.Node{
		Account: types.Account{
			Address:     nil, // todo
			PubKey:      hexPubKey,
			Balance:     balance,
			StakeAmount: stakeAmount,
		},
		URL:             nil, // todo
		SupportedChains: nsc, // just one for now
		IsAlive:         true,
	}

	return
}

// writes json nodepool for testing
func GenerateApplications() {
	RegisterPOKT()
	var result types.Applications
	fmt.Println()
	for i := 0; i < numberOfApplications; i++ {
		result = append(result, GenerateApplication())
	}
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = ioutil.WriteFile("tests/fixtures/JSON/randomApplicationPool.json", output, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GenerateApplication() (node types.Application) {
	randomSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randomSource)
	balance := GeneratePOKT(random.Int63())
	stakeAmount := GeneratePOKT(random.Int63())
	privateKey, _ := legacy.NewPrivateKey()
	pubKey := legacy.NewPublicKey(privateKey)
	pubKeyBytes := legacy.FromECDSAPub(pubKey)
	hexPubKey := types.AccountPublicKey(hex.EncodeToString(pubKeyBytes))
	node = types.Application{
		Account: types.Account{
			Address:     nil, // todo
			PubKey:      hexPubKey,
			Balance:     balance,
			StakeAmount: stakeAmount,
		},
		RequestedChains: nil, //todo
	}

	return
}

func RegisterPOKT() {
	err := sdk.RegisterDenom("pokt", sdk.NewDec(0))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GeneratePOKT(numberOf int64) types.POKT {
	test := types.POKT(sdk.NewCoin("pokt", sdk.NewInt(numberOf)))
	return test
}

func GenerateBlockHash() types.BlockID {
	seed := make([]byte, 10)
	rand.Read(seed)
	ranHash := legacy.SHA3FromBytes(seed)
	return types.BlockID(tdmnt.BlockID{
		Hash:        ranHash,
		PartsHeader: tdmnt.PartSetHeader{},
	})
}

func GenerateNonNativeBlockchain() []byte {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return legacy.SHA3FromBytes([]byte(tickers[r1.Intn(3)]))
}
