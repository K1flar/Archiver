package main

import (
	"archiver/internal/args"
	"archiver/internal/compressors/huffman"
	"archiver/internal/transformers/bwt"
	"archiver/internal/transformers/caesar"
	"archiver/internal/transformers/mirror"
	"archiver/internal/transformers/mtf"
	"archiver/pkg/archiver"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	DefaultOutputArchiveName = "archive.arc"
	DefaultOutputDirName     = "."
)

type Transformer interface {
	Transform() error
	Retransform() error
}

type Compressor interface {
	Compress() error
	Decompress() error
}

var Transformers = map[string]func(string, ...string) Transformer{
	"bwt":    func(s1 string, s2 ...string) Transformer { return bwt.New(s1) },
	"mtf":    func(s1 string, s2 ...string) Transformer { return mtf.New(s1) },
	"mirror": func(s1 string, s2 ...string) Transformer { return mirror.New(s1) },
	"caesar": func(s1 string, s2 ...string) Transformer {
		shift, _ := strconv.Atoi(s2[0])
		return caesar.New(s1, shift)
	},
}

var Compressors = map[string]func(string) Compressor{
	"huff": func(s string) Compressor { return huffman.New(s) },
}

func main() {
	s := time.Now()
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

		arc := archiver.New(outputFileName)
		err = arc.Archive(dirName)
		if err != nil {
			io.WriteString(os.Stderr, fmt.Sprintf("%s\n", err.Error()))
			return
		}
		for _, arg := range os.Args {
			if arg[0] != '-' {
				continue
			}
			flag := arg[1:]
			if _, ok := Transformers[flag]; ok {
				arg, _ := args.FindFlag(os.Args, flag)
				err = arc.Transform(Transformers[flag](outputFileName, arg))
				if err != nil {
					io.WriteString(os.Stderr, fmt.Sprintf("%s\n", err.Error()))
					return
				}
			}
			if _, ok := Compressors[flag]; ok {
				arc.Compress(Compressors[flag](outputFileName))
			}
		}
		e := time.Now()
		fmt.Printf("Done: %.2f sec", e.Sub(s).Seconds())

	case "unarchive":
		if len(os.Args) == 2 {
			io.WriteString(os.Stderr, "no archive\n")
			return
		}

		wd, err := os.Getwd()
		if err != nil {
			io.WriteString(os.Stderr, fmt.Sprintf("%s\n", err.Error()))
			return
		}

		archiveName := os.Args[2]
		if archiveName[0] != '/' {
			archiveName = path.Join(wd, archiveName)
		}

		outputDirName, exists := args.FindFlag(os.Args, "o")
		if !exists {
			outputDirName = DefaultOutputDirName
		}
		if outputDirName[0] != '/' {
			outputDirName = path.Join(wd, outputDirName)
		}

		arc := archiver.New(archiveName)
		for _, arg := range os.Args {
			if arg[0] != '-' {
				continue
			}
			flag := arg[1:]
			if _, ok := Transformers[flag]; ok {
				arg, _ := args.FindFlag(os.Args, flag)
				err = arc.Retransform(Transformers[flag](archiveName, arg))
				if err != nil {
					io.WriteString(os.Stderr, fmt.Sprintf("%s\n", err.Error()))
					return
				}
			}
			if _, ok := Compressors[flag]; ok {
				arc.Decompress(Compressors[flag](archiveName))
			}
		}
		err = arc.Unarchive(outputDirName)
		if err != nil {
			io.WriteString(os.Stderr, fmt.Sprintf("%s\n", err.Error()))
			return
		}
		e := time.Now()
		fmt.Printf("Done: %.2f sec", e.Sub(s).Seconds())
	default:
		io.WriteString(os.Stderr, "unknown command\n")
	}
}
