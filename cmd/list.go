package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zen-flo/todo-cli/internal/storage"
	"os"
	"sort"
)

// formatStatus возвращает цветной символ статуса задачи.
// ✅ зелёный — выполнено, ❌ красный — невыполнено
func formatStatus(completed bool) string {
	if completed {
		return "\033[32m✅\033[0m"
	}
	return "\033[31m❌\033[0m"
}

// listCmd — подкоманда "list", которая выводит все задачи.
// Поддерживает флаг --sort=name/date
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

		// Получаем флаг сортировки
		sortBy, _ := cmd.Flags().GetString("sort")

		// Сортировка задач
		switch sortBy {
		case "name":
			sort.Slice(tasks, func(i, j int) bool {
				return tasks[i].Title < tasks[j].Title
			})
		case "date":
			sort.Slice(tasks, func(i, j int) bool {
				return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
			})
		case "":
			// без сортировки
		default:
			fmt.Println("Неизвестный параметр сортировки. Используйте name или date.")
			os.Exit(1)
		}

		// Выводим все задачи с цветным статусом и выравниванием
		fmt.Println("Список задач:")
		for _, t := range tasks {
			fmt.Printf("%-4d %s %-20s %s\n",
				t.ID,
				formatStatus(t.Completed), // без %-2s
				t.Title,
				t.CreatedAt.Format("2006-01-02 15:04"),
			)
		}
	},
}

// init вызывается автоматически при старте программы.
// Здесь мы подключаем подкоманду "list" к rootCmd.
func init() {
	rootCmd.AddCommand(listCmd)
	// Флаг сортировки
	listCmd.Flags().StringP("sort", "s", "", "Сортировка: name или date")
}
