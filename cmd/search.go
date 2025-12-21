package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [слово или фраза]",
	Short: "Поиск задач по ключевому слову или фразе",
	Long: `Выполняет поиск по заголовкам и описаниям задач.

Поиск регистронезависимый и ищет вхождение указанного текста в любой части заголовка или описания.

Примеры:
  todo search "отчёт"
`,
	Args: cobra.MatchAll(
		cobra.MinimumNArgs(1), // минимум 1 аргумент
	),
	Run: func(cmd *cobra.Command, args []string) {
		err := mgr.Search(args[0])
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
