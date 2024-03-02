package caesar

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
	BufSize  = 4096
)

type CaesarTransformer struct {
	archive string
	shift   int
}

func New(archiveName string, shift int) *CaesarTransformer {
	return &CaesarTransformer{
		archive: archiveName,
		shift:   shift,
	}
}

func (t *CaesarTransformer) Transform() error {
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
	header.WriteString(fmt.Sprintf("%d;\n", t.shift))
	temp.WriteString(header.String())
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
		transformAlg := NewCaesarAlgorithm(t.shift)
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

func (t *CaesarTransformer) Retransform() error {
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

	shift, err := strconv.Atoi(strings.Split(header[:len(header)-1], ";")[0])
	if err != nil {
		return err
	}

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
		transformAlg := NewCaesarAlgorithm(shift)
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
