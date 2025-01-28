package service

import (
	"io"
	"lab/src/logger"
	"os"
)

func (ac *Action) FeatureDownload() error {
	logger.GetLogActivity().WriteLog("Carregando Hash do arquivo...")
	hashOriginalFile, err := repositoryFile.GetHashByName(ac.FileTarget)
	if err != nil {
		return err
	}

	logger.GetLogActivity().WriteLog("Carregando lista de chunk do arquivo...")
	list, err := repositoryChunk.GetChunkHashListByHashOriginalFile(hashOriginalFile)
	if err != nil {
		return err
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
