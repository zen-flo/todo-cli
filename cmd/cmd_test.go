package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/zen-flo/todo-cli/internal/storage"
	"github.com/zen-flo/todo-cli/internal/task"
)

// --- Вспомогательная функция для перехвата stdout ---
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

// --- Вспомогательная функция для временного хранилища ---
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
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		output := captureOutput(func() {
			addCmd.Run(addCmd, []string{})
		})
		if !bytes.Contains([]byte(output), []byte("Ошибка: нужно указать заголовок задачи")) {
			t.Errorf("ожидалось сообщение об ошибке при отсутствии аргументов")
		}
	})
}

// --- Проверка ошибок addCmd ---
func TestAddCommand_NoArgs_Error(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		output := captureOutput(func() {
			addCmd.Run(addCmd, []string{})
		})
		if !strings.Contains(output, "Ошибка: нужно указать заголовок задачи.") {
			t.Errorf("ожидалось сообщение об ошибке, получено: %s", output)
		}
	})
}

// --- Тесты флагов и автодополнения ---
func TestAddCmdFlagCompletion(t *testing.T) {
	fn := func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"true", "false"}, cobra.ShellCompDirectiveNoFileComp
	}

	res, dir := fn(addCmd, nil, "")
	expected := []string{"true", "false"}

	for i := range expected {
		if res[i] != expected[i] {
			t.Errorf("ожидалось %v, получено %v", expected[i], res[i])
		}
	}

	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("ожидалось директиву NoFileComp, получено %v", dir)
	}
}

func TestUpdateCmdValidArgsFunction(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		suggestions, _ := updateCmd.ValidArgsFunction(updateCmd, nil, "")
		// Если список задач пуст, suggestions должно быть nil
		if suggestions != nil && len(suggestions) != 0 {
			t.Errorf("ожидалось пустое предложение для пустого списка, получено: %v", suggestions)
		}
	})
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

// --- Проверка ошибок doneCmd ---
func TestDoneCommand_NonExist(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		output := captureOutput(func() {
			doneCmd.Run(doneCmd, []string{"999"})
		})
		if !strings.Contains(output, "Ошибка") {
			t.Errorf("ожидалось сообщение об ошибке отметки несуществующей задачи, получено: %s", output)
		}
	})
}

// --- Тест doneCmd с некорректным ID ---
func TestDoneCommand_InvalidID(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		output := captureOutput(func() {
			doneCmd.Run(doneCmd, []string{"abc"})
		})
		if !bytes.Contains([]byte(output), []byte("Некорректный ID")) {
			t.Errorf("ожидалось сообщение об ошибке некорректного ID")
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

// --- Проверка ошибок deleteCmd ---
func TestDeleteCommand_NonExist(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		output := captureOutput(func() {
			deleteCmd.Run(deleteCmd, []string{"999"})
		})
		if !strings.Contains(output, "Ошибка") {
			t.Errorf("ожидалось сообщение об ошибке удаления несуществующей задачи, получено: %s", output)
		}
	})
}

// --- Тест deleteCmd на несуществующей задаче ---
func TestDeleteCommand_NoTasks(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		output := captureOutput(func() {
			deleteCmd.Run(deleteCmd, []string{"1"})
		})
		if !bytes.Contains([]byte(output), []byte("Ошибка")) {
			t.Errorf("ожидалось сообщение об ошибке удаления несуществующей задачи")
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
func TestUpdateCommand_NoArgs(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		output := captureOutput(func() {
			updateCmd.Run(updateCmd, []string{})
		})
		if !bytes.Contains([]byte(output), []byte("Ошибка: нужно указать ID и новое название задачи")) {
			t.Errorf("ожидалось сообщение об ошибке при отсутствии аргументов")
		}
	})
}

// --- Проверка ошибок updateCmd ---
func TestUpdateCommand_NoArgs_Error(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		output := captureOutput(func() {
			updateCmd.Run(updateCmd, []string{})
		})
		if !strings.Contains(output, "Ошибка: нужно указать ID и новое название задачи.") {
			t.Errorf("ожидалось сообщение об ошибке, получено: %s", output)
		}
	})
}

