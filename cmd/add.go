package cmd

import (
	"fmt"
	"todo_cli/internal/manager"
	"todo_cli/internal/storage"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [заголовок] [описание]",
	Short: "Создание новой задачи",
	Long: `Создаёт новую задачу с указанным заголовком и опциональным описанием.

Заголовок является обязательным аргументом, описание — опциональным.
После создания задачу можно будет отредактировать командой edit.

Примеры:
  todo add "Купить продукты"
  todo add "Написать отчёт" "Подготовить отчёт для руководства"
`,
	Args: cobra.MatchAll(
		cobra.MinimumNArgs(1),
		cobra.MaximumNArgs(2),
	),
	Run: func(cmd *cobra.Command, args []string) {
		store := &storage.FileStorage{}
		filter := &manager.FilterTasks{}
		mgr := manager.NewManager(store, filter)
		if len(args[0]) == 0 {
			fmt.Print("укажите корректные данные для заголовка или описания задачи\n")
			return
		}
		title := args[0]
		data := make(map[string]string, 2)
		data["title"] = title
		if len(args[0]) > 1 {
			description := args[1]
			data["description"] = description
		}
		idTask, err := mgr.Create(data)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Printf("Задача #%d добавлена успешно\n", idTask)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
