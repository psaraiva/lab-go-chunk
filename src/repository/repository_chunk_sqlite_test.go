package repository

import (
	"database/sql"
	"lab/src/internal/entity"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
)

func setUpRepositoryChunkSqliteTest(t *testing.T) (*sql.DB, RepositoryChunkSqlite) {
	err := godotenv.Load("../../src/.env.test")
	if err != nil {
		panic("Error loading .env.test file")
	}

	db, err := sql.Open("sqlite3", os.Getenv("CONFIG_HOST_SQLITE"))
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	ddl := `
	CREATE TABLE files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		hash TEXT NOT NULL UNIQUE
	)`

	_, err = db.Exec(ddl)
	if err != nil {
		t.Fatalf("Failed to create table files: %v", err)
	}

	ddl = `
    CREATE TABLE chunks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        size INTEGER NOT NULL,
        file_id INTEGER NOT NULL,
        FOREIGN KEY (file_id) REFERENCES files(id)
    )`

	_, err = db.Exec(ddl)
	if err != nil {
		t.Fatalf("Failed to create table chunks: %v", err)
	}

	ddl = `
    CREATE TABLE chunk_hashes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        hash TEXT NOT NULL UNIQUE
    )`

	_, err = db.Exec(ddl)
	if err != nil {
		t.Fatalf("Failed to create table chunk_hashes: %v", err)
	}

	ddl = `
	CREATE TABLE chunks_has_chunk_hashes (
		chunk_id INTEGER NOT NULL,
		chunk_hash_id INTEGER NOT NULL,
		PRIMARY KEY (chunk_id, chunk_hash_id),
		FOREIGN KEY (chunk_id) REFERENCES chunks(id),
		FOREIGN KEY (chunk_hash_id) REFERENCES chunk_hashes(id)
	)`

	_, err = db.Exec(ddl)
	if err != nil {
		t.Fatalf("Failed to create table chunks_has_chunk_hashes: %v", err)
	}

	return db, RepositoryChunkSqlite{}
}

func setDownRepositoryChunkSqliteTest(t *testing.T) {
	err := os.Remove(os.Getenv("CONFIG_HOST_SQLITE"))
	if err != nil {
		t.Fatalf("Failed to remove test database: %v", err)
	}
}

