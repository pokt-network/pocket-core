package types

type Blockchain string

type Blockchains []Blockchain

func (b Blockchains) Contains(blockchain Blockchain) bool {
	for _, val := range b {
		if val == blockchain {
			return true
		}
	}
	return false
}
