package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd показывает базовую команду (тут только описание тк не указан Run) если команда передана без аргументов
var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "CLI приложение для управления задачами",
	Long: `

Простое приложения для создания и управления задачами.

Для справки вызовите:

todo -h`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.AddTemplateFunc("tr", translate)
	rootCmd.SetUsageTemplate(usageTemplate)
}

func translate(s string) string {
	translations := map[string]string{
		"Usage:":             "Использование:",
		"Available Commands": "Доступные команды",
		"Flags:":             "Флаги:",
		"help for":           "справка для",
		"default":            "по-умолчанию",
	}
	if t, ok := translations[s]; ok {
		return t
	}
	return s
}

const usageTemplate = `Использование:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [команда]{{end}}{{if gt (len .Aliases) 0}}

Псевдонимы:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Примеры:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Доступные команды:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Флаги:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Глобальные флаги:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Дополнительные команды помощи:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Для дополнительной информации используйте "{{.CommandPath}} [команда] --help"{{end}}
`
