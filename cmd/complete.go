package cmd

import (
	"fmt"
	"strconv"
	"todo_cli/internal/manager"
	"todo_cli/internal/storage"

	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete [ID задачи]",
	Short: "Отметить задачу как выполненную (установить статус 'completed')",
	Long: `Переводит задачу в статус "completed" (выполнена) и устанавливает дату завершения.

Используйте эту команду, когда задача полностью завершена.
Для выполнения команды необходимо передать ID задачи.

Примеры:
  todo complete 7
`,
	Run: func(cmd *cobra.Command, args []string) {
		store := &storage.FileStorage{}
		filter := &manager.FilterTasks{}
		mgr := manager.NewManager(store, filter)
		if len(args) < 1 {
			fmt.Print("не передан ID задачи\n")
			return
		}
		idTask, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("не верное значение для ID задачи: %v\n", args[0])
			return
		}
		err = mgr.Complete(idTask)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Printf("Задача #%d завершена\n", idTask)
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
