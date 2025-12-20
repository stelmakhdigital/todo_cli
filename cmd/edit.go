package cmd

import (
	"fmt"
	"strconv"
	"todo_cli/internal/manager"
	"todo_cli/internal/storage"

	"github.com/spf13/cobra"
)

var (
	title, description string
)

var editCmd = &cobra.Command{
	Use:   "edit [ID задачи]",
	Short: "Редактирование заголовка или описания задачи",
	Long: `Изменяет заголовок и/или описание существующей задачи.

Необходимо указать ID задачи и хотя бы один из флагов: --title или --description.
Можно изменить оба поля одновременно.

Примеры:
  todo edit 14 --title "Купить книгу по архитектуре облачных приложений"
  todo edit 5 --description "Новое описание задачи"
  todo edit 7 -t "Новый заголовок" -d "Новое описание"
`,
	Args: cobra.MatchAll(
		cobra.MinimumNArgs(1),
	),
	Run: func(cmd *cobra.Command, args []string) {
		store := &storage.FileStorage{}
		filter := &manager.FilterTasks{}
		mgr := manager.NewManager(store, filter)
		if title == "" && description == "" {
			fmt.Print("укажите значение для изменения заголовка или описания задачи\n")
			return
		}
		data := make(map[string]string, 2)
		if title != "" {
			data["title"] = title
		}
		if description != "" {
			data["description"] = description
		}
		idTask, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("не верное значение для ID задачи: %v\n", args[0])
			return
		}
		err = mgr.Edit(idTask, data)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().StringVarP(&title, "title", "t", "", "Новое название для заголовка задачи")
	editCmd.Flags().StringVarP(&description, "description", "d", "", "Новое описание для задачи")
}
