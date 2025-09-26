package storage

import (
	"os"
	"testing"
	"time"

	"github.com/zen-flo/todo-cli/internal/task"
)

// TestAddTask проверяет добавление задачи в JSONStore.
func TestAddTask(t *testing.T) {
	// Создаём временный файл, чтобы не трогать реальный tasks.json.
	tmpFile, err := os.CreateTemp("", "tasks_*.json")
	if err != nil {
		t.Fatalf("не удалось создать временный файл: %v", err)
	}

	// Гарантируем удаление временного файла после теста.
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	// Создаём новое хранилище, указывая путь к временному файлу.
	store := NewJSONStore(tmpFile.Name())

	// Создаём тестовую задачу.
	newTask := task.Task{
		ID:        1,
		Title:     "Тестовая задача",
		Completed: false,
		CreatedAt: time.Now(),
	}

	// Добавляем задачу в хранилище.
	if err := store.AddTask(newTask); err != nil {
		t.Fatalf("AddTask вернул ошибку: %v", err)
	}

	// Проверяем, что задача реально добавилась.
	tasks, err := store.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks вернул ошибку: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("ожидалось 1 задача, получили %d", len(tasks))
	}

	if tasks[0].Title != "Тестовая задача" {
		t.Errorf("ожидался заголовок 'Тестовая задача', получили '%s'", tasks[0].Title)
	}
}

// TestListTasks проверяет корректность загрузки задач из JSONStore.
func TestListTasks(t *testing.T) {
	// Создаём временный файл, чтобы не трогать реальный tasks.json.
	tmpFile, err := os.CreateTemp("", "tasks_*.json")
	if err != nil {
		t.Fatalf("не удалось создать временный файл: %v", err)
	}

	// Гарантируем удаление временного файла после теста.
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	// Создаём новое хранилище.
	store := NewJSONStore(tmpFile.Name())

	// Добавляем несколько тестовых задач.
	tasksToAdd := []task.Task{
		{ID: 1, Title: "Первая", CreatedAt: time.Now()},
		{ID: 2, Title: "Вторая", CreatedAt: time.Now().Add(time.Minute)},
		{ID: 3, Title: "Третья", Completed: true, CreatedAt: time.Now().Add(2 * time.Minute)},
	}

	for _, tt := range tasksToAdd {
		if err := store.AddTask(tt); err != nil {
			t.Fatalf("AddTask вернул ошибку: %v", err)
		}
	}

	// Получаем список задач.
	tasks, err := store.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks вернул ошибку: %v", err)
	}

	// Проверяем количество задач.
	if len(tasks) != len(tasksToAdd) {
		t.Fatalf("ожидалось %d задач, получено %d", len(tasksToAdd), len(tasks))
	}

	// Проверяем, что задачи совпадают по содержимому.
	for i, tt := range tasksToAdd {
		if tasks[i].Title != tt.Title {
			t.Errorf("ожидался заголовок %q, получено %q", tt.Title, tasks[i].Title)
		}
		if tasks[i].Completed != tt.Completed {
			t.Errorf("ошибка статуса: ожидалось %v, получено %v", tt.Completed, tasks[i].Completed)
		}
	}
}

