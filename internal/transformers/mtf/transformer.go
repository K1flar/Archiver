package mtf

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

type MTFTransformer struct {
	archive string
}

func New(archiveName string) *MTFTransformer {
	return &MTFTransformer{
		archive: archiveName,
	}
}

func (t *MTFTransformer) Transform() error {
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
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("%w", err)
		}
		transformAlg := NewMTFAlgorithm([]byte{})
		temp.Write(transformAlg.Transform(buf[:n]))
		header.WriteString(fmt.Sprintf("%s;", transformAlg.GetAlphabet()))
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

func (t *MTFTransformer) Retransform() error {
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

	alphabets := strings.Split(header[:len(header)-1], ";")

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
		if i >= len(alphabets) {
			return fmt.Errorf("mtf retransform: bad header")
		}
		alphabetCodes := strings.Split(alphabets[i], ",")
		alphabet := make([]byte, len(alphabetCodes))
		for i, c := range alphabetCodes {
			code, err := strconv.Atoi(c)
			if err != nil {
				return err
			}
			alphabet[i] = byte(code)
		}
		transformAlg := NewMTFAlgorithm(alphabet)
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
