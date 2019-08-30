package model

type AminoBuffer []byte

func (a AminoBuffer) Append(aminoBuffers ...AminoBuffer) AminoBuffer {
	result := a
	for _, amino := range aminoBuffers {
		result = append(result, []byte(amino)...)
	}
	return result
}
