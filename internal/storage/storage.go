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

type FileStorage struct{}

func checkExistsFile(fileName string) error {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		emptyTasks := []byte("[]")
		return os.WriteFile(fileName, emptyTasks, fileMode644)
	}
	return nil
}

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

func (fs *FileStorage) Clear(differentFileName *string) error {
	choiceNameFile := fileName
	if differentFileName != nil {
		choiceNameFile = *differentFileName
	}
	return os.Remove(choiceNameFile)
}
