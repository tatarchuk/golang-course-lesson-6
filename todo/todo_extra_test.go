// Extra tests for Task 1. The spec file todo_test.go must not be edited,
// so I put my own additional checks in this separate file. They check the
// things the spec does not: that the wrapped errors can still be inspected
// with errors.Is / errors.As (as the TODO comments in todo.go describe),
// that CreatedAt survives the round trip, and that a nil list is saved as
// an empty JSON array. They also cover the error branches, so coverage
// stays well above the required 80%.
package todo

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadTodos_MissingFile_WrapsErrNotExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")

	_, err := LoadTodos(path)
	if err == nil {
		t.Fatal("LoadTodos() error = nil, want an error for a missing file")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("errors.Is(err, os.ErrNotExist) = false, want true; err = %v", err)
	}
}

func TestLoadTodos_MalformedJSON_WrapsSyntaxError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte("{not valid json"), 0o644); err != nil {
		t.Fatalf("test setup error: %v", err)
	}

	_, err := LoadTodos(path)
	if err == nil {
		t.Fatal("LoadTodos() error = nil, want an error for malformed JSON")
	}
	var syntaxErr *json.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Errorf("errors.As(err, *json.SyntaxError) = false, want true; err = %v", err)
	}
}

func TestLoadTodos_ValidJSONButNotAList(t *testing.T) {
	// Valid JSON, but an object instead of an array of todos.
	path := filepath.Join(t.TempDir(), "object.json")
	if err := os.WriteFile(path, []byte(`{"id": 1}`), 0o644); err != nil {
		t.Fatalf("test setup error: %v", err)
	}

	if _, err := LoadTodos(path); err == nil {
		t.Error("LoadTodos() error = nil, want an error when the JSON is not a list")
	}
}

func TestSaveTodos_DirectoryDoesNotExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no-such-dir", "todos.json")

	err := SaveTodos(path, []Todo{{ID: 1, Title: "x"}})
	if err == nil {
		t.Fatal("SaveTodos() error = nil, want an error when the directory does not exist")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("errors.Is(err, os.ErrNotExist) = false, want true; err = %v", err)
	}
}

func TestSaveTodos_NilListIsSavedAsEmptyArray(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nil.json")

	if err := SaveTodos(path, nil); err != nil {
		t.Fatalf("SaveTodos() error = %v, want nil", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading saved file: %v", err)
	}
	if string(data) != "[]" {
		t.Errorf("file content = %q, want %q (a nil slice should be saved as an empty JSON array, not null)", string(data), "[]")
	}
}

func TestSaveAndLoadTodos_KeepsCreatedAt(t *testing.T) {
	// The spec tests only compare ID, Title and Done, so I check the
	// timestamp separately.
	path := filepath.Join(t.TempDir(), "time.json")
	created := time.Date(2026, 1, 10, 9, 30, 0, 0, time.UTC)

	if err := SaveTodos(path, []Todo{{ID: 1, Title: "x", CreatedAt: created}}); err != nil {
		t.Fatalf("SaveTodos() error = %v, want nil", err)
	}

	got, err := LoadTodos(path)
	if err != nil {
		t.Fatalf("LoadTodos() error = %v, want nil", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d todos, want 1", len(got))
	}
	if !got[0].CreatedAt.Equal(created) {
		t.Errorf("CreatedAt = %v, want %v", got[0].CreatedAt, created)
	}
}
