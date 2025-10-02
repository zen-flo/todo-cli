package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zen-flo/todo-cli/internal/storage"
)

// pendingCmd — подкоманда "pending", которая выводит только невыполненные задачи.
// Пример использования:
//
//	todo pending
var pendingCmd = &cobra.Command{
	Use:   "pending",                              // формат вызова
	Short: "Показать только невыполненные задачи", // краткое описание
	Run: func(cmd *cobra.Command, args []string) {
		// Создаём хранилище задач
		store := storage.NewJSONStore(tasksFile)

		// Получаем список всех задач
		tasks, err := store.ListTasks()
		if err != nil {
			fmt.Println("Ошибка при загрузке задач:", err)
			return
		}

		// Выводим все невыполненные задачи
		fmt.Println("Невыполненные задачи:")
		for _, t := range tasks {
			if !t.Completed {
				fmt.Printf("[%d] %s\n", t.ID, t.Title)
			}
		}
	},
}

// init автоматически вызывается при старте приложения.
// Здесь мы подключаем подкоманду "delete" к rootCmd.
func init() {
	rootCmd.AddCommand(pendingCmd)
}
