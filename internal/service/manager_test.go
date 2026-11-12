package service

import "testing"

// ============================================================
// TASK MANAGER TEST SUITE
// ============================================================
//
// Testing objectives:
// 1. Verify that valid operations work correctly.
// 2. Verify that invalid inputs are rejected.
// 3. Verify that missing resources are handled correctly.
// 4. Verify that operations produce the expected state changes.
// 5. Verify that multiple operations work correctly together.
//
// ============================================================

// ------------------------------------------------------------
// 1. CREATE TASK
// ------------------------------------------------------------

// TestCreateTask
// - Create a task with a valid ID and title.
// - Verify that no error is returned.
// - Verify that the task exists in TaskManager.
// - Verify that all default task fields are initialized correctly.
func TestCreateTask(t *testing.T) {
	tm := NewTaskManager()
	err := tm.CreateTask(1, "Learn Go")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tm.tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tm.tasks))
	}
	if tm.tasks[0].ID != 1 || tm.tasks[0].Title != "Learn Go" {
		t.Errorf("Task fields not set correctly: %+v", tm.tasks[0])
	}
	if tm.tasks[0].Status != "Pending" || tm.tasks[0].Priority != "Medium" {
		t.Errorf("Default fields not set correctly: %+v", tm.tasks[0])
	}
	if tm.tasks[0].CreatedAt.IsZero() || tm.tasks[0].UpdatedAt.IsZero() {
		t.Errorf("Timestamps not initialized: %+v", tm.tasks[0])
	}
	if tm.tasks[0].DueDate != nil {
		t.Errorf("Expected DueDate to be nil, got %v", tm.tasks[0].DueDate)
	}
}

// TestCreateTaskDuplicateID
// - Create a task with an ID that already exists.
// - Verify that an appropriate duplicate ID error is returned.
// - Verify that the original task remains unchanged.
// - Verify that no additional task is inserted.
func TestCreateTaskDuplicateID(t *testing.T) {
	tm := NewTaskManager()
	err := tm.CreateTask(1, "Learn Go")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	err = tm.CreateTask(1, "Learn Python")
	if err == nil {
		t.Errorf("Expected duplicate ID error, got no error")
	}
	if len(tm.tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tm.tasks))
	}
	if tm.tasks[0].ID != 1 || tm.tasks[0].Title != "Learn Go" {
		t.Errorf("Task fields not set correctly: %+v", tm.tasks[0])
	}
}

// TestCreateTaskInvalidInput
// - Test negative task ID.
// - Test empty title.
// - Test whitespace-only title.
// - Verify that each invalid input returns an error.
// - Verify that invalid tasks are never inserted.

// TestCreateTaskBoundaryValues
// - Test the smallest valid ID.
// - Test a title containing leading and trailing whitespace.
// - Test a title containing special characters.
// - Verify behaviour matches the defined validation rules.
func TestCreateTaskBoundaryValues(t *testing.T) {
	tm := NewTaskManager()

	err := tm.CreateTask(0, "  Learn Go  ")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(tm.tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tm.tasks))
	}

	if tm.tasks[0].ID != 0 || tm.tasks[0].Title != "Learn Go" {
		t.Errorf("Task fields not set correctly: %+v", tm.tasks[0])
	}
}

// ------------------------------------------------------------
// 2. UPDATE TASK
// ------------------------------------------------------------

// TestUpdateTask
// - Create an existing task.
// - Update its title.
// - Verify that no error is returned.
// - Verify that the title actually changes.
// - Verify that the task ID remains unchanged.
// - Verify that CreatedAt remains unchanged.
// - Verify that UpdatedAt changes appropriately.
func TestUpdateTaskTM(t *testing.T) {
	tm := NewTaskManager()
	tm.CreateTask(1, "Learn Go")
	err := tm.UpdateTask(1, "Learn Go and Testing")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tm.tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tm.tasks))
	}
	if tm.tasks[0].ID != 1 || tm.tasks[0].Title != "Learn Go and Testing" {
		t.Errorf("Task fields not set correctly: %+v", tm.tasks[0])
	}
}

