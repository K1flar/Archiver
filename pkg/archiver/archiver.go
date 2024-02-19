package archiver

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

const (
	TempFile = "temp.txt"
	BufSize  = 4096
)

type Archvier struct{}

func New() *Archvier {
	return &Archvier{}
}

func (a *Archvier) Archvie(inputDirName, outputFileName string) error {
	dir, err := os.OpenFile(inputDirName, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	defer dir.Close()

	out, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE, 0622)
	if err != nil {
		return err
	}
	defer out.Close()

	temp, err := os.OpenFile(TempFile, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer func() {
		temp.Close()
		os.Remove(TempFile)
	}()

	var headers strings.Builder
	err = a.archvieRecursive(dir, temp, &headers, path.Base(inputDirName))
	if err != nil {
		return err
	}

	_, err = out.WriteString(headers.String() + "\n")
	if err != nil {
		return err
	}

	buf := make([]byte, BufSize)
	temp.Seek(0, 0)
	for {
		n, err := temp.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("%w", err)
		}
		_, err = out.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Archvier) archvieRecursive(dir, out *os.File, headers *strings.Builder, curDir string) error {
	files, err := dir.ReadDir(0)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		headers.WriteString(fmt.Sprintf("%s/,%d;", curDir, 0))
	}

	for _, file := range files {
		if file.Name() == out.Name() {
			continue
		}
		filePath := path.Join(dir.Name(), file.Name())
		fileName := path.Base(filePath)
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		if file.IsDir() {
			err = a.archvieRecursive(f, out, headers, path.Join(curDir, fileName))
			if err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			headers.WriteString(fmt.Sprintf("%s,%d;", path.Join(curDir, fileName), len(data)))
			_, err = out.Write(data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
