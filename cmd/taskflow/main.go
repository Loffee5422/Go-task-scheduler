package main

import (
	"fmt"
	"os"

	"task_manager/internal/CLI"
	"task_manager/internal/service"
	"task_manager/internal/storage"
)

func main() {
	repo := storage.NewTaskRepository("tasks.json")
	manager, err := service.NewPersistentTaskManager(repo)
	if err != nil {
		fmt.Printf("Error initializing TaskManager: %v\n", err)
		return
	}

	commandLine := cli.NewCLI(manager)
	if err := commandLine.Run(os.Args[1:]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) > 1 && shouldSave(os.Args[1]) {
		if err := manager.Save(); err != nil {
			fmt.Printf("Error saving tasks: %v\n", err)
			os.Exit(1)
		}
	}
}

func shouldSave(command string) bool {
	switch command {
	case "add", "update", "delete", "complete", "reopen", "sort":
		return true
	default:
		return false
	}
}
