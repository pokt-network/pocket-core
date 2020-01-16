module github.com/pokt-network/pocket-core

go 1.13

require (
	github.com/btcsuite/btcd v0.0.0-20190824003749-130ea5bddde3 // indirect
	github.com/glacjay/goini v0.0.0-20161120062552-fd3024d87ee2 // indirect
	github.com/go-kit/kit v0.9.0
	github.com/julienschmidt/httprouter v1.2.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pokt-network/posmint v0.0.0-20200115125342-953ec9c4ec98
	github.com/prometheus/client_golang v1.1.0 // indirect
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/prometheus/procfs v0.0.4 // indirect
	github.com/rakyll/globalconf v0.0.0-20180912185831-87f8127c421f
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.32.7
	github.com/tendermint/tm-db v0.2.0
	github.com/wealdtech/go-merkletree v1.0.0
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55 // indirect
	gopkg.in/h2non/gock.v1 v1.0.15
)

replace github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.14
