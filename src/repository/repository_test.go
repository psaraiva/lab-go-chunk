package repository

import (
	"reflect"
	"testing"
)

func TestRepositoryMakeRepositoryFile(t *testing.T) {
	repo := MakeRepositoryFile(ENGINE_JSON)
	if reflect.TypeOf(repo) != reflect.TypeOf(RepositoryFileJson{}) {
		t.Fatalf("Expected RepositoryFileJson, got %v", repo)
	}

	repo = MakeRepositoryFile(ENGINE_SQLITE)
	if reflect.TypeOf(repo) != reflect.TypeOf(RepositoryFileSqlite{}) {
		t.Fatalf("Expected RepositoryFileSqlite, got %v", repo)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Expected panic, got %v", r)
		}
	}()

	MakeRepositoryFile("invalid")
}

func TestRepositoryMakeRepositoryChunk(t *testing.T) {
	repo := MakeRepositoryChunk(ENGINE_JSON)
	if reflect.TypeOf(repo) != reflect.TypeOf(RepositoryChunkJson{}) {
		t.Fatalf("Expected RepositoryChunkJson, got %v", repo)
	}

	repo = MakeRepositoryChunk(ENGINE_SQLITE)
	if reflect.TypeOf(repo) != reflect.TypeOf(RepositoryChunkSqlite{}) {
		t.Fatalf("Expected RepositoryChunkSqlite, got %v", repo)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Expected panic, got %v", r)
		}
	}()

	MakeRepositoryChunk("invalid")
}
