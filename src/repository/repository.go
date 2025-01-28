package repository

import (
	"database/sql"
	"lab/src/model"
	"os"
)

const (
	ENGINE_JSON   = "json"
	ENGINE_SQLITE = "sqlite"
)

type RepositoryFile interface {
	Create(model.File) (int64, error)
	GetHashByName(string) (string, error)
	IsExistsByHashFile(string) (bool, error)
	RemoveAll() error
	RemoveByHashFile(string) error
}

type RepositoryChunk interface {
	CountChunkHash(string) (int64, error)
	Create(model.Chunk) (int64, error)
	GetChunkHashListByHashOriginalFile(string) ([]string, error)
	RemoveAll() error
	RemoveByHashOriginalFile(string) ([]string, error)
}

type CountItem struct {
	Id    int64
	Total int
}

func MakeRepositoryFile(engine string) RepositoryFile {
	switch engine {
	case ENGINE_JSON:
		return RepositoryFileJson{}
	case ENGINE_SQLITE:
		return RepositoryFileSqlite{}
	}

	panic("Invalid engine to repository file: " + engine)
}

func MakeRepositoryChunk(engine string) RepositoryChunk {
	switch engine {
	case ENGINE_JSON:
		return RepositoryChunkJson{}
	case ENGINE_SQLITE:
		return RepositoryChunkSqlite{}
	}

	panic("Invalid engine to repository chunk: " + engine)
}

func getConectionSqlite() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", os.Getenv("CONFIG_HOST_SQLITE"))
	if err != nil {
		return nil, err
	}

	return db, ping(db)
}

func ping(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		defer db.Close()
	}
	return err
}
