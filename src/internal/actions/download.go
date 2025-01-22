package actions

import (
	"io"
	"lab/src/logger"
	"os"
)

func (ac *Action) FeatureDownload() error {
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

	logger.GetLogActivity().WriteLog("Gerando arquivo final...")
	return ac.generateFileByChunks(list, "")
}

func (ac *Action) generateFileByChunks(chunks []string, targetFolder string) error {
	out, err := os.Create(targetFolder + ac.FileTarget)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, file := range chunks {
		in, err := os.Open(os.Getenv("FOLDER_STORAGE") + "/" + file + ".bin")
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
