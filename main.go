package main

import "github.com/zen-flo/todo-cli/cmd"

// main — точка входа в приложение.
// Здесь мы просто вызываем cmd.Execute(), который подхватывает корневую команду cobra.
func main() {
	cmd.Execute()
}
