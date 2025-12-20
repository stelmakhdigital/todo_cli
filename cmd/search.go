package cmd

import (
	"fmt"
	"todo_cli/internal/manager"
	"todo_cli/internal/storage"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Поиск по задачам",
	Long: `

Ищет и выводит найденные задачи по вхождению слова или фразу по заголовку или описанию задачи.

Пример:

todo search "сделать"
`,
	Args: cobra.MatchAll(
		cobra.MinimumNArgs(1), // минимум 1 аргумент
	),
	Run: func(cmd *cobra.Command, args []string) {
		store := &storage.FileStorage{}
		filter := &manager.FilterTasks{}
		mgr := manager.NewManager(store, filter)
		err := mgr.Search(args[0])
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
