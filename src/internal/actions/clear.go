package actions

import (
	"lab/src/logger"
	"os"
	"path/filepath"
)

func (ac *Action) FeatureClear() error {
	err := repositoryChunkItem.RemoveAll()
	if err != nil {
		return err
	}

	err = repositoryFile.RemoveAll()
	if err != nil {
		return err
	}

	err = ac.cleanStorage()
	if err != nil {
		return err
	}

	err = ac.cleanTmp()
	if err != nil {
		return err
	}

	err = logger.GetLogError().ClearLog()
	if err != nil {
		return err
	}

	err = logger.GetLogActivity().ClearLog()
	if err != nil {
		return err
	}

	return nil
}

func (ac *Action) cleanStorage() error {
	err := filepath.Walk(os.Getenv("FOLDER_STORAGE"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return os.Remove(path)
	})
	return err
}

func (ac *Action) cleanTmp() error {
	err := filepath.Walk(os.Getenv("FOLDER_TMP"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return os.Remove(path)
	})
	return err
}
