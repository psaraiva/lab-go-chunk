package model

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

type File struct {
	Hash string
	Name string
}

func (f File) GenerateHashByOsFile(file *os.File) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes), nil
}
