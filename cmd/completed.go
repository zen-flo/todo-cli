package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zen-flo/todo-cli/internal/storage"
)

// completedCmd — подкоманда "completed", которая выводит только выполненные задачи.
// Пример использования:
//
//	todo completed
var completedCmd = &cobra.Command{
	Use:   "completed",                          // формат вызова
	Short: "Показать только выполненные задачи", // краткое описание
	Run: func(cmd *cobra.Command, args []string) {
		// Создаём хранилище задач
		store := storage.NewJSONStore("tasks.json")

		// Получаем список всех задач
		tasks, err := store.ListTasks()
		if err != nil {
			fmt.Println("Ошибка при загрузке задач:", err)
			return
		}

		// Выводим все выполненные задачи
		fmt.Println("Выполненные задачи:")
		for _, t := range tasks {
			if t.Completed {
				fmt.Printf("[%d] %s\n", t.ID, t.Title)
			}
		}
	},
}

// init автоматически вызывается при старте приложения.
// Здесь мы подключаем подкоманду "delete" к rootCmd.
func init() {
	rootCmd.AddCommand(completedCmd)
}
