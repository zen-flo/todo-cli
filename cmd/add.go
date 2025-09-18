package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zen-flo/todo-cli/internal/storage"
	"github.com/zen-flo/todo-cli/internal/task"
	"time"
)

var addCmd = &cobra.Command{
	Use:   "add [task title]",
	Short: "Добавить новую задачу",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		store := storage.NewJSONStore()
		newTask := task.Task{
			ID:        1, // пока фиксируем ID
			Title:     args[0],
			Completed: false,
			CreatedAt: time.Now(),
		}

		err := store.AddTask(newTask)
		if err != nil {
			fmt.Println("Ошибка при добавлении задачи:", err)
			return
		}

		fmt.Println("Добавлена задача:", newTask.Title)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
