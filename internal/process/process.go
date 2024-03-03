package process

import (
	"archiver/internal/flags"
	"fmt"
	"os"
	"path"
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
