package service

import (
	"lab/src/logger"
)

func (ac *Action) FeatureRemove() error {
	logger.GetLogActivity().WriteLog("Carregando Hash do arquivo...")
	hashOriginalFile, err := repositoryFile.GetHashByFileName(ac.FileTarget)
	if err != nil {
		return err
	}

	logger.GetLogActivity().WriteLog("Carregando lista de chunk do arquivo...")
	list, err := repositoryChunk.GetChunkHashListByHashOriginalFile(hashOriginalFile)
	if err != nil {
		return err
	}

	for _, hashChunk := range list {
		flag, err := repositoryChunk.IsChunkCanBeRemoved(hashChunk)
		if err != nil {
			return err
		}

		if !flag {
			continue
		}

		logger.GetLogActivity().WriteLog("Removendo chunk: " + hashChunk + ".bin")
		err = serviceStorage.RemoveFile(hashChunk + ".bin")
		if err != nil {
			return err
		}
	}

	logger.GetLogActivity().WriteLog("Removendo registro do arquivo de chunk...")
	err = repositoryFile.RemoveByHashFile(hashOriginalFile)
	if err != nil {
		return err
	}

	logger.GetLogActivity().WriteLog("Removendo registro do arquivo de hash...")
	return repositoryChunk.RemoveByHashOriginalFile(hashOriginalFile)
}
