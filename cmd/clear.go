package cmd

import (
	"fmt"
	"github.com/zen-flo/todo-cli/internal/task"

	"github.com/spf13/cobra"
	"github.com/zen-flo/todo-cli/internal/storage"
)

// clearCmd — подкоманда "clear", которая удаляет все завершённые задачи
// из списка. Используется для очистки списка задач от "мусора".
// Пример использования:
//
//	todo clear
var clearCmd = &cobra.Command{
	Use:   "clear",                          // формат вызова
	Short: "Удалить все завершённые задачи", // краткое описание
	Run: func(cmd *cobra.Command, args []string) {
		// Создаём хранилище задач
		store := storage.NewJSONStore("tasks.json")

		// Загружаем список задач
		tasks, err := store.ListTasks()
		if err != nil {
			fmt.Println("Ошибка при загрузке задач:", err)
			return
		}

		// Фильтруем только незавершённые
		active := make([]task.Task, 0)
		cleared := 0
		for _, t := range tasks {
			if !t.Completed {
				active = append(active, t)
			} else {
				cleared++
			}
		}

		// Перезаписываем список
		if err := store.OverwriteTasks(active); err != nil {
			fmt.Println("Ошибка при очистке завершённых задач:", err)
			return
		}

		// Сообщение пользователю
		if cleared > 0 {
			fmt.Printf("Удалено завершённых задач: %d\n", cleared)
		} else {
			fmt.Println("Нет завершённых задач для удаления.")
		}
	},
}

// init автоматически вызывается при старте приложения.
// Здесь мы подключаем подкоманду "clear" к rootCmd.
func init() {
	rootCmd.AddCommand(clearCmd)
}
