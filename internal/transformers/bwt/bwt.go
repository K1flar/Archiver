package bwt

import (
	"sort"
)

const AlphabetSize = 256

type BWTAlgorithm struct {
	code int
}

func NewBWTAlgorithm(code int) *BWTAlgorithm {
	return &BWTAlgorithm{
		code: code,
	}
}

func (a *BWTAlgorithm) Transform(bytes []byte) []byte {
	if len(bytes) == 0 {
		return []byte{}
	}
	res := make([]byte, len(bytes))
	rotations := make([][]byte, len(bytes))

	for i := range bytes {
		rotations[i] = []byte(string(bytes[i:]) + string(bytes[:i]))
	}

	sort.Slice(rotations, func(i, j int) bool { return string(rotations[i]) < string(rotations[j]) })

	for i, r := range rotations {
		if string(r) == string(bytes) {
			a.code = i
		}
		res[i] = r[len(r)-1]
	}

	return res
}

func (a *BWTAlgorithm) Retransform(bytes []byte) []byte {
	res := make([]byte, len(bytes))
	vec := make([]int, len(bytes))

	count := make([]int, AlphabetSize)
	for _, b := range bytes {
		count[b]++
	}

	sum := 0
	for i := range count {
		sum += count[i]
		count[i] = sum - count[i]
	}

	for i := range bytes {
		vec[count[bytes[i]]] = i
		count[bytes[i]]++
	}

	pos := vec[a.code]
	for i := range bytes {
		res[i] = bytes[pos]
		pos = vec[pos]
	}
	return res
}

func (a *BWTAlgorithm) GetCode() int {
	return a.code
}
