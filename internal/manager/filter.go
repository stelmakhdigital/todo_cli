package manager

import (
	"fmt"
	"strings"
	"todo_cli/internal/task"
)

type Filter interface {
	GetIndexByID(tasks []*task.Task, id int) *int
	GetTasksByStatus(tasks []*task.Task, status task.Status) ([]*task.Task, error)
	GetStatsTasksByStatus(tasks []*task.Task) map[string]interface{}
	GetTasksBySearchWord(tasks []*task.Task, word string) []*task.Task
}

type FilterTasks struct{}

func (f *FilterTasks) GetIndexByID(tasks []*task.Task, id int) *int {
	for index, value := range tasks {
		if value.ID == id {
			return &index
		}
	}
	return nil
}

func (f *FilterTasks) GetTasksByStatus(tasks []*task.Task, status task.Status) ([]*task.Task, error) {
	if !task.Status(status).Valid() {
		return nil, fmt.Errorf("ошибка валидации (%w): %s", task.ErrInvalidStatus, status)
	}
	filteredTasks := make([]*task.Task, 0, len(tasks))
	for _, value := range tasks {
		if value.Status == status {
			filteredTasks = append(filteredTasks, value)
		}
	}
	return filteredTasks, nil
}

func (f *FilterTasks) GetStatsTasksByStatus(tasks []*task.Task) map[string]interface{} {
	var allTask, completed, progress, pending int
	allTask = len(tasks)
	for _, value := range tasks {
		if value.Status == task.StatusPending {
			pending += 1
		}
		if value.Status == task.StatusProgress {
			progress += 1
		}
		if value.Status == task.StatusCompleted {
			completed += 1
		}
	}
	data := map[string]interface{}{
		"Всего задач:": allTask,
		"Выполнено":    completed,
		"В работе":     progress,
		"Ожидает":      pending,
	}
	return data
}

func (f *FilterTasks) GetTasksBySearchWord(tasks []*task.Task, word string) []*task.Task {
	foundTasks := make([]*task.Task, 0, len(tasks))
	word = strings.ToLower(word)
	for _, value := range tasks {
		title := strings.ToLower(value.Title)
		description := strings.ToLower(value.Description)
		if strings.Contains(title, word) || strings.Contains(description, word) {
			foundTasks = append(foundTasks, value)
		}
	}
	return foundTasks
}
