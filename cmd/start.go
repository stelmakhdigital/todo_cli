package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [ID задачи]",
	Short: "Начать выполнение задачи (установить статус 'in_progress')",
	Long: `Переводит задачу в статус "in_progress" (в работе).

Используйте эту команду, когда начинаете работать над задачей.
Для выполнения команды необходимо передать ID задачи.

Примеры:
  todo start 15
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
		err = mgr.Start(idTask)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Printf("Задача #%d переведена в статус 'in_progress' \n", idTask)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
