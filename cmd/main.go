package main

import (
	"archiver/internal/args"
	"archiver/pkg/archiver"
	"fmt"
	"io"
	"os"
	"path"
)

const DefaultOutputArchiveName = "archive.txt"

func main() {
	if len(os.Args) == 1 {
		io.WriteString(os.Stderr, "no command line arguments\n")
		return
	}

	switch os.Args[1] {
	case "archive":
		if len(os.Args) == 2 {
			io.WriteString(os.Stderr, "no input directory\n")
			return
		}

		wd, err := os.Getwd()
		if err != nil {
			io.WriteString(os.Stderr, fmt.Sprintf("%s\n", err.Error()))
			return
		}

		dirName := os.Args[2]
		if dirName[0] != '/' {
			dirName = path.Join(wd, dirName)
		}
		outputFileName, exists := args.FindFlag(os.Args, "o")
		if !exists {
			outputFileName = DefaultOutputArchiveName
		}
		if outputFileName[0] != '/' {
			outputFileName = path.Join(wd, outputFileName)
		}

		a := archiver.New()
		err = a.Archvie(dirName, outputFileName)
		if err != nil {
			io.WriteString(os.Stderr, fmt.Sprintf("%s\n", err.Error()))
			return
		}
	case "unarchive":
		if len(os.Args) == 2 {
			io.WriteString(os.Stderr, "no archvie\n")
			return
		}

		archiveName := os.Args[2]

		// unarchive()
		fmt.Println(archiveName)

	default:
		io.WriteString(os.Stderr, "unknown command\n")
	}
}
