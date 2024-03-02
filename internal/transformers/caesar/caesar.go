package caesar

const AlphabetSize = 256

type CaesarAlgorithm struct {
	shift int
}

func NewCaesarAlgorithm(shift int) *CaesarAlgorithm {
	if shift >= 0 {
		shift %= AlphabetSize
	} else {
		shift = AlphabetSize - (shift % AlphabetSize)
	}

	return &CaesarAlgorithm{
		shift: shift,
	}
}

func (a *CaesarAlgorithm) Transform(bytes []byte) []byte {
	return a.caesarShift(bytes, a.shift)
}

func (a *CaesarAlgorithm) Retransform(bytes []byte) []byte {
	return a.caesarShift(bytes, -a.shift)
}

func (a *CaesarAlgorithm) caesarShift(bytes []byte, shift int) []byte {
	res := make([]byte, len(bytes))
	for i, c := range bytes {
		res[i] = byte((int(c) + shift) % AlphabetSize)
	}
	return res
}
