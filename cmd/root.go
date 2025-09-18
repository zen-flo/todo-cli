package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "ToDo CLI — простой менеджер задач",
	Long:  "ToDo CLI позволяет добавлять, отмечgit add .\nать и удалять задачи прямо из терминала.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Используйте подкоманды, например: todo add \"купить хлеб\"")
	},
}

// Execute запускает CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
