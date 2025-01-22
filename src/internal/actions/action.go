package actions

import (
	"encoding/json"
	"fmt"
	"lab/src/interfaces"
	"lab/src/logger"
	"os"
)

const (
	ACTION_CLEAR    = "clear"
	ACTION_DOWNLOAD = "download"
	ACTION_REMOVE   = "remove"
	ACTION_UPLOAD   = "upload"
	CHUNK_SIZE      = 1 * 1024 // 1Kb
)

type Action struct {
	Type       string
	FileTarget string
	Hash       string
}

type hashItem struct {
	Hash string
	Name string
}

type chunkItem struct {
	HashFile string
	Chunk    []string
}

func MakeAction() Action {
	return Action{Type: "none", FileTarget: "none", Hash: "none"}
}

func Execute(action interfaces.ActionBase) error {
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

func (ac *Action) GetActionType() string {
	return ac.Type
}

func (ac *Action) getHashByFileName(fileName string) (string, error) {
	jsonFile, err := os.Open(os.Getenv("JSON_FILE_HASH"))
	if err != nil {
		return "", err
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	hashList := []hashItem{}
	err = decoder.Decode(&hashList)
	if err != nil {
		return "", err
	}

	for _, item := range hashList {
		if item.Name == fileName {
			return item.Hash, nil
		}
	}

	return "", fmt.Errorf("arquivo não encontrado")
}

func (ac *Action) getChunksByHash(hash string) ([]string, error) {
	jsonChunk, err := os.Open(os.Getenv("JSON_FILE_CHUNK"))
	if err != nil {
		return nil, err
	}
	defer jsonChunk.Close()

	decoder := json.NewDecoder(jsonChunk)
	chunkList := []chunkItem{}
	err = decoder.Decode(&chunkList)
	if err != nil {
		return nil, err
	}

	for _, item := range chunkList {
		if item.HashFile == hash {
			return item.Chunk, nil
		}
	}

	return nil, fmt.Errorf("arquivo não encontrado")
}
