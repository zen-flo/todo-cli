package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zen-flo/todo-cli/internal/storage"
)

// doneCmd — подкоманда "done", которая отмечает задачу как выполненную.
// Пример использования:
//
//	todo done 2  — пометит задачу с ID 2 как выполненную
var doneCmd = &cobra.Command{
	Use:   "done [task ID]",                  // формат вызова
	Short: "Отметить задачу как выполненную", // краткое описание
	Args:  cobra.ExactArgs(1),                // ожидаем ровно один аргумент — ID задачи
	Run: func(cmd *cobra.Command, args []string) {
		// Конвертируем аргумент в int (ID задачи)
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Некорректный ID задачи:", args[0])
			return
		}

		// Создаём хранилище задач
		store := storage.NewJSONStore(tasksFile)

		// Отмечаем задачу как выполненную
		err = store.MarkTaskDone(id)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		// Подтверждаем успешное выполнение
		fmt.Printf("Задача с ID %d отмечена как выполненная.\n", id)
	},
}

// init автоматически вызывается при старте приложения.
// Здесь мы подключаем подкоманду "done" к rootCmd.
func init() {
	doneCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		store := storage.NewJSONStore(tasksFile)
		tasks, err := store.ListTasks()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var suggestions []string
		for _, t := range tasks {
			if !t.Completed { // показываем только невыполненные
				suggestions = append(suggestions, fmt.Sprint(t.ID))
			}
		}
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	}

	rootCmd.AddCommand(doneCmd)
}
