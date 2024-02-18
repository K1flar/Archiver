package archiver

import (
	"fmt"
	"os"
	"path"
	"strings"
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

	var headers strings.Builder
	err = a.archvieRecursive(dir, out, &headers)
	if err != nil {
		return err
	}

	// TODO: подумать насчет headers
	data, err := os.ReadFile(outputFileName)
	if err != nil {
		return err
	}
	_, err = out.WriteAt([]byte(headers.String()), 0)
	if err != nil {
		return err
	}
	_, err = out.WriteAt(data, int64(len([]byte(headers.String()))))
	if err != nil {
		return err
	}

	return nil
}

func (a *Archvier) archvieRecursive(dir, out *os.File, headers *strings.Builder) error {
	files, err := dir.ReadDir(0)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := path.Join(dir.Name(), file.Name())
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		if file.IsDir() {
			err = a.archvieRecursive(f, out, headers)
			if err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			headers.WriteString(fmt.Sprintf("%s,%d;", filePath, len(data)))
			_, err = out.Write(data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