func TestRepositoryChunkSqliteCreate(t *testing.T) {
	db, repo := setUpRepositoryChunkSqliteTest(t)
	defer setDownRepositoryChunkSqliteTest(t)

	hash := "ABC123"
	_, err := db.Exec(`INSERT INTO files (name, hash) VALUES (?,?)`, "test.txt", hash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	id, err := repo.Create(entity.Chunk{
		HashOriginalFile: hash,
		HashList:         []string{"ABCDEF", "123456"},
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if id != 1 {
		t.Fatalf("Expected valid id, got %d", id)
	}
}

func TestRepositoryChunkSqliteCreateNotFoundFileHash(t *testing.T) {
	_, repo := setUpRepositoryChunkSqliteTest(t)
	defer setDownRepositoryChunkSqliteTest(t)

	_, err := repo.Create(entity.Chunk{
		HashOriginalFile: "123456789ABCDEFGHI",
		HashList:         []string{"ABCDEF", "123456"},
		Size:             1024,
	})

	expected := ErrorRecordNotFound
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepositoryChunkSqliteGetChunkHashListByHashOriginalFile(t *testing.T) {
	db, repo := setUpRepositoryChunkSqliteTest(t)
	defer setDownRepositoryChunkSqliteTest(t)

	hashFile := "123ABC123"
	_, err := db.Exec(`INSERT INTO files (name, hash) VALUES (?,?)`, "test.txt", hashFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedHashList := []string{"ABCDEF", "123456"}
	_, err = repo.Create(entity.Chunk{
		HashOriginalFile: hashFile,
		HashList:         expectedHashList,
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	chunkHashList, err := repo.GetChunkHashListByHashOriginalFile(hashFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(chunkHashList, expectedHashList) {
		t.Fatalf("Expected %v, got %v", expectedHashList, chunkHashList)
	}
}

func TestRepositoryChunkSqliteGetChunkHashListByHashOriginalFileNotFound(t *testing.T) {
	_, repo := setUpRepositoryChunkSqliteTest(t)
	defer setDownRepositoryChunkSqliteTest(t)

	_, err := repo.GetChunkHashListByHashOriginalFile("ABC456")
	expected := ErrorRecordNotFound
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepositoryChunkSqliteCountChunkHashZero(t *testing.T) {
	db, repo := setUpRepositoryChunkSqliteTest(t)
	defer setDownRepositoryChunkSqliteTest(t)

	_, err := db.Exec(`INSERT INTO files (name, hash) VALUES (?,?)`, "test_a.txt", "123ABC123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	countUsed, err := repo.CountUsedChunkHash("ABCDEF")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if countUsed != 0 {
		t.Fatalf("Expected count 1, got %d", countUsed)
	}
}

func TestRepositoryChunkSqliteCountChunkHashOne(t *testing.T) {
	db, repo := setUpRepositoryChunkSqliteTest(t)
	defer setDownRepositoryChunkSqliteTest(t)

	_, err := db.Exec(`INSERT INTO files (name, hash) VALUES (?,?)`, "test_a.txt", "123ABC123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create(entity.Chunk{
		HashOriginalFile: "123ABC123",
		HashList:         []string{"ABCDEF", "123456"},
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = db.Exec(`INSERT INTO files (name, hash) VALUES (?,?)`, "test_b.txt", "ABC123ABC")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	countUsed, err := repo.CountUsedChunkHash("ABCDEF")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if countUsed != 1 {
		t.Fatalf("Expected count 1, got %d", countUsed)
	}
}

func TestRepositoryChunkSqliteCountChunkHashTwo(t *testing.T) {
	db, repo := setUpRepositoryChunkSqliteTest(t)
	defer setDownRepositoryChunkSqliteTest(t)

	_, err := db.Exec(`INSERT INTO files (name, hash) VALUES (?,?)`, "test_a.txt", "123ABC123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create(entity.Chunk{
		HashOriginalFile: "123ABC123",
		HashList:         []string{"ABCDEF", "123456"},
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = db.Exec(`INSERT INTO files (name, hash) VALUES (?,?)`, "test_b.txt", "ABC123ABC")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create(entity.Chunk{
		HashOriginalFile: "ABC123ABC",
		HashList:         []string{"ABCDEF", "ABCDEF0123456"},
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	countUsed, err := repo.CountUsedChunkHash("ABCDEF")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if countUsed != 2 {
		t.Fatalf("Expected count 2, got %d", countUsed)
	}
}

func TestRepositoryChunkSqliteRemoveAll(t *testing.T) {
	db, repo := setUpRepositoryChunkSqliteTest(t)
	defer setDownRepositoryChunkSqliteTest(t)

	_, err := db.Exec(`INSERT INTO files (name, hash) VALUES (?,?)`, "test_a.txt", "123ABC123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create(entity.Chunk{
		HashOriginalFile: "123ABC123",
		HashList:         []string{"ABCDEF", "123456"},
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = db.Exec(`INSERT INTO files (name, hash) VALUES (?,?)`, "test_b.txt", "ABC123ABC")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create(entity.Chunk{
		HashOriginalFile: "ABC123ABC",
		HashList:         []string{"ABCDEF", "ABCDEF0123456"},
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repo.RemoveAll()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	count, err := repo.CountUsedChunkHash("ABCDEF")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 0 {
		t.Fatalf("Expected count 0, got %d", count)
	}

	repo.CountUsedChunkHash("123456")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 0 {
		t.Fatalf("Expected count 0, got %d", count)
	}

	repo.CountUsedChunkHash("ABCDEF0123456")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 0 {
		t.Fatalf("Expected count 0, got %d", count)
	}
}

func TestRepositoryChunkSqliteRemoveAllSequence(t *testing.T) {
	db, repo := setUpRepositoryChunkSqliteTest(t)
	defer setDownRepositoryChunkSqliteTest(t)

	_, err := db.Exec(`INSERT INTO files (name, hash) VALUES (?,?)`, "test_a.txt", "123ABC123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create(entity.Chunk{
		HashOriginalFile: "123ABC123",
		HashList:         []string{"ABCDEF", "123456"},
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = db.Exec(`INSERT INTO files (name, hash) VALUES (?,?)`, "test_b.txt", "ABC123ABC")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create(entity.Chunk{
		HashOriginalFile: "ABC123ABC",
		HashList:         []string{"ABCDEF", "ABCDEF0123456"},
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repo.RemoveAll()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	var seq int
	dml := `SELECT seq FROM sqlite_sequence WHERE name='chunk_hashes'`
	err = tx.QueryRow(dml).Scan(&seq)
	if err != sql.ErrNoRows {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx.Rollback()
}

func TestRepositoryChunkSqliteRemoveByHashOriginalFile(t *testing.T) {
	db, repo := setUpRepositoryChunkSqliteTest(t)
	defer setDownRepositoryChunkSqliteTest(t)

	originalHashFile := "123ABC123"
	_, err := db.Exec(`INSERT INTO files (name, hash) VALUES (?,?)`,
		"test_a.txt",
		originalHashFile)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedHashList := []string{"ABCDEF", "123456"}
	_, err = repo.Create(entity.Chunk{
		HashOriginalFile: originalHashFile,
		HashList:         expectedHashList,
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	chunkHashList, err := repo.RemoveByHashOriginalFile(originalHashFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(chunkHashList, expectedHashList) {
		t.Fatalf("Expected %v, got %v", expectedHashList, chunkHashList)
	}
}

func TestRepositoryChunkSqliteRemoveByHashOriginalFileNotFound(t *testing.T) {
	_, repo := setUpRepositoryChunkSqliteTest(t)
	defer setDownRepositoryChunkSqliteTest(t)

	_, err := repo.RemoveByHashOriginalFile("123ABC123")
	expected := ErrorRecordNotFound
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}
