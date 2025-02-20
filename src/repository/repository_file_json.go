package repository

import (
	"encoding/json"
	"fmt"
	"lab/src/internal/entity"
	"os"
)

type RepositoryFileJson struct{}

func (rfj RepositoryFileJson) Create(file entity.File) (int64, error) {
	var id int64
	fileJson, err := os.Open(os.Getenv("COLLECTION_FILE_JSON"))
	if err != nil {
		return id, err
	}
	defer fileJson.Close()

	decoder := json.NewDecoder(fileJson)
	fileList := []entity.File{}
	err = decoder.Decode(&fileList)
	if err != nil {
		return id, err
	}

	for _, item := range fileList {
		if item.Hash == file.Hash {
			return id, fmt.Errorf("arquivo já existe: %s", item.Name)
		}
	}

	fileList = append(fileList, entity.File{Hash: file.Hash, Name: file.Name})
	updatedJSON, err := json.MarshalIndent(fileList, "", "  ")
	if err != nil {
		return id, err
	}

	err = os.WriteFile(os.Getenv("COLLECTION_FILE_JSON"), updatedJSON, 0644)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (rfj RepositoryFileJson) GetHashByName(name string) (string, error) {
	jsonFile, err := os.Open(os.Getenv("COLLECTION_FILE_JSON"))
	if err != nil {
		return "", err
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	hashList := []entity.File{}
	err = decoder.Decode(&hashList)
	if err != nil {
		return "", err
	}

	for _, item := range hashList {
		if item.Name == name {
			return item.Hash, nil
		}
	}

	return "", fmt.Errorf("record not found")
}

func (rfj RepositoryFileJson) IsExistsByHash(hash string) (bool, error) {
	jsonHash, err := os.Open(os.Getenv("COLLECTION_FILE_JSON"))
	if err != nil {
		return false, err
	}
	defer jsonHash.Close()

	decoder := json.NewDecoder(jsonHash)
	fileList := []entity.File{}
	err = decoder.Decode(&fileList)
	if err != nil {
		return false, err
	}

	for _, item := range fileList {
		if item.Hash == hash {
			return true, nil
		}
	}

	return false, nil
}

func (rfj RepositoryFileJson) RemoveByHash(hash string) error {
	jsonHashFile, err := os.Open(os.Getenv("COLLECTION_FILE_JSON"))
	if err != nil {
		return err
	}
	defer jsonHashFile.Close()

	decoder := json.NewDecoder(jsonHashFile)
	fileList := []entity.File{}
	err = decoder.Decode(&fileList)
	if err != nil {
		return err
	}

	isDeleted := false
	for index, item := range fileList {
		if item.Hash == hash {
			isDeleted = true
			fileList = append(fileList[:index], fileList[index+1:]...)
			break
		}
	}

	if !isDeleted {
		return fmt.Errorf("record not found")
	}

	upJSON, err := json.MarshalIndent(fileList, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(os.Getenv("COLLECTION_FILE_JSON"), upJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (rfj RepositoryFileJson) RemoveAll() error {
	return os.WriteFile(os.Getenv("COLLECTION_FILE_JSON"), []byte("[]"), 0644)
}
