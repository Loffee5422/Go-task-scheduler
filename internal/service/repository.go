package service

type TaskRepository interface {
	CreateTask(id int, title string) error
	UpdateTask(id int, newTitle string) error
	DeleteTask(id int) error
	ListTasks() []Task
	FindTask(id int) (Task, error)
}
