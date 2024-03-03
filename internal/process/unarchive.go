package process

import (
	"archiver/internal/compressors/huffman"
	"archiver/internal/flags"
	"archiver/internal/transformers/bwt"
	"archiver/internal/transformers/caesar"
	"archiver/internal/transformers/mirror"
	"archiver/internal/transformers/mtf"
	"archiver/pkg/archiver"
	"fmt"
	"strconv"
	"time"
)

const DefaultOutputDirName = "."

type Retransformer interface {
	Retransform() error
}

type Decompressor interface {
	Decompress() error
}

var Retransformers = map[string]func(string, ...string) Retransformer{
	"bwt":    func(s1 string, s2 ...string) Retransformer { return bwt.New(s1) },
	"mtf":    func(s1 string, s2 ...string) Retransformer { return mtf.New(s1) },
	"mirror": func(s1 string, s2 ...string) Retransformer { return mirror.New(s1) },
	"caesar": func(s1 string, s2 ...string) Retransformer {
		shift, _ := strconv.Atoi(s2[0])
		return caesar.New(s1, shift)

	},
}

var Decompressors = map[string]func(string) Decompressor{
	"huff": func(s string) Decompressor { return huffman.New(s) },
}

func Unarchive(args []string) error {
	inputArchive, outputDir, err := checkIO("unarchive", args)
	if err != nil {
		return err
	}
	if len(outputDir) == 0 {
		outputDir = DefaultOutputDirName
	}

	s := time.Now()
	arc := archiver.New(inputArchive)

	for _, arg := range args {
		if arg[0] != '-' {
			continue
		}
		flag := arg[1:]
		if _, ok := Retransformers[flag]; ok {
			arg, _ := flags.FindFlag(args, flag)
			if err = arc.Retransform(Retransformers[flag](inputArchive, arg)); err != nil {
				return err
			}
		}
		if _, ok := Compressors[flag]; ok {
			if err = arc.Decompress(Decompressors[flag](inputArchive)); err != nil {
				return nil
			}
		}
	}

	err = arc.Unarchive(outputDir)
	if err != nil {
		return err
	}
	fmt.Printf("Archive done: %.2f sec\n", time.Since(s).Seconds())

	return nil
}
