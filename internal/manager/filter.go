package manager

import "todo_cli/internal/task"

type Filter interface {
	GetIndexByID(tasks []*task.Task, id int) *int
	// Search(tasks []*task.Task, word string) ([]*task.Task, error)
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
