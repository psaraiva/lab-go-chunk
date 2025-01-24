package repository

import (
	"encoding/json"
	"fmt"
	"lab/src/model"
	"os"
)

type RepositoryChunkJson struct{}

func (rcij RepositoryChunkJson) Create(chunk model.Chunk) error {
	jsonChunkFile, err := os.Open(os.Getenv("COLLECTION_JSON_CHUNK"))
	if err != nil {
		return err
	}
	defer jsonChunkFile.Close()

	decoder := json.NewDecoder(jsonChunkFile)
	chunkList := []model.Chunk{}
	err = decoder.Decode(&chunkList)
	if err != nil {
		return err
	}

	chunkList = append(chunkList, chunk)
	updatedJSON, err := json.MarshalIndent(chunkList, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(os.Getenv("COLLECTION_JSON_CHUNK"), updatedJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (rcij RepositoryChunkJson) GetChunkHashListByHashOriginalFile(hashOriginalFile string) ([]string, error) {
	jsonChunk, err := os.Open(os.Getenv("COLLECTION_JSON_CHUNK"))
	if err != nil {
		return nil, err
	}
	defer jsonChunk.Close()

	decoder := json.NewDecoder(jsonChunk)
	chunkList := []model.Chunk{}
	err = decoder.Decode(&chunkList)
	if err != nil {
		return nil, err
	}

	for _, chunk := range chunkList {
		if chunk.HashOriginalFile == hashOriginalFile {
			return chunk.HashList, nil
		}
	}

	return nil, fmt.Errorf("arquivo não encontrado")
}

// @WARNING, verificar o uso desse método nas demais implementações de repositoryChunk...
func (rcij RepositoryChunkJson) GetCountChunkMap(chunkList []model.Chunk) map[string]int {
	chunkCount := make(map[string]int)
	for _, item := range chunkList {
		for _, value := range item.HashList {
			chunkCount[value]++
		}
	}
	return chunkCount
}

func (rcij RepositoryChunkJson) IsChunkCanBeRemoved(chunk string) (bool, error) {
	jsonChunkFile, err := os.Open(os.Getenv("COLLECTION_JSON_CHUNK"))
	if err != nil {
		return false, err
	}
	defer jsonChunkFile.Close()

	decoder := json.NewDecoder(jsonChunkFile)
	chunkList := []model.Chunk{}
	err = decoder.Decode(&chunkList)
	if err != nil {
		return false, err
	}

	countList := rcij.GetCountChunkMap(chunkList)
	if countList[chunk] > 1 {
		return false, nil
	}

	return true, nil
}

func (rcij RepositoryChunkJson) RemoveAll() error {
	err := os.WriteFile(os.Getenv("COLLECTION_JSON_CHUNK"), []byte("[]"), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (rcij RepositoryChunkJson) RemoveByHashOriginalFile(hashOriginalFile string) error {
	jsonChunkFile, err := os.Open(os.Getenv("COLLECTION_JSON_CHUNK"))
	if err != nil {
		return err
	}
	defer jsonChunkFile.Close()

	decoder := json.NewDecoder(jsonChunkFile)
	chunkList := []model.Chunk{}
	err = decoder.Decode(&chunkList)
	if err != nil {
		return err
	}

	for index, item := range chunkList {
		if item.HashOriginalFile == hashOriginalFile {
			chunkList = append(chunkList[:index], chunkList[index+1:]...)
			break
		}
	}

	upJSON, err := json.MarshalIndent(chunkList, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(os.Getenv("COLLECTION_JSON_CHUNK"), upJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}
