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

func (fs *FileStorage) Save(tasks []*task.Task, newFileName *string) error {
	choiceNameFile := fileName
	if newFileName != nil {
		choiceNameFile = *newFileName
	}

	taskToString, err := DataToJson(&tasks)
	if err != nil {
		return fmt.Errorf("ошибка при преобразовании задачи: %w", err)
	}

	return os.WriteFile(choiceNameFile, taskToString, fileMode644)
}

func (fs *FileStorage) Load() ([]*task.Task, error) {
	var tasks []*task.Task

	stringToTasks, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить задачи: %w", err)
	}
	err = JsonToData(stringToTasks, &tasks)
	if err != nil {
		return nil, fmt.Errorf("не удалось преобразовать задачи: %w", err)
	}
	return tasks, nil
}
