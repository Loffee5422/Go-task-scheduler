package service

import (
	"errors"

	"task_manager/internal/storage"
	"task_manager/internal/task"
)

type Task = task.Task

type TaskManager struct {
	tasks []Task
	repo  *storage.TaskRepository
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: []Task{},
	}
}

// NewPersistentTaskManager creates a TaskManager
// and loads previously saved tasks from storage.
func NewPersistentTaskManager(
	repo *storage.TaskRepository,
) (*TaskManager, error) {
	if repo == nil {
		return nil, errors.New("repository cannot be nil")
	}

	// Read tasks from JSON.
	tasks, err := repo.Load()
	if err != nil {
		return nil, err
	}
	if tasks == nil {
		tasks = []Task{}
	}

	// Initialize TaskManager with loaded tasks.
	return &TaskManager{
		tasks: tasks,
		repo:  repo,
	}, nil
}

func (tm *TaskManager) Save() error {
	if tm == nil {
		return errors.New("task manager cannot be nil")
	}
	if tm.repo == nil {
		return errors.New("no repository set for TaskManager")
	}

	return tm.repo.Save(tm.tasks)
}

func (tm *TaskManager) CreateTask(id int, title string) error {
	return task.Create(&tm.tasks, id, title)
}

func (tm *TaskManager) UpdateTask(id int, newTitle string) error {
	return task.Update(&tm.tasks, id, newTitle)
}

func (tm *TaskManager) DeleteTask(id int) error {
	return task.Delete(&tm.tasks, id)
}

func (tm *TaskManager) ListTasks() []Task {
	tasks := make([]Task, len(tm.tasks))
	for i, current := range tm.tasks {
		tasks[i] = current
		if current.DueDate != nil {
			dueDate := *current.DueDate
			tasks[i].DueDate = &dueDate
		}
		if current.Tags != nil {
			tasks[i].Tags = append([]string(nil), current.Tags...)
		}
	}
	return tasks
}

func (tm *TaskManager) FindTask(id int) (Task, error) {
	found, err := task.Find(tm.tasks, id)
	if err != nil {
		return Task{}, err
	}
	return *found, nil
}

func (tm *TaskManager) CompleteTask(id int) error {
	return task.Complete(&tm.tasks, id)
}

func (tm *TaskManager) ReopenTask(id int) error {
	return task.Reopen(&tm.tasks, id)
}

func (tm *TaskManager) FilterTasksByStatus(status string) []Task {
	return task.FilterByStatus(tm.tasks, status)
}

func (tm *TaskManager) SortTasksByPriority() {
	task.SortByPriority(tm.tasks)
}

func (tm *TaskManager) ValidateStateTransition(id int, newStatus string) error {
	return task.ValidateStatusTransition(tm.tasks, id, newStatus)
}
