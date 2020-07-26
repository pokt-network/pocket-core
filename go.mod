module github.com/pokt-network/pocket-core

go 1.13

require (
	github.com/go-kit/kit v0.10.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/julienschmidt/httprouter v1.3.0
	github.com/pokt-network/posmint v0.0.0-20200726152834-2651cae46d68
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/iavl v0.12.4
	github.com/tendermint/tendermint v0.32.10
	github.com/tendermint/tm-db v0.2.0
	github.com/willf/bitset v1.1.10 // indirect
	github.com/willf/bloom v2.0.3+incompatible
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79
	gopkg.in/h2non/gock.v1 v1.0.15
)

replace github.com/tendermint/tendermint => github.com/pokt-network/tendermint v0.32.11-0.20200717190033-b520221b1537
