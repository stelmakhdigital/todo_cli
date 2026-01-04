package storage

import (
	"fmt"
	"os"
	"todo_cli/internal/task"
)

const fileName = "tasks.json"

const (
	permWrite = 2 << iota // 2
	permRead              // 4
)

const (
	owner = (permRead | permWrite) << 6 // 0600
	group = permRead << 3               // 0040
	other = permRead                    // 0004
)

const (
	fileMode644 = owner | group | other // 0644
)

// FileStorage реализует интерфейс Storage для работы с файловой системой.
// Сохраняет и загружает задачи в формате JSON.
type FileStorage struct{}

// checkExistsFile проверяет существование файла.
// Если файл не существует, создаёт его с пустым массивом задач "[]".
// Возвращает ошибку при проблемах с созданием файла.
func checkExistsFile(fileName string) error {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		emptyTasks := []byte("[]")
		return os.WriteFile(fileName, emptyTasks, fileMode644)
	}
	return nil
}

// Save сохраняет список задач в JSON-файл.
// Если newFileName = nil, использует дефолтное имя файла "tasks.json".
// Если newFileName указан, сохраняет в файл с указанным именем.
// Автоматически создаёт файл, если он не существует.
// Возвращает ошибку при проблемах с сериализацией или записью файла.
func (fs *FileStorage) Save(tasks []*task.Task, newFileName *string) error {
	choiceNameFile := fileName
	if newFileName != nil {
		choiceNameFile = *newFileName
	}
	err := checkExistsFile(choiceNameFile)
	if err != nil {
		return err
	}
	taskToString, err := DataToJson(&tasks)
	if err != nil {
		return fmt.Errorf("ошибка при преобразовании задачи: %w", err)
	}
	return os.WriteFile(choiceNameFile, taskToString, fileMode644)
}

// Load загружает список задач из JSON-файла.
// Если differentFileName = nil, загружает из дефолтного файла "tasks.json".
// Если differentFileName указан, загружает из файла с указанным именем.
// Автоматически создаёт файл с пустым списком, если он не существует.
// Возвращает список задач или ошибку при проблемах с чтением или десериализацией.
func (fs *FileStorage) Load(differentFileName *string) ([]*task.Task, error) {
	var tasks []*task.Task
	choiceNameFile := fileName
	if differentFileName != nil {
		choiceNameFile = *differentFileName
	}
	err := checkExistsFile(choiceNameFile)
	if err != nil {
		return nil, err
	}
	stringToTasks, err := os.ReadFile(choiceNameFile)
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить задачи: %w", err)
	}
	err = JsonToData(stringToTasks, &tasks)
	if err != nil {
		return nil, fmt.Errorf("не удалось преобразовать задачи: %w", err)
	}
	return tasks, nil
}

// Clear удаляет файл с задачами.
// Если differentFileName = nil, удаляет дефолтный файл "tasks.json".
// Если differentFileName указан, удаляет файл с указанным именем.
// Возвращает ошибку, если файл не существует или не может быть удалён.
func (fs *FileStorage) Clear(differentFileName *string) error {
	choiceNameFile := fileName
	if differentFileName != nil {
		choiceNameFile = *differentFileName
	}
	return os.Remove(choiceNameFile)
}
