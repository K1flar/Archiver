package archiver

import (
	"archiver/internal/transformers/bwt"
	"archiver/internal/transformers/mtf"
)

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

func (a *Archiver) NewBWTTransformer() *bwt.BWTTransformer {
	return bwt.New(a.archive)
}

func (a *Archiver) NewMTFAlgorithm() *mtf.MTFTransformer {
	return mtf.New(a.archive)
}

// type Transform func([]byte) []byte

// type Transformer interface {
// 	Transform([]byte) []byte
// }

// type Retransformer interface {
// 	Retransform([]byte) []byte
// }

// func (a *Archiver) Transform(transformer Transformer) error {
// 	return a.transformAlg(transformer.Transform)
// }

// func (a *Archiver) Retransform(retransformer Retransformer) error {
// 	return a.transformAlg(retransformer.Retransform)
// }

// func (a *Archiver) NewBWTTransformer() (*bwt.BWTTransformer, error) {
// 	f, err := os.Open(a.archive)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()

// 	scanner := bufio.NewScanner(f)
// 	scanner.Scan()
// 	// header := scanner.Text()
// 	// info

// 	// bwt.New()

// 	return nil, nil
// }

// func (a *Archiver) transformAlg(transform Transform) error {
// 	f, err := os.Open(a.archive)
// 	if err != nil {
// 		return err
// 	}

// 	temp, err := os.OpenFile(TempFile, os.O_RDWR|os.O_CREATE, 0777)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		temp.Close()
// 		os.Remove(TempFile)
// 	}()

// 	scanner := bufio.NewScanner(f)
// 	scanner.Scan()
// 	header := scanner.Text()

// 	_, err = temp.WriteString(header + "\n")
// 	if err != nil {
// 		return err
// 	}

// 	for scanner.Scan() {
// 		row := scanner.Text()
// 		temp.Write(transform([]byte(row)))
// 	}

// 	f.Close()
// 	f, err = os.OpenFile(a.archive, os.O_WRONLY|os.O_TRUNC, 0777)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	buf := make([]byte, BufSize)
// 	temp.Seek(0, 0)
// 	for {
// 		n, err := temp.Read(buf)
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			return fmt.Errorf("%w", err)
// 		}
// 		_, err = f.Write(buf[:n])
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
