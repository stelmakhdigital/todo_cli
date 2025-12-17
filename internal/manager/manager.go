package manager

import (
	"fmt"
	"todo_cli/internal/render"
	"todo_cli/internal/task"
)

type Storage interface {
	Save(tasks []*task.Task, newFileName *string) error
	Load() ([]*task.Task, error)
}

type ManagerTasks interface {
	Show(id int) error
	Create(data map[string]string) (*int, error)
	Edit(id int, data map[string]string) error
	Start(id int) error
	Complete(id int) error
	Delete(id int) error
}

// добавим зависимость для использования во внутренних методах
type Manager struct {
	store  Storage
	filter Filter
}

// конструктор
func NewManager(s Storage, f Filter) *Manager {
	return &Manager{
		store:  s,
		filter: f,
	}
}

func hasKeys(data map[string]string, keys ...string) bool {
	for _, key := range keys {
		if _, ok := data[key]; !ok {
			return false
		}
	}
	return true
}

func editTask(m *Manager, id int, data map[string]string) error {
	tasks, err := m.store.Load()
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	indexTask := m.filter.GetIndexByID(tasks, id)

	if indexTask == nil {
		fmt.Printf("задачи с #%d не существует...\n", id)
		return fmt.Errorf("не найдена задача с #%d", id)
	}

	if title, ok := data["title"]; ok {
		tasks[*indexTask].Title = title
	}

	if description, ok := data["description"]; ok {
		tasks[*indexTask].Description = description
	}

	if status, ok := data["status"]; ok {
		if !task.Status(string(status)).Valid() {
			return fmt.Errorf("неверный статус задачи: %v", status)
		}
		tasks[*indexTask].Status = task.Status(status)
	}

	err = m.store.Save(tasks, nil)
	if err != nil {
		return fmt.Errorf("ошибка при записи: %w", err)
	}

	return nil
}

func (m *Manager) Create(data map[string]string) (*int, error) {
	var idTask int = 0
	tasks, err := m.store.Load()
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании: %w", err)
	}
	for _, task := range tasks {
		if idTask < task.ID {
			idTask = task.ID
		}
	}
	idTask += 1

	if hasKeys(data, "Title", "Description") {
		newTask, err := task.NewTask(idTask, data["Title"], data["Description"], task.StatusPending.String())
		if err != nil {
			return nil, fmt.Errorf("ошибка при создании задачи: %w", err)
		}
		tasks = append(tasks, newTask)
		m.store.Save(tasks, nil)
		fmt.Printf("Задача #%d добавлена успешно\n", idTask)
		return &idTask, nil
	}

	return nil, fmt.Errorf("получены некорректные данные при создании задачи: %v", data)

}

func (m *Manager) Start(id int) error {
	data := map[string]string{"status": task.StatusProgress.String()}
	err := editTask(m, id, data)
	if err != nil {
		return fmt.Errorf("не удалось создать задачу: %w", err)
	}
	return nil
}

// Complete - реализует изменение статуса через внутреннюю функцию [changeStatus](#changeStatus)
//
// Менеджер: [Manager](#Manager)
func (m *Manager) Complete(id int) error {
	data := map[string]string{"status": task.StatusCompleted.String()}
	err := editTask(m, id, data)
	if err != nil {
		return fmt.Errorf("не удалось завершить задачу: %w", err)
	}
	return nil
}

func (m *Manager) Edit(id int, data map[string]string) error {
	err := editTask(m, id, data)
	if err != nil {
		return fmt.Errorf("не удалось отредактировать задачу: %w", err)
	}
	return nil
}

func (m *Manager) Delete(id int) error {
	tasks, err := m.store.Load()
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	indexTask := m.filter.GetIndexByID(tasks, id)

	if indexTask == nil {
		fmt.Printf("задачи с #%d не существует...\n", id)
		return fmt.Errorf("не найдена задача с #%d", id)
	}

	// удаляем из слайса tasks элемент под индексом indexTask
	if len(tasks) != 0 && *indexTask < len(tasks) {
		tasks = append(tasks[:*indexTask], tasks[*indexTask+1:]...)
	}

	err = m.store.Save(tasks, nil)
	if err != nil {
		return fmt.Errorf("ошибка при записи: %w", err)
	}

	return nil
}

func (m *Manager) Show(id int) error {
	tasks, err := m.store.Load()
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	indexTask := m.filter.GetIndexByID(tasks, id)

	if indexTask == nil {
		fmt.Printf("задачи с #%d не существует...\n", id)
		return fmt.Errorf("не найдена задача с #%d", id)
	}

	task := tasks[*indexTask]
	render.RenderOne(task)

	return nil
}
