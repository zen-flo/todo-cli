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
// NewJSONStore создаёт новый экземпляр JSONStore с указанным файлом.
// Потокобезопасный метод: внутренние операции синхронизированы мьютексом.
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

// AddTask добавляет новую задачу в хранилище.
// Потокобезопасный метод: использует мьютекс для синхронизации доступа.
// Возвращает ошибку, если не удалось сохранить задачу.
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

// ListTasks возвращает все задачи из хранилища.
// Потокобезопасный метод: использует мьютекс для синхронизации доступа.
// Возвращает слайс задач и ошибку, если не удалось загрузить данные.
func (s *JSONStore) ListTasks() ([]task.Task, error) {
	return s.loadTasks()
}

// UpdateTask изменяет название задачи с указанным ID.
// Потокобезопасный метод: использует мьютекс для синхронизации доступа.
// Если задача с таким ID не найдена, возвращает ошибку.
func (s *JSONStore) UpdateTask(id int, newTitle string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.loadTasks()
	if err != nil {
		return err
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Title = newTitle
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("задача с ID %d не найдена", id)
	}

	return s.saveTasks(tasks)
}

// DeleteTask удаляет задачу с указанным ID из хранилища.
// Потокобезопасный метод: использует мьютекс для синхронизации доступа.
// Если задача с таким ID не найдена, возвращает ошибку.
func (s *JSONStore) DeleteTask(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Загружаем все задачи
	tasks, err := s.loadTasks()
	if err != nil {
		return err
	}

	// Создаём новый слайс без удаляемой задачи
	var newTasks []task.Task
	found := false
	for _, t := range tasks {
		if t.ID == id {
			found = true
			continue // пропускаем удаляемую задачу
		}
		newTasks = append(newTasks, t)
	}

	if !found {
		return fmt.Errorf("задача с ID %d не найдена", id)
	}

	// Сохраняем обновлённый список задач
	return s.saveTasks(newTasks)
}

// MarkTaskDone отмечает задачу с указанным ID как выполненную.
// Потокобезопасный метод: использует мьютекс для синхронизации доступа.
// Если задача с таким ID не найдена, возвращает ошибку.
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
