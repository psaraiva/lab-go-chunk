package repository

import (
	"encoding/json"
	"io"
	"lab/src/internal/entity"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func setUpRepositoryFileJsonTest(t *testing.T) RepositoryFileJson {
	err := godotenv.Load("../../src/.env.test")
	if err != nil {
		panic("Error loading .env.test file")
	}

	err = os.WriteFile(os.Getenv("COLLECTION_FILE_JSON"), []byte("[]"), 0644)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	return RepositoryFileJson{}
}

func setDownRepositoryFileJsonTest(t *testing.T) {
	err := os.Remove(os.Getenv("COLLECTION_FILE_JSON"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestRepositoryFileJsonCreate(t *testing.T) {
	repo := setUpRepositoryFileJsonTest(t)
	defer setDownRepositoryFileJsonTest(t)

	expected := entity.File{
		Name: "test.txt",
		Hash: "123ABC456DEF789",
	}

	_, err := repo.Create(expected)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	file_content, err := os.Open(os.Getenv("COLLECTION_FILE_JSON"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer file_content.Close()

	bytes, err := io.ReadAll(file_content)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	files := []entity.File{}
	err = json.Unmarshal(bytes, &files)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %v", len(files))
	}

	if files[0].Name != expected.Name {
		t.Fatalf("Expected %v, got %v", expected.Name, files[0].Name)
	}

	if files[0].Hash != expected.Hash {
		t.Fatalf("Expected %v, got %v", expected.Hash, files[0].Hash)
	}
}

func TestRepositoryFileJsonGetHashByName(t *testing.T) {
	repo := setUpRepositoryFileJsonTest(t)
	defer setDownRepositoryFileJsonTest(t)

	name := "test.txt"
	expectedHash := "123ABC456DEF789"
	_, err := repo.Create(entity.File{
		Name: name,
		Hash: expectedHash,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	fileHash, err := repo.GetHashByName(name)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if fileHash != expectedHash {
		t.Fatalf("Expected %v, got %v", expectedHash, fileHash)
	}
}

func TestRepositoryFileJsonGetHashByNameNotFound(t *testing.T) {
	repo := setUpRepositoryFileJsonTest(t)
	defer setDownRepositoryFileJsonTest(t)

	_, err := repo.GetHashByName("test.txt")
	expected := ErrorRecordNotFound
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepositoryFileJsonIsExistsByHash(t *testing.T) {
	repo := setUpRepositoryFileJsonTest(t)
	defer setDownRepositoryFileJsonTest(t)

	name := "test.txt"
	hash := "123ABC456DEF789"
	_, err := repo.Create(entity.File{
		Name: name,
		Hash: hash,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	isExists, err := repo.IsExistsByHash(hash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !isExists {
		t.Fatalf("Expected %v, got %v", "exists", "not exists")
	}
}

func TestRepositoryFileJsonIsNotExistsByHash(t *testing.T) {
	repo := setUpRepositoryFileJsonTest(t)
	defer setDownRepositoryFileJsonTest(t)

	hash := "123ABC456DEF789"
	isExists, err := repo.IsExistsByHash(hash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if isExists {
		t.Fatalf("Expected %v, got %v", "exists", "not exists")
	}
}

func TestRepositoryFileJsonRemove(t *testing.T) {
	repo := setUpRepositoryFileJsonTest(t)
	defer setDownRepositoryFileJsonTest(t)

	name := "test.txt"
	hash := "123ABC456DEF789"
	_, err := repo.Create(entity.File{
		Name: name,
		Hash: hash,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	isExists, err := repo.IsExistsByHash(hash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !isExists {
		t.Fatalf("Expected %v, got %v", "exists", "not exists")
	}

	err = repo.RemoveByHash(hash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	isExists, err = repo.IsExistsByHash(hash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if isExists {
		t.Fatalf("Expected %v, got %v", "not exists", "exists")
	}
}

func TestRepositoryFileJsonRemoveNotFound(t *testing.T) {
	repo := setUpRepositoryFileJsonTest(t)
	defer setDownRepositoryFileJsonTest(t)

	hash := "123ABC456DEF789"
	err := repo.RemoveByHash(hash)
	expected := ErrorRecordNotFound
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepositoryFileJsonRemoveAll(t *testing.T) {
	repo := setUpRepositoryFileJsonTest(t)
	defer setDownRepositoryFileJsonTest(t)

	expected := 0
	_, err := repo.Create(entity.File{
		Name: "test_a.txt",
		Hash: "123ABC456DEF780",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create(entity.File{
		Name: "test_b.txt",
		Hash: "123ABC456DEF781",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create(entity.File{
		Name: "test_c.txt",
		Hash: "123ABC456DEF782",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repo.RemoveAll()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	file_content, err := os.Open(os.Getenv("COLLECTION_FILE_JSON"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer file_content.Close()

	bytes, err := io.ReadAll(file_content)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	files := []entity.File{}
	err = json.Unmarshal(bytes, &files)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(files) != expected {
		t.Fatalf("Expected %v file, got %v", expected, len(files))
	}
}

func TestRepositoryFileJsonRemoveAllEmptyList(t *testing.T) {
	repo := setUpRepositoryFileJsonTest(t)
	defer setDownRepositoryFileJsonTest(t)

	expected := 0
	err := repo.RemoveAll()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	file_content, err := os.Open(os.Getenv("COLLECTION_FILE_JSON"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer file_content.Close()

	bytes, err := io.ReadAll(file_content)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	files := []entity.File{}
	err = json.Unmarshal(bytes, &files)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(files) != expected {
		t.Fatalf("Expected %v file, got %v", expected, len(files))
	}
}
