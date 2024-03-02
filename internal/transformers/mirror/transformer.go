package mirror

import (
	"fmt"
	"io"
	"os"
)

const (
	TempFile = "temp.txt"
	BufSize  = 4096
)

type MirrorTransformer struct {
	archive string
}

func New(archiveName string) *MirrorTransformer {
	return &MirrorTransformer{
		archive: archiveName,
	}
}

func (t *MirrorTransformer) Transform() error {
	f, err := os.Open(t.archive)
	if err != nil {
		return err
	}

	temp, err := os.OpenFile(TempFile, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer func() {
		temp.Close()
		os.Remove(TempFile)
	}()

	buf := make([]byte, BufSize)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("%w", err)
		}
		transformAlg := NewMirrorAlgorithm()
		temp.Write(transformAlg.Transform(buf[:n]))
	}

	f.Close()
	f, err = os.OpenFile(t.archive, os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	temp.Seek(0, 0)
	for {
		n, err := temp.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("%w", err)
		}
		_, err = f.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *MirrorTransformer) Retransform() error {
	f, err := os.Open(t.archive)
	if err != nil {
		return err
	}

	temp, err := os.OpenFile(TempFile, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer func() {
		temp.Close()
		os.Remove(TempFile)
	}()

	buf := make([]byte, BufSize)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("%w", err)
		}
		transformAlg := NewMirrorAlgorithm()
		temp.Write(transformAlg.Retransform(buf[:n]))
	}

	f.Close()
	f, err = os.OpenFile(t.archive, os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	temp.Seek(0, 0)
	for {
		n, err := temp.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("%w", err)
		}
		_, err = f.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return nil
}
