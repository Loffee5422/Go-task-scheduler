package task

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

var (
	ErrNotFound          = errors.New("task not found")
	ErrConflict          = errors.New("task with this ID already exists")
	ErrBadRequest        = errors.New("bad task request")
	ErrInvalidTransition = errors.New("invalid state transition")
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
	Priority    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DueDate     *time.Time
	Tags        []string
}

func New(id int, title string) (Task, error) {
	now := time.Now().UTC()
	if id < 0 {
		return Task{}, fmt.Errorf("%w: task ID cannot be negative", ErrBadRequest)
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, fmt.Errorf("%w: title cannot be empty", ErrBadRequest)
	}

	return Task{
		ID:        id,
		Title:     title,
		Status:    "Pending",
		Priority:  "Medium",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func Find(tasks []Task, id int) (*Task, error) {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i], nil
		}
	}
	return nil, ErrNotFound
}

func Create(tasks *[]Task, id int, title string) error {
	if tasks == nil {
		return fmt.Errorf("%w: task collection cannot be nil", ErrBadRequest)
	}
	for _, existing := range *tasks {
		if existing.ID == id {
			return ErrConflict
		}
	}
	newTask, err := New(id, title)
	if err != nil {
		return err
	}
	*tasks = append(*tasks, newTask)
	return nil
}

func Update(tasks *[]Task, id int, title string) error {
	if tasks == nil {
		return fmt.Errorf("%w: task collection cannot be nil", ErrBadRequest)
	}
	if id < 0 {
		return fmt.Errorf("%w: task ID cannot be negative", ErrBadRequest)
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return fmt.Errorf("%w: title cannot be empty", ErrBadRequest)
	}

	current, err := Find(*tasks, id)
	if err != nil {
		return err
	}
	current.Title = title
	current.UpdatedAt = time.Now().UTC()
	return nil
}

func Delete(tasks *[]Task, id int) error {
	if tasks == nil {
		return fmt.Errorf("%w: task collection cannot be nil", ErrBadRequest)
	}
	if id < 0 {
		return fmt.Errorf("%w: task ID cannot be negative", ErrBadRequest)
	}
	for i := range *tasks {
		if (*tasks)[i].ID == id {
			*tasks = append((*tasks)[:i], (*tasks)[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func Complete(tasks *[]Task, id int) error {
	return updateStatus(tasks, id, "Completed")
}

func Reopen(tasks *[]Task, id int) error {
	return updateStatus(tasks, id, "Pending")
}

func FilterByStatus(tasks []Task, status string) []Task {
	status = strings.TrimSpace(status)
	filtered := make([]Task, 0)
	for _, current := range tasks {
		if current.Status == status {
			filtered = append(filtered, current)
		}
	}
	return filtered
}

func SortByPriority(tasks []Task) {
	priorityRank := map[string]int{"High": 0, "Medium": 1, "Low": 2}
	sort.SliceStable(tasks, func(i, j int) bool {
		left, leftKnown := priorityRank[tasks[i].Priority]
		right, rightKnown := priorityRank[tasks[j].Priority]
		if !leftKnown {
			left = len(priorityRank)
		}
		if !rightKnown {
			right = len(priorityRank)
		}
		return left < right
	})
}

func ValidateStatusTransition(tasks []Task, id int, newStatus string) error {
	current, err := Find(tasks, id)
	if err != nil {
		return err
	}
	if !validTransition(current.Status, strings.TrimSpace(newStatus)) {
		return fmt.Errorf("%w: %s", ErrInvalidTransition, ErrInvalidTransition)
	}
	return nil
}

func updateStatus(tasks *[]Task, id int, newStatus string) error {
	if tasks == nil {
		return fmt.Errorf("%w: task collection cannot be nil", ErrBadRequest)
	}
	if id < 0 {
		return fmt.Errorf("%w: task ID cannot be negative", ErrBadRequest)
	}
	current, err := Find(*tasks, id)
	if err != nil {
		return err
	}
	if !validTransition(current.Status, newStatus) {
		return fmt.Errorf("%w: %s", ErrInvalidTransition, ErrInvalidTransition)
	}
	current.Status = newStatus
	current.UpdatedAt = time.Now().UTC()
	return nil
}

func validTransition(current, next string) bool {
	switch current {
	case "Pending":
		return next == "In Progress" || next == "Completed"
	case "In Progress":
		return next == "Completed"
	case "Completed":
		return next == "Pending"
	default:
		return false
	}
}

// These aliases keep the package's existing internal tests focused on implementation details.
func newTask(id int, title string) (Task, error) {
	return New(id, title)
}

func findTask(tasks []Task, id int) (*Task, error) {
	return Find(tasks, id)
}

func updateTask(tasks *[]Task, id int, title string) error {
	return Update(tasks, id, title)
}

func deleteTask(tasks *[]Task, id int) error {
	return Delete(tasks, id)
}

func listTasks(tasks []Task) {
	for _, current := range tasks {
		fmt.Println(current.ID, current.Title, current.Status)
	}
}
