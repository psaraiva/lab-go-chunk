package service

import (
	"errors"
	"fmt"
	"lab/src/interfaces"
	"lab/src/logger"
	"lab/src/repository"
	"os"
)

const (
	ACTION_CLEAR         = "clear"
	ACTION_DOWNLOAD      = "download"
	ACTION_REMOVE        = "remove"
	ACTION_UPLOAD        = "upload"
	STRING_EMPTY_DEFAULT = "none"
)

type Action struct {
	Type       string
	FileTarget string
	Hash       string
}

var (
	// repository
	repositoryFile  repository.RepositoryFile
	repositoryChunk repository.RepositoryChunk
	// service
	serviceTemporaryArea interfaces.ServiceTemporaryArea
	serviceStorage       interfaces.ServiceStorage
	// error
	errorActionDefault = errors.New("ocorreu um erro nessa operação")
	errorFileNotFound  = errors.New("file not found")
)

func MakeAction() Action {
	repositoryFile = repository.MakeRepositoryFile(os.Getenv("ENGINE_COLLECTION"))
	repositoryChunk = repository.MakeRepositoryChunk(os.Getenv("ENGINE_COLLECTION"))
	serviceTemporaryArea = MakeServiceTemporaryArea()
	serviceStorage = MakeServiceStorage()

	return Action{
		Type:       STRING_EMPTY_DEFAULT,
		FileTarget: STRING_EMPTY_DEFAULT,
		Hash:       STRING_EMPTY_DEFAULT,
	}
}

func Execute(action interfaces.ServiceAction) error {
	switch action.GetActionType() {
	case ACTION_CLEAR:
		logger.GetLogActivity().WriteLog("Iniciando rotina de limpeza.")
		err := action.FeatureClear()
		if err != nil {
			return err
		}

		logger.GetLogActivity().WriteLog("Rotina de limpeza realizada com sucesso!")
	case ACTION_UPLOAD:
		logger.GetLogActivity().WriteLog("Iniciando rotina de carregamento de arquivo.")
		err := action.FeatureUpload()
		if err != nil {
			return err
		}

		logger.GetLogActivity().WriteLog("Arquivo carregado com sucesso!")
	case ACTION_DOWNLOAD:
		logger.GetLogActivity().WriteLog("Iniciando rotina de descarregamento de arquivo.")
		err := action.FeatureDownload()
		if err != nil {
			return err
		}

		logger.GetLogActivity().WriteLog("Arquivo descarregado com sucesso!")
	case ACTION_REMOVE:
		logger.GetLogActivity().WriteLog("Iniciando rotina de remoção de arquivo.")
		err := action.FeatureRemove()
		if err != nil {
			return err
		}

		logger.GetLogActivity().WriteLog("Arquivo removido com sucesso!")
	}

	return nil
}

func (ac Action) GetActionType() string {
	return ac.Type
}

func (ac Action) LogErrorWrite(action string, msg string) {
	logger.GetLogError().WriteLog(fmt.Sprintf("%s: %s", action, msg))
}
