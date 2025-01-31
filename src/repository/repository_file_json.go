package repository

import (
	"encoding/json"
	"fmt"
	"lab/src/model"
	"os"
)

type RepositoryFileJson struct{}

func (rfj RepositoryFileJson) Create(file model.File) (int64, error) {
	var id int64
	jsonHash, err := os.Open(os.Getenv("COLLECTION_JSON_FILE"))
	if err != nil {
		return id, err
	}
	defer jsonHash.Close()

	decoder := json.NewDecoder(jsonHash)
	fileList := []model.File{}
	err = decoder.Decode(&fileList)
	if err != nil {
		return id, err
	}

	for _, item := range fileList {
		if item.Hash == file.Hash {
			return id, fmt.Errorf("arquivo já existe: %s", item.Name)
		}
	}

	fileList = append(fileList, model.File{Hash: file.Hash, Name: file.Name})
	updatedJSON, err := json.MarshalIndent(fileList, "", "  ")
	if err != nil {
		return id, err
	}

	err = os.WriteFile(os.Getenv("COLLECTION_JSON_FILE"), updatedJSON, 0644)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (rfj RepositoryFileJson) GetHashByName(name string) (string, error) {
	jsonFile, err := os.Open(os.Getenv("COLLECTION_JSON_FILE"))
	if err != nil {
		return "", err
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	hashList := []model.File{}
	err = decoder.Decode(&hashList)
	if err != nil {
		return "", err
	}

	for _, item := range hashList {
		if item.Name == name {
			return item.Hash, nil
		}
	}

	return "", fmt.Errorf("arquivo não encontrado")
}

func (rfj RepositoryFileJson) IsExistsByHash(hash string) (bool, error) {
	jsonHash, err := os.Open(os.Getenv("COLLECTION_JSON_FILE"))
	if err != nil {
		return false, err
	}
	defer jsonHash.Close()

	decoder := json.NewDecoder(jsonHash)
	fileList := []model.File{}
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
	jsonHashFile, err := os.Open(os.Getenv("COLLECTION_JSON_FILE"))
	if err != nil {
		return err
	}
	defer jsonHashFile.Close()

	decoder := json.NewDecoder(jsonHashFile)
	fileList := []model.File{}
	err = decoder.Decode(&fileList)
	if err != nil {
		return err
	}

	for index, item := range fileList {
		if item.Hash == hash {
			fileList = append(fileList[:index], fileList[index+1:]...)
			break
		}
	}

	upJSON, err := json.MarshalIndent(fileList, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(os.Getenv("COLLECTION_JSON_FILE"), upJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (rfj RepositoryFileJson) RemoveAll() error {
	err := os.WriteFile(os.Getenv("COLLECTION_JSON_FILE"), []byte("[]"), 0644)
	if err != nil {
		return err
	}
	return nil
}