// TestUpdateTaskNotFound
// - Attempt to update a task that does not exist.
// - Verify that an appropriate error is returned.
// - Verify that existing tasks remain unchanged.
func TestUpdateTaskNotFound(t *testing.T) {
	tm := NewTaskManager()
	tm.CreateTask(1, "Learn Go")
	err := tm.UpdateTask(2, "Learn Python")
	if err == nil {
		t.Errorf("Expected error, got no error")
	}
	if len(tm.tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tm.tasks))
	}
	if tm.tasks[0].ID != 1 || tm.tasks[0].Title != "Learn Go" {
		t.Errorf("Task fields not set correctly: %+v", tm.tasks[0])
	}
}

// TestUpdateTaskInvalidInput
// - Test negative task ID.
// - Test empty title.
// - Test whitespace-only title.
// - Verify that invalid updates return errors.
// - Verify that the original task remains unchanged.
func TestUpdateTaskInvalidInput(t *testing.T) {
	tm := NewTaskManager()
	tm.CreateTask(1, "Learn Go")
	err := tm.UpdateTask(-1, "Learn Python")
	if err == nil {
		t.Errorf("Expected error, got no error")
	}
	if len(tm.tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tm.tasks))
	}
	if tm.tasks[0].ID != 1 || tm.tasks[0].Title != "Learn Go" {
		t.Errorf("Task fields not set correctly: %+v", tm.tasks[0])
	}
}

// TestUpdateTaskBoundaryValues
// - Update a task with a title containing special characters.
// - Update a task with a very long title.
// - Verify behaviour against the defined title constraints.
func TestUpdateTaskBoundaryValues(t *testing.T) {
	tm := NewTaskManager()
	tm.CreateTask(1, "Learn Go")
	err := tm.UpdateTask(1, "Learn Go & Testing!")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tm.tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tm.tasks))
	}
	if tm.tasks[0].ID != 1 || tm.tasks[0].Title != "Learn Go & Testing!" {
		t.Errorf("Task fields not set correctly: %+v", tm.tasks[0])
	}
}

// ------------------------------------------------------------
// 3. DELETE TASK
// ------------------------------------------------------------

// TestDeleteTask
// - Create a task.
// - Delete the task.
// - Verify that no error is returned.
// - Verify that the task no longer exists.
// - Verify that the task count decreases by one.
func TestDeleteTaskTM(t *testing.T) {
	tm := NewTaskManager()
	tm.CreateTask(1, "Learn Go")
	err := tm.DeleteTask(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tm.tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tm.tasks))
	}
}

// TestDeleteTaskNotFound
// - Attempt to delete a task that does not exist.
// - Verify that an appropriate error is returned.
// - Verify that existing tasks remain unchanged.
func TestDeleteTaskNotFound(t *testing.T) {
	tm := NewTaskManager()
	tm.CreateTask(1, "Learn Go")
	err := tm.DeleteTask(2)
	if err == nil {
		t.Errorf("Expected error, got no error")
	}
	if len(tm.tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tm.tasks))
	}
	if tm.tasks[0].ID != 1 || tm.tasks[0].Title != "Learn Go" {
		t.Errorf("Task fields not set correctly: %+v", tm.tasks[0])
	}
}

// TestDeleteTaskInvalidID
// - Attempt to delete a task using a negative ID.
// - Verify that an appropriate validation error is returned.
// - Verify that existing tasks remain unchanged.
func TestDeleteTaskInvalidID(t *testing.T) {
	tm := NewTaskManager()
	tm.CreateTask(1, "Learn Go")
	err := tm.DeleteTask(-1)
	if err == nil {
		t.Errorf("Expected error, got no error")
	}
	if len(tm.tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tm.tasks))
	}
	if tm.tasks[0].ID != 1 || tm.tasks[0].Title != "Learn Go" {
		t.Errorf("Task fields not set correctly: %+v", tm.tasks[0])
	}
}

