package service

import (
	"fmt"
	"lab/src/logger"
)

func (ac *Action) FeatureRemove() error {
	logger.GetLogActivity().WriteLog("Obtendo Hash do arquivo...")
	hashOriginalFile, err := repositoryFile.GetHashByName(ac.FileTarget)
	if err != nil {
		return err
	}

	if hashOriginalFile == "" {
		logger.GetLogActivity().WriteLog(fmt.Sprintf("Arquivo não encontrado: %s.", ac.FileTarget))
		return fmt.Errorf("arquivo não encontrado: %s", ac.FileTarget)
	}

	logger.GetLogActivity().WriteLog("Removendo registro(s) da coleção chunk...")
	hashList, err := repositoryChunk.RemoveByHashOriginalFile(hashOriginalFile)
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
		return err
	}

	logger.GetLogActivity().WriteLog("Removendo registro da coleção file...")
	err = repositoryFile.RemoveByHashFile(hashOriginalFile)
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
		return err
	}

	for _, hash := range hashList {
		logger.GetLogActivity().WriteLog("Removendo chunk: " + hash + ".bin")
		err = serviceStorage.RemoveFile(hash + ".bin")
		if err != nil {
			logger.GetLogError().WriteLog(err.Error())
			return err
		}
	}

	return err
}
