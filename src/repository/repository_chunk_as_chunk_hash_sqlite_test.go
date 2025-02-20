package repository

import (
	"database/sql"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setUpRepostoryChunkAsChunkHashSqliteTest(t *testing.T) (*sql.DB, RepositoryChunkHasChunkHashSqlite) {
	err := godotenv.Load("../../src/.env.test")
	if err != nil {
		panic("Error loading .env.test file")
	}

	db, err := sql.Open("sqlite3", os.Getenv("CONFIG_HOST_SQLITE"))
	assert.NoError(t, err)

	ddl := `
    CREATE TABLE chunks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        size INTEGER NOT NULL,
        file_id INTEGER NOT NULL)`

	_, err = db.Exec(ddl)
	assert.NoError(t, err)

	ddl = `
    CREATE TABLE chunk_hashes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        hash TEXT NOT NULL UNIQUE)`

	_, err = db.Exec(ddl)
	assert.NoError(t, err)

	ddl = `
    CREATE TABLE chunks_has_chunk_hashes (
        chunk_id INTEGER NOT NULL,
        chunk_hash_id INTEGER NOT NULL,
        PRIMARY KEY (chunk_id, chunk_hash_id),
        FOREIGN KEY (chunk_id) REFERENCES chunks(id),
        FOREIGN KEY (chunk_hash_id) REFERENCES chunk_hashes(id))`

	_, err = db.Exec(ddl)
	assert.NoError(t, err)
	return db, RepositoryChunkHasChunkHashSqlite{}
}

func setDownRepostoryChunkAsChunkHashSqliteTest(t *testing.T) {
	assert.NoError(t, os.Remove(os.Getenv("CONFIG_HOST_SQLITE")))
}

func TestRepositoryChunkHasChunkHashSqliteCreate(t *testing.T) {
	db, repo := setUpRepostoryChunkAsChunkHashSqliteTest(t)
	defer setDownRepostoryChunkAsChunkHashSqliteTest(t)

	db.Exec(`INSERT INTO chunks (file_id, size) VALUES (1, 1024)`)
	db.Exec(`INSERT INTO chunk_hashes (hash) VALUES ('ABCD1234')`)

	tx, err := db.Begin()
	assert.NoError(t, err)
	assert.NoError(t, repo.Create(1, 1, tx))

	var count int64
	assert.NoError(
		t,
		tx.QueryRow(`
            SELECT COUNT(chunk_hash_id)
              FROM chunks_has_chunk_hashes`).Scan(&count))
	assert.Equal(t, count, int64(1))
	assert.NoError(t, tx.Rollback())
}

func TestRepositoryChunkHasChunkHashSqliteRemoveAllWithTransaction(t *testing.T) {
	db, repo := setUpRepostoryChunkAsChunkHashSqliteTest(t)
	defer setDownRepostoryChunkAsChunkHashSqliteTest(t)

	db.Exec(`INSERT INTO chunks (file_id, size) VALUES (1, 1024)`)
	db.Exec(`INSERT INTO chunks (file_id, size) VALUES (2, 1024)`)
	db.Exec(`INSERT INTO chunk_hashes (hash) VALUES ('ABCD1234')`)
	db.Exec(`INSERT INTO chunk_hashes (hash) VALUES ('EFGH5678')`)

	tx, err := db.Begin()
	assert.NoError(t, err)
	assert.NoError(t, repo.Create(1, 1, tx))
	assert.NoError(t, repo.Create(2, 2, tx))
	assert.NoError(t, repo.RemoveAllWithTransaction(tx))

	var count int64
	err = db.QueryRow(`SELECT COUNT(chunk_hash_id) FROM chunks_has_chunk_hashes`).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, count, int64(0))
	assert.NoError(t, tx.Rollback())
}

func TestRepositoryChunkHasChunkHashSqliteRemoveByChunkId(t *testing.T) {
	db, repo := setUpRepostoryChunkAsChunkHashSqliteTest(t)
	defer setDownRepostoryChunkAsChunkHashSqliteTest(t)

	db.Exec(`INSERT INTO chunks (file_id, size) VALUES (1, 1024)`)
	db.Exec(`INSERT INTO chunk_hashes (hash) VALUES ('ABCD1234')`)

	tx, err := db.Begin()
	assert.NoError(t, err)
	assert.NoError(t, repo.Create(1, 1, tx))
	assert.NoError(t, repo.RemoveByChunkId(1, tx))

	var count int64
	assert.NoError(
		t,
		tx.QueryRow(`
          SELECT COUNT(chunk_hash_id)
            FROM chunks_has_chunk_hashes
           WHERE chunk_id = 1`).Scan(&count))
	assert.Equal(t, count, int64(0))
	assert.NoError(t, tx.Rollback())
}

