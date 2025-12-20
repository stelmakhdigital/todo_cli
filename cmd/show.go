package cmd

import (
	"fmt"
	"strconv"
	"todo_cli/internal/manager"
	"todo_cli/internal/storage"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [ID задачи]",
	Short: "Просмотр детальной информации о задаче по её ID",
	Long: `Отображает полную информацию о задаче: заголовок, описание, статус, дату создания и завершения.

Для просмотра необходимо передать ID задачи в качестве аргумента.

Примеры:
  todo show 12
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
		err = mgr.Show(idTask)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
