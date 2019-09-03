package main

import "github.com/pokt-network/pocket-core/tests/fixtures"

func main() {
	GenerateAccounts()
}

func GenerateAccounts() {
	fixtures.GenerateAliveNodes()
	fixtures.GenerateDevelopers()
}
