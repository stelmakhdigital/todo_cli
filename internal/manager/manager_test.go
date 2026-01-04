//go:build !production

package manager

import (
	"errors"
	"testing"
	"todo_cli/internal/task"
	"todo_cli/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Save(tasks []*task.Task, newFileName *string) error {
	args := m.Called(tasks, newFileName)
	return args.Error(0)
}

func (m *MockStorage) Load(differentFileName *string) ([]*task.Task, error) {
	args := m.Called(differentFileName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*task.Task), args.Error(1)
}

func (m *MockStorage) Clear(differentFileName *string) error {
	args := m.Called(differentFileName)
	return args.Error(0)
}

type MockFilter struct {
	mock.Mock
}

func (m *MockFilter) GetIndexByID(tasks []*task.Task, id int) *int {
	args := m.Called(tasks, id)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*int)
}

func (m *MockFilter) GetTasksByStatus(tasks []*task.Task, status task.Status) ([]*task.Task, error) {
	args := m.Called(tasks, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*task.Task), args.Error(1)
}

func (m *MockFilter) GetTasksBySearchWord(tasks []*task.Task, word string) []*task.Task {
	args := m.Called(tasks, word)
	return args.Get(0).([]*task.Task)
}

func (m *MockFilter) GetStatsTasksByStatus(tasks []*task.Task) map[string]interface{} {
	args := m.Called(tasks)
	return args.Get(0).(map[string]interface{})
}

type MockRender struct {
	mock.Mock
}

func (m *MockRender) RenderDetailed(t *task.Task) {
	m.Called(t)
}

func (m *MockRender) RenderList(tasks []*task.Task) {
	m.Called(tasks)
}

func (m *MockRender) RenderMap(data map[string]interface{}) {
	m.Called(data)
}

func TestCreate(t *testing.T) {
	tasksEmpty := testutil.EmptyTasks()
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err)

	tests := []struct {
		name        string
		data        map[string]string
		loadTasks   []*task.Task
		loadErr     error
		saveErr     error
		expectedID  *int
		expectedErr bool
	}{
		{"создание первой задачи", map[string]string{"title": "New Task", "description": "New Desc"}, tasksEmpty, nil, nil, testutil.IntPtr(1), false},
		{"создание задачи когда уже есть задачи", map[string]string{"title": "Task 7", "description": "Desc 7"}, tasksMany, nil, nil, testutil.IntPtr(7), false},
		{"ошибка при загрузке", map[string]string{"title": "Task", "description": "Desc"}, nil, errors.New("load error"), nil, nil, true},
		{"нет title в данных", map[string]string{"description": "Desc"}, tasksEmpty, nil, nil, nil, true},
		{"нет description в данных", map[string]string{"title": "Task"}, tasksEmpty, nil, nil, nil, true},
		{"пустые данные", map[string]string{}, tasksEmpty, nil, nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			mockFilter := new(MockFilter)
			mockRender := new(MockRender)

			mockStorage.On("Load", mock.Anything).Return(tt.loadTasks, tt.loadErr)
			if !tt.expectedErr && tt.loadErr == nil && tt.data["title"] != "" && tt.data["description"] != "" {
				mockStorage.On("Save", mock.Anything, mock.Anything).Return(tt.saveErr)
				mockRender.On("RenderDetailed", mock.Anything).Return()
			}

			manager := NewManager(mockStorage, mockFilter, mockRender)
			result, err := manager.Create(tt.data)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, result)
			}
		})
	}
}

