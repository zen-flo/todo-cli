package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zen-flo/todo-cli/internal/storage"
	"github.com/zen-flo/todo-cli/internal/task"
	"os"
	"sort"
	"time"
)

// formatStatus возвращает цветной символ статуса задачи.
// ✅ зелёный — выполнено, ❌ красный — невыполнено
func formatStatus(completed bool) string {
	if completed {
		return "\033[32m✅\033[0m"
	}
	return "\033[31m❌\033[0m"
}

// formatTaskTitle форматирует название задачи, добавляет значки и подсветку.
// 🔥 — важная задача, ⏰ — просроченная (старше 7 дней и не выполнена)
func formatTaskTitle(t task.Task) string {
	title := t.Title

	if t.Important {
		title = "🔥 " + title
	}

	if !t.Completed && time.Since(t.CreatedAt) > 7*24*time.Hour {
		// Просрочена — жёлтый цвет
		title = "\033[33m" + title + "\033[0m"
	}

	return title
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

		if len(tasks) == 0 {
			fmt.Println("Список задач пуст. Добавьте новую с помощью: todo add \"Название задачи\"")
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

		// Вывод заголовков таблицы
		fmt.Printf("%-4s %-7s %-25s %s\n", "ID", "Status", "Title", "CreatedAt")
		fmt.Println("------------------------------------------------------------")

		for _, t := range tasks {
			fmt.Printf("%-4d %-7s %-30s %s\n",
				t.ID,
				formatStatus(t.Completed),
				formatTaskTitle(t),
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
