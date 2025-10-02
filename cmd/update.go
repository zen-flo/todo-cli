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
	Args: func(cmd *cobra.Command, args []string) error { // ожидаем ровно два аргумента: ID и новый заголовок
		if len(args) < 2 {
			fmt.Println("Ошибка: нужно указать ID и новое название задачи.")
			return fmt.Errorf("недостаточно аргументов")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Проверяем количество аргументов прямо в Run
		if len(args) < 2 {
			fmt.Println("Ошибка: нужно указать ID и новое название задачи.")
			return
		}
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
		store := storage.NewJSONStore(tasksFile)

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
	rootCmd.AddCommand(updateCmd)

	// Флаг важности
	updateCmd.Flags().BoolP("important", "i", false, "Сделать задачу важной")

	// Автодополнение для аргументов
	updateCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if tasksFile == "" {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		store := storage.NewJSONStore(tasksFile)
		tasks, err := store.ListTasks()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		// Если последний аргумент --important, предлагаем true/false
		if len(args) > 0 && args[len(args)-1] == "--important" {
			return []string{"true", "false"}, cobra.ShellCompDirectiveNoFileComp
		}

		// Показываем ID всех задач
		var suggestions []string
		for _, t := range tasks {
			suggestions = append(suggestions, fmt.Sprint(t.ID))
		}
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	}
}
