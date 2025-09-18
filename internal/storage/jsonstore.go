package storage

import (
	"github.com/zen-flo/todo-cli/internal/task"
)

type JSONStore struct {
	// Тут позже будет путь к файлу и задачи
}

func NewJSONStore() *JSONStore {
	return &JSONStore{}
}

func (s *JSONStore) AddTask(t task.Task) error {
	// Пока заглушка
	return nil
}

func (s *JSONStore) ListTasks() ([]task.Task, error) {
	return []task.Task{}, nil
}

func (s *JSONStore) UpdateTask(t task.Task) error {
	return nil
}

func (s *JSONStore) DeleteTask(id int) error {
	return nil
}
