package repository

import (
	"lab/src/model"
)

const (
	ENGINE_JSON   = "json"
	ENGINE_SQLITE = "sqlite"
)

type RepositoryFile interface {
	Create(model.File) error
	GetHashByFileName(string) (string, error)
	IsExistsByHashFile(string) (bool, error)
	RemoveAll() error
	RemoveByHashFile(string) error
}

type RepositoryChunk interface {
	Create(model.Chunk) error
	GetChunkHashListByHashOriginalFile(string) ([]string, error)
	GetCountChunkMap([]model.Chunk) map[string]int
	IsChunkCanBeRemoved(string) (bool, error)
	RemoveAll() error
	RemoveByHashOriginalFile(string) error
}

func MakeRepositoryFile(engine string) RepositoryFile {
	if ENGINE_JSON == engine {
		return RepositoryFileJson{}
	}

	panic("Invalid engine to repository file: " + engine)
}

func MakeRepositoryChunk(engine string) RepositoryChunkJson {
	if ENGINE_JSON == engine {
		return RepositoryChunkJson{}
	}

	panic("Invalid engine to repository chunk: " + engine)
}
