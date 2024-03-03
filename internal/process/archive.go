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

const DefaultOutputArchiveName = "archive.arc"

type Transformer interface {
	Transform() error
}

type Compressor interface {
	Compress() error
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

func Archive(args []string) error {
	inputDir, outputArchive, err := checkIO("archive", args)
	if err != nil {
		return err
	}
	if len(outputArchive) == 0 {
		outputArchive = DefaultOutputArchiveName
	}

	s := time.Now()
	arc := archiver.New(outputArchive)
	err = arc.Archive(inputDir)
	if err != nil {
		return err
	}

	for _, arg := range args {
		if arg[0] != '-' {
			continue
		}
		flag := arg[1:]
		if _, ok := Transformers[flag]; ok {
			arg, _ := flags.FindFlag(args, flag)
			if err = arc.Transform(Transformers[flag](outputArchive, arg)); err != nil {
				return err
			}
		}
		if _, ok := Compressors[flag]; ok {
			if err = arc.Compress(Compressors[flag](outputArchive)); err != nil {
				return nil
			}
		}
	}
	fmt.Printf("Archive done: %.2f sec\n", time.Since(s).Seconds())

	return nil
}
