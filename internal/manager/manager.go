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
	List(status task.Status) error
	Create(data map[string]string) (*int, error)
	Edit(id int, data map[string]string) error
	Start(id int) error
	Complete(id int) error
	Delete(id int) error
	Stats() error
	Search(word string) error
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

func editTask(m *Manager, tasks []*task.Task, id int, data map[string]string) (*int, error) {
	indexTask := m.filter.GetIndexByID(tasks, id)
	if indexTask == nil {
		return nil, fmt.Errorf("не найдена задача с #%d", id)
	}
	if title, ok := data["title"]; ok {
		tasks[*indexTask].Title = title
	}
	if description, ok := data["description"]; ok {
		tasks[*indexTask].Description = description
	}
	if status, ok := data["status"]; ok {
		if !task.Status(string(status)).Valid() {
			return nil, fmt.Errorf("неверный статус задачи: %v", status)
		}
		tasks[*indexTask].Status = task.Status(status)
	}
	err := m.store.Save(tasks, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при записи: %w", err)
	}
	return indexTask, nil
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

	if hasKeys(data, "title", "description") {
		newTask, err := task.NewTask(idTask, data["title"], data["description"], task.StatusPending.String())
		if err != nil {
			return nil, fmt.Errorf("ошибка при создании задачи: %w", err)
		}
		tasks = append(tasks, newTask)
		m.store.Save(tasks, nil)
		render.RenderDetailed(newTask)
		return &idTask, nil
	}
	return nil, fmt.Errorf("получены некорректные данные при создании задачи: %v", data)
}

func (m *Manager) Start(id int) error {
	data := map[string]string{"status": task.StatusProgress.String()}
	tasks, err := m.store.Load()
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	indexTask, err := editTask(m, tasks, id, data)
	if err != nil {
		return fmt.Errorf("не удалось создать задачу: %w", err)
	}
	render.RenderDetailed(tasks[*indexTask])
	return nil
}

// Complete - реализует изменение статуса через внутреннюю функцию [changeStatus](#changeStatus)
//
// Менеджер: [Manager](#Manager)
func (m *Manager) Complete(id int) error {
	data := map[string]string{"status": task.StatusCompleted.String()}
	tasks, err := m.store.Load()
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	indexTask, err := editTask(m, tasks, id, data)
	if err != nil {
		return fmt.Errorf("не удалось завершить задачу: %w", err)
	}
	render.RenderDetailed(tasks[*indexTask])
	return nil
}

func (m *Manager) Edit(id int, data map[string]string) error {
	tasks, err := m.store.Load()
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	indexTask, err := editTask(m, tasks, id, data)
	if err != nil {
		return fmt.Errorf("не удалось отредактировать задачу: %w", err)
	}
	render.RenderDetailed(tasks[*indexTask])
	return nil
}

func (m *Manager) Delete(id int) error {
	tasks, err := m.store.Load()
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	indexTask := m.filter.GetIndexByID(tasks, id)

	if indexTask == nil {
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
		return fmt.Errorf("не найдена задача с #%d", id)
	}
	render.RenderDetailed(tasks[*indexTask])
	return nil
}

func (m *Manager) List(status string) error {
	tasks, err := m.store.Load()
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	if status == "all" {
		render.RenderList(tasks)
	} else {
		if !task.Status(status).Valid() {
			return fmt.Errorf("передан некорректный статус для фильтрации: %s", status)
		}
		tasksFiltered, err := m.filter.GetTasksByStatus(tasks, task.Status(status))
		if err != nil {
			return fmt.Errorf("ошибка при фильтрации задач: %w", err)
		}
		render.RenderList(tasksFiltered)
	}
	return nil
}

func (m *Manager) Stats() error {
	tasks, err := m.store.Load()
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	data := m.filter.GetStatsTasksByStatus(tasks)
	render.RenderMap(data)
	return nil
}

func (m *Manager) Search(word string) error {
	tasks, err := m.store.Load()
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	foundTasks := m.filter.GetTasksBySearchWord(tasks, word)

	if len(foundTasks) >= 1 {
		render.RenderList(foundTasks)
	} else {
		return fmt.Errorf("задачи по фразе %s - не найдены", word)
	}
	return nil
}
