package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"task_manager/internal/service"
)

type CLI struct {
	manager *service.TaskManager
	output  io.Writer
}

func NewCLI(manager *service.TaskManager) *CLI {
	return NewCLIWithWriter(manager, os.Stdout)
}

func NewCLIWithWriter(manager *service.TaskManager, output io.Writer) *CLI {
	if output == nil {
		output = io.Discard
	}
	return &CLI{
		manager: manager,
		output:  output,
	}
}

func (cli *CLI) Run(args []string) error {
	if len(args) == 0 {
		return errors.New("no command provided") // No command provided, do nothing.
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "add":
		return cli.handleAdd(commandArgs)
	case "list":
		return cli.handleList()
	case "update":
		return cli.handleUpdate(commandArgs)
	case "delete":
		return cli.handleDelete(commandArgs)
	case "find":
		return cli.handleFind(commandArgs)
	case "complete":
		return cli.handleComplete(commandArgs)
	case "reopen":
		return cli.handleReopen(commandArgs)
	case "filter":
		return cli.handleFilter(commandArgs)
	case "sort":
		return cli.handleSort(commandArgs)
	case "help":
		return cli.handleHelp()
	default:
		return errors.New("unknown command")
	}
}

func (cli *CLI) handleAdd(args []string) error {
	if len(args) < 2 {
		return errors.New("usage: add <id> <title>")
	}
	id, err := parseID(args[0])
	if err != nil {
		return err
	}
	return cli.manager.CreateTask(id, strings.Join(args[1:], " "))
}

func (cli *CLI) handleList() error {
	tasks := cli.manager.ListTasks()
	for _, task := range tasks {
		fmt.Fprintln(cli.output, strconv.Itoa(task.ID)+": "+task.Title)
	}
	return nil
}

func (cli *CLI) handleUpdate(args []string) error {
	if len(args) < 2 {
		return errors.New("usage: update <id> <new_title>")
	}
	id, err := parseID(args[0])
	if err != nil {
		return err
	}
	newTitle := strings.Join(args[1:], " ")
	return cli.manager.UpdateTask(id, newTitle)
}

func (cli *CLI) handleDelete(args []string) error {
	if len(args) < 1 {
		return errors.New("usage: delete <id>")
	}
	id, err := parseID(args[0])
	if err != nil {
		return err
	}
	return cli.manager.DeleteTask(id)
}

func (cli *CLI) handleFind(args []string) error {
	if len(args) < 1 {
		return errors.New("usage: find <id>")
	}
	id, err := parseID(args[0])
	if err != nil {
		return err
	}
	task, err := cli.manager.FindTask(id)
	if err != nil {
		return err
	}
	fmt.Fprintln(cli.output, strconv.Itoa(task.ID)+": "+task.Title)
	return nil
}

func (cli *CLI) handleComplete(args []string) error {
	if len(args) < 1 {
		return errors.New("usage: complete <id>")
	}
	id, err := parseID(args[0])
	if err != nil {
		return err
	}
	return cli.manager.CompleteTask(id)
}

func (cli *CLI) handleReopen(args []string) error {
	if len(args) < 1 {
		return errors.New("usage: reopen <id>")
	}
	id, err := parseID(args[0])
	if err != nil {
		return err
	}
	return cli.manager.ReopenTask(id)
}

func (cli *CLI) handleFilter(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: filter <status>")
	}
	for _, task := range cli.manager.FilterTasksByStatus(args[0]) {
		fmt.Fprintln(cli.output, strconv.Itoa(task.ID)+": "+task.Title+" ["+task.Status+"]")
	}
	return nil
}

func (cli *CLI) handleSort(args []string) error {
	if len(args) != 0 {
		return errors.New("usage: sort")
	}
	cli.manager.SortTasksByPriority()
	return cli.handleList()
}

func (cli *CLI) handleHelp() error {
	fmt.Fprintln(cli.output, "Available commands:")
	fmt.Fprintln(cli.output, "  add <id> <title>       - Add a new task")
	fmt.Fprintln(cli.output, "  list                   - List all tasks")
	fmt.Fprintln(cli.output, "  update <id> <new_title> - Update a task")
	fmt.Fprintln(cli.output, "  delete <id>            - Delete a task")
	fmt.Fprintln(cli.output, "  find <id>              - Find a task")
	fmt.Fprintln(cli.output, "  complete <id>          - Complete a task")
	fmt.Fprintln(cli.output, "  reopen <id>            - Reopen a completed task")
	fmt.Fprintln(cli.output, "  filter <status>        - List tasks by status")
	fmt.Fprintln(cli.output, "  sort                   - Sort tasks by priority")
	fmt.Fprintln(cli.output, "  help                   - Show this help message")
	return nil
}

// RunCommand executes a command when the command name and arguments are separate.
func (cli *CLI) RunCommand(command string, args []string) error {
	return cli.Run(append([]string{command}, args...))
}

func parseID(value string) (int, error) {
	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.New("task ID must be an integer")
	}
	return id, nil
}
