package mtf

import "fmt"

const AlphabetSize = 256

type MTFAlgorithm struct {
	alphabet []byte
}

func NewMTFAlgorithm(alphabet []byte) *MTFAlgorithm {
	return &MTFAlgorithm{
		alphabet: alphabet,
	}
}

func (t *MTFAlgorithm) Transform(bytes []byte) []byte {
	res := make([]byte, len(bytes))
	alphabet := make([]int, AlphabetSize)

	for _, b := range bytes {
		alphabet[b]++
	}

	for i, count := range alphabet {
		if count != 0 {
			t.alphabet = append(t.alphabet, byte(i))
		}
	}

	alphabetCopy := make([]byte, len(t.alphabet))
	copy(alphabetCopy, t.alphabet)

	for i, b := range bytes {
		res[i] = byte(t.findIndexInAlphabet(b))
		t.moveToFront(int(res[i]))
	}

	t.alphabet = alphabetCopy

	return res
}

func (t *MTFAlgorithm) Retransform(bytes []byte) []byte {
	res := make([]byte, len(bytes))
	for i, b := range bytes {
		res[i] = t.alphabet[b]
		t.moveToFront(int(b))
	}
	return res
}

func (t *MTFAlgorithm) findIndexInAlphabet(b byte) int {
	for i, c := range t.alphabet {
		if c == b {
			return i
		}
	}
	return -1
}

func (t *MTFAlgorithm) moveToFront(index int) {
	if index > 0 {
		temp := t.alphabet[index]
		copy(t.alphabet[1:index+1], t.alphabet[:index])
		t.alphabet[0] = temp
	}
}

func (t *MTFAlgorithm) GetAlphabet() string {
	res := ""
	for _, ch := range t.alphabet {
		res += fmt.Sprintf("%d,", ch)
	}
	if len(res) != 0 {
		res = res[:len(res)-1]
	}
	return res
}
