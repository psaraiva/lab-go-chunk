package service

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

type Storage struct{}

func MakeServiceStorage() Storage {
	return Storage{}
}

func (s Storage) CreateFile(data []byte, name string) error {
	file, err := os.Create(fmt.Sprintf("%s/%s", os.Getenv("SERVICE_FILE_STORAGE"), name))
	if err != nil {
		return err
	}
	defer file.Close()

	for _, value := range data {
		err = binary.Write(file, binary.LittleEndian, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Storage) Clear() error {
	err := filepath.Walk(os.Getenv("SERVICE_FILE_STORAGE"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return os.Remove(path)
	})
	return err
}

func (s Storage) GetFile(name string) (*os.File, error) {
	return os.Open(os.Getenv("SERVICE_FILE_STORAGE") + "/" + name)
}

func (s Storage) RemoveFile(name string) error {
	return os.Remove(os.Getenv("SERVICE_FILE_STORAGE") + "/" + name)
}
