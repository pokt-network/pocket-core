package types

type Blockchain AminoBuffer

type Blockchains map[string]struct{}

func (bc Blockchain) String() string {
	return AminoBuffer(bc).String()
}

func (bcs Blockchains) GetChainURL(blockchain string) (string, error) {
	// todo
	return "8.8.8.8", nil
}

func (bcs Blockchains) Contains(blockchain string) bool {
	// todo
	for k := range bcs {
		if k == blockchain {
			return true
		}
	}
	return false
}
