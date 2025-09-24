package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zen-flo/todo-cli/internal/task"
	"os"
	"sync"
)

// JSONStore — реализация интерфейса Storage.
// Задачи хранятся в JSON-файле на диске.
type JSONStore struct {
	FilePath string     // путь к файлу с задачами
	mu       sync.Mutex // мьютекс для защиты при параллельном доступе
}

// NewJSONStore — конструктор JSONStore.
func NewJSONStore(filePath string) *JSONStore {
	return &JSONStore{FilePath: filePath}
}

// loadTasks — приватный метод, загружает задачи из JSON-файла.
func (s *JSONStore) loadTasks() ([]task.Task, error) {
	data, err := os.ReadFile(s.FilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []task.Task{}, nil // если файла нет — возвращаем пустой список
		}
		return nil, err
	}

	var tasks []task.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// saveTasks — приватный метод, сохраняет список задач в JSON-файл.
func (s *JSONStore) saveTasks(tasks []task.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.FilePath, data, 0644)
}

// AddTask — добавляет новую задачу в список и сохраняет её в JSON.
func (s *JSONStore) AddTask(t task.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Загружаем все задачи
	tasks, err := s.loadTasks()
	if err != nil {
		return err
	}

	// Находим максимальный ID
	maxID := 0
	for _, existing := range tasks {
		if existing.ID > maxID {
			maxID = existing.ID
		}
	}

	// Присваиваем новый ID
	t.ID = maxID + 1

	// Добавляем задачу
	tasks = append(tasks, t)

	// Сохраняем обратно
	return s.saveTasks(tasks)
}

// ListTasks — возвращает все задачи из JSON-файла.
func (s *JSONStore) ListTasks() ([]task.Task, error) {
	return s.loadTasks()
}

// UpdateTask — обновляет задачу по ID.
func (s *JSONStore) UpdateTask(t task.Task) error {
	tasks, err := s.loadTasks()
	if err != nil {
		return err
	}

	for i, existing := range tasks {
		if existing.ID == t.ID {
			tasks[i] = t
			return s.saveTasks(tasks)
		}
	}
	return errors.New("task not found")
}

// DeleteTask — удаляет задачу по ID.
func (s *JSONStore) DeleteTask(id int) error {
	tasks, err := s.loadTasks()
	if err != nil {
		return err
	}

	for i, existing := range tasks {
		if existing.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return s.saveTasks(tasks)
		}
	}
	return errors.New("task not found")
}

// MarkTaskDone отмечает задачу с указанным ID как выполненную.
func (s *JSONStore) MarkTaskDone(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Загружаем все задачи
	tasks, err := s.loadTasks()
	if err != nil {
		return err
	}

	// Ищем задачу по ID
	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Completed = true
			found = true
			break
		}
	}

	// Если задача не найдена — сообщаем пользователю
	if !found {
		return fmt.Errorf("задача с ID %d не найдена", id)
	}

	// Сохраняем обновлённый список задач
	return s.saveTasks(tasks)
}