func TestRepositoryChunkHasChunkHashSqliteRemoveByChunkIdNotFound(t *testing.T) {
	db, repo := setUpRepostoryChunkAsChunkHashSqliteTest(t)
	defer setDownRepostoryChunkAsChunkHashSqliteTest(t)

	tx, err := db.Begin()
	assert.NoError(t, err)
	assert.Equal(
		t,
		repo.RemoveByChunkId(1, tx).Error(),
		ErrorRecordNotFound.Error())
	assert.NoError(t, tx.Rollback())
}

func TestRepositoryChunkHasChunkHashSqliteCountChunkHashByChunkId(t *testing.T) {
	db, repo := setUpRepostoryChunkAsChunkHashSqliteTest(t)
	defer setDownRepostoryChunkAsChunkHashSqliteTest(t)

	db.Exec(`INSERT INTO chunks (file_id, size) VALUES (1, 1024)`)
	db.Exec(`INSERT INTO chunks (file_id, size) VALUES (2, 1024)`)
	db.Exec(`INSERT INTO chunk_hashes (hash) VALUES ('ABCD1234')`)
	db.Exec(`INSERT INTO chunk_hashes (hash) VALUES ('EFGH5678')`)

	tx, err := db.Begin()
	assert.NoError(t, err)
	assert.NoError(t, repo.Create(1, 1, tx))
	assert.NoError(t, repo.Create(1, 2, tx))
	assert.NoError(t, repo.Create(2, 1, tx))
	assert.NoError(t, tx.Commit())

	count, err := repo.CountChunkHashByChunkId(1)
	assert.NoError(t, err)
	assert.Equal(t, len(count), 2)
	assert.Equal(t, count[0].Id, int64(1))
	assert.Equal(t, count[0].Total, 2)
	assert.Equal(t, count[1].Id, int64(2))
	assert.Equal(t, count[1].Total, 1)
}

func TestRepositoryChunkHasChunkHashSqliteCountChunkHashByChunkIdNotFound(t *testing.T) {
	_, repo := setUpRepostoryChunkAsChunkHashSqliteTest(t)
	defer setDownRepostoryChunkAsChunkHashSqliteTest(t)

	count, err := repo.CountChunkHashByChunkId(1)
	assert.NoError(t, err)
	assert.Equal(t, len(count), 0)
}

func TestRepositoryChunkHasChunkHashSqliteGetChunkHashIdsByChunkId(t *testing.T) {
	db, repo := setUpRepostoryChunkAsChunkHashSqliteTest(t)
	defer setDownRepostoryChunkAsChunkHashSqliteTest(t)

	db.Exec(`INSERT INTO chunks (file_id, size) VALUES (1, 1024)`)
	db.Exec(`INSERT INTO chunk_hashes (hash) VALUES ('ABCD1234')`)
	db.Exec(`INSERT INTO chunk_hashes (hash) VALUES ('EFGH5678')`)

	tx, err := db.Begin()
	assert.NoError(t, err)
	assert.NoError(t, repo.Create(1, 1, tx))
	assert.NoError(t, repo.Create(1, 2, tx))
	assert.NoError(t, tx.Commit())

	ids, err := repo.GetChunkHashIdsByChunkId(1)
	assert.NoError(t, err)
	assert.Equal(t, len(ids), 2)
	assert.Equal(t, ids[0], int64(1))
	assert.Equal(t, ids[1], int64(2))
}

func TestRepositoryChunkHasChunkHashSqliteGetChunkHashIdsByChunkIdNotFound(t *testing.T) {
	_, repo := setUpRepostoryChunkAsChunkHashSqliteTest(t)
	defer setDownRepostoryChunkAsChunkHashSqliteTest(t)

	ids, err := repo.GetChunkHashIdsByChunkId(1)
	assert.NoError(t, err)
	assert.Equal(t, len(ids), 0)
}
