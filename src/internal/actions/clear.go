package actions

import (
	"lab/src/logger"
	"os"
	"path/filepath"
)

func (ac *Action) FeatureClear() error {
	err := ac.resetChunkCollection()
	if err != nil {
		return err
	}

	err = ac.resetHashCollection()
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

func (ac *Action) resetChunkCollection() error {
	err := os.WriteFile(os.Getenv("JSON_FILE_CHUNK"), []byte("[]"), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (ac *Action) resetHashCollection() error {
	err := os.WriteFile(os.Getenv("JSON_FILE_HASH"), []byte("[]"), 0644)
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
