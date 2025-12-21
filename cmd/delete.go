package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [ID задачи]",
	Short: "Удаление задачи по её ID",
	Long: `Полностью удаляет задачу из списка.

Внимание: операция необратима! Удалённую задачу невозможно восстановить.
Для удаления необходимо передать ID задачи.

Примеры:
  todo delete 8
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Print("не передан ID задачи\n")
			return
		}
		idTask, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("не верное значение для ID задачи: %v\n", args[0])
			return
		}
		err = mgr.Delete(idTask)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Printf("задача с #%d успешно удалена\n", idTask)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
