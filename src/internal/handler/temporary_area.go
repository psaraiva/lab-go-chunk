package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type TemporaryArea struct{}

func MakeServiceTemporaryArea() TemporaryArea {
	return TemporaryArea{}
}

func (ta TemporaryArea) Clear() error {
	err := filepath.Walk(os.Getenv("SERVICE_FILE_TMP"), func(path string, info os.FileInfo, err error) error {
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

func (ta TemporaryArea) GetFile(name string) (*os.File, error) {
	return os.Open(os.Getenv("SERVICE_FILE_TMP") + "/" + name)
}

func (ta TemporaryArea) CreateFileByFileSource(name string, source *os.File) error {
	destinationFile, err := os.Create(os.Getenv("SERVICE_FILE_TMP") + "/" + name)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, source)
	if err != nil {
		return err
	}

	return nil
}

func (ta TemporaryArea) RemoveFile(name string) error {
	return os.Remove(fmt.Sprintf("%s/%s", os.Getenv("SERVICE_FILE_TMP"), name))
}
