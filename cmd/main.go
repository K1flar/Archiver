package main

import (
	"archiver/internal/args"
	"fmt"
	"io"
	"os"
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

		dirName := os.Args[2]
		outputFileName, exists := args.FindFlag(os.Args, "o")
		if !exists {
			outputFileName = DefaultOutputArchiveName
		}

		// archive()
		fmt.Println(dirName, outputFileName)
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