func TestStart(t *testing.T) {
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err)

	tests := []struct {
		name        string
		taskID      int
		loadTasks   []*task.Task
		loadErr     error
		indexResult *int
		saveErr     error
		expectedErr bool
	}{
		{"успешный старт задачи #1", 1, tasksMany, nil, testutil.IntPtr(0), nil, false},
		{"задача не найдена", 99, tasksMany, nil, nil, nil, true},
		{"ошибка при загрузке", 1, nil, errors.New("load error"), nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			mockFilter := new(MockFilter)
			mockRender := new(MockRender)

			mockStorage.On("Load", mock.Anything).Return(tt.loadTasks, tt.loadErr)
			mockFilter.On("GetIndexByID", mock.Anything, tt.taskID).Return(tt.indexResult)
			if tt.indexResult != nil && tt.loadErr == nil {
				mockStorage.On("Save", mock.Anything, mock.Anything).Return(tt.saveErr)
				mockRender.On("RenderDetailed", mock.Anything).Return()
			}

			manager := NewManager(mockStorage, mockFilter, mockRender)
			err := manager.Start(tt.taskID)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestComplete(t *testing.T) {
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err)

	tests := []struct {
		name        string
		taskID      int
		loadTasks   []*task.Task
		loadErr     error
		indexResult *int
		saveErr     error
		expectedErr bool
	}{
		{"успешное завершение задачи #1", 1, tasksMany, nil, testutil.IntPtr(0), nil, false},
		{"задача не найдена", 99, tasksMany, nil, nil, nil, true},
		{"ошибка при загрузке", 1, nil, errors.New("load error"), nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			mockFilter := new(MockFilter)
			mockRender := new(MockRender)

			mockStorage.On("Load", mock.Anything).Return(tt.loadTasks, tt.loadErr)
			mockFilter.On("GetIndexByID", mock.Anything, tt.taskID).Return(tt.indexResult)
			if tt.indexResult != nil && tt.loadErr == nil {
				mockStorage.On("Save", mock.Anything, mock.Anything).Return(tt.saveErr)
				mockRender.On("RenderDetailed", mock.Anything).Return()
			}

			manager := NewManager(mockStorage, mockFilter, mockRender)
			err := manager.Complete(tt.taskID)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEdit(t *testing.T) {
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err)

	tests := []struct {
		name        string
		taskID      int
		data        map[string]string
		loadTasks   []*task.Task
		loadErr     error
		indexResult *int
		saveErr     error
		expectedErr bool
	}{
		{"редактирование title", 1, map[string]string{"title": "Updated"}, tasksMany, nil, testutil.IntPtr(0), nil, false},
		{"редактирование description", 1, map[string]string{"description": "New desc"}, tasksMany, nil, testutil.IntPtr(0), nil, false},
		{"редактирование status", 1, map[string]string{"status": "completed"}, tasksMany, nil, testutil.IntPtr(0), nil, false},
		{"некорректный status", 1, map[string]string{"status": "invalid"}, tasksMany, nil, testutil.IntPtr(0), nil, true},
		{"задача не найдена", 99, map[string]string{"title": "Test"}, tasksMany, nil, nil, nil, true},
		{"ошибка при загрузке", 1, map[string]string{"title": "Test"}, nil, errors.New("load error"), nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			mockFilter := new(MockFilter)
			mockRender := new(MockRender)

			mockStorage.On("Load", mock.Anything).Return(tt.loadTasks, tt.loadErr)
			mockFilter.On("GetIndexByID", mock.Anything, tt.taskID).Return(tt.indexResult)
			if tt.indexResult != nil && tt.loadErr == nil && tt.data["status"] != "invalid" {
				mockStorage.On("Save", mock.Anything, mock.Anything).Return(tt.saveErr)
				mockRender.On("RenderDetailed", mock.Anything).Return()
			}

			manager := NewManager(mockStorage, mockFilter, mockRender)
			err := manager.Edit(tt.taskID, tt.data)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err)

	tests := []struct {
		name        string
		taskID      int
		loadTasks   []*task.Task
		loadErr     error
		indexResult *int
		saveErr     error
		expectedErr bool
	}{
		{"успешное удаление задачи #1", 1, tasksMany, nil, testutil.IntPtr(0), nil, false},
		{"задача не найдена", 99, tasksMany, nil, nil, nil, true},
		{"ошибка при загрузке", 1, nil, errors.New("load error"), nil, nil, true},
		{"ошибка при сохранении", 1, tasksMany, nil, testutil.IntPtr(0), errors.New("save error"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			mockFilter := new(MockFilter)
			mockRender := new(MockRender)

			mockStorage.On("Load", mock.Anything).Return(tt.loadTasks, tt.loadErr)
			mockFilter.On("GetIndexByID", mock.Anything, tt.taskID).Return(tt.indexResult)
			if tt.indexResult != nil && tt.loadErr == nil {
				mockStorage.On("Save", mock.Anything, mock.Anything).Return(tt.saveErr)
			}

			manager := NewManager(mockStorage, mockFilter, mockRender)
			err := manager.Delete(tt.taskID)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestShow(t *testing.T) {
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err)

	tests := []struct {
		name        string
		taskID      int
		loadTasks   []*task.Task
		loadErr     error
		indexResult *int
		expectedErr bool
	}{
		{"успешный показ задачи #1", 1, tasksMany, nil, testutil.IntPtr(0), false},
		{"задача не найдена", 99, tasksMany, nil, nil, true},
		{"ошибка при загрузке", 1, nil, errors.New("load error"), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			mockFilter := new(MockFilter)
			mockRender := new(MockRender)

			mockStorage.On("Load", mock.Anything).Return(tt.loadTasks, tt.loadErr)
			mockFilter.On("GetIndexByID", mock.Anything, tt.taskID).Return(tt.indexResult)
			if tt.indexResult != nil {
				mockRender.On("RenderDetailed", mock.Anything).Return()
			}

			manager := NewManager(mockStorage, mockFilter, mockRender)
			err := manager.Show(tt.taskID)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestList(t *testing.T) {
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err)

	tests := []struct {
		name          string
		status        string
		loadTasks     []*task.Task
		loadErr       error
		filteredTasks []*task.Task
		filterErr     error
		expectedErr   bool
	}{
		{"список всех задач", "all", tasksMany, nil, nil, nil, false},
		{"фильтр по pending", "pending", tasksMany, nil, tasksMany[:2], nil, false},
		{"невалидный статус", "invalid", tasksMany, nil, nil, nil, true},
		{"ошибка при загрузке", "all", nil, errors.New("load error"), nil, nil, true},
		{"ошибка при фильтрации", "pending", tasksMany, nil, nil, errors.New("filter error"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			mockFilter := new(MockFilter)
			mockRender := new(MockRender)

			mockStorage.On("Load", mock.Anything).Return(tt.loadTasks, tt.loadErr)
			if tt.status != "all" && tt.status != "invalid" && tt.loadErr == nil {
				mockFilter.On("GetTasksByStatus", mock.Anything, task.Status(tt.status)).Return(tt.filteredTasks, tt.filterErr)
			}
			if tt.loadErr == nil && !tt.expectedErr {
				mockRender.On("RenderList", mock.Anything).Return()
			}

			manager := NewManager(mockStorage, mockFilter, mockRender)
			err := manager.List(tt.status)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStats(t *testing.T) {
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err)
	stats := testutil.StatsTask(6, 3, 1, 2)

	tests := []struct {
		name        string
		loadTasks   []*task.Task
		loadErr     error
		stats       map[string]interface{}
		expectedErr bool
	}{
		{"успешная статистика", tasksMany, nil, stats, false},
		{"ошибка при загрузке", nil, errors.New("load error"), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			mockFilter := new(MockFilter)
			mockRender := new(MockRender)

			mockStorage.On("Load", mock.Anything).Return(tt.loadTasks, tt.loadErr)
			if tt.loadErr == nil {
				mockFilter.On("GetStatsTasksByStatus", mock.Anything).Return(tt.stats)
				mockRender.On("RenderMap", tt.stats).Return()
			}

			manager := NewManager(mockStorage, mockFilter, mockRender)
			err := manager.Stats()

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err)

	tests := []struct {
		name        string
		word        string
		loadTasks   []*task.Task
		loadErr     error
		foundTasks  []*task.Task
		expectedErr bool
	}{
		{"найдены задачи", "pending", tasksMany, nil, tasksMany[:2], false},
		{"задачи не найдены", "notfound", tasksMany, nil, []*task.Task{}, true},
		{"ошибка при загрузке", "test", nil, errors.New("load error"), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			mockFilter := new(MockFilter)
			mockRender := new(MockRender)

			mockStorage.On("Load", mock.Anything).Return(tt.loadTasks, tt.loadErr)
			if tt.loadErr == nil {
				mockFilter.On("GetTasksBySearchWord", mock.Anything, tt.word).Return(tt.foundTasks)
				if len(tt.foundTasks) > 0 {
					mockRender.On("RenderList", tt.foundTasks).Return()
				}
			}

			manager := NewManager(mockStorage, mockFilter, mockRender)
			err := manager.Search(tt.word)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHasKeys(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]string
		keys     []string
		expected bool
	}{
		{"все ключи есть", map[string]string{"a": "1", "b": "2"}, []string{"a", "b"}, true},
		{"один ключ отсутствует", map[string]string{"a": "1"}, []string{"a", "b"}, false},
		{"пустая карта", map[string]string{}, []string{"a"}, false},
		{"пустой список ключей", map[string]string{"a": "1"}, []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasKeys(tt.data, tt.keys...)
			assert.Equal(t, tt.expected, result)
		})
	}
}
