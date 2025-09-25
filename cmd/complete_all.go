package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zen-flo/todo-cli/internal/storage"
)

// completeAllCmd — подкоманда "complete-all", которая отмечает
// все задачи как выполненные.
// Пример использования:
//
//	todo complete-all
var completeAllCmd = &cobra.Command{
	Use:   "complete-all",                        // формат вызова
	Short: "Отметить все задачи как выполненные", // краткое описание
	Run: func(cmd *cobra.Command, args []string) {
		// Создаём хранилище задач
		store := storage.NewJSONStore("tasks.json")

		// Загружаем все задачи
		tasks, err := store.ListTasks()
		if err != nil {
			fmt.Println("Ошибка при загрузке задач:", err)
			return
		}

		// Обновляем статус всех задач
		updated := 0
		for i := range tasks {
			if !tasks[i].Completed {
				tasks[i].Completed = true
				updated++
			}
		}

		// Сохраняем изменения
		if err := store.OverwriteTasks(tasks); err != nil {
			fmt.Println("Ошибка при обновлении задач:", err)
			return
		}

		if updated > 0 {
			fmt.Printf("Отмечено как выполненные задач: %d\n", updated)
		} else {
			fmt.Println("Все задачи уже выполнены.")
		}
	},
}

// init автоматически вызывается при старте приложения.
// Здесь мы подключаем подкоманду "complete-all" к rootCmd.
func init() {
	rootCmd.AddCommand(completeAllCmd)
}
