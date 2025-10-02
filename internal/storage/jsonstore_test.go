package storage

import (
	"os"
	"sort"
	"testing"
	"time"

	"github.com/zen-flo/todo-cli/internal/task"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// TestAddTask проверяет добавление задачи с учетом Important и Completed.
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
		Important: true,
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

	if tasks[0].Title != "Тестовая задача" || !tasks[0].Important || tasks[0].Completed {
		t.Errorf("неверные данные задачи: %+v", tasks[0])
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

	// Добавляем несколько тестовых задач
	tasksToAdd := []task.Task{
		{ID: 1, Title: "Первая", CreatedAt: time.Now(), Important: true},
		{ID: 2, Title: "Вторая", CreatedAt: time.Now().Add(time.Minute)},
		{ID: 3, Title: "Третья", Completed: true, CreatedAt: time.Now().Add(2 * time.Minute)},
	}

	for _, tt := range tasksToAdd {
		if err := store.AddTask(tt); err != nil {
			t.Fatalf("AddTask вернул ошибку: %v", err)
		}
	}

	// Получаем все задачи.
	allTasks, err := store.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks вернул ошибку: %v", err)
	}

	// Проверяем количество задач.
	if len(allTasks) != len(tasksToAdd) {
		t.Fatalf("ожидалось %d задач, получено %d", len(tasksToAdd), len(allTasks))
	}

	// Проверяем, что задачи совпадают по содержимому.
	for i, tt := range tasksToAdd {
		if allTasks[i].Title != tt.Title {
			t.Errorf("ожидался заголовок %q, получено %q", tt.Title, allTasks[i].Title)
		}
		if allTasks[i].Completed != tt.Completed {
			t.Errorf("ошибка статуса: ожидалось %v, получено %v", tt.Completed, allTasks[i].Completed)
		}
	}

	// Проверка фильтра important
	var importantTasks []task.Task
	for _, t := range allTasks {
		if t.Important {
			importantTasks = append(importantTasks, t)
		}
	}
	if len(importantTasks) != 1 || importantTasks[0].ID != 1 {
		t.Errorf("фильтр important не работает, получили %+v", importantTasks)
	}

	// Проверка фильтра completed
	var completedTasks []task.Task
	for _, t := range allTasks {
		if t.Completed {
			completedTasks = append(completedTasks, t)
		}
	}
	if len(completedTasks) != 1 || completedTasks[0].ID != 3 {
		t.Errorf("фильтр completed не работает, получили %+v", completedTasks)
	}

	// Проверка сортировки по имени
	nameSorted := make([]task.Task, len(allTasks))
	copy(nameSorted, allTasks)
	c := collate.New(language.Russian)
	sort.Slice(allTasks, func(i, j int) bool {
		return c.CompareString(nameSorted[i].Title, nameSorted[j].Title) < 0
	})
	if nameSorted[0].Title != "Первая" || nameSorted[2].Title != "Третья" {
		t.Errorf("сортировка по имени некорректна: %+v", nameSorted)
	}

	// Проверка сортировки по дате
	dateSorted := make([]task.Task, len(allTasks))
	copy(dateSorted, allTasks)
	sort.Slice(dateSorted, func(i, j int) bool {
		return dateSorted[i].CreatedAt.Before(dateSorted[j].CreatedAt)
	})
	if dateSorted[0].ID != 1 || dateSorted[2].ID != 3 {
		t.Errorf("сортировка по дате некорректна: %+v", dateSorted)
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
		Important: false,
		CreatedAt: time.Now(),
	}

	if err := store.AddTask(initialTask); err != nil {
		t.Fatalf("AddTask вернул ошибку: %v", err)
	}

	// Обновляем title и important
	if err := store.UpdateTask(1, "Новое название", true); err != nil {
		t.Fatalf("UpdateTask вернул ошибку: %v", err)
	}

	// Загружаем обновлённый список.
	// Проверяем, что название изменилось.
	tasks, _ := store.ListTasks()
	if tasks[0].Title != "Новое название" || !tasks[0].Important {
		t.Errorf("обновление задачи не сработало: %+v", tasks[0])
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
		Important: false,
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

// --- Тест ListTasks() на несуществующем файле ---
func TestJSONStore_ReadNonExist(t *testing.T) {
	store := NewJSONStore("/non/exist/file.json")
	tasks, err := store.ListTasks()
	if err != nil || len(tasks) != 0 {
		t.Errorf("ожидалось пустой список без ошибки, получили: %v", err)
	}
}

// --- Тест UpdateTask() на несуществующем файле ---
func TestJSONStore_UpdateNonExist(t *testing.T) {
	store := NewJSONStore("/non/exist/file.json")
	err := store.UpdateTask(999, "Title", false)
	if err == nil {
		t.Errorf("ожидалось, что обновление несуществующей задачи вернет ошибку")
	}
}

// --- Тест DeleteTask() на несуществующем файле ---
func TestJSONStore_DeleteNonExist(t *testing.T) {
	store := NewJSONStore("/non/exist/file.json")
	err := store.DeleteTask(999)
	if err == nil {
		t.Errorf("ожидалось, что удаление несуществующей задачи вернет ошибку")
	}
}
