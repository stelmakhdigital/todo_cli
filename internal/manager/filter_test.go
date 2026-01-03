package manager

import (
	"testing"
	"todo_cli/internal/task"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func intPtr(v int) *int {
	return &v
}

func fixtureTwoTasks() ([]*task.Task, error) {
	tasks := []*task.Task{}
	newTaskOne, err := task.NewTask(1, "test task 1", "test description", task.StatusPending.String())
	if err != nil {
		return nil, err
	}
	tasks = append(tasks, newTaskOne)
	newTaskTwo, err := task.NewTask(2, "test task 2", "test description", task.StatusProgress.String())
	if err != nil {
		return nil, err
	}
	tasks = append(tasks, newTaskTwo)

	return tasks, nil
}

func fixtureMixedStatusTasks() ([]*task.Task, error) {
	// 2 задачи StatusPending
	// 1 задача StatusProgress
	// 3 задачи StatusCompleted
	tasks := []*task.Task{}

	t1, err := task.NewTask(1, "pending task 1", "description", task.StatusPending.String())
	if err != nil {
		return nil, err
	}
	tasks = append(tasks, t1)

	t2, err := task.NewTask(2, "pending task 2", "description", task.StatusPending.String())
	if err != nil {
		return nil, err
	}
	tasks = append(tasks, t2)

	t3, err := task.NewTask(3, "progress task", "description", task.StatusProgress.String())
	if err != nil {
		return nil, err
	}
	tasks = append(tasks, t3)

	t4, err := task.NewTask(4, "completed task 1", "description", task.StatusCompleted.String())
	if err != nil {
		return nil, err
	}
	tasks = append(tasks, t4)

	t5, err := task.NewTask(5, "completed task 2", "description", task.StatusCompleted.String())
	if err != nil {
		return nil, err
	}
	tasks = append(tasks, t5)

	t6, err := task.NewTask(6, "completed task 3", "description", task.StatusCompleted.String())
	if err != nil {
		return nil, err
	}
	tasks = append(tasks, t6)

	return tasks, nil
}

func TestGetIndexByID(t *testing.T) {
	tasksEmpty := []*task.Task{}
	tasksTwo, err := fixtureTwoTasks()
	if err != nil {
		t.Errorf("ошибка во время содания задач. %v", err)
	}
	tests := []struct {
		name     string
		tasks    []*task.Task
		indexID  int
		expected *int
	}{
		{"пустой список задач", tasksEmpty, 1, nil},
		{"задача с ID=1 на индексе 0", tasksTwo, 1, intPtr(0)},
		{"задача с ID=2 на индексе 1", tasksTwo, 2, intPtr(1)},
		{"несуществующий ID", tasksTwo, 99, nil},
	}
	filter := &FilterTasks{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filter.GetIndexByID(tt.tasks, tt.indexID)
			assert.Equal(t, tt.expected, result, tt.name)
		})
	}
}

func TestGetTasksByStatus(t *testing.T) {
	tasksEmpty := []*task.Task{}
	tasksMixed, err := fixtureMixedStatusTasks()
	require.NoError(t, err, "не должно быть ошибки при создании задач")

	tests := []struct {
		name          string
		tasks         []*task.Task
		status        task.Status
		expectedCount int
		expectedErr   error
	}{
		{
			name:          "пустой список задач",
			tasks:         tasksEmpty,
			status:        task.StatusPending,
			expectedCount: 0,
			expectedErr:   nil,
		},
		{
			name:          "фильтр по статусу pending - 2 задачи",
			tasks:         tasksMixed,
			status:        task.StatusPending,
			expectedCount: 2,
			expectedErr:   nil,
		},
		{
			name:          "фильтр по статусу in_progress - 1 задача",
			tasks:         tasksMixed,
			status:        task.StatusProgress,
			expectedCount: 1,
			expectedErr:   nil,
		},
		{
			name:          "фильтр по статусу completed - 3 задачи",
			tasks:         tasksMixed,
			status:        task.StatusCompleted,
			expectedCount: 3,
			expectedErr:   nil,
		},
		{
			name:          "невалидный статус - ошибка",
			tasks:         tasksMixed,
			status:        task.Status("invalid_status"),
			expectedCount: 0,
			expectedErr:   task.ErrInvalidStatus,
		},
	}

	filter := &FilterTasks{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := filter.GetTasksByStatus(tt.tasks, tt.status)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, tt.expectedCount)

				for _, tsk := range result {
					assert.Equal(t, tt.status, tsk.Status)
				}
			}
		})
	}
}
