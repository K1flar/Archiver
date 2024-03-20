package bwt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	TempFile = "temp.txt"
	BufSize  = 1024
)

type BWTTransformer struct {
	archive string
}

func New(archiveName string) *BWTTransformer {
	return &BWTTransformer{
		archive: archiveName,
	}
}

func (t *BWTTransformer) Transform() error {
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

	var header strings.Builder
	buf := make([]byte, BufSize)
	f.Seek(0, 0)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("%w", err)
		}
		transformAlg := NewBWTAlgorithm(0)
		temp.Write(transformAlg.Transform(buf[:n]))
		header.WriteString(fmt.Sprintf("%d;", transformAlg.GetCode()))
	}

	f.Close()
	f, err = os.OpenFile(t.archive, os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(header.String() + "\n")
	if err != nil {
		return err
	}

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

func (t *BWTTransformer) Retransform() error {
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

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	header := scanner.Text()

	codes := strings.Split(header[:len(header)-1], ";")

	i := 0
	buf := make([]byte, BufSize)
	f.Seek(int64(len(header)+1), 0)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("%w", err)
		}
		if i >= len(codes) {
			return fmt.Errorf("bwt retransform: bad header")
		}
		code, err := strconv.Atoi(codes[i])
		if err != nil {
			return err
		}
		transformAlg := NewBWTAlgorithm(code)
		temp.Write(transformAlg.Retransform(buf[:n]))
		i++
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
