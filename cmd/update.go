package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zen-flo/todo-cli/internal/storage"
)

// updateCmd — подкоманда "update", которая изменяет название задачи по ID.
// Пример использования:
//
//	todo update 2 "Новое название задачи"
var updateCmd = &cobra.Command{
	Use:   "update [task ID] [new title]",   // формат вызова
	Short: "Изменить название задачи по ID", // краткое описание
	Args:  cobra.ExactArgs(2),               // ожидаем ровно два аргумента: ID и новый заголовок
	Run: func(cmd *cobra.Command, args []string) {
		// Конвертируем аргумент в int (ID задачи)
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Некорректный ID задачи:", args[0])
			return
		}

		newTitle := args[1]

		// Считываем флаг, что задача важная
		important, _ := cmd.Flags().GetBool("important")

		// Создаём хранилище задач
		store := storage.NewJSONStore("tasks.json")

		// Обновляем название задачи через публичный метод UpdateTask
		err = store.UpdateTask(id, newTitle, important)
		if err != nil {
			fmt.Println("Ошибка при обновлении задачи:", err)
			return
		}

		// Подтверждаем успешное обновление
		fmt.Printf("Задача с ID %d успешно обновлена.\n", id)
	},
}

// init автоматически вызывается при старте приложения.
// Здесь мы подключаем подкоманду "update" к rootCmd.
func init() {
	updateCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		store := storage.NewJSONStore("tasks.json")
		tasks, err := store.ListTasks()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var suggestions []string
		for _, t := range tasks {
			suggestions = append(suggestions, fmt.Sprint(t.ID))
		}
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	}

	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolP("important", "i", false, "Сделать задачу важной")
}
