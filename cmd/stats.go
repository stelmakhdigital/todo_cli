package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Просмотр статистики по задачам",
	Long: `Отображает статистику по всем задачам в разбивке по статусам.

Показывает общее количество задач и количество задач для каждого статуса:
pending (ожидает), in_progress (в работе), completed (выполнено).

Примеры:
  todo stats
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := mgr.Stats()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
