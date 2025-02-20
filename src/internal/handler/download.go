package handler

import (
	"fmt"
	"io"
	"lab/src/logger"
	"lab/src/repository"
	"os"
)

func (ac *Action) FeatureDownload() error {
	logger.GetLogActivity().WriteLog("Carregando Hash do arquivo...")
	hashOriginalFile, err := repositoryFile.GetHashByName(ac.FileTarget)
	if err != nil {
		ac.LogErrorWrite(
			ACTION_DOWNLOAD,
			fmt.Sprintf(
				"file not found by name: %s",
				ac.FileTarget))

		if err == repository.ErrorRecordNotFound {
			return fmt.Errorf("arquivo não encontrado: %s", ac.FileTarget)
		}

		return errorActionDefault
	}

	logger.GetLogActivity().WriteLog("Carregando lista de chunk do arquivo...")
	list, err := repositoryChunk.GetChunkHashListByHashOriginalFile(hashOriginalFile)
	if err != nil {
		ac.LogErrorWrite(
			ACTION_DOWNLOAD,
			fmt.Sprintf(
				" fail to get chunk list by original file hash - %s",
				hashOriginalFile))

		if err == repository.ErrorRecordNotFound {
			return fmt.Errorf("arquivo não encontrado: %s", ac.FileTarget)
		}

		return errorActionDefault
	}

	logger.GetLogActivity().WriteLog("Gerando arquivo final...")
	return ac.generateFileByChunkHashList(list, "")
}

func (ac *Action) generateFileByChunkHashList(chunks []string, targetFolder string) error {
	out, err := os.Create(targetFolder + ac.FileTarget)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, file := range chunks {
		in, err := serviceStorage.GetFile(file + ".bin")
		if err != nil {
			return err
		}
		defer in.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
	}

	return nil
}
