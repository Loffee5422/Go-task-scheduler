package cli

import (
	"bytes"
	"strings"
	"testing"

	"task_manager/internal/service"
)

func TestRunAddAndFind(t *testing.T) {
	manager := service.NewTaskManager()
	commandLine := NewCLIWithWriter(manager, &bytes.Buffer{})

	if err := commandLine.Run([]string{"add", "1", "Learn", "Go"}); err != nil {
		t.Fatalf("add failed: %v", err)
	}
	found, err := manager.FindTask(1)
	if err != nil {
		t.Fatalf("find failed: %v", err)
	}
	if found.Title != "Learn Go" {
		t.Fatalf("expected joined title %q, got %q", "Learn Go", found.Title)
	}
}

func TestRunValidationErrors(t *testing.T) {
	commandLine := NewCLIWithWriter(service.NewTaskManager(), &bytes.Buffer{})
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"missing command", nil, "no command provided"},
		{"unknown command", []string{"unknown"}, "unknown command"},
		{"invalid ID", []string{"add", "not-an-id", "title"}, "task ID must be an integer"},
		{"missing title", []string{"add", "1"}, "usage: add <id> <title>"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := commandLine.Run(test.args)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected error containing %q, got %v", test.want, err)
			}
		})
	}
}

func TestRunCommandDelegatesToRun(t *testing.T) {
	manager := service.NewTaskManager()
	commandLine := NewCLIWithWriter(manager, &bytes.Buffer{})

	if err := commandLine.Run([]string{"add", "1", "Task"}); err != nil {
		t.Fatalf("RunCommand failed: %v", err)
	}
	if _, err := manager.FindTask(1); err != nil {
		t.Fatalf("task was not created: %v", err)
	}
}

func TestUpdateTask(t *testing.T) {
	manager := service.NewTaskManager()
	commandLine := NewCLIWithWriter(manager, &bytes.Buffer{})

	if err := commandLine.Run([]string{"add", "1", "Old title"}); err != nil {
		t.Fatal(err)
	}
	if err := commandLine.Run([]string{"update", "1", "New", "title"}); err != nil {
		t.Fatalf("update failed: %v", err)
	}
	task, err := manager.FindTask(1)
	if err != nil {
		t.Fatal(err)
	}
	if task.Title != "New title" {
		t.Fatalf("expected updated title %q, got %q", "New title", task.Title)
	}
}

func TestDeleteTask(t *testing.T) {
	manager := service.NewTaskManager()
	commandLine := NewCLIWithWriter(manager, &bytes.Buffer{})

	if err := commandLine.Run([]string{"add", "1", "Task"}); err != nil {
		t.Fatal(err)
	}
	if err := commandLine.Run([]string{"delete", "1"}); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if _, err := manager.FindTask(1); err == nil {
		t.Fatal("expected deleted task to be absent")
	}
}

func TestCompleteTask(t *testing.T) {
	manager := service.NewTaskManager()
	commandLine := NewCLIWithWriter(manager, &bytes.Buffer{})

	if err := commandLine.Run([]string{"add", "1", "Task"}); err != nil {
		t.Fatal(err)
	}
	if err := commandLine.Run([]string{"complete", "1"}); err != nil {
		t.Fatalf("complete failed: %v", err)
	}
	task, err := manager.FindTask(1)
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != "Completed" {
		t.Fatalf("expected Completed, got %q", task.Status)
	}
}

func TestReopenTask(t *testing.T) {
	manager := service.NewTaskManager()
	commandLine := NewCLIWithWriter(manager, &bytes.Buffer{})

	if err := commandLine.Run([]string{"add", "1", "Task"}); err != nil {
		t.Fatal(err)
	}
	if err := commandLine.Run([]string{"complete", "1"}); err != nil {
		t.Fatal(err)
	}
	if err := commandLine.Run([]string{"reopen", "1"}); err != nil {
		t.Fatalf("reopen failed: %v", err)
	}
	task, err := manager.FindTask(1)
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != "Pending" {
		t.Fatalf("expected Pending, got %q", task.Status)
	}
}

func TestFilterAndSortCommands(t *testing.T) {
	var output bytes.Buffer
	manager := service.NewTaskManager()
	commandLine := NewCLIWithWriter(manager, &output)

	if err := commandLine.Run([]string{"add", "1", "Medium"}); err != nil {
		t.Fatal(err)
	}
	if err := commandLine.Run([]string{"add", "2", "High"}); err != nil {
		t.Fatal(err)
	}
	if err := commandLine.Run([]string{"filter", "Pending"}); err != nil {
		t.Fatalf("filter failed: %v", err)
	}
	if got := output.String(); !strings.Contains(got, "1: Medium [Pending]") ||
		!strings.Contains(got, "2: High [Pending]") {
		t.Fatalf("unexpected filter output: %q", got)
	}
	output.Reset()
	if err := commandLine.Run([]string{"sort"}); err != nil {
		t.Fatalf("sort failed: %v", err)
	}
	if got := output.String(); got != "1: Medium\n2: High\n" {
		t.Fatalf("unexpected sort output: %q", got)
	}
}

func TestFindNonexistentTask(t *testing.T) {
	commandLine := NewCLIWithWriter(service.NewTaskManager(), &bytes.Buffer{})

	err := commandLine.Run([]string{"find", "99"})
	if err == nil || err.Error() != "task not found" {
		t.Fatalf("expected task not found, got %v", err)
	}
}

func TestDuplicateID(t *testing.T) {
	commandLine := NewCLIWithWriter(service.NewTaskManager(), &bytes.Buffer{})

	if err := commandLine.Run([]string{"add", "1", "First"}); err != nil {
		t.Fatal(err)
	}
	err := commandLine.Run([]string{"add", "1", "Second"})
	if err == nil || err.Error() != "task with this ID already exists" {
		t.Fatalf("expected duplicate ID error, got %v", err)
	}
}

func TestListOutputEmptyAndPopulated(t *testing.T) {
	var output bytes.Buffer
	manager := service.NewTaskManager()
	commandLine := NewCLIWithWriter(manager, &output)

	if err := commandLine.Run([]string{"list"}); err != nil {
		t.Fatal(err)
	}
	if output.String() != "" {
		t.Fatalf("expected empty list output, got %q", output.String())
	}
	if err := commandLine.Run([]string{"add", "1", "First"}); err != nil {
		t.Fatal(err)
	}
	if err := commandLine.Run([]string{"add", "2", "Second"}); err != nil {
		t.Fatal(err)
	}
	output.Reset()
	if err := commandLine.Run([]string{"list"}); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), "1: First\n2: Second\n"; got != want {
		t.Fatalf("expected list output %q, got %q", want, got)
	}
}

func TestHelpOutput(t *testing.T) {
	var output bytes.Buffer
	commandLine := NewCLIWithWriter(service.NewTaskManager(), &output)

	if err := commandLine.Run([]string{"help"}); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"Available commands:",
		"add <id> <title>",
		"update <id> <new_title>",
		"delete <id>",
		"find <id>",
		"complete <id>",
		"reopen <id>",
		"filter <status>",
		"sort",
	} {
		if !strings.Contains(output.String(), expected) {
			t.Errorf("help output missing %q: %q", expected, output.String())
		}
	}
}
