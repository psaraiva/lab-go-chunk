package repository

import (
	"encoding/json"
	"fmt"
	"lab/src/model"
	"os"
)

type RepositoryChunkItemJson struct{}

func (rcij RepositoryChunkItemJson) Create(chunkItem model.ChunkItem) error {
	jsonChunkFile, err := os.Open(os.Getenv("JSON_FILE_CHUNK"))
	if err != nil {
		return err
	}
	defer jsonChunkFile.Close()

	decoder := json.NewDecoder(jsonChunkFile)
	chunkList := []model.ChunkItem{}
	err = decoder.Decode(&chunkList)
	if err != nil {
		return err
	}

	chunkList = append(chunkList, chunkItem)
	updatedJSON, err := json.MarshalIndent(chunkList, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(os.Getenv("JSON_FILE_CHUNK"), updatedJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (rcij RepositoryChunkItemJson) GetChunkHashListByHashFile(hashFile string) ([]string, error) {
	jsonChunk, err := os.Open(os.Getenv("JSON_FILE_CHUNK"))
	if err != nil {
		return nil, err
	}
	defer jsonChunk.Close()

	decoder := json.NewDecoder(jsonChunk)
	chunkList := []model.ChunkItem{}
	err = decoder.Decode(&chunkList)
	if err != nil {
		return nil, err
	}

	for _, chunk := range chunkList {
		if chunk.HashFile == hashFile {
			return chunk.HashList, nil
		}
	}

	return nil, fmt.Errorf("arquivo não encontrado")
}

// @WARNING, verificar o uso desse método nas demais implementações de repositoryChunkItem...
func (rcij RepositoryChunkItemJson) GetCountChunkMap(chunkList []model.ChunkItem) map[string]int {
	chunkCount := make(map[string]int)
	for _, item := range chunkList {
		for _, value := range item.HashList {
			chunkCount[value]++
		}
	}
	return chunkCount
}

func (rcij RepositoryChunkItemJson) IsChunkCanBeRemoved(chunk string) (bool, error) {
	jsonChunkFile, err := os.Open(os.Getenv("JSON_FILE_CHUNK"))
	if err != nil {
		return false, err
	}
	defer jsonChunkFile.Close()

	decoder := json.NewDecoder(jsonChunkFile)
	chunkList := []model.ChunkItem{}
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

func (rcij RepositoryChunkItemJson) RemoveAll() error {
	err := os.WriteFile(os.Getenv("JSON_FILE_CHUNK"), []byte("[]"), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (rcij RepositoryChunkItemJson) RemoveByHashFile(hashFile string) error {
	jsonChunkFile, err := os.Open(os.Getenv("JSON_FILE_CHUNK"))
	if err != nil {
		return err
	}
	defer jsonChunkFile.Close()

	decoder := json.NewDecoder(jsonChunkFile)
	chunkList := []model.ChunkItem{}
	err = decoder.Decode(&chunkList)
	if err != nil {
		return err
	}

	for index, item := range chunkList {
		if item.HashFile == hashFile {
			chunkList = append(chunkList[:index], chunkList[index+1:]...)
			break
		}
	}

	upJSON, err := json.MarshalIndent(chunkList, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(os.Getenv("JSON_FILE_CHUNK"), upJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}
