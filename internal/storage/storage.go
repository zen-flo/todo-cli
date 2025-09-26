package storage

import "github.com/zen-flo/todo-cli/internal/task"

// Storage — интерфейс для работы с задачами.
// Это позволит легко подменять хранилище (например, JSON → SQLite).
type Storage interface {
	AddTask(t task.Task) error       // добавить задачу
	ListTasks() ([]task.Task, error) // получить список задач
	UpdateTask(t task.Task) error    // обновить задачу
	DeleteTask(id int) error         // удалить задачу
}
