package task

import (
	"errors"
	"fmt"
	"time"
	"unicode/utf8"
)

// type-alias
type Status string

var (
	ErrInvalidStatus = errors.New("некорректный статус")
	ErrInvalidID     = errors.New("некорректный ID")
	ErrTaskNotFound  = errors.New("задача не найдена")
	ErrTaskTitle     = errors.New("пустое название")
)

const (
	StatusPending   Status = "pending"
	StatusProgress  Status = "in_progress"
	StatusCompleted Status = "completed"
)

// теги структур json или yaml задаются через тильда кавычки
// omitempty - пропустить если нету значения
// "-" - исключить вообще
type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created,omitempty"`
	CompletedAt *time.Time `json:"completed,omitempty"`
}

// это метод - функция с получателем (receiver)
// реализует интерфейс для вывода значения с преобразованием типов
// fmt.PrintLn(StatusPending) вызовет StatusPending.String() и выведет "pending"
func (status Status) String() string {
	return string(status)
}

func (status Status) Valid() bool {
	switch status {
	case StatusPending, StatusProgress, StatusCompleted:
		return true
	}
	return false
}

// NewTask - конструктор указывается через New префикс. У конструктора нет именованных аргументов -
// аргументы должны передаваться в том же порядке
// %w позволяет обернуть ошибку для error.Is() проверки
func NewTask(id int, title, description, status string) (*Task, error) {
	if !Status(status).Valid() {
		return nil, fmt.Errorf("ошибка валидации (%w): %s", ErrInvalidStatus, status)
	}
	// проверяем кол-во символов по Unicode через руны
	if utf8.RuneCountInString(title) <= 1 {
		return nil, fmt.Errorf("ошибка в названии задачи (%w)", ErrTaskTitle)
	}
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      Status(status),
		CreatedAt:   time.Now(),
	}, nil
}
