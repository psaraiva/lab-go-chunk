package repository

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func setUpRepostoryChunkHashSqliteTest(t *testing.T) (*sql.DB, RepositoryChunkHashSqlite) {
	err := godotenv.Load("../../src/.env.test")
	if err != nil {
		panic("Error loading .env.test file")
	}

	db, err := sql.Open("sqlite3", os.Getenv("CONFIG_HOST_SQLITE"))
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	ddl := `
    CREATE TABLE chunk_hashes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        hash TEXT NOT NULL UNIQUE
    )`

	_, err = db.Exec(ddl)
	if err != nil {
		t.Fatalf("Failed to create table chunk_hashes: %v", err)
	}

	return db, RepositoryChunkHashSqlite{}
}

func setDownRepostoryChunkHashSqliteTest(t *testing.T) {
	err := os.Remove(os.Getenv("CONFIG_HOST_SQLITE"))
	if err != nil {
		t.Fatalf("Failed to remove test database: %v", err)
	}
}

func TestRepostoryChunkHashSqliteCreate(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTest(t)
	defer setDownRepostoryChunkHashSqliteTest(t)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	id, err := repo.Create("123ABC", tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if id != 1 {
		t.Fatalf("Expected valid id, got %d", id)
	}
}

func TestRepostoryChunkHashSqliteCreateUniqueConstraintHash(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTest(t)
	defer setDownRepostoryChunkHashSqliteTest(t)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	hash := "12ABC456"
	_, err = repo.Create(hash, tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create(hash, tx)
	expected := errors.New("UNIQUE constraint failed: chunk_hashes.hash")
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepostoryChunkHashSqliteGetIdByHash(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTest(t)
	defer setDownRepostoryChunkHashSqliteTest(t)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	hash := "123789ABC"
	expectedId, err := repo.Create(hash, tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	id, err := repo.GetIdByHash(hash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if id != expectedId {
		t.Fatalf("Expected id %d, got %d", expectedId, id)
	}
}

func TestRepostoryChunkHashSqliteGetIdByHashNotFound(t *testing.T) {
	_, repo := setUpRepostoryChunkHashSqliteTest(t)
	defer setDownRepostoryChunkHashSqliteTest(t)

	_, err := repo.GetIdByHash("ABC123456")
	expected := ErrorRecordNotFound
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepostoryChunkHashSqliteGetHashById(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTest(t)
	defer setDownRepostoryChunkHashSqliteTest(t)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	expected := "ABC123DEF"
	id, err := repo.Create(expected, tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx.Commit()

	hash, err := repo.GetHashById(id)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if hash != expected {
		t.Fatalf("Expected hash '%s', got %s", expected, hash)
	}
}

func TestRepostoryChunkHashSqliteGetHashByIdNotFound(t *testing.T) {
	_, repo := setUpRepostoryChunkHashSqliteTest(t)
	defer setDownRepostoryChunkHashSqliteTest(t)

	expected := ErrorRecordNotFound
	_, err := repo.GetHashById(7777)
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepostoryChunkHashSqliteRemoveAllWithTransaction(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTest(t)
	defer setDownRepostoryChunkHashSqliteTest(t)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	_, err = repo.Create("ABC123", tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create("123456789ABCDEFG", tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx.Commit()

	tx, err = db.Begin()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repo.RemoveAllWithTransaction(tx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var count int
	query := `SELECT COUNT(id) FROM chunk_hashes`
	err = tx.QueryRow(query).Scan(&count)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count > 0 {
		t.Fatalf("Expected no rows, got %v", count)
	}

	var seq int
	query = `SELECT seq FROM sqlite_sequence WHERE name='chunk_hashes'`
	err = tx.QueryRow(query).Scan(&seq)
	if err != sql.ErrNoRows {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx.Rollback()
}

func TestRepostoryChunkHashSqliteRemoveByChunkHashWithTransactionIds(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTest(t)
	defer setDownRepostoryChunkHashSqliteTest(t)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	id1, err := repo.Create("ABC123", tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	id2, err := repo.Create("ABC456789123DEF", tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx.Commit()
	ids := []int64{id1, id2}

	tx, err = db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repo.RemoveByIdsWithTransaction(ids, tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx.Commit()

	expected := ErrorRecordNotFound
	_, err = repo.GetHashById(id1)
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}

	_, err = repo.GetHashById(id2)
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepostoryChunkHashSqliteRemoveByChunkHashIdsWithTransactionNotFound(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTest(t)
	defer setDownRepostoryChunkHashSqliteTest(t)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	ids := []int64{7777}
	err = repo.RemoveByIdsWithTransaction(ids, tx)
	expected := ErrorRecordNotFound
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepostoryChunkHashSqliteRemoveByChunkHashIdWithTransaction(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTest(t)
	defer setDownRepostoryChunkHashSqliteTest(t)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	id, err := repo.Create("ABC123456", tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx.Commit()
	tx, err = db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repo.RemoveByIdWithTransaction(id, tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx.Commit()

	_, err = repo.GetHashById(id)
	expected := ErrorRecordNotFound
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepostoryChunkHashSqliteRemoveByChunkHashIdWithTransactionNotFound(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTest(t)
	defer setDownRepostoryChunkHashSqliteTest(t)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	err = repo.RemoveByIdWithTransaction(7777, tx)
	expected := ErrorRecordNotFound
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}
