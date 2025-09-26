package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zen-flo/todo-cli/internal/storage"
	"github.com/zen-flo/todo-cli/internal/task"
	"time"
)

// addCmd — подкоманда "add", которая создаёт новую задачу.
// Пример использования:
//
//	todo add "Купить хлеб"
var addCmd = &cobra.Command{
	Use:   "add [task title]",      // формат вызова
	Short: "Добавить новую задачу", // краткое описание
	Args:  cobra.ExactArgs(1),      // ожидаем ровно один аргумент — название задачи
	Run: func(cmd *cobra.Command, args []string) {
		// Создаём новое хранилище задач.
		// Указываем путь к файлу.
		store := storage.NewJSONStore("tasks.json")

		// Формируем новую задачу.
		// ID сейчас фиксированный (1) — это временное решение,
		// позже JSONStore будет генерировать уникальные ID.
		newTask := task.Task{
			ID:        1, // пока фиксируем ID
			Title:     args[0],
			Completed: false,
			CreatedAt: time.Now(),
		}

		// Пытаемся добавить задачу в хранилище.
		err := store.AddTask(newTask)
		if err != nil {
			fmt.Println("Ошибка при добавлении задачи:", err)
			return
		}

		// Если всё ок — выводим сообщение пользователю.
		fmt.Println("Добавлена задача:", newTask.Title)
	},
}

// init — автоматически вызывается при запуске.
// Здесь мы подключаем подкоманду add к rootCmd.
func init() {
	rootCmd.AddCommand(addCmd)
}
