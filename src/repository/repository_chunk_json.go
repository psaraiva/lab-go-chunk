package repository

import (
	"encoding/json"
	"fmt"
	"lab/src/internal/entity"
	"os"
)

type RepositoryChunkJson struct{}

func (rcj RepositoryChunkJson) Create(chunk entity.Chunk) (int64, error) {
	var id int64
	jsonChunkFile, err := os.Open(os.Getenv("COLLECTION_CHUNK_JSON"))
	if err != nil {
		return id, err
	}
	defer jsonChunkFile.Close()

	decoder := json.NewDecoder(jsonChunkFile)
	chunkList := []entity.Chunk{}
	err = decoder.Decode(&chunkList)
	if err != nil {
		return id, err
	}

	chunkList = append(chunkList, chunk)
	updatedJSON, err := json.MarshalIndent(chunkList, "", "  ")
	if err != nil {
		return id, err
	}

	err = os.WriteFile(os.Getenv("COLLECTION_CHUNK_JSON"), updatedJSON, 0644)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (rcj RepositoryChunkJson) GetChunkHashListByHashOriginalFile(hashOriginalFile string) ([]string, error) {
	jsonChunk, err := os.Open(os.Getenv("COLLECTION_CHUNK_JSON"))
	if err != nil {
		return nil, err
	}
	defer jsonChunk.Close()

	decoder := json.NewDecoder(jsonChunk)
	chunkList := []entity.Chunk{}
	err = decoder.Decode(&chunkList)
	if err != nil {
		return nil, err
	}

	for _, chunk := range chunkList {
		if chunk.HashOriginalFile == hashOriginalFile {
			return chunk.HashList, nil
		}
	}

	return nil, fmt.Errorf("record not found")
}

func (rcj RepositoryChunkJson) getCountChunkMap(chunkList []entity.Chunk) map[string]int64 {
	chunkCount := make(map[string]int64)
	for _, item := range chunkList {
		for _, value := range item.HashList {
			chunkCount[value]++
		}
	}
	return chunkCount
}

func (rcj RepositoryChunkJson) CountUsedChunkHash(hash string) (int64, error) {
	jsonChunkFile, err := os.Open(os.Getenv("COLLECTION_CHUNK_JSON"))
	if err != nil {
		return 0, err
	}
	defer jsonChunkFile.Close()

	decoder := json.NewDecoder(jsonChunkFile)
	chunkList := []entity.Chunk{}
	err = decoder.Decode(&chunkList)
	if err != nil {
		return 0, err
	}

	countList := rcj.getCountChunkMap(chunkList)
	return countList[hash], nil
}

func (rcj RepositoryChunkJson) RemoveAll() error {
	err := os.WriteFile(os.Getenv("COLLECTION_CHUNK_JSON"), []byte("[]"), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (rcj RepositoryChunkJson) RemoveByHashOriginalFile(hashOriginalFile string) ([]string, error) {
	var hashList []string
	jsonChunkFile, err := os.Open(os.Getenv("COLLECTION_CHUNK_JSON"))
	if err != nil {
		return hashList, err
	}
	defer jsonChunkFile.Close()

	decoder := json.NewDecoder(jsonChunkFile)
	chunkList := []entity.Chunk{}
	err = decoder.Decode(&chunkList)
	if err != nil {
		return hashList, err
	}

	chunkTarget := entity.Chunk{}
	for index, item := range chunkList {
		if item.HashOriginalFile == hashOriginalFile {
			chunkTarget = chunkList[index]
			chunkList = append(chunkList[:index], chunkList[index+1:]...)
			break
		}
	}

	upJSON, err := json.MarshalIndent(chunkList, "", "  ")
	if err != nil {
		return hashList, err
	}

	err = os.WriteFile(os.Getenv("COLLECTION_CHUNK_JSON"), upJSON, 0644)
	if err != nil {
		return hashList, err
	}

	for _, hash := range chunkTarget.HashList {
		reps, err := rcj.isChunkHashCanBeRemoved(hash)
		if err != nil {
			return hashList, err
		}

		if !reps {
			continue
		}

		hashList = append(hashList, hash)
	}

	return hashList, nil
}

func (rcj RepositoryChunkJson) isChunkHashCanBeRemoved(hash string) (bool, error) {
	total, err := rcj.CountUsedChunkHash(hash)
	if err != nil {
		return false, err
	}
	return total == 0, nil
}
