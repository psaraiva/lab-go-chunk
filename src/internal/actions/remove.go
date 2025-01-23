package actions

import (
	"encoding/json"
	"lab/src/logger"
	"lab/src/models"
	"os"
)

func (ac *Action) FeatureRemove() error {
	logger.GetLogActivity().WriteLog("Carregando Hash do arquivo...")
	hash, err := ac.getHashByFileName(ac.FileTarget)
	if err != nil {
		return err
	}

	logger.GetLogActivity().WriteLog("Carregando lista de chunk do arquivo...")
	list, err := ac.getChunksByHash(hash)
	if err != nil {
		return err
	}

	for _, hashChunk := range list {
		flag, err := ac.isChunkCanBeRemoved(hashChunk)
		if err != nil {
			return err
		}

		if !flag {
			continue
		}

		logger.GetLogActivity().WriteLog("Removendo chunk: " + hashChunk + ".bin")
		err = os.Remove(os.Getenv("FOLDER_STORAGE") + "/" + hashChunk + ".bin")
		if err != nil {
			return err
		}
	}

	logger.GetLogActivity().WriteLog("Removendo registro do arquivo de chunk...")
	err = ac.removeHashFileToChunkCollection(hash)
	if err != nil {
		return err
	}

	logger.GetLogActivity().WriteLog("Removendo registro do arquivo de hash...")
	return ac.removeHashToCollection(hash)
}

func (ac *Action) isChunkCanBeRemoved(chunk string) (bool, error) {
	jsonChunkFile, err := os.Open(os.Getenv("JSON_FILE_CHUNK"))
	if err != nil {
		return false, err
	}
	defer jsonChunkFile.Close()

	decoder := json.NewDecoder(jsonChunkFile)
	chunkList := []models.ChunkItem{}
	err = decoder.Decode(&chunkList)
	if err != nil {
		return false, err
	}

	countList := ac.getCountChunkMap(chunkList)
	if countList[chunk] > 1 {
		return false, nil
	}

	return true, nil
}

func (ac *Action) getCountChunkMap(chunkList []models.ChunkItem) map[string]int {
	chunkCount := make(map[string]int)
	for _, item := range chunkList {
		for _, value := range item.HashList {
			chunkCount[value]++
		}
	}
	return chunkCount
}

func (ac *Action) removeHashFileToChunkCollection(hashFile string) error {
	jsonChunkFile, err := os.Open(os.Getenv("JSON_FILE_CHUNK"))
	if err != nil {
		return err
	}
	defer jsonChunkFile.Close()

	decoder := json.NewDecoder(jsonChunkFile)
	chunkList := []models.ChunkItem{}
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

func (ac *Action) removeHashToCollection(hashString string) error {
	jsonHashFile, err := os.Open(os.Getenv("JSON_FILE_HASH"))
	if err != nil {
		return err
	}
	defer jsonHashFile.Close()

	decoder := json.NewDecoder(jsonHashFile)
	fileList := []models.File{}
	err = decoder.Decode(&fileList)
	if err != nil {
		return err
	}

	for index, item := range fileList {
		if item.Hash == hashString {
			fileList = append(fileList[:index], fileList[index+1:]...)
			break
		}
	}

	upJSON, err := json.MarshalIndent(fileList, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(os.Getenv("JSON_FILE_HASH"), upJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}
