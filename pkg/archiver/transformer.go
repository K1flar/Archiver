package archiver

type Transformer interface {
	Transform() error
}

type Retransformer interface {
	Retransform() error
}

func (a *Archiver) Transform(transformer Transformer) error {
	return transformer.Transform()
}

func (a *Archiver) Retransform(retransformer Retransformer) error {
	return retransformer.Retransform()
}
