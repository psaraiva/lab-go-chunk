package entity

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type Chunk struct {
	HashOriginalFile string
	HashList         []string
	Size             int
}

func (c Chunk) GenerateChunkByOsFile(hashFile string, chunkSize int, file *os.File) (Chunk, error) {
	defer file.Close()
	chunk := Chunk{}

	if chunkSize < 1024 { // 1kb
		return chunk, fmt.Errorf("falha na configuração de: Chunk Size")
	}

	var chunks []string
	buf := make([]byte, chunkSize)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return chunk, err
		}

		if n == 0 {
			break
		}

		chunks = append(chunks, c.GenerateHash(buf, n))
	}

	chunk.HashOriginalFile = hashFile
	chunk.HashList = chunks
	chunk.Size = chunkSize
	return chunk, nil
}

func (c Chunk) GenerateHash(buf []byte, ref int) string {
	return c.generateHashMd5(buf, ref)
}

func (c Chunk) generateHashMd5(buf []byte, ref int) string {
	chunkHash := md5.Sum(buf[:ref])
	return hex.EncodeToString(chunkHash[:])
}
