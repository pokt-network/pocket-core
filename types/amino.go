package types

type AminoBuffer []byte

func (a AminoBuffer) Append(aminoBuffers ...AminoBuffer) AminoBuffer {
	result := a
	for _, amino := range aminoBuffers {
		result = append(result, []byte(amino)...)
	}
	return result
}

func (a AminoBuffer) String() string {
	//todo
	return string(a)
}
