package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version переменные заполняются при сборке через -ldflags
	Version   = "dev"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Показать версию приложения",
	Long:  `Отображает версию приложения и дату сборки`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Версия: %s\n", Version)
		fmt.Printf("Дата: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
