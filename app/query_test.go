package app

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pokt-network/posmint/crypto/keys"
	"github.com/tendermint/tendermint/node"
)

// TODO list
// - start tm node
// - find ways to inject data into the tmnode...
var _ = Describe("App", func() {
	var (
		tmNode  *node.Node
		pswrd   string = "SOME_W$RD_P4W000RDD"
		nodeUri = ":8080"
		dbKeybase keys.KeyPair
		err error
	)
	BeforeEach(func() {
		dbKeybase, err = keys.NewInMemory().Create(pswrd)
		if err != nil {
			panic(err)
		}
		MakeCodec()

		InitGenesis()

		setcoinbasePassphrase(pswrd)
		setTmNode(nodeUri)
		tmNode = InitTendermint("", "")
	})
	Describe("Query block", func() {
		It("Gets node height", func() {
			Expect(tmNode).To(Equal(tmNode))
			Expect(dbKeybase).To(Equal(dbKeybase))
		})
	})
})
