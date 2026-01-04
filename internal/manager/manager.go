package manager

import (
	"fmt"
	"todo_cli/internal/render"
	"todo_cli/internal/task"
)

type Storage interface {
	Save(tasks []*task.Task, newFileName *string) error
	Load(differentFileName *string) ([]*task.Task, error)
	Clear(differentFileName *string) error
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
	render render.Render
}

// конструктор
func NewManager(s Storage, f Filter, r render.Render) *Manager {
	return &Manager{
		store:  s,
		filter: f,
		render: r,
	}
}

// hasKeys проверяет наличие всех указанных ключей в карте data.
// Возвращает true, если все ключи присутствуют, иначе false.
func hasKeys(data map[string]string, keys ...string) bool {
	for _, key := range keys {
		if _, ok := data[key]; !ok {
			return false
		}
	}
	return true
}

// editTask изменяет поля задачи по её ID и сохраняет изменения в хранилище.
// Принимает менеджер, список задач, ID задачи и карту с новыми данными (title, description, status).
// Возвращает индекс изменённой задачи или ошибку, если задача не найдена или данные невалидны.
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

// Create создаёт новую задачу со статусом "pending".
// Принимает карту data с обязательными ключами "title" и "description".
// Автоматически назначает новый уникальный ID (максимальный существующий + 1).
// Сохраняет задачу в хранилище и выводит детальную информацию.
// Возвращает указатель на ID созданной задачи или ошибку при невалидных данных.
func (m *Manager) Create(data map[string]string) (*int, error) {
	var idTask int = 0
	tasks, err := m.store.Load(nil)
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
		m.render.RenderDetailed(newTask)
		return &idTask, nil
	}
	return nil, fmt.Errorf("получены некорректные данные при создании задачи: %v", data)
}

// Start переводит задачу в статус "in_progress" (в работе).
// Находит задачу по ID, изменяет её статус и сохраняет изменения.
// Выводит детальную информацию об обновлённой задаче.
// Возвращает ошибку, если задача не найдена или произошла ошибка при сохранении.
func (m *Manager) Start(id int) error {
	data := map[string]string{"status": task.StatusProgress.String()}
	tasks, err := m.store.Load(nil)
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	indexTask, err := editTask(m, tasks, id, data)
	if err != nil {
		return fmt.Errorf("не удалось создать задачу: %w", err)
	}
	m.render.RenderDetailed(tasks[*indexTask])
	return nil
}

// Complete переводит задачу в статус "completed" (выполнена).
// Находит задачу по ID, изменяет её статус на завершённый и сохраняет изменения.
// Выводит детальную информацию об обновлённой задаче.
// Возвращает ошибку, если задача не найдена или произошла ошибка при сохранении.
func (m *Manager) Complete(id int) error {
	data := map[string]string{"status": task.StatusCompleted.String()}
	tasks, err := m.store.Load(nil)
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	indexTask, err := editTask(m, tasks, id, data)
	if err != nil {
		return fmt.Errorf("не удалось завершить задачу: %w", err)
	}
	m.render.RenderDetailed(tasks[*indexTask])
	return nil
}

// Edit изменяет данные существующей задачи.
// Принимает ID задачи и карту data с новыми значениями (title, description, status).
// Можно изменять как одно поле, так и несколько одновременно.
// Выводит детальную информацию об обновлённой задаче.
// Возвращает ошибку, если задача не найдена, статус невалиден или ошибка при сохранении.
func (m *Manager) Edit(id int, data map[string]string) error {
	tasks, err := m.store.Load(nil)
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	indexTask, err := editTask(m, tasks, id, data)
	if err != nil {
		return fmt.Errorf("не удалось отредактировать задачу: %w", err)
	}
	m.render.RenderDetailed(tasks[*indexTask])
	return nil
}

// Delete удаляет задачу по её ID.
// Находит задачу в списке, удаляет её из слайса и сохраняет изменения в хранилище.
// Возвращает ошибку, если задача не найдена или произошла ошибка при сохранении.
func (m *Manager) Delete(id int) error {
	tasks, err := m.store.Load(nil)
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

// Show выводит детальную информацию о конкретной задаче.
// Загружает задачи из хранилища, находит задачу по ID и отображает её данные.
// Возвращает ошибку, если задача не найдена или произошла ошибка при загрузке.
func (m *Manager) Show(id int) error {
	tasks, err := m.store.Load(nil)
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	indexTask := m.filter.GetIndexByID(tasks, id)
	if indexTask == nil {
		return fmt.Errorf("не найдена задача с #%d", id)
	}
	m.render.RenderDetailed(tasks[*indexTask])
	return nil
}

// List выводит список задач с опциональной фильтрацией по статусу.
// Если status = "all", выводит все задачи без фильтрации.
// Иначе фильтрует задачи по указанному статусу (pending, in_progress, completed).
// Возвращает ошибку, если передан некорректный статус или ошибка при загрузке.
func (m *Manager) List(status string) error {
	tasks, err := m.store.Load(nil)
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	if status == "all" {
		m.render.RenderList(tasks)
	} else {
		if !task.Status(status).Valid() {
			return fmt.Errorf("передан некорректный статус для фильтрации: %s", status)
		}
		tasksFiltered, err := m.filter.GetTasksByStatus(tasks, task.Status(status))
		if err != nil {
			return fmt.Errorf("ошибка при фильтрации задач: %w", err)
		}
		m.render.RenderList(tasksFiltered)
	}
	return nil
}

// Stats выводит статистику по задачам.
// Собирает и отображает количество задач по каждому статусу и общее количество.
// Формат вывода: всего задач, выполнено, в работе, ожидает.
// Возвращает ошибку, если произошла ошибка при загрузке задач.
func (m *Manager) Stats() error {
	tasks, err := m.store.Load(nil)
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	data := m.filter.GetStatsTasksByStatus(tasks)
	m.render.RenderMap(data)
	return nil
}

// Search выполняет поиск задач по ключевому слову.
// Ищет совпадения в заголовке и описании задач (регистронезависимый поиск).
// Выводит список найденных задач.
// Возвращает ошибку, если задачи не найдены или произошла ошибка при загрузке.
func (m *Manager) Search(word string) error {
	tasks, err := m.store.Load(nil)
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	foundTasks := m.filter.GetTasksBySearchWord(tasks, word)

	if len(foundTasks) >= 1 {
		m.render.RenderList(foundTasks)
	} else {
		return fmt.Errorf("задачи по фразе %s - не найдены", word)
	}
	return nil
}