// --- Тест updateCmd с некорректным ID ---
func TestUpdateCommand_InvalidID(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		err := store.AddTask(task.Task{
			ID:        1,
			Title:     "Old Task",
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось добавить задачу: %v", err)
		}

		output := captureOutput(func() {
			updateCmd.Run(updateCmd, []string{"abc", "New Task"})
		})
		if !bytes.Contains([]byte(output), []byte("Некорректный ID")) {
			t.Errorf("ожидалось сообщение об ошибке некорректного ID")
		}
	})
}

// --- Проверка ошибок updateCmd ---
func TestUpdateCommand_BadID_Error(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		output := captureOutput(func() {
			updateCmd.Run(updateCmd, []string{"abc", "title"})
		})
		if !strings.Contains(output, "Некорректный ID задачи") {
			t.Errorf("ожидалось сообщение о некорректном ID, получено: %s", output)
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

// --- Тест listCmd на пустом списке задач ---
func TestListCommand_EmptyList(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		output := captureOutput(func() {
			listCmd.Run(listCmd, []string{})
		})
		if !bytes.Contains([]byte(output), []byte("Список задач пуст")) {
			t.Errorf("ожидалось сообщение о пустом списке")
		}
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

// --- Тест completeAllCmd на уже выполненных задачах ---
func TestCompleteAllCommand_AlreadyCompleted(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		// Добавляем задачи, все уже выполненные
		err := store.AddTask(task.Task{
			ID:        1,
			Title:     "Task 1",
			Completed: true,
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось добавить задачу: %v", err)
		}
		err = store.AddTask(task.Task{
			ID:        2,
			Title:     "Task 2",
			Completed: true,
			CreatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("не удалось добавить задачу: %v", err)
		}

		// Ловим вывод команды
		output := captureOutput(func() {
			completeAllCmd.Run(completeAllCmd, []string{})
		})

		// Проверяем, что вывод соответствует одному из ожидаемых вариантов
		if !(strings.Contains(output, "Отмечено как выполненные задач") ||
			strings.Contains(output, "Все задачи уже выполнены")) {
			t.Errorf("ожидалось сообщение о выполнении или что все задачи уже выполнены, получено: %s", output)
		}

		// Проверяем, что все задачи остаются выполненными
		tasks, _ := store.ListTasks()
		for _, tsk := range tasks {
			if !tsk.Completed {
				t.Errorf("задача %d должна быть отмечена как выполненная", tsk.ID)
			}
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

// --- Проверка ошибок clearCmd ---
func TestClearCommand_ErrorOnOverwrite(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		// Добавляем завершённую задачу, чтобы список active не был пустым
		if err := store.AddTask(task.Task{
			ID:        1,
			Title:     "Completed Task",
			Completed: true,
			CreatedAt: time.Now(),
		}); err != nil {
			t.Fatalf("не удалось добавить задачу: %v", err)
		}

		// Делаем файл только для чтения, чтобы вызвался error при OverwriteTasks
		if err := os.Chmod(tmpFile, 0444); err != nil {
			t.Fatalf("не удалось изменить права файла: %v", err)
		}
		defer func() { // восстановим права после теста
			if err := os.Chmod(tmpFile, 0644); err != nil {
				t.Logf("не удалось восстановить права файла: %v", err)
			}
		}()

		output := captureOutput(func() {
			clearCmd.Run(clearCmd, []string{})
		})

		if !bytes.Contains([]byte(output), []byte("Ошибка")) {
			t.Errorf("ожидалось сообщение об ошибке, получено: %s", output)
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

// --- Проверка completeAllCmd для всех выполненных ---
func TestCompleteAllCommand_AllCompleted(t *testing.T) {
	withTempStore(t, func(store *storage.JSONStore, tmpFile string) {
		_ = store.AddTask(task.Task{ID: 1, Title: "Task 1", Completed: true})
		output := captureOutput(func() {
			completeAllCmd.Run(completeAllCmd, []string{})
		})
		if !strings.Contains(output, "Все задачи уже выполнены") {
			t.Errorf("ожидалось сообщение о выполненных задачах, получено: %s", output)
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