// TestDeleteTask проверяет корректность удаления задачи по ID из JSONStore.
func TestDeleteTask(t *testing.T) {
	// Создаём временный файл для теста.
	tmpFile, err := os.CreateTemp("", "tasks_*.json")
	if err != nil {
		t.Fatalf("не удалось создать временный файл: %v", err)
	}

	// Гарантируем удаление файла после теста.
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	// Инициализируем хранилище.
	store := NewJSONStore(tmpFile.Name())

	// Добавляем несколько задач.
	tasksToAdd := []task.Task{
		{ID: 1, Title: "Первая", CreatedAt: time.Now()},
		{ID: 2, Title: "Вторая", CreatedAt: time.Now()},
		{ID: 3, Title: "Третья", CreatedAt: time.Now()},
	}

	for _, tt := range tasksToAdd {
		if err := store.AddTask(tt); err != nil {
			t.Fatalf("AddTask вернул ошибку: %v", err)
		}
	}

	// Удаляем задачу с ID = 2.
	if err := store.DeleteTask(2); err != nil {
		t.Fatalf("DeleteTask вернул ошибку: %v", err)
	}

	// Получаем обновлённый список.
	tasks, err := store.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks вернул ошибку: %v", err)
	}

	// Проверяем, что количество задач уменьшилось.
	if len(tasks) != 2 {
		t.Fatalf("ожидалось 2 задачи после удаления, получено %d", len(tasks))
	}

	// Проверяем, что задачи с ID 2 больше нет.
	for _, tt := range tasks {
		if tt.ID == 2 {
			t.Errorf("задача с ID=2 всё ещё существует после удаления")
		}
	}
}

// TestUpdateTask проверяет корректность обновления задачи по ID.
func TestUpdateTask(t *testing.T) {
	// Создаём временный файл для теста.
	tmpFile, err := os.CreateTemp("", "tasks_*.json")
	if err != nil {
		t.Fatalf("не удалось создать временный файл: %v", err)
	}

	// Гарантируем удаление временного файла после теста.
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	// Инициализируем хранилище.
	store := NewJSONStore(tmpFile.Name())

	// Добавляем тестовую задачу.
	initialTask := task.Task{
		ID:        1,
		Title:     "Старое название",
		Completed: false,
		CreatedAt: time.Now(),
	}

	if err := store.AddTask(initialTask); err != nil {
		t.Fatalf("AddTask вернул ошибку: %v", err)
	}

	// Обновляем заголовок задачи.
	newTitle := "Новое название"
	if err := store.UpdateTask(1, newTitle); err != nil {
		t.Fatalf("UpdateTask вернул ошибку: %v", err)
	}

	// Загружаем обновлённый список.
	tasks, err := store.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks вернул ошибку: %v", err)
	}

	// Проверяем, что в списке осталась только одна задача.
	if len(tasks) != 1 {
		t.Fatalf("ожидалась 1 задача, получено %d", len(tasks))
	}

	// Проверяем, что название изменилось.
	if tasks[0].Title != newTitle {
		t.Errorf("ожидалось новое название %q, получено %q", newTitle, tasks[0].Title)
	}
}

// TestMarkDone проверяет корректность отметки задачи как выполненной.
func TestMarkDone(t *testing.T) {
	// Создаём временный файл для теста.
	tmpFile, err := os.CreateTemp("", "tasks_*.json")
	if err != nil {
		t.Fatalf("не удалось создать временный файл: %v", err)
	}

	// Удаляем файл после теста.
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	// Инициализируем хранилище.
	store := NewJSONStore(tmpFile.Name())

	// Добавляем задачу, которая изначально не выполнена.
	taskToAdd := task.Task{
		ID:        1,
		Title:     "Проверить MarkDone",
		Completed: false,
		CreatedAt: time.Now(),
	}

	if err := store.AddTask(taskToAdd); err != nil {
		t.Fatalf("AddTask вернул ошибку: %v", err)
	}

	// Отмечаем задачу как выполненную.
	if err := store.MarkTaskDone(1); err != nil {
		t.Fatalf("MarkDone вернул ошибку: %v", err)
	}

	// Получаем список задач.
	tasks, err := store.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks вернул ошибку: %v", err)
	}

	// Проверяем, что задача осталась одна.
	if len(tasks) != 1 {
		t.Fatalf("ожидалась 1 задача, получено %d", len(tasks))
	}

	// Проверяем, что задача действительно отмечена как выполненная.
	if !tasks[0].Completed {
		t.Errorf("ожидалось, что задача будет выполнена, но Completed=false")
	}
}
