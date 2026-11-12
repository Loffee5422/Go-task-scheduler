# Go Project: Build a Task Manager CLI

Beginner → Intermediate

Go fundamentals

Estimated: 20–30 hours

The objective is to build a command-line task manager using Go, starting with a minimal working application and gradually extending it into a maintainable, production-quality CLI tool.

I will structure this as an engineering project brief rather than a coding tutorial. You will get:

- Functional requirements.
    
- Technical constraints.
    
- Development milestones.
    
- Acceptance criteria.
    
- Architectural hints.
    
- Testing challenges.
    
- Optional advanced extensions.
    

I will not provide the implementation or complete code. You should make the architectural decisions, write the code, debug it, and validate the results yourself.

# 1. Project Overview

TaskFlow CLI

Your project

# A local task management system

Build a terminal-based application that allows users to create, organize, update, search, and track tasks without requiring a database server or graphical interface.

![Taskline: manage your tasks with style](https://images.openai.com/static-rsc-4/4rgFmcf4UAUrOW49L5w1cbVtQkeLg5NtiSRuK-BeSfDaUpO7Vcw4vxDoZ53RnJkis6MER7dSY3WJIbX7IVdizdL620dN7b-1FAl8agMrynljqmn4AIUynP8685PHPgQKjVZka8Fuj9wP14uBQp5kI_8AUN_Tb5N2OPdRDZ1s3qQ?purpose=inline)

[medevel.com](https://medevel.com/taskline/?utm_source=chatgpt.com)

Core capabilities

Create tasks

Create and validate tasks.

List tasks

Filter and sort tasks.

Update tasks

Modify existing task properties.

Track progress

Complete and reopen tasks.

Search

Find tasks using keywords.

Persistence

Preserve data between executions.

## Learning objectives

By completing this project, you should understand:

|Go concept|Practical application|
|---|---|
|Structs|Task data models|
|Methods|Task behavior|
|Slices|Managing collections|
|Maps|Indexing and lookup|
|Interfaces|Decoupling storage|
|Pointers|Modifying shared data|
|Error handling|Robust CLI behavior|
|Packages|Modular application design|
|File I/O|Persistent storage|
|JSON|Data serialization|
|Goroutines|Optional concurrent operations|
|Testing|Reliable application behavior|
|Modules|Dependency management|

The important objective is not simply learning Go syntax. It is learning how Go applications are structured.

# 2. Define the Product Requirements

Before writing code, establish what the application should do.

## 2.1 Task data model

Every task should contain the following information.

|Field|Type|Required|Description|
|---|---|---|---|
|ID|Integer|Yes|Unique identifier|
|Title|String|Yes|Task name|
|Description|String|No|Additional information|
|Status|Enum-like string|Yes|Pending / In Progress / Completed|
|Priority|Enum-like string|Yes|Low / Medium / High|
|CreatedAt|Time|Yes|Creation timestamp|
|UpdatedAt|Time|Yes|Last modification|
|DueDate|Time|No|Deadline|
|Tags|String collection|No|Task categories|

Design challenge:

Think about how you would represent these fields in Go.

Consider:

- Should status be a string or a custom type?
    
- Should priority be a custom type?
    
- Should DueDate be a pointer?
    
- Should timestamps use `time.Time`?
    
- How will you prevent invalid status values?
    

Hint: Go's type system can help you prevent invalid application states. Do not rely entirely on runtime string comparisons.

## 2.2 Functional requirements

Your application must support these operations.

### A. Create

Example:

```
taskflow add "Learn Go"
```

Additional functionality:

- Set priority.
    
- Add a description.
    
- Set a due date.
    
- Add tags.
    

Expected behavior:

- Generate a unique ID.
    
- Set creation and modification timestamps.
    
- Assign sensible defaults.
    
- Validate user input.
    

Challenge: What should happen if the title is empty or the supplied priority is invalid?

### B. List

Example:

```
taskflow list
```

Additional examples:

```
taskflow list --status pending
taskflow list --priority high
taskflow list --sort due-date
taskflow list --tag learning
```

Expected behavior:

- Display all tasks by default.
    
- Support multiple filters.
    
- Support sorting.
    
- Handle empty collections.
    
- Display readable output.
    

Challenge: How will you combine multiple filters?

For example, a user might request:

- Pending tasks.
    
- High priority.
    
- Due within seven days.
    

Think about whether filtering should happen in the CLI layer or the task management layer.

### C. Update

Example:

```
taskflow update 12 --title "Learn Go concurrency"
```

Supported fields:

- Title.
    
- Description.
    
- Priority.
    
- Due date.
    
- Tags.
    
- Status.
    

Expected behavior:

- Find the correct task.
    
- Validate the requested changes.
    
- Preserve fields that were not modified.
    
- Update the modification timestamp.
    

Challenge: How do you distinguish between a field that was not provided and a field that was explicitly cleared?

This is an important data-modeling problem.

### D. Complete

Example:

```
taskflow complete 12
```

Additional functionality:

```
taskflow reopen 12
```

Expected behavior:

- Change task status.
    
- Reject invalid task IDs.
    
- Prevent unnecessary modifications.
    
- Persist the result.
    

### E. Delete

Example:

```
taskflow delete 12
```

Consider supporting:

```
taskflow delete 12 --force
```

Challenge: Should deletion be permanent, or should you implement soft deletion?

For the first version, permanent deletion is sufficient.

### F. Search

Example:

```
taskflow search "golang"
```

Search should support:

- Title.
    
- Description.
    
- Tags.
    

Consider case-insensitive matching.

# 3. CLI Interface Design

A good CLI should be predictable.

Your commands should follow consistent naming conventions and argument structures.

Here is the proposed interface:

|Command|Purpose|
|---|---|
|`add`|Create a task|
|`list`|List tasks|
|`show`|Display task details|
|`update`|Modify a task|
|`complete`|Complete a task|
|`reopen`|Reopen a task|
|`delete`|Delete a task|
|`search`|Search tasks|
|`stats`|Show task statistics|
|`help`|Show documentation|
|`version`|Show application version|

## CLI design challenge

Design the following behaviors yourself:

1. What happens when the user enters an unknown command?
    
2. What happens when a required argument is missing?
    
3. How should optional flags behave?
    
4. Should commands support aliases?
    
5. How should the application communicate errors?
    
6. Should output be human-readable, machine-readable, or both?
    

Hint: Look at how `git`, `docker`, and `kubectl` structure their command-line interfaces.

Pay particular attention to:

- Subcommands.
    
- Flags.
    
- Help messages.
    
- Exit codes.
    
- Consistent error messages.
    

# 4. Architecture and Package Design

This is one of the most important parts of the project.

Avoid putting the entire application inside `main.go`.

You should aim for a small modular application, not an unnecessarily complicated enterprise architecture.

## Suggested architecture

## Application architecture

CLI Layer

Commands · Flags · Input · Output

Application / Service Layer

Task operations · Validation · Business rules

Domain Layer

Task · Status · Priority · Domain rules

Storage Layer

Repository contract · JSON persistence

JSON file

Initial storage implementation

In-memory

Testing implementation

## Suggested project structure

This is a structural target, not a requirement to create every package immediately.

```
taskflow/
│
├── cmd/
│   └── taskflow/
│
├── internal/
│   ├── task/
│   ├── service/
│   ├── storage/
│   └── cli/
│
├── test/
│
├── go.mod
├── README.md
└── Makefile
```

### Your architectural responsibilities

Domain

- Define what a task is.
    
- Define valid task states.
    
- Define task-related rules.
    

Service

- Coordinate task operations.
    
- Enforce business rules.
    
- Handle interactions with storage.
    

Storage

- Save tasks.
    
- Load tasks.
    
- Update tasks.
    
- Delete tasks.
    

CLI

- Interpret user commands.
    
- Parse arguments.
    
- Call application operations.
    
- Display results.
    

### Important architectural questions

Try answering these before implementation:

1. Should the domain package know anything about JSON?
    
2. Should the CLI directly modify task objects?
    
3. Should the service know about command-line flags?
    
4. Should the storage layer decide whether a task is high priority?
    
5. How would you replace JSON storage with SQLite without rewriting the application?
    

Hint: Interfaces are useful when they represent a meaningful boundary between components. Avoid creating interfaces for every struct or function.

# 5. Development Milestones

Work through these milestones sequentially. Do not move to the next milestone until the acceptance criteria are satisfied.

Estimated time assumes approximately 2–3 hours per session.

## Milestone 1 — Go fundamentals and domain model

Foundation

Estimated: 3–4 hours

Objective: Build the domain model independently of CLI and persistence.

Tasks:

- Create the Go module.
    
- Define the Task model.
    
- Define task statuses.
    
- Define priority levels.
    
- Establish validation rules.
    
- Implement task-related methods.
    
- Define how IDs are assigned.
    
- Define timestamp behavior.
    

Challenges:

- How do Go structs differ from Python classes?
    
- How do value receivers differ from pointer receivers?
    
- When should you use pointers?
    
- How does Go handle zero values?
    
- How can custom types improve correctness?
    

Acceptance criteria:

- Task objects can be created.
    
- Invalid task data can be rejected.
    
- Task state transitions follow defined rules.
    
- Domain logic does not depend on storage or CLI packages.
    

## Milestone 2 — In-memory task management

Core functionality

Estimated: 3–4 hours

Objective: Make the task manager functional without persistent storage.

Tasks:

- Add tasks.
    
- Retrieve tasks by ID.
    
- List all tasks.
    
- Update tasks.
    
- Complete tasks.
    
- Reopen tasks.
    
- Delete tasks.
    
- Search tasks.
    

Challenges:

- Should you use a slice or a map?
    
- What are the performance implications?
    
- How do you handle duplicate IDs?
    
- What happens when an ID does not exist?
    
- How do you avoid accidentally modifying unrelated tasks?
    

Acceptance criteria:

All operations work correctly in memory.

You should be able to create 100 tasks, retrieve individual tasks, update them, and delete them without inconsistent state.

## Milestone 3 — Persistent storage

Data engineering

Estimated: 3–4 hours

Objective: Ensure tasks survive application restarts.

Initially, use JSON file persistence.

Tasks:

- Design a JSON storage format.
    
- Save tasks to disk.
    
- Load tasks during application startup.
    
- Handle missing storage files.
    
- Handle malformed JSON.
    
- Handle file permission errors.
    
- Prevent accidental data loss during writes.
    

Challenges:

- What happens if the application crashes during a write?
    
- Should you overwrite the entire file or modify individual records?
    
- How do you handle corrupted files?
    
- Should an empty file represent an empty task collection?
    
- Where should the application store its data?
    

Hint: Investigate Go's `encoding/json`, `os`, `path/filepath`, and `io` packages.

Important: Think about atomic file replacement. A partially written JSON file should not destroy the user's previous data.

Acceptance criteria:

1. Create tasks.
    
2. Exit the application.
    
3. Restart the application.
    
4. Verify that previously created tasks still exist.
    

## Milestone 4 — Build the CLI

User interface

Estimated: 4–5 hours

Objective: Make the application usable from a terminal.

Tasks:

- Implement command parsing.
    
- Implement command routing.
    
- Support positional arguments.
    
- Support optional flags.
    
- Implement help output.
    
- Implement version output.
    
- Display readable task tables.
    
- Handle invalid input.
    
- Implement meaningful exit codes.
    

Challenges:

- Should you use Go's standard `flag` package or a CLI framework?
    
- What are the trade-offs between simplicity and extensibility?
    
- How will you support nested commands?
    
- How will you handle quoted strings?
    
- How should errors be communicated?
    

Hint: Compare the standard library's `flag` package with Cobra.

Do not add a framework just because it is popular. Make a decision based on your requirements.

Acceptance criteria:

Every core operation must be executable through the terminal.

## Milestone 5 — Filtering, sorting and statistics

Intermediate

Estimated: 3–4 hours

Objective: Add functionality that makes the application genuinely useful.

Tasks:

- Filter by status.
    
- Filter by priority.
    
- Filter by tags.
    
- Filter by deadline.
    
- Sort by creation date.
    
- Sort by due date.
    
- Sort by priority.
    
- Search titles and descriptions.
    
- Calculate task statistics.
    

Example statistics:

```
TASK STATISTICS

Total tasks:       45
Pending:           12
In Progress:        8
Completed:         25

Completion Rate: 55.6%

Overdue:            3
High Priority:      7
```

Challenges:

- How should multiple filters interact?
    
- Should sorting happen before or after filtering?
    
- How should overdue tasks be calculated?
    
- How should you handle tasks without deadlines?
    
- Should completion rate include deleted tasks?
    

Acceptance criteria:

FilteringAcceptance criteria:

- Multiple filters work together.
    
- Sorting produces deterministic results.
    
- Statistics accurately reflect stored tasks.
    
- Edge cases are handled correctly.
    

## Milestone 6 — Testing and reliability

Engineering quality

Estimated: 3–4 hours

Objective: Make the application reliable and prevent regressions.

Tasks:

- Write unit tests for domain logic.
    
- Write tests for task operations.
    
- Test storage operations.
    
- Test malformed input.
    
- Test invalid task IDs.
    
- Test missing files.
    
- Test corrupted JSON.
    
- Test command-line behavior.
    
- Add integration tests.
    
- Measure test coverage.
    

Challenges:

1. How do you test code that reads and writes files without modifying your real task data?
    
2. How can you test error conditions?
    
3. What should be mocked?
    
4. When should you use an in-memory implementation?
    
5. How do you prevent tests from depending on execution order?
    

Hint: Explore:

- `testing`
    
- `testing.T`
    
- `t.Run`
    
- `t.TempDir()`
    
- Table-driven tests
    
- `go test ./...`
    
- `go test -race ./...`
    
- `go test -cover ./...`
    

Acceptance criteria:

- All tests pass.
    
- Tests are deterministic.
    
- Tests do not depend on your personal filesystem.
    
- Important business rules have direct test coverage.
    
- No data corruption occurs during normal application operations.
    

## Milestone 7 — Error handling and application hardening

Production readiness

Estimated: 2–3 hours

Objective: Handle unexpected conditions without crashing or corrupting data.

Tasks:

- Define consistent error handling.
    
- Introduce meaningful error messages.
    
- Distinguish expected and unexpected errors.
    
- Wrap errors with additional context.
    
- Validate all external input.
    
- Handle interrupted operations.
    
- Improve file write safety.
    
- Ensure resources are closed correctly.
    

Challenges:

Consider these scenarios:

|Scenario|Expected behavior|
|---|---|
|Invalid task ID|Clear error|
|Missing task|Not-found error|
|Invalid priority|Validation error|
|Corrupted JSON|Storage error|
|Permission denied|Filesystem error|
|Disk full|Write failure|
|Empty title|Validation error|
|Invalid date|Input error|

Hint: Investigate:

- `errors.New`
    
- `fmt.Errorf`
    
- `%w`
    
- `errors.Is`
    
- `errors.As`
    
- `defer`
    
- `os.Rename`
    

Think about the difference between returning an error and terminating the entire application.

Acceptance criteria:

The application must fail gracefully. Errors should retain enough context to identify their source without exposing unnecessary internal implementation details.

# 6. Advanced Features

These are optional extensions. Complete the core application first.

## Level 1 — Improve usability

Intermediate

|Feature|Learning objective|
|---|---|
|Colored terminal output|Terminal formatting|
|Interactive task selection|CLI interaction|
|Due-date reminders|Time handling|
|Task categories|Data modeling|
|Recurring tasks|Scheduling logic|
|Export to CSV|Data transformation|
|Import from JSON|Validation and parsing|
|Configuration file|Application configuration|

Potential challenge:

How would you allow users to configure the default task priority without hardcoding it?

Explore environment variables, configuration files, and command-line flags.

## Level 2 — Concurrency

Go-specific learning

Introduce concurrency only when there is a meaningful use case.

Potential features:

- Background overdue-task checking.
    
- Periodic automatic backups.
    
- Concurrent task searching.
    
- Background synchronization.
    
- Parallel import validation.
    

Challenges:

1. What happens when two goroutines modify the same task collection?
    
2. How do you prevent data races?
    
3. When should you use a mutex?
    
4. When should you use channels?
    
5. How do you stop background goroutines?
    
6. What is the role of `context.Context`?
    

Hint: Investigate:

- Goroutines.
    
- Channels.
    
- `sync.Mutex`.
    
- `sync.RWMutex`.
    
- `sync.WaitGroup`.
    
- `context.Context`.
    
- Race detector.
    

Important: A task manager does not inherently need concurrency. Do not introduce goroutines simply to demonstrate that you know how to use them.

## Level 3 — Database integration

Advanced

Replace JSON storage with SQLite.

Requirements:

- Design a relational schema.
    
- Define primary keys.
    
- Define appropriate indexes.
    
- Implement database operations.
    
- Handle transactions.
    
- Handle database errors.
    
- Preserve existing task functionality.
    
- Migrate existing JSON data into SQLite.
    

Challenges:

- What makes a database transaction atomic?
    
- Why are indexes useful?
    
- What is the difference between SQL queries and prepared statements?
    
- How should you handle concurrent database access?
    
- How do you migrate data without losing existing records?
    

Hint: Explore Go's `database/sql` package and SQLite drivers.

## Level 4 — Performance and benchmarking

Advanced

Your task manager will probably be fast enough for ordinary usage.

However, use this milestone to learn Go's performance tooling.

Tasks:

- Benchmark task creation.
    
- Benchmark task retrieval.
    
- Benchmark filtering.
    
- Benchmark sorting.
    
- Benchmark JSON serialization.
    
- Profile CPU usage.
    
- Profile memory allocation.
    
- Identify unnecessary allocations.
    
- Compare slice-based and map-based lookup.
    

Useful commands:

```
go test -bench=. -benchmem ./...
go test -cpuprofile=cpu.prof ./...
go test -memprofile=mem.prof ./...
```

Challenges:

- When is a map faster than a slice?
    
- What is the cost of repeated memory allocation?
    
- How does Go's garbage collector affect performance?
    
- How do you identify an actual bottleneck instead of guessing?
    

Hint: Investigate `pprof` and Go's benchmarking framework.

# 7. Recommended Development Workflow

Treat this as a real engineering project.

Do not write the entire application in one sitting.

Use the following workflow for each milestone.

## Development lifecycle

1. Understand requirements
    
    Identify expected behavior, constraints, and edge cases.
    
2. Design
    
    Decide which package owns the functionality and define its interfaces.
    
3. Implement
    
    Write the smallest working implementation.
    
4. Test
    
    Verify normal cases, boundary conditions, and failure cases.
    
5. Review
    
    Check readability, complexity, naming, and package dependencies.
    
6. Refactor
    
    Remove unnecessary complexity without changing behavior.
    
7. Commit
    
    Save a working, independently understandable milestone.
    

## Git workflow

Use Git from the beginning.

Suggested milestones:

```
feat/domain-model
feat/in-memory-storage
feat/json-persistence
feat/cli-commands
feat/filtering-sorting
test/core-functionality
feat/error-handling
feat/advanced-features
```

Suggested commit conventions:

```
feat: add task creation
feat: implement JSON persistence
feat: support task filtering
test: add task service tests
fix: handle corrupted storage files
refactor: separate CLI from business logic
```

## Go commands you should know

|Command|Purpose|
|---|---|
|`go mod init`|Initialize module|
|`go mod tidy`|Manage dependencies|
|`go run`|Execute application|
|`go build`|Compile application|
|`go test`|Run tests|
|`go test ./...`|Test all packages|
|`go vet ./...`|Detect suspicious code|
|`gofmt`|Format source code|
|`go doc`|Read documentation|
|`go get`|Add or update dependencies|

Additional challenge:

Can you create a Makefile that allows you to execute the most important development operations with short commands?

# 8. Technical Research Checklist

Do not blindly copy implementations from tutorials. Use official documentation to understand the underlying concepts.

These are the primary resources you should consult.

![Golang and why it matters. Go, or as its easily google-able… | by James O'Toole | Medium](https://images.openai.com/static-rsc-4/StbaYV9H0kiDpA5LLjDcstrfAmmR5o4qXWbK03fzjS0DnD9qaXyouHOLuvAF4wJZv99FRoqZydC2Xe-aL_6h73ybrrWto0cydN44j5L_PZB7aWdKP-MOVU1Er7Fdif8DnC8GCq-qy1kL6V3Z-404zOkDMbnf9Cjdos_E4DTQLR0?purpose=inline)

Go Documentation

Primary reference

[Official documentation](https://go.dev/doc/?utm_source=chatgpt.com)

Go Tour

Learn syntax, types, methods, interfaces, and concurrency.

[Go Tour](https://go.dev/tour/?utm_source=chatgpt.com)

Effective Go

Learn idiomatic Go conventions and programming practices.

[Effective Go](https://go.dev/doc/effective_go?utm_source=chatgpt.com)

Go Standard Library

Explore packages for JSON, filesystem operations, time, testing, and synchronization.

[Standard library](https://pkg.go.dev/std?utm_source=chatgpt.com)

Go by Example

Practical examples of Go language features.

[Go by Example](https://gobyexample.com/?utm_source=chatgpt.com)

Cobra CLI

Understand how professional Go CLI applications organize commands.

[Cobra documentation](https://cobra.dev/?utm_source=chatgpt.com)

# 9. Final Acceptance Criteria

Your project is complete when it satisfies the following requirements.

Project completion

# 0%

0 of 13 acceptance criteria satisfied

Domain model

Valid task representation and state transitions.

CRUD operations

Create, retrieve, update, and delete tasks.

Persistence

Tasks survive application restarts.

CLI

All core operations are accessible from the terminal.

Filtering

Multiple filters and sorting work correctly.

Search

Search works across task fields.

Statistics

Task statistics are accurate.

Error handling

Invalid input and storage failures are handled.

Unit tests

Domain and service behavior is tested.

Integration tests

CLI and persistence work together.

Code quality

Formatting, vetting, and package boundaries are sound.

Documentation

README includes installation, usage, and examples.

Build

Application builds into an executable binary.

# 10. Self-Assessment

After finishing the project, you should be able to answer these questions without looking up the answers.

Go fundamentals

1. Why does Go use explicit error returns instead of exceptions?
    
2. What is the difference between a value receiver and pointer receiver?
    
3. How do interfaces work in Go?
    
4. What is the difference between an array and a slice?
    
5. How does Go's garbage collector work?
    
6. What are zero values?
    
7. How does Go handle package visibility?
    

Software architecture

8. Why should business logic be separated from CLI parsing?
    
9. What is the purpose of a repository interface?
    
10. When is an abstraction unnecessary?
    
11. How do you prevent circular package dependencies?
    
12. How would you replace JSON with PostgreSQL without changing business logic?
    

Concurrency

13. What is a goroutine?
    
14. What is a data race?
    
15. When should you use a mutex instead of a channel?
    
16. What does `context.Context` solve?
    
17. How does the Go race detector work?
    

Engineering quality

18. How do you design deterministic tests?
    
19. How do you benchmark a Go function?
    
20. How do you identify memory allocation bottlenecks?
    
21. How do you handle partial failures in file persistence?
    
22. What makes an application production-ready?
    

# 11. Suggested Time Allocation

A reasonable intensive schedule is:

|Milestone|Hours|
|---|---|
|Domain model|3|
|In-memory management|3|
|JSON persistence|4|
|CLI implementation|4|
|Filtering and statistics|3|
|Testing|4|
|Error handling|2|
|Documentation and final review|2|
|Total|25 hours|

Advanced extensions are separate and should not be included in the initial 25-hour target.

# 12. Rules for This Learning Project

I recommend following these constraints throughout development.

1. Do not copy complete implementations. Understand the underlying Go concepts before using them.
    
2. Use the standard library first. Introduce dependencies only when they provide meaningful value.
    
3. Avoid premature abstraction. Do not build a generic task engine for a simple task manager.
    
4. Test incrementally. Do not wait until the application is complete.
    
5. Use meaningful errors. Do not silently ignore failures.
    
6. Keep package boundaries explicit. Avoid circular dependencies and unnecessary cross-package access.
    
7. Use Git checkpoints. Each milestone should leave the application in a working state.
    
8. Prioritize correctness over performance. Optimize only after measuring.
    
9. Write idiomatic Go. Prefer simple, explicit code over unnecessarily sophisticated patterns.
    
10. Be able to explain every abstraction. If you cannot justify why it exists, reconsider it.
    

My recommendation: Complete Milestones 1–7 before touching concurrency or databases. A well-designed JSON-backed CLI with comprehensive tests will teach you substantially more about practical Go engineering than immediately building a complicated concurrent application.

Your final objective is to be able to build a small Go application independently, understand its architecture, diagnose its failures, and justify your design decisions.