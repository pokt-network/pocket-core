package crypto

type Signature []byte // todo

func MockVerifySignature(publicKey, messageHash, signature []byte) bool{
	//todo
	if len(signature) == 0 || len(publicKey) ==0 || len(messageHash) == 0 {
		return false
	}
	return true
}
