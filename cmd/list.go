package cmd

import (
	"fmt"
	"todo_cli/internal/manager"
	"todo_cli/internal/storage"

	"github.com/spf13/cobra"
)

var status string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Просмотр всех задач или фильтрация по статусу",
	Long: `Отображает список задач с возможностью фильтрации по статусу.

По умолчанию показывает все задачи. Используйте флаг --status для фильтрации.
Доступные статусы: pending, in_progress, completed.

Примеры:
  todo list
  todo list --status completed
  todo list -s in_progress
`,
	Run: func(cmd *cobra.Command, args []string) {
		store := &storage.FileStorage{}
		filter := &manager.FilterTasks{}
		mgr := manager.NewManager(store, filter)
		err := mgr.List(status)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&status, "status", "s", "all", "Название статуса для фильтра")
}