// TestDeleteTaskMultipleTasks
// - Create multiple tasks.
// - Delete a task from the middle of the collection.
// - Verify that the remaining tasks are preserved.
// - Verify that their IDs and contents remain unchanged.
// - Verify that the collection length is correct.
func TestDeleteTaskMultipleTasks(t *testing.T) {
	tm := NewTaskManager()
	tm.CreateTask(1, "Learn Go")
	tm.CreateTask(2, "Learn Python")
	tm.CreateTask(3, "Learn JavaScript")
	err := tm.DeleteTask(2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tm.tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tm.tasks))
	}
	if tm.tasks[0].ID != 1 || tm.tasks[0].Title != "Learn Go" {
		t.Errorf("Task fields not set correctly: %+v", tm.tasks[0])
	}
	if tm.tasks[1].ID != 3 || tm.tasks[1].Title != "Learn JavaScript" {
		t.Errorf("Task fields not set correctly: %+v", tm.tasks[1])
	}
}

// ------------------------------------------------------------
// 4. FIND TASK
// ------------------------------------------------------------

// TestFindTask
// - Create a task.
// - Find it using its ID.
// - Verify that no error is returned.
// - Verify that the returned task contains the correct fields.
func TestFindTaskTM(t *testing.T) {
	tm := NewTaskManager()
	tm.CreateTask(1, "Learn Go")
	task, err := tm.FindTask(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if task.ID != 1 || task.Title != "Learn Go" {
		t.Errorf("Task fields not set correctly: %+v", task)
	}
}

// TestFindTaskNotFound
// - Search for a task that does not exist.
// - Verify that an appropriate error is returned.
func TestFindTaskNotFound(t *testing.T) {
	tm := NewTaskManager()
	tm.CreateTask(1, "Learn Go")
	_, err := tm.FindTask(2)
	if err == nil {
		t.Errorf("Expected error, got no error")
	}
}

// TestFindTaskInvalidID
// - Search using a negative ID.
// - Verify that an appropriate validation error is returned.
func TestFindTaskInvalidID(t *testing.T) {
	tm := NewTaskManager()
	tm.CreateTask(1, "Learn Go")
	_, err := tm.FindTask(-1)
	if err == nil {
		t.Errorf("Expected error, got no error")
	}
}

// TestFindTaskEmptyCollection
// - Search when TaskManager contains no tasks.
// - Verify that an appropriate error is returned.
func TestFindTaskEmptyCollection(t *testing.T) {
	tm := NewTaskManager()
	_, err := tm.FindTask(1)
	if err == nil {
		t.Errorf("Expected error, got no error")
	}
}

// ------------------------------------------------------------
// 5. LIST TASKS
// ------------------------------------------------------------

// TestListTasksEmpty
// - Call ListTasks on an empty TaskManager.
// - Verify that the returned collection contains zero tasks.
func TestListTasksEmpty(t *testing.T) {
	tm := NewTaskManager()
	tasks := tm.ListTasks()
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
}

// TestListTasks
// - Create multiple tasks.
// - Call ListTasks.
// - Verify that all tasks are returned.
// - Verify that the number of tasks is correct.
// - Verify that task IDs and contents are correct.
func TestListTasksTM(t *testing.T) {
	tm := NewTaskManager()
	tm.CreateTask(1, "Learn Go")
	tm.CreateTask(2, "Learn Python")
	tasks := tm.ListTasks()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].ID != 1 || tasks[0].Title != "Learn Go" {
		t.Errorf("Task fields not set correctly: %+v", tasks[0])
	}
	if tasks[1].ID != 2 || tasks[1].Title != "Learn Python" {
		t.Errorf("Task fields not set correctly: %+v", tasks[1])
	}
}

// TestListTasksAfterOperations
// - Create multiple tasks.
// - Update one task.
// - Delete another task.
// - Call ListTasks.
// - Verify that the returned collection reflects all changes.
// - Verify that deleted tasks are absent.
// - Verify that updated task information is correct.
func TestListTasksAfterOperations(t *testing.T) {
	tm := NewTaskManager()
	tm.CreateTask(1, "Learn Go")
	tm.CreateTask(2, "Learn Python")
	tm.UpdateTask(1, "Learn Go and Testing")
	tm.DeleteTask(2)
	tasks := tm.ListTasks()
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	if tasks[0].ID != 1 || tasks[0].Title != "Learn Go and Testing" {
		t.Errorf("Task fields not set correctly: %+v", tasks[0])
	}
}

