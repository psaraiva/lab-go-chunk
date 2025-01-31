package repository

import (
	"database/sql"
	"errors"
	"lab/src/model"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func setUpRepositoryFileSqliteTestDB(t *testing.T) (*sql.DB, RepositoryFileSqlite) {
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

	return db, RepositoryFileSqlite{}
}

func setDownRepositoryFileSqliteTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`DROP TABLE files`)
	if err != nil {
		t.Fatalf("Failed to drop table files: %v", err)
	}

	db.Close()
	err = os.Remove(os.Getenv("CONFIG_HOST_SQLITE"))
	if err != nil {
		t.Fatalf("Failed to remove test database: %v", err)
	}
}

func TestRepositoryFileSqliteFileCreate(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	id, err := repo.Create(model.File{
		Name: "test.txt",
		Hash: "123ABC456DEF789"},
	)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if id != 1 {
		t.Fatalf("Expected valid id, got %d", id)
	}
}

func TestRepositoryFileSqliteCreateUniqueConstraintName(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	_, err := repo.Create(model.File{
		Name: "test.txt",
		Hash: "ABC123",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create(model.File{
		Name: "test.txt",
		Hash: "123456789ABC",
	})

	expected := errors.New("UNIQUE constraint failed: files.name")
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepositoryFileSqliteCreateUniqueConstraintHash(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	_, err := repo.Create(model.File{
		Name: "test_a.txt",
		Hash: "123ABC456",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create(model.File{
		Name: "test_b.txt",
		Hash: "123ABC456",
	})

	expected := errors.New("UNIQUE constraint failed: files.hash")
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepositoryFileSqliteGetHashByName(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	expectedName := "test.txt"
	expectedHash := "123ABC"

	_, err := repo.Create(model.File{
		Name: expectedName,
		Hash: expectedHash,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	hash, err := repo.GetHashByName(expectedName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if hash != expectedHash {
		t.Fatalf("Expected hash '%s', got %s", expectedHash, hash)
	}
}

func TestRepositoryFileSqliteGetHashByNameNotFound(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	hash := "ABC123"
	_, err := repo.Create(model.File{
		Name: "test.txt",
		Hash: hash,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.GetHashByName(hash)
	expected := errors.New("record not found")
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepositoryFileSqliteIsExistsByHash(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	hash := "ABC123"
	_, err := repo.Create(model.File{
		Name: "test.txt",
		Hash: hash,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	exists, err := repo.IsExistsByHash(hash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !exists {
		t.Fatalf("Expected file to exist, but it does not")
	}
}

func TestRepositoryFileSqliteIsExistsByHashNotFound(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	name := "test.txt"
	_, err := repo.Create(model.File{
		Name: name,
		Hash: "5f5adedeea13569a610a771521f66274",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	exists, err := repo.IsExistsByHash(name)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if exists {
		t.Fatalf("Expected file not to exist, but it does exists")
	}
}

func TestRepositoryFileSqliteRemoveByHash(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	hash := "ABCDEF123456"
	_, err := repo.Create(model.File{
		Name: "test.txt",
		Hash: hash,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repo.RemoveByHash(hash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	exists, err := repo.IsExistsByHash(hash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if exists {
		t.Fatalf("Expected file to be removed, but it still exists")
	}
}

func TestRepositoryFileSqliteRemoveByHashNotFound(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	_, err := repo.Create(model.File{
		Name: "test.txt",
		Hash: "ABCDEF123456"},
	)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repo.RemoveByHash("123ABC456")
	expected := errors.New("record not found")
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepositoryFileSqliteRemoveAll(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	_, err := repo.Create(model.File{
		Name: "test_a.txt",
		Hash: "5f5adedeea13569a610a771521f66274"},
	)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create(model.File{
		Name: "test_b.txt",
		Hash: "69e13300af627698d1b16901d82a28ce"},
	)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repo.RemoveAll()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var count int
	query := `SELECT COUNT(id) FROM files`
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 0 {
		t.Fatalf("Expected no rows, got %v", count)
	}

	var seq int
	query = `SELECT seq FROM sqlite_sequence WHERE name='files'`
	err = db.QueryRow(query).Scan(&seq)
	if err != sql.ErrNoRows {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestRepositoryFileSqliteGetIdByHash(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	hash := "ABCDEF123456"
	expectedId, err := repo.Create(model.File{
		Name: "test.txt",
		Hash: hash},
	)

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

func TestRepositoryFileSqliteGetIdByHashNotFound(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	_, err := repo.GetIdByHash("test.txt")
	expected := errors.New("record not found")
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}
