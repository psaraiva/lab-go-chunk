package actions

import (
	"encoding/json"
	"lab/src/logger"
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
	chunkList := []chunkItem{}
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

func (ac *Action) getCountChunkMap(chunkList []chunkItem) map[string]int {
	chunkCount := make(map[string]int)
	for _, item := range chunkList {
		for _, value := range item.Chunk {
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
	chunkList := []chunkItem{}
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
	hashList := []hashItem{}
	err = decoder.Decode(&hashList)
	if err != nil {
		return err
	}

	for index, item := range hashList {
		if item.Hash == hashString {
			hashList = append(hashList[:index], hashList[index+1:]...)
			break
		}
	}

	upJSON, err := json.MarshalIndent(hashList, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(os.Getenv("JSON_FILE_HASH"), upJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}
