//go:build !production

package storage

import (
	"os"
	"path/filepath"
	"testing"
	"todo_cli/internal/task"
	"todo_cli/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataToJson(t *testing.T) {
	tasksEmpty := testutil.EmptyTasks()
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err)

	tests := []struct {
		name        string
		data        interface{}
		expectedErr bool
	}{
		{"пустой список задач", &tasksEmpty, false},
		{"список из нескольких задач", &tasksMany, false},
		{"простая структура", &map[string]string{"key": "value"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result []byte
			var err error
			switch v := tt.data.(type) {
			case *[]*task.Task:
				result, err = DataToJson(v)
			case *map[string]string:
				result, err = DataToJson(v)
			}
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Greater(t, len(result), 0)
			}
		})
	}
}

func TestJsonToData(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		expectedErr bool
	}{
		{"валидный JSON задач", `[{"id":1,"title":"Test","description":"Desc","status":"pending"}]`, false},
		{"пустой массив", `[]`, false},
		{"невалидный JSON", `{invalid json}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tasks []*task.Task
			err := JsonToData([]byte(tt.jsonData), &tasks)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileStorage_Save(t *testing.T) {
	tmpDir := t.TempDir()
	tasksEmpty := testutil.EmptyTasks()
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err)

	tests := []struct {
		name     string
		tasks    []*task.Task
		fileName *string
	}{
		{"сохранение пустого списка", tasksEmpty, testutil.StrPtr(filepath.Join(tmpDir, "empty.json"))},
		{"сохранение нескольких задач", tasksMany, testutil.StrPtr(filepath.Join(tmpDir, "many.json"))},
		{"сохранение с fileName=nil", tasksMany, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FileStorage{}

			// для nil в fileName меняем рабочую директорию
			if tt.fileName == nil {
				oldDir, _ := os.Getwd()
				defer os.Chdir(oldDir)
				os.Chdir(tmpDir)
			}

			err := fs.Save(tt.tasks, tt.fileName)
			assert.NoError(t, err)

			// чекаем что файл создан
			checkFile := fileName
			if tt.fileName != nil {
				checkFile = *tt.fileName
			}
			_, err = os.Stat(checkFile)
			assert.NoError(t, err)

			// очистка
			if tt.fileName != nil {
				os.Remove(*tt.fileName)
			}
		})
	}
}

func TestFileStorage_Load(t *testing.T) {
	tmpDir := t.TempDir()
	tasksMany, err := testutil.ManyTasks()
	require.NoError(t, err)

	// создаём тестовый файл
	testFile := filepath.Join(tmpDir, "test.json")
	fs := &FileStorage{}
	err = fs.Save(tasksMany, &testFile)
	require.NoError(t, err)

	tests := []struct {
		name          string
		fileName      *string
		expectedCount int
		expectedErr   bool
	}{
		{"загрузка существующего файла", &testFile, 6, false},
		{"загрузка несуществующего файла", testutil.StrPtr(filepath.Join(tmpDir, "new.json")), 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fs.Load(tt.fileName)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}
		})
	}

	// очистка
	os.Remove(testFile)
}

func TestFileStorage_Clear(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "clear.json")

	// создаем файл для удаления
	fs := &FileStorage{}
	err := fs.Save(testutil.EmptyTasks(), &testFile)
	require.NoError(t, err)

	tests := []struct {
		name        string
		fileName    *string
		expectedErr bool
	}{
		{"удаление существующего файла", &testFile, false},
		{"удаление несуществующего файла", testutil.StrPtr(filepath.Join(tmpDir, "notexist.json")), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fs.Clear(tt.fileName)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// чекаем что файл удалён
				_, err := os.Stat(*tt.fileName)
				assert.True(t, os.IsNotExist(err))
			}
		})
	}
}

func TestCheckExistsFile(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "existing.json")
	newFile := filepath.Join(tmpDir, "new.json")

	// создание существующего файл
	err := os.WriteFile(existingFile, []byte("test"), 0644)
	require.NoError(t, err)

	tests := []struct {
		name     string
		fileName string
	}{
		{"существующий файл", existingFile},
		{"новый файл", newFile},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkExistsFile(tt.fileName)
			assert.NoError(t, err)

			// проверка что файл существует
			_, err = os.Stat(tt.fileName)
			assert.NoError(t, err)
		})
	}
	// очистка
	os.Remove(existingFile)
	os.Remove(newFile)
}
