package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zen-flo/todo-cli/internal/storage"
)

// listCmd — подкоманда "list", которая выводит все задачи.
// Пример использования:
//
//	todo list
var listCmd = &cobra.Command{
	Use:   "list",                // формат вызова
	Short: "Показать все задачи", // краткое описание
	Run: func(cmd *cobra.Command, args []string) {
		// Создаём хранилище задач
		store := storage.NewJSONStore("tasks.json")

		// Получаем список всех задач
		tasks, err := store.ListTasks()
		if err != nil {
			fmt.Println("Ошибка при загрузке задач:", err)
			return
		}

		// Если список пуст — выводим сообщение
		if len(tasks) == 0 {
			fmt.Println("Список задач пуст.")
			return
		}

		// Выводим все задачи с их статусом (выполнена/не выполнена)
		for _, t := range tasks {
			status := "❌"
			if t.Completed {
				status = "✅"
			}
			fmt.Printf("[%s] %d: %s\n", status, t.ID, t.Title)
		}
	},
}

// init вызывается автоматически при старте программы.
// Здесь мы подключаем подкоманду "list" к rootCmd.
func init() {
	rootCmd.AddCommand(listCmd)
}
