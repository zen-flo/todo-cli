package cmd

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/zen-flo/todo-cli/internal/storage"
	"github.com/zen-flo/todo-cli/internal/task"
)

// captureOutput ловит вывод в stdout
// captureOutput ловит вывод в stdout
func captureOutput(f func()) string {
	var buf bytes.Buffer

	// Сохраняем старый stdout
	old := os.Stdout

	// Создаём pipe
	r, w, err := os.Pipe()
	if err != nil {
		panic("не удалось создать pipe: " + err.Error())
	}

	// Перенаправляем stdout
	os.Stdout = w

	// Выполняем функцию
	f()

	// Закрываем writer
	if err := w.Close(); err != nil {
		panic("не удалось закрыть writer: " + err.Error())
	}

	// Восстанавливаем stdout
	os.Stdout = old

	// Читаем всё из reader
	if _, err := buf.ReadFrom(r); err != nil {
		panic("не удалось прочитать из pipe: " + err.Error())
	}

	return buf.String()
}

// helper: создает временный JSONStore и подменяет tasksFile
func withTempStore(t *testing.T, f func(store *storage.JSONStore, tmpFile string)) {
	tmpFile, err := os.CreateTemp("", "tasks_*.json")
	if err != nil {
		t.Fatalf("не удалось создать временный файл: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("не удалось удалить временный файл: %v", err)
		}
	}()

	orig := tasksFile
	tasksFile = tmpFile.Name()
	defer func() { tasksFile = orig }()

	store := storage.NewJSONStore(tasksFile)
	f(store, tmpFile.Name())
}

// --- Тест команды addCmd ---
func TestAddCommand(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		if err := addCmd.Flags().Set("important", "true"); err != nil {
			t.Fatalf("не удалось установить флаг: %v", err)
		}

		captureOutput(func() {
			addCmd.Run(addCmd, []string{"Test Task"})
		})

		tasks, _ := store.ListTasks()
		if len(tasks) != 1 {
			t.Fatalf("ожидалось 1 задача, получено %d", len(tasks))
		}
		if tasks[0].Title != "Test Task" || !tasks[0].Important {
			t.Errorf("неверные данные задачи: %+v", tasks[0])
		}
	})
}

// --- Тест addCmd без аргументов ---
func TestAddCommand_NoArgs(t *testing.T) {
	output := captureOutput(func() {
		addCmd.Run(addCmd, []string{})
	})
	if !bytes.Contains([]byte(output), []byte("нужно указать заголовок задачи")) {
		t.Errorf("ожидалось сообщение об ошибке при отсутствии аргументов")
	}
}

// --- Тест команды doneCmd ---
func TestDoneCommand(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		err := store.AddTask(task.Task{
			ID:        1,
			Title:     "Task 1",
			Completed: true,
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось выполнить AddTask: %v", err)
		}

		captureOutput(func() {
			doneCmd.Run(doneCmd, []string{"1"})
		})

		tasks, _ := store.ListTasks()
		if !tasks[0].Completed {
			t.Errorf("задача должна быть отмечена как выполненная")
		}
	})
}

// --- Тест команды deleteCmd ---
func TestDeleteCommand(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		err := store.AddTask(task.Task{
			ID:        1,
			Title:     "Task 1",
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось выполнить AddTask: %v", err)
		}
		captureOutput(func() {
			deleteCmd.Run(deleteCmd, []string{"1"})
		})

		tasks, _ := store.ListTasks()
		if len(tasks) != 0 {
			t.Errorf("задача должна быть удалена")
		}
	})
}

// --- Тест на неправильный ввод для deleteCmd ---
func TestDeleteCommand_InvalidID(t *testing.T) {
	output := captureOutput(func() {
		deleteCmd.Run(deleteCmd, []string{"abc"})
	})
	if !bytes.Contains([]byte(output), []byte("Некорректный ID")) {
		t.Errorf("ожидалось сообщение об ошибке для некорректного ID")
	}
}

// --- Тест команды updateCmd ---
func TestUpdateCommand(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		err := store.AddTask(task.Task{
			ID:        1,
			Title:     "Old Title",
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось выполнить AddTask: %v", err)
		}
		if err := updateCmd.Flags().Set("important", "true"); err != nil {
			t.Fatalf("не удалось установить флаг: %v", err)
		}

		captureOutput(func() {
			updateCmd.Run(updateCmd, []string{"1", "New Title"})
		})

		tasks, _ := store.ListTasks()
		if tasks[0].Title != "New Title" || !tasks[0].Important {
			t.Errorf("обновление задачи не сработало: %+v", tasks[0])
		}
	})
}

// --- Тест updateCmd без аргументов ---
func TestUpdateCommand_MissingArgs(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		output := captureOutput(func() {
			// вызываем Run с пустым args
			updateCmd.Run(updateCmd, []string{})
		})

		if !bytes.Contains([]byte(output), []byte("Ошибка: нужно указать ID и новое название задачи")) {
			t.Errorf("ожидалось сообщение о нехватке аргументов, получено: %s", output)
		}
	})
}

