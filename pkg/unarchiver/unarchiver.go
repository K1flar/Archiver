package unarchiver

import (
	"bufio"
	"os"
	"path"
	"strconv"
	"strings"
)

const BufSize = 4096

type Unarchiver struct{}

func New() *Unarchiver {
	return &Unarchiver{}
}

func (ua *Unarchiver) Unarchive(archiveName, outputDirName string) error {
	arch, err := os.OpenFile(archiveName, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	defer arch.Close()

	scanner := bufio.NewScanner(arch)
	scanner.Scan()
	header := scanner.Text()
	arch.Seek(int64(len(header)+1), 0)

	filesInfo := strings.Split(header[:len(header)-1], ";")
	for _, fi := range filesInfo {
		fileInfo := strings.Split(fi, ",")
		filePath := fileInfo[0]
		fileSize, err := strconv.Atoi(fileInfo[1])
		if err != nil {
			return err
		}
		fileData := make([]byte, fileSize)
		absoluteFilePath := path.Join(outputDirName, filePath)

		var isDir bool
		if len(strings.Split(filePath, ".")) == 1 && fileSize == 0 {
			isDir = true
			err = os.MkdirAll(absoluteFilePath, 0777)
			if err != nil {
				return err
			}
		}

		_, err = arch.Read(fileData)
		if err != nil {
			return err
		}

		if _, err = os.Stat(absoluteFilePath); os.IsNotExist(err) {
			err = os.MkdirAll(path.Dir(absoluteFilePath), 0777)
			if err != nil {
				return err
			}
		}

		if !isDir {
			err = os.WriteFile(absoluteFilePath, fileData, 0777)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
