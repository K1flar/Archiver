package main

import (
	"archiver/internal/process"
	"fmt"
	"io"
	"os"
)

var Processes = map[string]func([]string) error{
	"archive":   process.Archive,
	"unarchive": process.Unarchive,
}

func main() {
	if len(os.Args) == 1 {
		io.WriteString(os.Stderr, "no command line arguments\n")
		return
	}

	if process, ok := Processes[os.Args[1]]; ok {
		if err := process(os.Args[2:]); err != nil {
			io.WriteString(os.Stderr, fmt.Sprintf("%s\n", err.Error()))
		}
	} else {
		io.WriteString(os.Stderr, "unknown command\n")
	}

}
