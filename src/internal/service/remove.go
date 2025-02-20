package service

import (
	"fmt"
	"lab/src/logger"
	"lab/src/repository"
)

func (ac *Action) FeatureRemove() error {
	logger.GetLogActivity().WriteLog("Obtendo Hash do arquivo...")
	hashOriginalFile, err := repositoryFile.GetHashByName(ac.FileTarget)
	if err != nil {
		ac.LogErrorWrite(
			ACTION_REMOVE,
			fmt.Sprintf(
				" file not found by name -  %s",
				ac.FileTarget))

		if err == repository.ErrorRecordNotFound {
			return fmt.Errorf("arquivo não encontrado: %s", ac.FileTarget)
		}

		return errorActionDefault
	}

	if hashOriginalFile == "" {
		ac.LogErrorWrite(
			ACTION_REMOVE,
			" hash original file not found")

		return fmt.Errorf("arquivo não encontrado: %s", ac.FileTarget)
	}

	logger.GetLogActivity().WriteLog("Removendo registro(s) da coleção chunk...")
	hashList, err := repositoryChunk.RemoveByHashOriginalFile(hashOriginalFile)
	if err != nil {
		ac.LogErrorWrite(
			ACTION_REMOVE,
			fmt.Sprintf(" fail remove by original file - %s", err.Error()))

		if err == repository.ErrorRecordNotFound {
			return fmt.Errorf("arquivo não encontrado: %s", ac.FileTarget)
		}

		return errorActionDefault
	}

	logger.GetLogActivity().WriteLog("Removendo registro da coleção file...")
	err = repositoryFile.RemoveByHash(hashOriginalFile)
	if err != nil {
		ac.LogErrorWrite(
			ACTION_REMOVE,
			fmt.Sprintf(" fail remove by hash - %s", err.Error()))

		if err == repository.ErrorRecordNotFound {
			return fmt.Errorf("arquivo não encontrado: %s", ac.FileTarget)
		}

		return errorActionDefault
	}

	for _, hash := range hashList {
		logger.GetLogActivity().WriteLog("Removendo chunk: " + hash + ".bin")
		err = serviceStorage.RemoveFile(hash + ".bin")
		if err != nil {
			ac.LogErrorWrite(
				ACTION_REMOVE,
				fmt.Sprintf(" fail remove file to storage - %s", err.Error()))
			return errorActionDefault
		}
	}

	return err
}
