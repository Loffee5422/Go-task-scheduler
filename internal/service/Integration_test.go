package service

import (
	"os"
	"path/filepath"
	"reflect"
	"task_manager/internal/storage"
	"testing"
)

// newTestRepository creates an isolated JSON repository.
// Each test gets its own temporary directory.
func newTestRepository(t *testing.T) *storage.TaskRepository {
	t.Helper()

	filename := filepath.Join(t.TempDir(), "tasks.json")
	return storage.NewTaskRepository(filename)
}

// reloadManager simulates restarting the application.
// A new TaskManager reads its state from the same repository.
func reloadManager(
	t *testing.T,
	repo *storage.TaskRepository,
) *TaskManager {
	t.Helper()

	tm, err := NewPersistentTaskManager(repo)
	if err != nil {
		t.Fatalf("Failed to reload TaskManager: %v", err)
	}
	return tm
}

// TestPersistentInitialization verifies that a missing JSON file
// is treated as an empty task collection.
func TestPersistentInitialization(t *testing.T) {
	repo := newTestRepository(t)

	tm, err := NewPersistentTaskManager(repo)
	if err != nil {
		t.Fatalf("Initialization failed: %v", err)
	}

	if tm == nil {
		t.Fatal("Expected non-nil TaskManager")
	}

	if len(tm.ListTasks()) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tm.ListTasks()))
	}
}

// TestPersistentSaveAndReload verifies that task data survives
// creating a completely new TaskManager instance.
func TestPersistentSaveAndReload(t *testing.T) {
	repo := newTestRepository(t)
	tm := reloadManager(t, repo)

	if err := tm.CreateTask(1, "Learn Go"); err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	if err := tm.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Simulate application restart.
	reloaded := reloadManager(t, repo)

	original := tm.ListTasks()
	restored := reloaded.ListTasks()

	if !reflect.DeepEqual(original, restored) {
		t.Errorf(
			"Tasks changed after reload:\nOriginal: %+v\nRestored: %+v",
			original,
			restored,
		)
	}

	if len(restored) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(restored))
	}

	if restored[0].ID != 1 || restored[0].Title != "Learn Go" {
		t.Errorf("Incorrect restored task: %+v", restored[0])
	}
}

// TestPersistentMultipleTasks verifies that multiple tasks
// are saved and loaded without losing their order or data.
func TestPersistentMultipleTasks(t *testing.T) {
	repo := newTestRepository(t)
	tm := reloadManager(t, repo)

	tasks := []struct {
		id    int
		title string
	}{
		{1, "Learn Go"},
		{2, "Write Tests"},
		{3, "Build CLI"},
	}

	for _, tc := range tasks {
		if err := tm.CreateTask(tc.id, tc.title); err != nil {
			t.Fatalf("CreateTask(%d) failed: %v", tc.id, err)
		}
	}

	if err := tm.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	reloaded := reloadManager(t, repo)
	got := reloaded.ListTasks()

	if len(got) != len(tasks) {
		t.Fatalf("Expected %d tasks, got %d", len(tasks), len(got))
	}

	for i, expected := range tasks {
		if got[i].ID != expected.id || got[i].Title != expected.title {
			t.Errorf("Task %d mismatch: %+v", i, got[i])
		}
	}
}

// TestPersistentUpdateAndReload verifies that task modifications
// are persisted, including UpdatedAt.
func TestPersistentUpdateAndReload(t *testing.T) {
	repo := newTestRepository(t)
	tm := reloadManager(t, repo)

	if err := tm.CreateTask(1, "Learn Go"); err != nil {
		t.Fatal(err)
	}

	if err := tm.Save(); err != nil {
		t.Fatal(err)
	}

	if err := tm.UpdateTask(1, "Learn Rust"); err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}

	updated, err := tm.FindTask(1)
	if err != nil {
		t.Fatal(err)
	}

	if err := tm.Save(); err != nil {
		t.Fatal(err)
	}

	reloaded := reloadManager(t, repo)

	got, err := reloaded.FindTask(1)
	if err != nil {
		t.Fatalf("Updated task not found after reload: %v", err)
	}

	if got.Title != "Learn Rust" {
		t.Errorf("Expected Learn Rust, got %q", got.Title)
	}

	if !got.UpdatedAt.Equal(updated.UpdatedAt) {
		t.Errorf(
			"UpdatedAt changed after reload: expected %v, got %v",
			updated.UpdatedAt,
			got.UpdatedAt,
		)
	}
}

// TestPersistentDeleteAndReload verifies that deleted tasks
// do not reappear after restarting the application.
func TestPersistentDeleteAndReload(t *testing.T) {
	repo := newTestRepository(t)
	tm := reloadManager(t, repo)

	for _, tc := range []struct {
		id    int
		title string
	}{
		{1, "Task One"},
		{2, "Task Two"},
	} {
		if err := tm.CreateTask(tc.id, tc.title); err != nil {
			t.Fatal(err)
		}
	}

	if err := tm.Save(); err != nil {
		t.Fatal(err)
	}

	if err := tm.DeleteTask(1); err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}

	if err := tm.Save(); err != nil {
		t.Fatal(err)
	}

	reloaded := reloadManager(t, repo)

	if _, err := reloaded.FindTask(1); err == nil {
		t.Error("Deleted task reappeared after reload")
	}

	if _, err := reloaded.FindTask(2); err != nil {
		t.Errorf("Remaining task disappeared: %v", err)
	}

	if len(reloaded.ListTasks()) != 1 {
		t.Errorf("Expected 1 remaining task")
	}
}

