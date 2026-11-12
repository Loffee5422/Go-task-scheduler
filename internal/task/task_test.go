package task

import (
	"testing"
)

// ==========================================
// Test: Create Task
// ==========================================

func TestNewTask(t *testing.T) {
	task, err := newTask(1, "Learn Go")

	if err != nil {
		t.Fatalf("newTask failed: %v", err)
	}

	// Verify default properties
	if task.ID != 1 {
		t.Errorf("Expected ID 1, got %d", task.ID)
	}

	if task.Title != "Learn Go" {
		t.Errorf("Expected title 'Learn Go', got '%s'", task.Title)
	}

	if task.Status != "Pending" {
		t.Errorf("Expected status 'Pending', got '%s'", task.Status)
	}

	if task.Priority != "Medium" {
		t.Errorf("Expected priority 'Medium', got '%s'", task.Priority)
	}

	if task.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be initialized")
	}

	if task.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be initialized")
	}

	if task.DueDate != nil {
		t.Errorf("Expected DueDate to be nil, got %v", task.DueDate)
	}

	if len(task.Tags) != 0 {
		t.Errorf("Expected empty tags, got %v", task.Tags)
	}
}

// ==========================================
// Test: Find Task
// ==========================================

func TestFindTask(t *testing.T) {
	task1, _ := newTask(1, "Learn Go")
	task2, _ := newTask(2, "Learn Testing")

	tasks := []Task{task1, task2}

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "find existing first task",
			id:      1,
			wantErr: false,
		},
		{
			name:    "find existing second task",
			id:      2,
			wantErr: false,
		},
		{
			name:    "find nonexistent task",
			id:      999,
			wantErr: true,
		},
		{
			name:    "find task in empty collection",
			id:      1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentTasks := tasks

			if tt.name == "find task in empty collection" {
				currentTasks = []Task{}
			}

			found, err := findTask(currentTasks, tt.id)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Expected error: %v, got: %v",
					tt.wantErr, err)
			}

			if err == nil && found.ID != tt.id {
				t.Errorf("Expected ID %d, got %d",
					tt.id, found.ID)
			}
		})
	}
}

// ==========================================
// Test: Update Task
// ==========================================

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		newTitle  string
		wantTitle string
		wantErr   bool
	}{
		{
			name:      "update existing task",
			id:        1,
			newTitle:  "Learn Go Testing",
			wantTitle: "Learn Go Testing",
			wantErr:   false,
		},
		{
			name:      "update nonexistent task",
			id:        999,
			newTitle:  "Invalid Update",
			wantTitle: "Learn Go",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := newTask(1, "Learn Go")
			if err != nil {
				t.Fatalf("newTask failed: %v", err)
			}

			tasks := []Task{task}

			err = updateTask(&tasks, tt.id, tt.newTitle)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Expected error: %v, got: %v",
					tt.wantErr, err)
			}

			if tasks[0].Title != tt.wantTitle {
				t.Errorf("Expected title '%s', got '%s'",
					tt.wantTitle, tasks[0].Title)
			}
		})
	}
}

// ==========================================
// Test: Delete Task
// ==========================================

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name       string
		initialIDs []int
		deleteID   int
		wantIDs    []int
		wantErr    bool
	}{
		{
			name:       "delete only task",
			initialIDs: []int{1},
			deleteID:   1,
			wantIDs:    []int{},
			wantErr:    false,
		},
		{
			name:       "delete middle task",
			initialIDs: []int{1, 2, 3},
			deleteID:   2,
			wantIDs:    []int{1, 3},
			wantErr:    false,
		},
		{
			name:       "delete first task",
			initialIDs: []int{1, 2, 3},
			deleteID:   1,
			wantIDs:    []int{2, 3},
			wantErr:    false,
		},
		{
			name:       "delete last task",
			initialIDs: []int{1, 2, 3},
			deleteID:   3,
			wantIDs:    []int{1, 2},
			wantErr:    false,
		},
		{
			name:       "delete nonexistent task",
			initialIDs: []int{1, 2, 3},
			deleteID:   999,
			wantIDs:    []int{1, 2, 3},
			wantErr:    true,
		},
		{
			name:       "delete from empty collection",
			initialIDs: []int{},
			deleteID:   1,
			wantIDs:    []int{},
			wantErr:    true,
		},
		{
			name:       "delete ID from different collection",
			initialIDs: []int{10, 20, 30},
			deleteID:   2,
			wantIDs:    []int{10, 20, 30},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			var tasks []Task

			for _, id := range tt.initialIDs {
				task, err := newTask(id, "Test Task")
				if err != nil {
					t.Fatalf("newTask failed: %v", err)
				}

				tasks = append(tasks, task)
			}

			// Act
			err := deleteTask(&tasks, tt.deleteID)

			// Assert: error behaviour
			if (err != nil) != tt.wantErr {
				t.Fatalf("Expected error: %v, got: %v",
					tt.wantErr, err)
			}

			// Assert: resulting collection
			if len(tasks) != len(tt.wantIDs) {
				t.Fatalf("Expected %d tasks, got %d",
					len(tt.wantIDs), len(tasks))
			}

			// Assert: correct elements and ordering
			for i, expectedID := range tt.wantIDs {
				if tasks[i].ID != expectedID {
					t.Errorf("Position %d: expected ID %d, got %d",
						i, expectedID, tasks[i].ID)
				}
			}
		})
	}
}

// ==========================================
// Test: List Tasks
// ==========================================

func TestListTasks(t *testing.T) {
	task1, err := newTask(1, "Learn Go")
	if err != nil {
		t.Fatalf("newTask failed: %v", err)
	}

	task2, err := newTask(2, "Learn Testing")
	if err != nil {
		t.Fatalf("newTask failed: %v", err)
	}

	tests := []struct {
		name  string
		tasks []Task
	}{
		{
			name:  "list multiple tasks",
			tasks: []Task{task1, task2},
		},
		{
			name:  "list empty collection",
			tasks: []Task{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Currently verifies that listing executes.
			// Does not yet verify printed output.
			listTasks(tt.tasks)
		})
	}
}
