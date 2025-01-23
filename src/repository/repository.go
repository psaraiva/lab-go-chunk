package repository

import (
	"lab/src/model"
)

const (
	ENGINE_JSON   = "json"
	ENGINE_SQLITE = "sqlite"
)

type RepositoryFile interface {
	Create(file model.File) error
	GetHashByFileName(fileName string) (string, error)
	IsExistsByHashFile(hashFile string) (bool, error)
	RemoveAll() error
	RemoveByHashFile(hashFile string) error
}

type RepositoryChunkItem interface {
	Create(model.ChunkItem) error
	GetChunkHashListByHashFile(hashFile string) ([]string, error)
	GetCountChunkMap(chunkList []model.ChunkItem) map[string]int
	IsChunkCanBeRemoved(chunk string) (bool, error)
	RemoveAll() error
	RemoveByHashFile(hashFile string) error
}

func MakeRepositoryFile(engine string) RepositoryFile {
	if ENGINE_JSON == engine {
		return RepositoryFileJson{}
	}

	panic("Invalid engine to repository file" + engine)
}

func MakeRepositoryChunkItem(engine string) RepositoryChunkItem {
	if ENGINE_JSON == engine {
		return RepositoryChunkItemJson{}
	}

	panic("Invalid engine to repository chunk item, " + engine)
}
