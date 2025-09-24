package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd — это корневая команда CLI.
// К ней будут добавляться все подкоманды (например, add, list, done).
var rootCmd = &cobra.Command{
	Use:   "todo",                              // имя исполняемой команды
	Short: "ToDo CLI — простой менеджер задач", // краткое описание
	Long: `Todo CLI — это минималистичный менеджер задач.
Позволяет добавлять, просматривать, отмечать и удалять задачи прямо из терминала.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Используйте подкоманды, например: todo add \"купить хлеб\"")
	},
}

// Execute — функция, которая запускает корневую команду.
// Если возникнет ошибка, приложение завершится с кодом 1.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
