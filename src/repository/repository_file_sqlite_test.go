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

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS files (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE,
        hash TEXT NOT NULL UNIQUE
    )`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("Failed to create file table: %v", err)
	}

	return db, RepositoryFileSqlite{}
}

func setDownRepositoryFileSqliteTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`DROP TABLE files`)
	if err != nil {
		t.Fatalf("Failed to drop files table: %v", err)
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

	file := model.File{Name: "test.txt", Hash: "5f5adedeea13569a610a771521f66274"}
	id, err := repo.Create(file)
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

	file := model.File{Name: "test.txt", Hash: "5f5adedeea13569a610a771521f66274"}
	_, _ = repo.Create(file)

	file = model.File{Name: "test.txt", Hash: "69e13300af627698d1b16901d82a28ce"}
	_, err := repo.Create(file)

	expected := errors.New("UNIQUE constraint failed: files.name")
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepositoryFileSqliteCreateUniqueConstraintHash(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	file := model.File{Name: "testA.txt", Hash: "5f5adedeea13569a610a771521f66274"}
	_, _ = repo.Create(file)

	file = model.File{Name: "testB.txt", Hash: "5f5adedeea13569a610a771521f66274"}
	_, err := repo.Create(file)

	expected := errors.New("UNIQUE constraint failed: files.hash")
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepositoryFileSqliteGetHashByName(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	expected := "5f5adedeea13569a610a771521f66274"
	file := model.File{Name: "test.txt", Hash: expected}
	_, err := repo.Create(file)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	hash, err := repo.GetHashByName(file.Name)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if hash != expected {
		t.Fatalf("Expected hash '%s', got %s", expected, hash)
	}
}

func TestRepositoryFileSqliteGetHashByNameNotFound(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	file := model.File{Name: "test.txt", Hash: "5f5adedeea13569a610a771521f66274"}
	_, err := repo.Create(file)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := errors.New("record not found")
	_, err = repo.GetHashByName(file.Hash)
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepositoryFileSqliteIsExistsByHash(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	file := model.File{Name: "test.txt", Hash: "5f5adedeea13569a610a771521f66274"}
	_, err := repo.Create(file)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	exists, err := repo.IsExistsByHash(file.Hash)
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

	file := model.File{Name: "test.txt", Hash: "5f5adedeea13569a610a771521f66274"}
	_, err := repo.Create(file)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	exists, err := repo.IsExistsByHash(file.Name)
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

	file := model.File{Name: "test.txt", Hash: "5f5adedeea13569a610a771521f66274"}
	_, err := repo.Create(file)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repo.RemoveByHash(file.Hash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	exists, err := repo.IsExistsByHash(file.Hash)
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

	file := model.File{Name: "test.txt", Hash: "5f5adedeea13569a610a771521f66274"}
	_, err := repo.Create(file)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := errors.New("record not found")
	err = repo.RemoveByHash(file.Name)
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepositoryFileSqliteRemoveAll(t *testing.T) {
	db, repo := setUpRepositoryFileSqliteTestDB(t)
	defer setDownRepositoryFileSqliteTestDB(t, db)

	file1 := model.File{Name: "test_a.txt", Hash: "5f5adedeea13569a610a771521f66274"}
	_, err := repo.Create(file1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	file2 := model.File{Name: "test_b.txt", Hash: "69e13300af627698d1b16901d82a28ce"}
	_, err = repo.Create(file2)
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

	if count > 0 {
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

	file := model.File{Name: "test.txt", Hash: "5f5adedeea13569a610a771521f66274"}
	expectedId, err := repo.Create(file)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	id, err := repo.GetIdByHash(file.Hash)
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

	file := model.File{Name: "test.txt", Hash: "5f5adedeea13569a610a771521f66274"}
	_, err := repo.Create(file)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := errors.New("record not found")
	_, err = repo.GetIdByHash(file.Name)
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}