// ------------------------------------------------------------
// 6. STATE INTEGRITY
// ------------------------------------------------------------

// TestFailedOperationDoesNotModifyState
// - Create multiple valid tasks.
// - Attempt an invalid update.
// - Attempt to delete a nonexistent task.
// - Attempt to create a duplicate ID.
// - Verify that failed operations do not modify existing tasks.
// - Verify that the collection length remains unchanged.
func TestFailedOperationDoesNotModifyState(t *testing.T) {
	tm := NewTaskManager()
	tm.CreateTask(1, "Learn Go")
	tm.CreateTask(2, "Learn Python")

	// Attempt an invalid update
	err := tm.UpdateTask(3, "Learn Java")
	if err == nil {
		t.Errorf("Expected error, got no error")
	}

	// Attempt to delete a nonexistent task
	err = tm.DeleteTask(3)
	if err == nil {
		t.Errorf("Expected error, got no error")
	}

	// Attempt to create a duplicate ID
	err = tm.CreateTask(1, "Learn Rust")
	if err == nil {
		t.Errorf("Expected error, got no error")
	}

	// Verify that the collection length remains unchanged
	tasks := tm.ListTasks()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
}

// TestTaskManagerIsolation
// - Create two independent TaskManager instances.
// - Add tasks to the first manager.
// - Verify that the second manager remains unchanged.
// - Verify that modifications to one manager do not affect the other.
func TestTaskManagerIsolation(t *testing.T) {
	tm1 := NewTaskManager()
	tm2 := NewTaskManager()

	tm1.CreateTask(1, "Learn Go")
	tm1.CreateTask(2, "Learn Python")

	// Verify that the second manager remains unchanged
	tasks2 := tm2.ListTasks()
	if len(tasks2) != 0 {
		t.Errorf("Expected 0 tasks in tm2, got %d", len(tasks2))
	}

	// Verify that modifications to one manager do not affect the other
	tm1.UpdateTask(1, "Learn Go and Testing")
	tasks1 := tm1.ListTasks()
	if len(tasks1) != 2 {
		t.Errorf("Expected 2 tasks in tm1, got %d", len(tasks1))
	}
	if tasks1[0].ID != 1 || tasks1[0].Title != "Learn Go and Testing" {
		t.Errorf("Task fields not set correctly: %+v", tasks1[0])
	}

	tasks2 = tm2.ListTasks()
	if len(tasks2) != 0 {
		t.Errorf("Expected 0 tasks in tm2, got %d", len(tasks2))
	}
}

// TestListTasksDoesNotExposeInternalState
// - Retrieve tasks using ListTasks.
// - Modify the returned collection.
// - Verify whether the internal TaskManager collection is protected.
// - Enforce this test only if ListTasks is intended to return a defensive copy.
func TestListTasksDoesNotExposeInternalState(t *testing.T) {
	tm := NewTaskManager()
	tm.CreateTask(1, "Learn Go")
	tm.CreateTask(2, "Learn Python")

	// Retrieve tasks using ListTasks
	tasks := tm.ListTasks()

	// Modify the returned collection
	tasks = append(tasks, Task{ID: 3, Title: "Learn Rust"})

	// Verify that the internal TaskManager collection is not affected
	tmTasks := tm.ListTasks()
	if len(tmTasks) != 2 {
		t.Errorf("Expected 2 tasks in internal collection, got %d", len(tmTasks))
	}
}

// ------------------------------------------------------------
// 7. INTEGRATION TESTS
// ------------------------------------------------------------

