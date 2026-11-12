# TaskFlow CLI

TaskFlow is a local task manager written in Go. Tasks are stored in `tasks.json`
and survive between command invocations.

## Run

```text
go run ./cmd/taskflow <command> [arguments]
```

## Commands

```text
add <id> <title>        Create a task
list                    List all tasks
find <id>               Show one task
update <id> <title>     Change a task title
delete <id>             Delete a task
complete <id>           Mark a task completed
reopen <id>             Move a completed task back to pending
filter <status>         List tasks with a status
sort                    Sort tasks by priority and list them
help                    Show command help
```

Supported statuses are `Pending`, `In Progress`, and `Completed`. New tasks
start as `Pending` with `Medium` priority.

Mutating commands save successful changes to `tasks.json`. Read-only commands
do not rewrite the file. A failed command or failed save exits with a non-zero
status.

## Test

```text
go test ./...
```
