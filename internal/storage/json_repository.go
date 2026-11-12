package storage

import (
	"encoding/json"
	"os"
	"task_manager/internal/task"
)

type Task = task.Task

// TaskRepository stores the filename used for persistence.
// It is a struct, NOT an interface.
type TaskRepository struct {
	filename string
}

// NewTaskRepository creates a repository linked to a JSON file.
// Example: repo := NewTaskRepository("tasks.json")
func NewTaskRepository(filename string) *TaskRepository {
	return &TaskRepository{
		filename: filename,
	}
}

// Load retrieves tasks from the repository's JSON file.
// If the file doesn't exist, return an empty task list.
// Other errors (e.g. malformed JSON) are returned to the caller.
func (r *TaskRepository) Load() ([]Task, error) {
	tasks, err := LoadTasksFromFile(r.filename)
	if os.IsNotExist(err) {
		return []Task{}, nil
	}
	return tasks, err
}

// Save persists the provided tasks to the repository's file.
// The existing file is overwritten with the new task list.
func (r *TaskRepository) Save(tasks []Task) error {
	return SaveTasksToFile(r.filename, tasks)
}

// LoadTasksFromFile reads and decodes tasks from a JSON file.
func LoadTasksFromFile(filename string) ([]Task, error) {
	// Read the file into a byte slice.
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Convert JSON bytes into Go task objects.
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	if tasks == nil {
		tasks = []Task{}
	}

	return tasks, nil
}

// SaveTasksToFile encodes tasks as JSON and writes them to disk.
func SaveTasksToFile(filename string, tasks []Task) error {
	// Convert Go objects into formatted JSON bytes.
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON bytes to the file.
	// 0644 = owner can read/write; others can read.
	return os.WriteFile(filename, data, 0644)
}

// FileExists checks whether a file or directory is accessible.
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// CreateFileIfNotExists initializes a missing JSON file.
// It writes [] so the file contains valid JSON.
func CreateFileIfNotExists(filename string) error {
	// O_EXCL prevents overwriting a file created by another process.
	file, err := os.OpenFile(
		filename,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL,
		0644,
	)

	if os.IsExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	_, writeErr := file.Write([]byte("[]"))
	closeErr := file.Close()

	if writeErr != nil {
		return writeErr
	}
	return closeErr
}
