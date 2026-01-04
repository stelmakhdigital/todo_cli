//go:build !production

package manager

import (
	"testing"
	"todo_cli/internal/task"
	"todo_cli/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIndexByID(t *testing.T) {
	tasksEmpty := testutil.EmptyTasks()
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err, "не должно быть ошибки при создании задач")

	tests := []struct {
		name     string
		tasks    []*task.Task
		indexID  int
		expected *int
	}{
		{"пустой список задач", tasksEmpty, 1, nil},
		{"задача с ID=1 на индексе 0", tasksMany, 1, testutil.IntPtr(0)},
		{"задача с ID=2 на индексе 1", tasksMany, 2, testutil.IntPtr(1)},
		{"несуществующий ID", tasksMany, 99, nil},
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
	tasksEmpty := testutil.EmptyTasks()
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err, "не должно быть ошибки при создании задач")

	tests := []struct {
		name          string
		tasks         []*task.Task
		status        task.Status
		expectedCount int
		expectedErr   error
	}{
		{"пустой список задач", tasksEmpty, task.StatusPending, 0, nil},
		{"фильтр по статусу pending - 2 задачи", tasksMany, task.StatusPending, 2, nil},
		{"фильтр по статусу in_progress - 1 задача", tasksMany, task.StatusProgress, 1, nil},
		{"фильтр по статусу completed - 3 задачи", tasksMany, task.StatusCompleted, 3, nil},
		{"невалидный статус - ошибка", tasksMany, task.Status("invalid_status"), 0, task.ErrInvalidStatus},
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

func TestGetTasksBySearchWord(t *testing.T) {
	tasksEmpty := testutil.EmptyTasks()
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err, "не должно быть ошибки при создании задачи")
	tests := []struct {
		name          string
		tasks         []*task.Task
		searchWord    string
		expectedCount int
	}{
		{"пустой список задач", tasksEmpty, "someone", 0},
		{"поиск по фразе pending", tasksMany, "pending", 2},
		{"поиск по фразе TaSk", tasksMany, "TaSk", 6},
		{"поиск по фразе not_found", tasksMany, "not_found", 0},
		{"поиск по фразе scri", tasksMany, "scri", 6},
	}
	filter := &FilterTasks{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filter.GetTasksBySearchWord(tt.tasks, tt.searchWord)
			assert.Len(t, result, tt.expectedCount)
		})
	}
}

func TestGetStatsTasksByStatus(t *testing.T) {
	tasksEmpty := testutil.EmptyTasks()
	tasksSingle, err := testutil.SingleTask()
	require.NoError(t, err, "не должно быть ошибки при создании задачи")
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err, "не должно быть ошибки при создании задачи")
	tests := []struct {
		name          string
		tasks         []*task.Task
		expectedStats map[string]interface{}
	}{
		{"пустой список задач", tasksEmpty, testutil.StatsTask(0, 0, 0, 0)},
		{"одна задача в статусе pending", tasksSingle, testutil.StatsTask(1, 0, 0, 1)},
		{"несколько задач с разными статусами (2-pending; 1-progress; 3-completed)", tasksMany, testutil.StatsTask(6, 3, 1, 2)},
	}
	filter := &FilterTasks{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filter.GetStatsTasksByStatus(tt.tasks)
			assert.Equal(t, tt.expectedStats, result)
		})
	}
}
