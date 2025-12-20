package cmd

import (
	"fmt"
	"todo_cli/internal/manager"
	"todo_cli/internal/storage"

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
		store := &storage.FileStorage{}
		filter := &manager.FilterTasks{}
		mgr := manager.NewManager(store, filter)
		err := mgr.Stats()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
