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
	Short: "Просмотр всех задач или списка задач по статусу",
	Long: `

Показывает список всех задач.
Если передать в аргумент статус (допускается pending|in_progress|completed) то выведет только задачи с этим статусом

Пример:

todo list --status completed
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