// --- Тест команды listCmd ---
func TestListCommand(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		err := store.AddTask(task.Task{
			ID:        1,
			Title:     "Task A",
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось выполнить AddTask: %v", err)
		}
		err = store.AddTask(task.Task{
			ID:        2,
			Title:     "Task B",
			Completed: true,
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось выполнить AddTask: %v", err)
		}

		captureOutput(func() {
			if err := listCmd.Flags().Set("filter", "completed"); err != nil {
				t.Fatalf("не удалось установить флаг: %v", err)
			}
			listCmd.Run(listCmd, []string{})
		})
	})
}

// --- Тест на ошибочный фильтр/сортировку ---
func TestListCommand_InvalidSort(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		if err := listCmd.Flags().Set("sort", "bad"); err != nil {
			t.Fatalf("не удалось установить флаг: %v", err)
		}

		output := captureOutput(func() {
			listCmd.Run(listCmd, []string{})
		})

		if !bytes.Contains([]byte(output), []byte("Ошибка: неизвестный способ сортировки")) {
			t.Errorf("ожидалось сообщение об ошибке сортировки, получено: %s", output)
		}
	})
}

// --- Тест команды pendingCmd ---
func TestPendingCommand(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		err := store.AddTask(task.Task{
			ID:        1,
			Title:     "Task 1",
			Completed: false,
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось выполнить AddTask: %v", err)
		}
		err = store.AddTask(task.Task{
			ID:        2,
			Title:     "Task 2",
			Completed: true,
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось выполнить AddTask: %v", err)
		}

		output := captureOutput(func() {
			pendingCmd.Run(pendingCmd, []string{})
		})

		if !bytes.Contains([]byte(output), []byte("Task 1")) {
			t.Errorf("ожидалось, что будет выведена только невыполненная задача")
		}
		if bytes.Contains([]byte(output), []byte("Task 2")) {
			t.Errorf("невыполненная задача должна быть исключена")
		}
	})
}

// --- Тест команды completedCmd ---
func TestCompletedCommand(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		err := store.AddTask(task.Task{
			ID:        1,
			Title:     "Task 1",
			Completed: true,
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось выполнить AddTask: %v", err)
		}
		err = store.AddTask(task.Task{
			ID:        2,
			Title:     "Task 2",
			Completed: false,
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось выполнить AddTask: %v", err)
		}

		output := captureOutput(func() {
			completedCmd.Run(completedCmd, []string{})
		})

		if !bytes.Contains([]byte(output), []byte("Task 1")) {
			t.Errorf("ожидалось, что будет выведена выполненная задача")
		}
		if bytes.Contains([]byte(output), []byte("Task 2")) {
			t.Errorf("невыполненная задача не должна выводиться")
		}
	})
}

// --- Тест команды clearCmd ---
func TestClearCommand(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		err := store.AddTask(task.Task{
			ID:        1,
			Title:     "Task 1",
			Completed: true,
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось выполнить AddTask: %v", err)
		}
		err = store.AddTask(task.Task{
			ID:        2,
			Title:     "Task 2",
			Completed: false,
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось выполнить AddTask: %v", err)
		}

		captureOutput(func() {
			clearCmd.Run(clearCmd, []string{})
		})

		tasks, _ := store.ListTasks()
		if len(tasks) != 1 || tasks[0].ID != 2 {
			t.Errorf("завершённая задача не была удалена")
		}
	})
}

// --- Тест clearCmd, когда нечего удалять ---
func TestClearCommand_EmptyList(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		output := captureOutput(func() {
			clearCmd.Run(clearCmd, []string{})
		})
		if !bytes.Contains([]byte(output), []byte("Нет завершённых задач")) {
			t.Errorf("ожидалось сообщение при пустом списке")
		}
	})
}

// --- Тест команды completeAllCmd ---
func TestCompleteAllCommand(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		err := store.AddTask(task.Task{
			ID:        1,
			Title:     "Task 1",
			Completed: false,
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось выполнить AddTask: %v", err)
		}
		err = store.AddTask(task.Task{
			ID:        2,
			Title:     "Task 2",
			Completed: false,
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось выполнить AddTask: %v", err)
		}

		captureOutput(func() {
			completeAllCmd.Run(completeAllCmd, []string{})
		})

		tasks, _ := store.ListTasks()
		for _, tsk := range tasks {
			if !tsk.Completed {
				t.Errorf("задача %d должна быть отмечена как выполненная", tsk.ID)
			}
		}
	})
}

// --- Тест команды searchCmd ---
func TestSearchCommand(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		err := store.AddTask(task.Task{
			ID:        1,
			Title:     "Buy Milk",
			Completed: false,
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось выполнить AddTask: %v", err)
		}
		err = store.AddTask(task.Task{
			ID:        2,
			Title:     "Read Book",
			Completed: false,
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось выполнить AddTask: %v", err)
		}

		output := captureOutput(func() {
			searchCmd.Run(searchCmd, []string{"milk"})
		})

		if !bytes.Contains([]byte(output), []byte("Buy Milk")) {
			t.Errorf("задача с ключевым словом 'milk' должна быть найдена")
		}
		if bytes.Contains([]byte(output), []byte("Read Book")) {
			t.Errorf("задача без ключевого слова не должна выводиться")
		}
	})
}
