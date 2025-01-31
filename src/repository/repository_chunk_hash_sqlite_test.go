package repository

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func setUpRepostoryChunkHashSqliteTestDB(t *testing.T) (*sql.DB, RepositoryChunkHashSqlite) {
	err := godotenv.Load("../../src/.env.test")
	if err != nil {
		panic("Error loading .env.test file")
	}

	db, err := sql.Open("sqlite3", os.Getenv("CONFIG_HOST_SQLITE"))
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS chunk_hashes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        hash TEXT NOT NULL UNIQUE
    )`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("Failed to create chunk hashes table: %v", err)
	}

	return db, RepositoryChunkHashSqlite{}
}

func setDownRepostoryChunkHashSqliteTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`DROP TABLE chunk_hashes`)
	if err != nil {
		t.Fatalf("Failed to drop chunk hashes table: %v", err)
	}

	db.Close()
	err = os.Remove(os.Getenv("CONFIG_HOST_SQLITE"))
	if err != nil {
		t.Fatalf("Failed to remove test database: %v", err)
	}
}

func TestRepostoryChunkHashSqliteCreate(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTestDB(t)
	defer setDownRepostoryChunkHashSqliteTestDB(t, db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	id, err := repo.Create("5f5adedeea13569a610a771521f66274", tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if id != 1 {
		t.Fatalf("Expected valid id, got %d", id)
	}
}

func TestRepostoryChunkHashSqliteCreateUniqueConstraintHash(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTestDB(t)
	defer setDownRepostoryChunkHashSqliteTestDB(t, db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	hash := "5f5adedeea13569a610a771521f66274"
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
	db, repo := setUpRepostoryChunkHashSqliteTestDB(t)
	defer setDownRepostoryChunkHashSqliteTestDB(t, db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	hash := "5f5adedeea13569a610a771521f66274"
	expectedId, err := repo.Create(hash, tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx.Commit()

	id, err := repo.GetIdByHash(hash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if id != expectedId {
		t.Fatalf("Expected id %d, got %d", expectedId, id)
	}
}

func TestRepostoryChunkHashSqliteGetIdByHashNotFound(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTestDB(t)
	defer setDownRepostoryChunkHashSqliteTestDB(t, db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	_, err = repo.Create("5f5adedeea13569a610a771521f66274", tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx.Commit()

	expected := errors.New("record not found")
	_, err = repo.GetIdByHash("test")
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepostoryChunkHashSqliteGetHashById(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTestDB(t)
	defer setDownRepostoryChunkHashSqliteTestDB(t, db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	expected := "5f5adedeea13569a610a771521f66274"
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
	db, repo := setUpRepostoryChunkHashSqliteTestDB(t)
	defer setDownRepostoryChunkHashSqliteTestDB(t, db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	_, err = repo.Create("5f5adedeea13569a610a771521f66274", tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx.Commit()

	expected := errors.New("record not found")
	_, err = repo.GetHashById(10)
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepostoryChunkHashSqliteRemoveAllWithTransaction(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTestDB(t)
	defer setDownRepostoryChunkHashSqliteTestDB(t, db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	_, err = repo.Create("5f5adedeea13569a610a771521f66274", tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create("69e13300af627698d1b16901d82a28ce", tx)
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

func TestRepostoryChunkHashSqliteRemoveByChunkHashIds(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTestDB(t)
	defer setDownRepostoryChunkHashSqliteTestDB(t, db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	id1, err := repo.Create("5f5adedeea13569a610a771521f66274", tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	id2, err := repo.Create("69e13300af627698d1b16901d82a28ce", tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx.Commit()
	ids := []int64{id1, id2}

	tx, err = db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repo.RemoveByIds(ids, tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx.Commit()

	expected := errors.New("record not found")
	_, err = repo.GetHashById(id1)
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}

	_, err = repo.GetHashById(id2)
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepostoryChunkHashSqliteRemoveByChunkHashIdsNotFound(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTestDB(t)
	defer setDownRepostoryChunkHashSqliteTestDB(t, db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	ids := []int64{10}
	err = repo.RemoveByIds(ids, tx)
	expected := errors.New("record not found")
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepostoryChunkHashSqliteRemoveByChunkHashId(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTestDB(t)
	defer setDownRepostoryChunkHashSqliteTestDB(t, db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	id, err := repo.Create("5f5adedeea13569a610a771521f66274", tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx.Commit()
	tx, err = db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repo.RemoveById(id, tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tx.Commit()

	expected := errors.New("record not found")
	_, err = repo.GetHashById(id)
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepostoryChunkHashSqliteRemoveByChunkHashIdNotFound(t *testing.T) {
	db, repo := setUpRepostoryChunkHashSqliteTestDB(t)
	defer setDownRepostoryChunkHashSqliteTestDB(t, db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer tx.Rollback()

	err = repo.RemoveById(10, tx)
	expected := errors.New("record not found")
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}
