package service

import (
	"lab/src/logger"
)

func (ac *Action) FeatureClear() error {
	err := repositoryChunk.RemoveAll()
	if err != nil {
		return err
	}

	err = repositoryFile.RemoveAll()
	if err != nil {
		return err
	}

	err = serviceStorage.Clear()
	if err != nil {
		return err
	}

	err = serviceTemporaryArea.Clear()
	if err != nil {
		return err
	}

	err = logger.GetLogError().Clear()
	if err != nil {
		return err
	}

	err = logger.GetLogActivity().Clear()
	if err != nil {
		return err
	}

	return nil
}
