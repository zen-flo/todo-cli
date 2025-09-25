package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zen-flo/todo-cli/internal/storage"
)

// searchCmd — подкоманда "search", которая ищет задачи по ключевому слову
// в их названии (Title). Поиск нечувствителен к регистру.
// Пример использования:
//
//	todo search хлеб
var searchCmd = &cobra.Command{
	Use:   "search [keyword]",                // формат вызова
	Short: "Найти задачи по ключевому слову", // краткое описание
	Args:  cobra.ExactArgs(1),                // ожидаем ровно один аргумент — слово для поиска
	Run: func(cmd *cobra.Command, args []string) {
		keyword := strings.ToLower(args[0]) // приводим к нижнему регистру для нечувствительного поиска

		// Создаём хранилище задач
		store := storage.NewJSONStore("tasks.json")

		// Загружаем все задачи
		tasks, err := store.ListTasks()
		if err != nil {
			fmt.Println("Ошибка при загрузке задач:", err)
			return
		}

		// Фильтруем задачи по ключевому слову
		fmt.Printf("Результаты поиска по \"%s\":\n", args[0])
		found := false
		for _, t := range tasks {
			if strings.Contains(strings.ToLower(t.Title), keyword) {
				status := "❌"
				if t.Completed {
					status = "✅"
				}
				fmt.Printf("[%s] %d: %s\n", status, t.ID, t.Title)
				found = true
			}
		}

		// Если задач не найдено — выводим сообщение
		if !found {
			fmt.Println("Задачи не найдены.")
		}
	},
}

// init автоматически вызывается при старте приложения.
// Здесь мы подключаем подкоманду "search" к rootCmd.
func init() {
	rootCmd.AddCommand(searchCmd)
}
