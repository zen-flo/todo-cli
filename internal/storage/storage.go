package storage

import "github.com/zen-flo/todo-cli/internal/task"

// Storage описывает интерфейс для сохранения задач
type Storage interface {
	AddTask(t task.Task) error
	ListTasks() ([]task.Task, error)
	UpdateTask(t task.Task) error
	DeleteTask(id int) error
}
