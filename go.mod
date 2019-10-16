module github.com/pokt-network/pocket-core

go 1.12

require (
	github.com/btcsuite/btcd v0.0.0-20190824003749-130ea5bddde3 // indirect
	github.com/cosmos/cosmos-sdk v0.37.1
	github.com/ethereum/go-ethereum v1.9.3
	github.com/google/flatbuffers v0.0.0-20190424190944-bf9eb67ab937
	github.com/gorilla/mux v1.7.0
	github.com/h2non/gock v1.0.15
	github.com/julienschmidt/httprouter v1.2.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/prometheus/client_golang v1.1.0 // indirect
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/prometheus/procfs v0.0.4 // indirect
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/cobra v0.0.5
	github.com/tendermint/tendermint v0.32.5
	github.com/tendermint/tm-db v0.2.0
	golang.org/x/crypto v0.0.0-20190426145343-a29dc8fdc734
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/appengine v1.4.0 // indirect
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55 // indirect
)

replace github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.14
