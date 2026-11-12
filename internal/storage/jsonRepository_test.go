package storage

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)



func TestLoadTasksFromFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "tasks.json")

	_, err := LoadTasksFromFile(filename)
	if err == nil {
		t.Fatal("Expected error when loading nonexistent file")
	}
}

func TestFileExists(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "tasks.json")

	if FileExists(filename) {
		t.Fatal("Expected FileExists to return false")
	}

	if err := os.WriteFile(filename, []byte("[]"), 0644); err != nil {
		t.Fatal(err)
	}

	if !FileExists(filename) {
		t.Fatal("Expected FileExists to return true")
	}
}

func TestSaveTasksToFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "tasks.json")

	tasks := []Task{
		{Title: "Task 1", Description: "Description 1"},
		{Title: "Task 2", Description: "Description 2"},
	}

	if err := SaveTasksToFile(filename, tasks); err != nil {
		t.Fatalf("SaveTasksToFile failed: %v", err)
	}

	if !FileExists(filename) {
		t.Fatal("Expected file to be created")
	}
}

func TestMalformedJSON(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "malformed.json")

	malformedJSON := `[{"name": "Task 1"`

	if err := os.WriteFile(filename, []byte(malformedJSON), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadTasksFromFile(filename)
	if err == nil {
		t.Fatal("Expected error when loading malformed JSON")
	}
}

func TestSaveAndLoadTasks(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "tasks.json")

	tasks := []Task{
		{Title: "Task 1", Description: "Description 1"},
		{Title: "Task 2", Description: "Description 2"},
	}

	if err := SaveTasksToFile(filename, tasks); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loadedTasks, err := LoadTasksFromFile(filename)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if !reflect.DeepEqual(tasks, loadedTasks) {
		t.Errorf("Data mismatch:\nExpected: %#v\nGot: %#v", tasks, loadedTasks)
	}
}

func TestSaveEmptyTasks(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "tasks.json")

	tasks := []Task{}

	if err := SaveTasksToFile(filename, tasks); err != nil {
		t.Fatalf("Save empty tasks failed: %v", err)
	}

	loadedTasks, err := LoadTasksFromFile(filename)
	if err != nil {
		t.Fatalf("Load empty tasks failed: %v", err)
	}

	if len(loadedTasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(loadedTasks))
	}
}

func TestCreateFileIfNotExists(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "tasks.json")

	if err := CreateFileIfNotExists(filename); err != nil {
		t.Fatalf("CreateFileIfNotExists failed: %v", err)
	}

	if !FileExists(filename) {
		t.Fatal("Expected file to exist after creation")
	}

	// Calling it again should not return an error.
	if err := CreateFileIfNotExists(filename); err != nil {
		t.Fatalf("Second creation failed: %v", err)
	}
}
