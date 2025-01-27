package service

import (
	"lab/src/interfaces"
	"lab/src/logger"
	"lab/src/repository"
	"os"
)

const (
	ACTION_CLEAR    = "clear"
	ACTION_DOWNLOAD = "download"
	ACTION_REMOVE   = "remove"
	ACTION_UPLOAD   = "upload"
)

type Action struct {
	Type       string
	FileTarget string
	Hash       string
}

var repositoryFile repository.RepositoryFile
var repositoryChunk repository.RepositoryChunk
var serviceTemporaryArea interfaces.ServiceTemporaryArea
var serviceStorage interfaces.ServiceStorage

func MakeAction() Action {
	repositoryFile = repository.MakeRepositoryFile(os.Getenv("ENGINE_COLLECTION"))
	repositoryChunk = repository.MakeRepositoryChunk(os.Getenv("ENGINE_COLLECTION"))
	serviceTemporaryArea = MakeServiceTemporaryArea()
	serviceStorage = MakeServiceStorage()

	return Action{
		Type:       "none",
		FileTarget: "none",
		Hash:       "none",
	}
}

func Execute(action interfaces.ServiceAction) error {
	switch action.GetActionType() {
	case ACTION_CLEAR:
		logger.GetLogActivity().WriteLog("Iniciando rotina de limpeza.")
		err := action.FeatureClear()
		if err != nil {
			logger.GetLogError().WriteLog(err.Error())
			return err
		}

		logger.GetLogActivity().WriteLog("Rotina de limpeza realizada com sucesso!")
	case ACTION_UPLOAD:
		logger.GetLogActivity().WriteLog("Iniciando rotina de carregamento de arquivo.")
		err := action.FeatureUpload()
		if err != nil {
			logger.GetLogError().WriteLog(err.Error())
			return err
		}

		logger.GetLogActivity().WriteLog("Arquivo carregado com sucesso!")
	case ACTION_DOWNLOAD:
		logger.GetLogActivity().WriteLog("Iniciando rotina de descarregamento de arquivo.")
		err := action.FeatureDownload()
		if err != nil {
			logger.GetLogError().WriteLog(err.Error())
			return err
		}

		logger.GetLogActivity().WriteLog("Arquivo descarregado com sucesso!")
	case ACTION_REMOVE:
		logger.GetLogActivity().WriteLog("Iniciando rotina de remoção de arquivo.")
		err := action.FeatureRemove()
		if err != nil {
			logger.GetLogError().WriteLog(err.Error())
			return err
		}

		logger.GetLogActivity().WriteLog("Arquivo removido com sucesso!")
	}

	return nil
}

func (ac Action) GetActionType() string {
	return ac.Type
}