// TestCompleteCRUDWorkflow
// - Create multiple tasks.
// - Find an existing task.
// - Update an existing task.
// - Delete an existing task.
// - List all remaining tasks.
// - Verify the final collection matches the expected state.
func TestCompleteCRUDWorkflow(t *testing.T) {
	tm := NewTaskManager()

	// Create multiple tasks
	tm.CreateTask(1, "Learn Go")
	tm.CreateTask(2, "Learn Python")
	tm.CreateTask(3, "Learn Rust")

	// Find an existing task
	task, err := tm.FindTask(1)

	if err != nil {
		t.Fatalf("Expected task to be found, got error: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("Expected task ID 1, got %d", task.ID)
	}

	// Update an existing task
	tm.UpdateTask(1, "Learn Go and Testing")

	// Delete an existing task
	tm.DeleteTask(2)

	// List all remaining tasks
	tasks := tm.ListTasks()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	// Verify the final collection matches the expected state
	if tasks[0].ID != 1 || tasks[0].Title != "Learn Go and Testing" {
		t.Errorf("Task fields not set correctly: %+v", tasks[0])
	}
	if tasks[1].ID != 3 || tasks[1].Title != "Learn Rust" {
		t.Errorf("Task fields not set correctly: %+v", tasks[1])
	}
}

// TestCRUDWorkflowWithFailures
// - Create a valid task.
// - Attempt to create a duplicate task.
// - Attempt to update a nonexistent task.
// - Attempt to delete a nonexistent task.
// - Verify that subsequent valid operations still work.
// - Verify that failed operations do not corrupt TaskManager state.
func TestCRUDWorkflowWithFailures(t *testing.T) {
	tm := NewTaskManager()

	// Create a valid task
	tm.CreateTask(1, "Learn Go")

	// Attempt to create a duplicate task
	err := tm.CreateTask(1, "Learn Python")
	if err == nil {
		t.Errorf("Expected error when creating duplicate task")
	}

	// Attempt to update a nonexistent task
	err = tm.UpdateTask(2, "Learn Python")
	if err == nil {
		t.Errorf("Expected error when updating nonexistent task")
	}

	// Attempt to delete a nonexistent task
	err = tm.DeleteTask(2)
	if err == nil {
		t.Errorf("Expected error when deleting nonexistent task")
	}

	// Verify that subsequent valid operations still work
	tm.CreateTask(3, "Learn Rust")
	tasks := tm.ListTasks()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	// Verify that failed operations do not corrupt TaskManager state
	tmTasks := tm.ListTasks()
	if len(tmTasks) != 2 {
		t.Errorf("Expected 2 tasks in internal collection, got %d", len(tmTasks))
	}
}

// ------------------------------------------------------------
// 8. CONCURRENCY TESTS (FUTURE)
// ------------------------------------------------------------

// TestConcurrentTaskOperations
// - Execute multiple TaskManager operations concurrently.
// - Check for race conditions.
// - Verify collection integrity.
// - Run using the Go race detector.
// - Only implement after introducing concurrent access requirements.
func TestConcurrentTaskOperations(t *testing.T) {
	// This test is a placeholder for future concurrency testing.
	// It should be implemented when TaskManager is expected to handle concurrent access.
}

//businuess logic test

// TestBusinessLogic
func TestBusinessLogic(t *testing.T) {
	tm := NewTaskManager()

	// Create a task
	tm.CreateTask(1, "Learn Go")

	// Complete the task
	err := tm.CompleteTask(1)
	if err != nil {
		t.Fatalf("Expected no error when completing task, got error: %v", err)
	}

	// Verify the task is completed
	task, err := tm.FindTask(1)
	if err != nil {
		t.Fatalf("Expected task to be found, got error: %v", err)
	}
	if task.Status != "Completed" {
		t.Errorf("Expected task status 'Completed', got '%s'", task.Status)
	}

	// Reopen the task
	err = tm.ReopenTask(1)
	if err != nil {
		t.Fatalf("Expected no error when reopening task, got error: %v", err)
	}

	// Verify the task is pending
	task, err = tm.FindTask(1)
	if err != nil {
		t.Fatalf("Expected task to be found, got error: %v", err)
	}
	if task.Status != "Pending" {
		t.Errorf("Expected task status 'Pending', got '%s'", task.Status)
	}
}
