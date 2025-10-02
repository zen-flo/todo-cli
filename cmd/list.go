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

// printTasksTable выводит задачи в виде таблицы с выравниванием и цветным статусом.
func printTasksTable(tasks []task.Task) {
	// Заголовок таблицы
	fmt.Printf("\033[36m%-4s %-7s %-20s %-16s\033[0m\n", "ID", "STATUS", "TITLE", "CREATED AT")
	fmt.Println("-----------------------------------------------")

	// Строки таблицы
	for _, t := range tasks {
		fmt.Printf("%-4d %-7s %-30s %-16s\n",
			t.ID,
			formatStatus(t.Completed),
			formatTaskTitle(t),
			t.CreatedAt.Format("2006-01-02 15:04"),
		)
	}
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
		// Получаем флаг сортировки и проверяем
		sortBy, _ := cmd.Flags().GetString("sort") //сортировка: name, date
		if sortBy != "" && sortBy != "title" && sortBy != "created" {
			fmt.Printf("Ошибка: неизвестный способ сортировки: %s\n", sortBy)
			return
		}

		// Создаём хранилище задач
		store := storage.NewJSONStore(tasksFile)

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

		// Получаем флаги
		filter, _ := cmd.Flags().GetString("filter")         // фильтр: all, pending, completed
		importantOnly, _ := cmd.Flags().GetBool("important") // фильтр: важные. Да/нет.

		// Фильтрация по статусу
		filtered := make([]task.Task, 0)
		for _, t := range tasks {
			if importantOnly && !t.Important {
				continue
			}
			switch filter {
			case "pending":
				if !t.Completed {
					filtered = append(filtered, t)
				}
			case "completed":
				if t.Completed {
					filtered = append(filtered, t)
				}
			default:
				filtered = append(filtered, t)
			}
		}

		// Сортировка задач
		switch sortBy {
		case "name":
			sort.Slice(filtered, func(i, j int) bool { return filtered[i].Title < filtered[j].Title })
		case "date":
			sort.Slice(filtered, func(i, j int) bool { return filtered[i].CreatedAt.Before(filtered[j].CreatedAt) })
		case "":
		default:
			fmt.Println("Неизвестный параметр сортировки. Используйте name или date.")
			os.Exit(1)
		}

		// Вывод заголовков таблицы
		printTasksTable(filtered)
	},
}

// init вызывается автоматически при старте программы.
// Здесь мы подключаем подкоманду "list" к rootCmd.
func init() {
	rootCmd.AddCommand(listCmd)

	// Флаги
	listCmd.Flags().StringP("sort", "s", "", "Сортировка: name или date")
	listCmd.Flags().StringP("filter", "f", "all", "Фильтр: all, pending, completed")
	listCmd.Flags().BoolP("important", "i", false, "Показать только важные задачи")

	// Автодополнение для флага --sort
	_ = listCmd.RegisterFlagCompletionFunc("sort", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"name", "date"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Автодополнение для флага --filter
	_ = listCmd.RegisterFlagCompletionFunc("filter", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"all", "pending", "completed"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Автодополнение для флага --important
	_ = listCmd.RegisterFlagCompletionFunc("important", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"true", "false"}, cobra.ShellCompDirectiveNoFileComp
	})
}
