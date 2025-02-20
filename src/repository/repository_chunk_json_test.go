package repository

import (
	"encoding/json"
	"io"
	"lab/src/internal/entity"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func setUpRepositoryChunkJsonTest(t *testing.T) RepositoryChunkJson {
	err := godotenv.Load("../../src/.env.test")
	if err != nil {
		panic("Error loading .env.test file")
	}

	err = os.WriteFile(os.Getenv("COLLECTION_CHUNK_JSON"), []byte("[]"), 0644)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	return RepositoryChunkJson{}
}

func setDownRepositoryChunkJsonTest(t *testing.T) {
	err := os.Remove(os.Getenv("COLLECTION_CHUNK_JSON"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestRepositoryChunkJsonCreate(t *testing.T) {
	repo := setUpRepositoryChunkJsonTest(t)
	defer setDownRepositoryChunkJsonTest(t)

	hashOriginalFile := "ABCDEF123456789"
	HashList := []string{"123ABC456DEF789", "987GHI654JKL321"}
	size := 1024
	_, err := repo.Create(entity.Chunk{
		HashOriginalFile: hashOriginalFile,
		HashList:         HashList,
		Size:             size,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	file_content, err := os.Open(os.Getenv("COLLECTION_CHUNK_JSON"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer file_content.Close()

	bytes, err := io.ReadAll(file_content)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	files := []entity.Chunk{}
	err = json.Unmarshal(bytes, &files)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %v", len(files))
	}

	if files[0].HashOriginalFile != hashOriginalFile {
		t.Fatalf("Expected %v, got %v", hashOriginalFile, files[0].HashOriginalFile)
	}

	if len(files[0].HashList) != 2 {
		t.Fatalf("Expected 2, got %v", len(files[0].HashList))
	}

	if files[0].HashList[0] != HashList[0] {
		t.Fatalf("Expected %v, got %v", HashList[0], files[0].HashList[0])
	}

	if files[0].HashList[1] != HashList[1] {
		t.Fatalf("Expected %v, got %v", HashList[1], files[0].HashList[1])
	}

	if files[0].Size != size {
		t.Fatalf("Expected %v, got %v", size, files[0].Size)
	}
}

func TestRepositoryChunkJsonGetChunkHashListByHashOriginalFile(t *testing.T) {
	repo := setUpRepositoryChunkJsonTest(t)
	defer setDownRepositoryChunkJsonTest(t)

	hashOriginalFile := "ABCDEF123456789"
	expectedHashList := []string{"123ABC456DEF789"}
	size := 1024
	_, err := repo.Create(entity.Chunk{
		HashOriginalFile: hashOriginalFile,
		HashList:         expectedHashList,
		Size:             size,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	hashList, err := repo.GetChunkHashListByHashOriginalFile(hashOriginalFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(hashList) != len(expectedHashList) {
		t.Fatalf("Expected %v, got %v", len(expectedHashList), len(hashList))
	}

	if hashList[0] != expectedHashList[0] {
		t.Fatalf("Expected %v, got %v", expectedHashList[0], hashList[0])
	}
}

func TestRepositoryChunkJsonGetChunkHashListByHashOriginalFileNotFound(t *testing.T) {
	repo := setUpRepositoryChunkJsonTest(t)
	defer setDownRepositoryChunkJsonTest(t)

	_, err := repo.GetChunkHashListByHashOriginalFile("ASCDEF123456789")
	expected := ErrorRecordNotFound
	if err.Error() != expected.Error() {
		t.Fatalf("Expected error %v, got %v", expected, err)
	}
}

func TestRepositoryChunkJsonCountUsedChunkHashZero(t *testing.T) {
	repo := setUpRepositoryChunkJsonTest(t)
	defer setDownRepositoryChunkJsonTest(t)

	count, err := repo.CountUsedChunkHash("ABC123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 0 {
		t.Fatalf("Expected 0, got %v", count)
	}
}

func TestRepositoryChunkJsonCountUsedChunkHashOne(t *testing.T) {
	repo := setUpRepositoryChunkJsonTest(t)
	defer setDownRepositoryChunkJsonTest(t)

	expectedHash := "ABC123"
	HashList := []string{expectedHash, "987GHI654JKL321"}
	_, err := repo.Create(entity.Chunk{
		HashOriginalFile: "ABCDEF123456789",
		HashList:         HashList,
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	count, err := repo.CountUsedChunkHash(expectedHash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 1 {
		t.Fatalf("Expected 1 got %v", count)
	}
}

func TestRepositoryChunkJsonCountUsedChunkHashTwo(t *testing.T) {
	repo := setUpRepositoryChunkJsonTest(t)
	defer setDownRepositoryChunkJsonTest(t)

	expectedHash := "ABC123"
	HashList := []string{expectedHash, "987GHI654JKL321"}
	_, err := repo.Create(entity.Chunk{
		HashOriginalFile: "ABCDEF123456789",
		HashList:         HashList,
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	HashList = []string{expectedHash, "ASDF1231"}
	_, err = repo.Create(entity.Chunk{
		HashOriginalFile: "ABCDEF123450021",
		HashList:         HashList,
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	count, err := repo.CountUsedChunkHash(expectedHash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 2 {
		t.Fatalf("Expected 2, got %v", count)
	}
}

func TestRepositoryChunkJsonRemoveAll(t *testing.T) {
	repo := setUpRepositoryChunkJsonTest(t)
	defer setDownRepositoryChunkJsonTest(t)

	_, err := repo.Create(entity.Chunk{
		HashOriginalFile: "ABCDEF123456789",
		HashList:         []string{"123ABC456DEF789", "987GHI654JKL321"},
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create(entity.Chunk{
		HashOriginalFile: "ABCDEF123456700",
		HashList:         []string{"123ABC456DEF789", "987GHI654JKL000"},
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repo.RemoveAll()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	file_content, err := os.Open(os.Getenv("COLLECTION_CHUNK_JSON"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer file_content.Close()

	bytes, err := io.ReadAll(file_content)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	chunks := []entity.Chunk{}
	err = json.Unmarshal(bytes, &chunks)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(chunks) != 0 {
		t.Fatalf("Expected 0 chunk, got %v", len(chunks))
	}
}

func TestRepositoryChunkJsonRemoveByHashOriginalFileSimple(t *testing.T) {
	repo := setUpRepositoryChunkJsonTest(t)
	defer setDownRepositoryChunkJsonTest(t)

	hashOriginalFile := "ABCDEF123456789"
	_, err := repo.Create(entity.Chunk{
		HashOriginalFile: hashOriginalFile,
		HashList:         []string{"123ABC456DEF789", "987GHI654JKL321"},
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	removedHashList, err := repo.RemoveByHashOriginalFile(hashOriginalFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(removedHashList) != 2 {
		t.Fatalf("Expected %v, got %v", 2, len(removedHashList))
	}
}

func TestRepositoryChunkJsonRemoveByHashOriginalFileComplex(t *testing.T) {
	repo := setUpRepositoryChunkJsonTest(t)
	defer setDownRepositoryChunkJsonTest(t)

	hashShare := "123ABC456DEF789"
	hashOriginalFile := "ABC123"
	expectedRemovedHash := "987GHI654JKL320"
	_, err := repo.Create(entity.Chunk{
		HashOriginalFile: hashOriginalFile,
		HashList:         []string{hashShare, expectedRemovedHash},
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.Create(entity.Chunk{
		HashOriginalFile: "ABCDEF45646",
		HashList:         []string{hashShare, "987GHI654JKL300"},
		Size:             1024,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	removedHashList, err := repo.RemoveByHashOriginalFile(hashOriginalFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(removedHashList) != 1 {
		t.Fatalf("Expected %v, got %v", 1, len(removedHashList))
	}

	if removedHashList[0] != expectedRemovedHash {
		t.Fatalf("Expected %v, got %v", expectedRemovedHash, removedHashList[0])
	}
}
