package process

import (
	"archiver/internal/flags"
	"fmt"
	"os"
	"path"
)

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

func checkIO(processName string, args []string) (string, string, error) {
	if len(args) == 0 {
		return "", "", fmt.Errorf("%s: no input file", processName)
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	var input, output string
	input = args[0]
	if input[0] != '/' {
		input = path.Join(wd, input)
	}

	output, _ = flags.FindFlag(args, "o")
	if len(output) != 0 && output[0] != '/' {
		output = path.Join(wd, output)
	}

	return input, output, nil
}

func getSize(count int64) string {
	switch {
	case count >= GB:
		return fmt.Sprintf("%.2fGb", float64(count)/GB)
	case count >= MB:
		return fmt.Sprintf("%.2fMb", float64(count)/MB)
	case count >= KB:
		return fmt.Sprintf("%.2fKb", float64(count)/KB)
	}
	return fmt.Sprintf("%db", count)
}

func pow(n int, p uint) int {
	if p == 0 {
		return 1
	}
	return n * pow(n, p-1)
}
