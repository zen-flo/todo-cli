package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zen-flo/todo-cli/internal/storage"
	"github.com/zen-flo/todo-cli/internal/task"
	"time"
)

// addCmd — подкоманда "add", которая создаёт новую задачу.
// Пример использования:
//
//	todo add "Купить хлеб"
var addCmd = &cobra.Command{
	Use:   "add [task title]",      // формат вызова
	Short: "Добавить новую задачу", // краткое описание
	Args:  cobra.ExactArgs(1),      // ожидаем ровно один аргумент — название задачи
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Ошибка: нужно указать заголовок задачи.")
			return
		}
		// Создаём новое хранилище задач.
		// Указываем путь к файлу.
		store := storage.NewJSONStore(tasksFile)

		// Считываем флаг, что задача важная
		important, _ := cmd.Flags().GetBool("important")

		// Формируем новую задачу.
		// ID сейчас фиксированный (1) — это временное решение,
		// позже JSONStore будет генерировать уникальные ID.
		newTask := task.Task{
			// ID присваивается автоматически внутри AddTask
			Title:     args[0],
			Completed: false,
			Important: important,
			CreatedAt: time.Now(),
		}

		// Пытаемся добавить задачу в хранилище.
		err := store.AddTask(newTask)
		if err != nil {
			fmt.Println("Ошибка при добавлении задачи:", err)
			return
		}

		// Если всё ок — выводим сообщение пользователю.
		fmt.Println("Добавлена задача:", newTask.Title)
	},
}

// init — автоматически вызывается при запуске.
// Здесь мы подключаем подкоманду add к rootCmd.
func init() {
	rootCmd.AddCommand(addCmd)

	// Флаг важности
	addCmd.Flags().BoolP("important", "i", false, "Отметить задачу как важную")

	// Для bool-флага автодополнение пустое, чтобы не ломать shell completion
	_ = addCmd.RegisterFlagCompletionFunc("important", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"true", "false"}, cobra.ShellCompDirectiveNoFileComp
	})
}
