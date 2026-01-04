//go:build !production

package testutil

import "todo_cli/internal/task"

func EmptyTasks() []*task.Task {
	return []*task.Task{}
}

func SingleTask() ([]*task.Task, error) {
	tasks := []*task.Task{}
	task, err := task.NewTask(1, "single task", "description", task.StatusPending.String())
	if err != nil {
		return nil, err
	}
	tasks = append(tasks, task)
	return tasks, nil
}

func ManyTasks() ([]*task.Task, error) {
	// 2 задачи StatusPending
	// 1 задача StatusProgress
	// 3 задачи StatusCompleted
	tasks := []*task.Task{}
	samples := []struct {
		id          int
		title       string
		description string
		status      string
	}{
		{1, "pending task 1", "description", task.StatusPending.String()},
		{2, "pending task 2", "description", task.StatusPending.String()},
		{3, "progress task", "description", task.StatusProgress.String()},
		{4, "completed task 1", "description", task.StatusCompleted.String()},
		{5, "completed task 2", "description", task.StatusCompleted.String()},
		{6, "completed task 3", "description", task.StatusCompleted.String()},
	}
	for _, values := range samples {
		task, err := task.NewTask(values.id, values.title, values.description, values.status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func StatsTask(all, completed, progress, pending int) map[string]interface{} {
	return map[string]interface{}{
		"Всего задач:": all,
		"Выполнено":    completed,
		"В работе":     progress,
		"Ожидает":      pending,
	}
}
