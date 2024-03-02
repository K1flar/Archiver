package mirror

type MirrorAlgorithm struct{}

func NewMirrorAlgorithm() *MirrorAlgorithm {
	return &MirrorAlgorithm{}
}

func (a *MirrorAlgorithm) Transform(bytes []byte) []byte {
	return a.transform(bytes)
}

func (a *MirrorAlgorithm) Retransform(bytes []byte) []byte {
	return a.transform(bytes)
}

func (a *MirrorAlgorithm) transform(bytes []byte) []byte {
	res := make([]byte, len(bytes))
	for i := 0; i < len(bytes)/2; i++ {
		res[i] = bytes[len(bytes)-i-1]
		res[len(res)-i-1] = bytes[i]
	}
	return res
}
