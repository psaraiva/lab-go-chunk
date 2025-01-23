package actions

import (
	"lab/src/logger"
	"os"
)

func (ac *Action) FeatureRemove() error {
	logger.GetLogActivity().WriteLog("Carregando Hash do arquivo...")
	hashFile, err := repositoryFile.GetHashByFileName(ac.FileTarget)
	if err != nil {
		return err
	}

	logger.GetLogActivity().WriteLog("Carregando lista de chunk do arquivo...")
	list, err := repositoryChunkItem.GetChunkHashListByHashFile(hashFile)
	if err != nil {
		return err
	}

	for _, hashChunk := range list {
		flag, err := repositoryChunkItem.IsChunkCanBeRemoved(hashChunk)
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
	err = repositoryFile.RemoveByHashFile(hashFile)
	if err != nil {
		return err
	}

	logger.GetLogActivity().WriteLog("Removendo registro do arquivo de hash...")
	return repositoryChunkItem.RemoveByHashFile(hashFile)
}
