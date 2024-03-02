package archiver

import "archiver/internal/compressors/huffman"

type Compressor interface {
	Compress() error
}

type Decompressor interface {
	Decompress() error
}

func (a *Archiver) Compress(compressor Compressor) error {
	return compressor.Compress()
}

func (a *Archiver) Decompress(decompressor Decompressor) error {
	return decompressor.Decompress()
}

func (a *Archiver) NewHuffmanComressor() *huffman.HuffmanCompressor {
	return huffman.New(a.archive)
}
