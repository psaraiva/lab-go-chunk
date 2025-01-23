package model

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type ChunkItem struct {
	HashFile string
	HashList []string
	Size     int
}

func (ci ChunkItem) GenerateChunkByOsFile(hashFile string, chunkSize int, file *os.File) (ChunkItem, error) {
	defer file.Close()
	chunkItem := ChunkItem{}

	if chunkSize < 1024 { // 1kb
		return chunkItem, fmt.Errorf("falha na configuração de: Chunk Size")
	}

	var chunks []string
	buf := make([]byte, chunkSize)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return chunkItem, err
		}

		if n == 0 {
			break
		}

		chunks = append(chunks, ci.GenerateHash(buf, n))
	}

	chunkItem.HashFile = hashFile
	chunkItem.HashList = chunks
	chunkItem.Size = chunkSize
	return chunkItem, nil
}

func (ci ChunkItem) GenerateHash(buf []byte, ref int) string {
	return ci.generateHashMd5(buf, ref)
}

func (ci ChunkItem) generateHashMd5(buf []byte, ref int) string {
	chunkHash := md5.Sum(buf[:ref])
	return hex.EncodeToString(chunkHash[:])
}