// TestPersistentCompleteAndReopen verifies that status
// transitions survive saving and reloading.
func TestPersistentCompleteAndReopen(t *testing.T) {
	repo := newTestRepository(t)
	tm := reloadManager(t, repo)

	if err := tm.CreateTask(1, "Learn Go"); err != nil {
		t.Fatal(err)
	}

	if err := tm.CompleteTask(1); err != nil {
		t.Fatalf("CompleteTask failed: %v", err)
	}

	if err := tm.Save(); err != nil {
		t.Fatal(err)
	}

	reloaded := reloadManager(t, repo)

	completed, err := reloaded.FindTask(1)
	if err != nil {
		t.Fatal(err)
	}

	if completed.Status != "Completed" {
		t.Errorf("Expected Completed, got %q", completed.Status)
	}

	if err := reloaded.ReopenTask(1); err != nil {
		t.Fatalf("ReopenTask failed: %v", err)
	}

	if err := reloaded.Save(); err != nil {
		t.Fatal(err)
	}

	reloadedAgain := reloadManager(t, repo)

	reopened, err := reloadedAgain.FindTask(1)
	if err != nil {
		t.Fatal(err)
	}

	if reopened.Status != "Pending" {
		t.Errorf("Expected Pending, got %q", reopened.Status)
	}
}

// TestPersistentDuplicateID verifies that duplicate ID
// validation still works after loading saved tasks.
func TestPersistentDuplicateID(t *testing.T) {
	repo := newTestRepository(t)
	tm := reloadManager(t, repo)

	if err := tm.CreateTask(1, "Original"); err != nil {
		t.Fatal(err)
	}

	if err := tm.Save(); err != nil {
		t.Fatal(err)
	}

	reloaded := reloadManager(t, repo)

	err := reloaded.CreateTask(1, "Duplicate")
	if err == nil {
		t.Error("Expected duplicate ID error after reload")
	}

	if len(reloaded.ListTasks()) != 1 {
		t.Error("Duplicate creation changed task count")
	}
}

// TestPersistentEmptyList verifies that deleting all tasks
// and saving does not restore previously deleted tasks.
func TestPersistentEmptyList(t *testing.T) {
	repo := newTestRepository(t)
	tm := reloadManager(t, repo)

	if err := tm.CreateTask(1, "Temporary"); err != nil {
		t.Fatal(err)
	}

	if err := tm.Save(); err != nil {
		t.Fatal(err)
	}

	if err := tm.DeleteTask(1); err != nil {
		t.Fatal(err)
	}

	if err := tm.Save(); err != nil {
		t.Fatal(err)
	}

	reloaded := reloadManager(t, repo)

	if len(reloaded.ListTasks()) != 0 {
		t.Errorf(
			"Expected empty task list, got %+v",
			reloaded.ListTasks(),
		)
	}
}

// TestPersistentMalformedJSON verifies that corrupted data
// produces an error instead of silently resetting tasks.
func TestPersistentMalformedJSON(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "tasks.json")

	invalidJSON := []byte(`[{"ID":1,"Title":"Broken"`)

	if err := os.WriteFile(filename, invalidJSON, 0644); err != nil {
		t.Fatal(err)
	}

	repo := storage.NewTaskRepository(filename)

	tm, err := NewPersistentTaskManager(repo)

	if err == nil {
		t.Error("Expected error for malformed JSON")
	}

	if tm != nil {
		t.Error("Expected nil TaskManager after load failure")
	}
}

// TestPersistentSaveFailure verifies that errors from
// the filesystem are propagated to the caller.
func TestPersistentSaveFailure(t *testing.T) {
	dir := t.TempDir()

	// Parent directory does not exist.
	filename := filepath.Join(dir, "missing", "tasks.json")

	repo := storage.NewTaskRepository(filename)
	tm := reloadManager(t, repo)

	if err := tm.CreateTask(1, "Learn Go"); err != nil {
		t.Fatal(err)
	}

	err := tm.Save()
	if err == nil {
		t.Error("Expected Save to fail for invalid path")
	}

	// Explicit Save failure must not remove the in-memory task.
	if _, err := tm.FindTask(1); err != nil {
		t.Error("Task disappeared from memory after failed Save")
	}

	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		t.Errorf("Expected no output file, got Stat error: %v", err)
	}
}

// TestTaskManagerWithoutRepository verifies that the
// original in-memory TaskManager remains usable.
func TestTaskManagerWithoutRepository(t *testing.T) {
	tm := NewTaskManager()

	if err := tm.CreateTask(1, "Learn Go"); err != nil {
		t.Fatal(err)
	}

	if err := tm.Save(); err == nil {
		t.Error("Expected error when no repository is configured")
	}

	if _, err := tm.FindTask(1); err != nil {
		t.Error("In-memory task should remain available")
	}
}

// TestPersistentUnsavedChanges verifies that mutations
// remain in memory until Save is explicitly called.
func TestPersistentUnsavedChanges(t *testing.T) {
	repo := newTestRepository(t)
	tm := reloadManager(t, repo)

	if err := tm.CreateTask(1, "Saved Task"); err != nil {
		t.Fatal(err)
	}

	if err := tm.Save(); err != nil {
		t.Fatal(err)
	}

	// Modify memory without calling Save.
	if err := tm.CreateTask(2, "Unsaved Task"); err != nil {
		t.Fatal(err)
	}

	reloaded := reloadManager(t, repo)

	if _, err := reloaded.FindTask(1); err != nil {
		t.Error("Previously saved task is missing")
	}

	if _, err := reloaded.FindTask(2); err == nil {
		t.Error("Unsaved task unexpectedly persisted")
	}
}
