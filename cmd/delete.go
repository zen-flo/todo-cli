package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zen-flo/todo-cli/internal/storage"
)

// deleteCmd — подкоманда "delete", которая удаляет задачу по ID.
// Пример использования:
//
//	todo delete 2  — удалит задачу с ID 2
var deleteCmd = &cobra.Command{
	Use:   "delete [task ID]",     // формат вызова
	Short: "Удалить задачу по ID", // краткое описание
	Args:  cobra.ExactArgs(1),     // ожидаем ровно один аргумент — ID задачи
	Run: func(cmd *cobra.Command, args []string) {
		// Конвертируем аргумент в int (ID задачи)
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Некорректный ID задачи:", args[0])
			return
		}

		// Создаём хранилище задач
		store := storage.NewJSONStore("tasks.json")

		// Удаляем задачу через публичный метод
		err = store.DeleteTask(id)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		// Подтверждаем успешное удаление
		fmt.Printf("Задача с ID %d успешно удалена.\n", id)
	},
}

// init автоматически вызывается при старте приложения.
// Здесь мы подключаем подкоманду "delete" к rootCmd.
func init() {
	rootCmd.AddCommand(deleteCmd)
}
