package manager

import (
	"fmt"
	"strings"
	"todo_cli/internal/storage"
	"todo_cli/internal/task"
	"unicode/utf8"
)

type ManagerTasks interface {
	GetAll() ([]*task.Task, error)
	// Show
	// Create
	// Start
	// Complete
	// Delete
	Render(tasks []*task.Task)
}

func GetAll() ([]*task.Task, error) {
	var tasks []*task.Task
	stringToTasks, err := storage.Load()
	if err != nil {
		return nil, fmt.Errorf("менеджеру не удалось загрузить задачи: %w", err)
	}
	err = storage.JsonToData(stringToTasks, &tasks)
	if err != nil {
		return nil, fmt.Errorf("менеджеру не удалось преобразовать задачи: %w", err)
	}
	return tasks, nil
}

func Create() {

}

func Render(tasks []*task.Task) {
	var columnMax int = 0

	for _, value := range tasks {
		if columnMax < utf8.RuneCountInString(value.Title) {
			columnMax = utf8.RuneCountInString(value.Title)
		}
	}
	columnMax += 5
	fmt.Printf("%-4s | %-*s | %-12s | %-15s\n", "ID", columnMax, "Название", "Статус", "Создана")
	fmt.Println(strings.Repeat("-", columnMax+40))

	for _, task := range tasks {
		fmt.Printf("%-4d | %-*s | %-12s | %-15s\n",
			task.ID, columnMax, task.Title, task.Status,
			task.CreatedAt.Format("02.01.2006"))

	}
}
